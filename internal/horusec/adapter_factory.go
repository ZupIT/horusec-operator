package horusec

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/controllers"
)

type AdapterFactory struct {
	scheme *runtime.Scheme
	svc    *Service
}

func NewAdapterFactory(scheme *runtime.Scheme, svc *Service) *AdapterFactory {
	return &AdapterFactory{scheme: scheme, svc: svc}
}

func (a *AdapterFactory) CreateHorusecPlatformAdapter(ctx context.Context, key client.ObjectKey) (controllers.HorusecPlatformAdapter, error) {
	svc := a.svc
	scheme := a.scheme
	resource, err := svc.LookupHorusecPlatform(ctx, key)
	if err != nil {
		return nil, err
	}
	return &Adapter{scheme: scheme, svc: svc, resource: resource}, err
}
