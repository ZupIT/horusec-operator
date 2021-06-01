package horusec

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
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

func (s *Service) ListAuthDeployments(ctx context.Context, namespace string) (*v1.DeploymentList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(auth.Labels),
	}
	list := &v1.DeploymentList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Auth deployments: %w", err)
	}
	return list, nil
}

func (s *Service) ListAuthServiceAccounts(ctx context.Context, namespace string) (*core.ServiceAccountList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(auth.Labels),
	}
	list := &core.ServiceAccountList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Auth service accounts: %w", err)
	}
	return list, nil
}
