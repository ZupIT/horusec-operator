package horusec

import (
	"context"

	"github.com/ZupIT/horusec-operator/controllers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AdapterFactory struct {
	builder ResourceBuilder
	svc     *Service
}

func NewAdapterFactory(builder ResourceBuilder, svc *Service) *AdapterFactory {
	return &AdapterFactory{builder: builder, svc: svc}
}

func (a *AdapterFactory) CreateHorusecPlatformAdapter(
	ctx context.Context, key client.ObjectKey) (controllers.HorusecPlatformAdapter, error) {
	svc := a.svc
	builder := a.builder
	resource, err := svc.LookupHorusecPlatform(ctx, key)
	if err != nil {
		return nil, err
	}
	return &Adapter{builder: builder, svc: svc, resource: resource}, err
}
