package horusec

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/internal/inventory"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

type Service struct {
	client k8s.Client
	log    logr.Logger
	scheme *runtime.Scheme
}

func NewService(client k8s.Client, scheme *runtime.Scheme) *Service {
	return &Service{
		client: client,
		log:    ctrl.Log.WithName("services").WithName("Horusec"),
		scheme: scheme,
	}
}

func (s *Service) LookupResourceAdapter(ctx context.Context, key k8s.ObjectKey) (*Adapter, error) {
	r := new(v2alpha1.HorusecPlatform)
	err := s.client.Get(ctx, key, r)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return &Adapter{resource: r, svc: s}, nil
}

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
		if err := s.client.Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err)
		}
		s.log.Info(fmt.Sprintf("%T %q deleted", obj, obj.GetName()))
	}

	return nil
}

func (s *Service) ListAuthDeployments(ctx context.Context, namespace string) (*v1.DeploymentList, error) {
	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(map[string]string{
			"app.kubernetes.io/name":       "auth",
			"app.kubernetes.io/managed-by": "horusec",
		}),
	}
	list := &v1.DeploymentList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Auth deployments: %w", err)
	}
	return list, nil
}
