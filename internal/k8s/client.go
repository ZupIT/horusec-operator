// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"context"
	"fmt"
	"strings"

	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/tracing"
)

type Client struct{ client.Client }

func NewClient(client client.Client) *Client {
	return &Client{Client: client}
}

func (d *Client) Apply(ctx context.Context, objects Objects) error {
	for _, obj := range objects.ToBeDeleted() {
		if err := d.delete(ctx, obj); err != nil {
			return err
		}
	}

	for _, obj := range objects.ToBeUpdated() {
		if err := d.update(ctx, obj); err != nil {
			return err
		}
	}

	for _, obj := range objects.ToBeCreated() {
		if err := d.create(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (d *Client) GetHorus(ctx context.Context, namespacedName types.NamespacedName) (*v2alpha1.HorusecPlatform, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	horus := new(v2alpha1.HorusecPlatform)
	if err := d.Get(ctx, namespacedName, horus); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to lookup resource: %w", err))
	}

	horus, err := v2alpha1.MergeWithDefaultValues(horus)
	if err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to merge default values: %w", err))
	}

	return horus, nil
}

func (d *Client) UpdateHorusStatus(ctx context.Context, horus *v2alpha1.HorusecPlatform) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := span.Logger()
	defer span.Finish()

	conditions := make([]interface{}, 0, 0)
	for _, condition := range horus.Status.Conditions {
		conditions = append(conditions, "condition."+condition.Type, string(condition.Status))
		span.SetTag("horus.condition."+condition.Type, condition.Status)
	}
	span.SetTag("horus.condition.state", horus.Status.State)

	err := d.Status().Update(ctx, horus)
	if err != nil {
		return span.HandleError(err)
	}
	log.V(1).
		WithValues("condition.state", horus.Status.State).
		WithValues(conditions...).
		Info(fmt.Sprintf("%T %q status updated", horus, horus.GetName()))
	return nil
}

func (d *Client) ListAutoscalingByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]autoscaling.HorizontalPodAutoscaler, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &autoscaling.HorizontalPodAutoscalerList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s Autoscaling: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) ListDeploymentsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]apps.Deployment, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &apps.DeploymentList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s deployments: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) ListIngressByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]networkingv1.Ingress, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &networkingv1.IngressList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s ingress: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) ListPodsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.Pod, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &core.PodList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s pods: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) ListJobsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]batch.Job, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &batch.JobList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s jobs: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) ListServiceAccountsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.ServiceAccount, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &core.ServiceAccountList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s service accounts: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) ListServicesByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.Service, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &core.ServiceList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s services: %w", owner.GetName(), err))
	}
	return list.Items, nil
}

func (d *Client) delete(ctx context.Context, obj client.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger().V(1)

	deleteOptions := []client.DeleteOption{client.PropagationPolicy(meta.DeletePropagationBackground)}
	if err := d.Delete(ctx, obj, deleteOptions...); err != nil {
		return span.HandleError(fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err))
	}
	log.Info(fmt.Sprintf("%s deleted", wrap(obj).String()))
	return nil
}

func (d *Client) update(ctx context.Context, obj client.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger().V(1)

	if err := d.Update(ctx, obj); err != nil {
		return span.HandleError(fmt.Errorf("failed to update %T %q: %w", obj, obj.GetName(), err))
	}
	log.Info(fmt.Sprintf("%s updated", wrap(obj).String()))
	return nil
}

func (d *Client) create(ctx context.Context, obj client.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger().V(1)

	if err := d.Create(ctx, obj); err != nil {
		return span.HandleError(fmt.Errorf("failed to create %T %q: %w", obj, obj.GetName(), err))
	}
	log.Info(fmt.Sprintf("%s created", wrap(obj).String()))
	return nil
}

type wrapper struct {
	obj client.Object
}

func wrap(obj client.Object) *wrapper {
	return &wrapper{obj: obj}
}

func (o *wrapper) kind() string {
	ss := strings.Split(fmt.Sprintf("%T", o.obj), ".")
	return ss[len(ss)-1]
}

func (o *wrapper) String() string {
	kind := o.kind()
	name := o.obj.GetName()
	return fmt.Sprintf("%s %q", strings.ToLower(kind), name)
}
