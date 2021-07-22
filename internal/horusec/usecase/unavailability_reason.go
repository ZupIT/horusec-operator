package usecase

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/api/v2alpha1/condition"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/internal/tracing"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type UnavailabilityReason struct {
	logs   KubernetesLogs
	client KubernetesClient
}

func NewUnavailabilityReason(logs KubernetesLogs, client KubernetesClient) *UnavailabilityReason {
	return &UnavailabilityReason{logs: logs, client: client}
}

func (u *UnavailabilityReason) EnsureUnavailabilityReason(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if !resource.AnyStatusConditionFalseOrUnknown() {
		return operation.RequeueAfter(30*time.Second, nil)
	}

	containers, err := u.listCrashingContainers(ctx, resource)
	if err != nil {
		return nil, err
	}

	changed := false
	for _, container := range containers.FilterCrashLoopBackOff() {
		contype := condition.ComponentMap[container.component]
		logs, err := u.logs.PreviousContainerLogs(ctx, container.pod, container.name)
		if err != nil {
			return nil, err
		}

		if msg := u.searchForDatabaseErrors(bytes.NewReader(logs)); msg != "" &&
			resource.SetStatusCondition(condition.False(contype, condition.DatabaseReason(msg))) {
			changed = true
		}

		if msg := u.searchForBrokerErrors(bytes.NewReader(logs)); msg != "" &&
			resource.SetStatusCondition(condition.False(contype, condition.BrokerReason(msg))) {
			changed = true
		}
	}

	for _, container := range containers.FilterConfigError() {
		var reason *condition.Reason
		contype := condition.ComponentMap[container.component]
		msg := container.StateWaitingMessage()

		if container.HasSecretError() {
			reason = condition.SecretReason(msg)
		} else {
			reason = condition.ConfigReason(msg)
		}

		if resource.SetStatusCondition(condition.False(contype, reason)) {
			changed = true
		}
	}

	if changed {
		return operation.RequeueOnErrorOrStop(u.client.UpdateHorusStatus(ctx, resource))
	}

	return operation.RequeueAfter(10*time.Second, nil)
}

func (u *UnavailabilityReason) listCrashingContainers(ctx context.Context, resource *v2alpha1.HorusecPlatform) (containerStatuses, error) {
	status, err := u.podStatusOf(ctx, resource)
	if err != nil {
		return nil, err
	}

	crashingContainers := make([]*containerStatus, 0, 0)
	for component, conditionType := range condition.ComponentMap {
		if !resource.IsStatusConditionTrue(conditionType) {
			containers := status.OfComponent(component).CrashingContainers()
			crashingContainers = append(crashingContainers, containers...)
		}
	}

	return crashingContainers, nil
}

func (u *UnavailabilityReason) searchForDatabaseErrors(logs io.Reader) string {
	scanner := bufio.NewScanner(logs)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "{ERROR_DATABASE}") {
			compRegEx := regexp.MustCompile(`error="(.*?)"`)
			match := compRegEx.FindStringSubmatch(text)
			if len(match) == 0 {
				return ""
			}
			return match[len(match)-1]
		}
	}

	return ""
}

func (u *UnavailabilityReason) searchForBrokerErrors(logs io.Reader) string {
	scanner := bufio.NewScanner(logs)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "{ERROR_BROKER}") {
			match := strings.Split(text, "panic: {ERROR_BROKER} ")
			if len(match) == 0 {
				return ""
			}
			return match[len(match)-1]
		}
	}

	return ""
}

func (u *UnavailabilityReason) podStatusOf(ctx context.Context, resource *v2alpha1.HorusecPlatform) (podStatuses, error) {
	pods, err := u.client.ListPodsByOwner(ctx, resource)
	if err != nil {
		return nil, err
	}

	items := make([]*podStatus, 0, len(pods))
	for _, pod := range pods {
		item := pod
		items = append(items, &podStatus{item: &item})
	}
	return items, nil
}

type (
	podStatus       struct{ item *core.Pod }
	podStatuses     []*podStatus
	containerStatus struct {
		component string
		pod       types.NamespacedName
		name      string
		state     core.ContainerState
	}
	containerStatuses []*containerStatus
)

func newPodContainer(pod *core.Pod, container core.ContainerStatus) *containerStatus {
	return &containerStatus{
		component: pod.Labels["app.kubernetes.io/component"],
		pod:       types.NamespacedName{Namespace: pod.GetNamespace(), Name: pod.GetName()},
		name:      container.Name,
		state:     container.State,
	}
}

func (p *podStatus) IsCrashing() bool {
	for _, status := range p.item.Status.ContainerStatuses {
		container := newPodContainer(p.item, status)
		if container.StateWaitingReason() != "" {
			return true
		}
	}
	return false
}

func (p *podStatus) CrashingContainers() containerStatuses {
	containers := make([]*containerStatus, 0, 0)
	for _, status := range p.item.Status.ContainerStatuses {
		container := newPodContainer(p.item, status)
		if container.StateWaitingReason() != "" {
			containers = append(containers, container)
		}
	}
	return containers
}

func (p podStatuses) IsCrashing() bool {
	for _, status := range p {
		if status.IsCrashing() {
			return true
		}
	}
	return false
}

func (p podStatuses) OfComponent(component string) podStatuses {
	items := make([]*podStatus, 0)
	for _, status := range p {
		pod := status.item
		if c, ok := pod.Labels["app.kubernetes.io/component"]; ok && c == component {
			items = append(items, &podStatus{item: pod})
		}
	}
	return items
}

func (p podStatuses) CrashingContainers() containerStatuses {
	containers := make([]*containerStatus, 0, 0)
	for _, pod := range p {
		if pod.IsCrashing() {
			containers = append(containers, pod.CrashingContainers()...)
		}
	}
	return containers
}

func (p *containerStatus) StateWaitingReason() string {
	waiting := p.state.Waiting
	if waiting != nil && (waiting.Reason == "CrashLoopBackOff" || waiting.Reason == "CreateContainerConfigError") {
		return waiting.Reason
	}
	return ""
}

func (p *containerStatus) StateWaitingMessage() string {
	waiting := p.state.Waiting
	if waiting != nil {
		return waiting.Message
	}
	return ""
}

func (p *containerStatus) HasSecretError() bool {
	return regexp.
		MustCompile(`(?i)\b(secret)\b`).
		MatchString(p.StateWaitingMessage())
}

func (p containerStatuses) FilterCrashLoopBackOff() containerStatuses {
	containers := make([]*containerStatus, 0, 0)
	for _, container := range p {
		if container.StateWaitingReason() == "CrashLoopBackOff" {
			containers = append(containers, container)
		}
	}
	return containers
}

func (p containerStatuses) FilterConfigError() containerStatuses {
	containers := make([]*containerStatus, 0, 0)
	for _, container := range p {
		if container.StateWaitingReason() == "CreateContainerConfigError" {
			containers = append(containers, container)
		}
	}
	return containers
}
