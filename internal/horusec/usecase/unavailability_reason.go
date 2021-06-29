package usecase

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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
	log := span.Logger()
	defer span.Finish()

	if !resource.AnyStatusConditionFalseOrUnknown() {
		return operation.RequeueAfter(30*time.Second, nil)
	}

	containers, err := u.listCrashingContainers(ctx, resource)
	if err != nil {
		return nil, err
	}

	changed := false
	for _, container := range containers {
		logs, err := u.logs.PreviousContainerLogs(ctx, container.pod, container.name)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(logs)
		if msg := u.searchForDatabaseErrors(reader); msg != "" {
			log.V(0).
				WithValues("pod", container.pod).
				WithValues("container", container.name).
				Info("an error with database was found")

			conditionType := condition.ComponentMap[container.component]
			if resource.SetStatusCondition(condition.False(conditionType, condition.DatabaseReason(msg))) {
				changed = true
			}
		}
	}

	if changed {
		return operation.RequeueOnErrorOrStop(u.client.UpdateHorusStatus(ctx, resource))
	}

	return operation.RequeueAfter(10*time.Second, nil)
}

func (u *UnavailabilityReason) listCrashingContainers(ctx context.Context, resource *v2alpha1.HorusecPlatform) ([]*podContainer, error) {
	status, err := u.podStatusOf(ctx, resource)
	if err != nil {
		return nil, err
	}

	crashingContainers := make([]*podContainer, 0, 0)
	for component, conditionType := range condition.ComponentMap {
		if !resource.IsStatusConditionTrue(conditionType) {
			containers := status.ofComponent(component).CrashingContainers()
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
			var compRegEx = regexp.MustCompile(`error="(.*?)"`)
			match := compRegEx.FindStringSubmatch(text)
			if len(match) == 0 {
				return ""
			}
			return match[len(match)-1]
		}
	}

	return ""
}

func (u *UnavailabilityReason) podStatusOf(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*podStatuses, error) {
	pods, err := u.client.ListPodsByOwner(ctx, resource)
	if err != nil {
		return nil, err
	}

	items := make([]*podStatus, 0, len(pods))
	for _, pod := range pods {
		item := pod
		items = append(items, &podStatus{item: &item})
	}
	return &podStatuses{items: items}, nil
}

type podStatus struct{ item *core.Pod }

func (p podStatus) IsCrashing() bool {
	for _, status := range p.item.Status.ContainerStatuses {
		waiting := status.State.Waiting
		if waiting != nil && waiting.Reason == "CrashLoopBackOff" {
			return true
		}
	}
	return false
}

func (p podStatus) CrashingContainers() []*podContainer {
	containers := make([]*podContainer, 0, 0)
	for _, status := range p.item.Status.ContainerStatuses {
		waiting := status.State.Waiting
		if waiting != nil && waiting.Reason == "CrashLoopBackOff" {
			containers = append(containers, &podContainer{
				component: p.item.Labels["app.kubernetes.io/component"],
				pod:       types.NamespacedName{Namespace: p.item.GetNamespace(), Name: p.item.GetName()},
				name:      status.Name,
			})
		}
	}
	return containers
}

type podStatuses struct {
	items []*podStatus
}

func (p *podStatuses) IsCrashing() bool {
	for _, status := range p.items {
		if status.IsCrashing() {
			return true
		}
	}
	return false
}

func (p *podStatuses) String() string {
	pods := make([]string, 0, len(p.items))
	for _, pod := range p.items {
		phase := func() string {
			if pod.IsCrashing() {
				return "CrashLoopBackOff"
			}
			return string(pod.item.Status.Phase)
		}()
		name := pod.item.GetName()
		namespace := pod.item.GetNamespace()
		pods = append(pods, fmt.Sprintf("%s/%s[%s]", namespace, name, phase))
	}
	return strings.Join(pods, ", ")
}

func (p *podStatuses) ofComponent(component string) *podStatuses {
	items := make([]*podStatus, 0)
	for _, status := range p.items {
		pod := status.item
		if c, ok := pod.Labels["app.kubernetes.io/component"]; ok && c == component {
			items = append(items, &podStatus{item: pod})
		}
	}
	return &podStatuses{items: items}
}

func (p *podStatuses) CrashingContainers() []*podContainer {
	containers := make([]*podContainer, 0, 0)
	for _, pod := range p.items {
		if pod.IsCrashing() {
			containers = append(containers, pod.CrashingContainers()...)
		}
	}
	return containers
}

type podContainer struct {
	component string
	pod       types.NamespacedName
	name      string
}
