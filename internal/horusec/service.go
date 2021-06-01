package horusec

import (
	"context"
	"fmt"

	autoScalingV2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/api/networking/v1beta1"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/inventory"
)

type Service struct {
	client k8s.Client
	log    logr.Logger
}

func NewService(client k8s.Client) *Service {
	return &Service{
		client: client,
		log:    ctrl.Log.WithName("services").WithName("Horusec"),
	}
}

func (s *Service) LookupHorusecPlatform(ctx context.Context, key k8s.ObjectKey) (*v2alpha1.HorusecPlatform, error) {
	r := new(v2alpha1.HorusecPlatform)
	err := s.client.Get(ctx, key, r)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return r, nil
}

func (s *Service) UpdateHorusecPlatformStatus(ctx context.Context, resource *v2alpha1.HorusecPlatform) error {
	err := s.client.Status().Update(ctx, resource)
	if err != nil {
		return err
	}
	s.log.Info(fmt.Sprintf("%T %q status updated", resource, resource.GetName()))
	return nil
}

//nolint:funlen // to improve in the future
func (s *Service) Apply(ctx context.Context, inv inventory.Object) error {
	for _, obj := range inv.Create {
		if err := s.client.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create %T %q: %w", obj, obj.GetName(), err)
		}
		s.log.Info(fmt.Sprintf("%T %q created", obj, obj.GetName()))
	}

	for _, obj := range inv.Update {
		if err := s.client.Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update %T %q: %w", obj, obj.GetName(), err)
		}
		s.log.Info(fmt.Sprintf("%T %q updated", obj, obj.GetName()))
	}

	for _, obj := range inv.Delete {
		if err := s.client.Delete(ctx, obj); err != nil {
			return fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err)
		}
		s.log.Info(fmt.Sprintf("%T %q deleted", obj, obj.GetName()))
	}

	return nil
}

func (s *Service) ListDeployments(ctx context.Context,
	namespace string, matchingLabels map[string]string) (*v1.DeploymentList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(matchingLabels),
	}
	list := &v1.DeploymentList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list %s deployments: %w", matchingLabels["app.kubernetes.io/name"], err)
	}
	return list, nil
}

func (s *Service) ListAutoscaling(ctx context.Context,
	namespace string, matchingLabels map[string]string) (*autoScalingV2beta2.HorizontalPodAutoscalerList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(matchingLabels),
	}
	list := &autoScalingV2beta2.HorizontalPodAutoscalerList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list %s Autoscaling: %w", matchingLabels["app.kubernetes.io/name"], err)
	}
	return list, nil
}

func (s *Service) ListServiceAccounts(
	ctx context.Context, namespace, name string, labels map[string]string) (*core.ServiceAccountList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(labels),
	}

	list := &core.ServiceAccountList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list %s service accounts: %w", name, err)
	}

	return list, nil
}

func (s *Service) ListIngress(
	ctx context.Context, namespace, name string, labels map[string]string) (*v1beta1.IngressList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(labels),
	}

	list := &v1beta1.IngressList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list %s ingress: %w", name, err)
	}

	return list, nil
}
