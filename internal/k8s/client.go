package k8s

import (
	"context"
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/tracing"
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	for _, obj := range objects.ToBeUpdated() {
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
	defer span.Finish()

	err := d.Status().Update(ctx, horus)
	if err != nil {
		return span.HandleError(err)
	}
	span.Logger().Info(fmt.Sprintf("%T %q status updated", horus, horus.GetName()))
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

func (d *Client) ListIngressByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]networking.Ingress, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
		client.MatchingLabels(owner.GetDefaultLabel()),
	}
	list := &networking.IngressList{}
	if err := d.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s ingress: %w", owner.GetName(), err))
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
	log := span.Logger()

	deleteOptions := []client.DeleteOption{client.PropagationPolicy(meta.DeletePropagationBackground)}
	if err := d.Delete(ctx, obj, deleteOptions...); err != nil {
		return span.HandleError(fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err))
	}
	log.Info(fmt.Sprintf("%T %q deleted", obj, obj.GetName()))
	return nil
}

func (d *Client) update(ctx context.Context, obj client.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	if err := d.Update(ctx, obj); err != nil {
		return span.HandleError(fmt.Errorf("failed to update %T %q: %w", obj, obj.GetName(), err))
	}
	log.Info(fmt.Sprintf("%T %q updated", obj, obj.GetName()))
	return nil
}

func (d *Client) create(ctx context.Context, obj client.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	if err := d.Create(ctx, obj); err != nil {
		return span.HandleError(fmt.Errorf("failed to create %T %q: %w", obj, obj.GetName(), err))
	}
	log.Info(fmt.Sprintf("%T %q created", obj, obj.GetName()))
	return nil
}
