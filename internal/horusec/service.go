package horusec

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
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
