package horusec

import (
	"context"

	"github.com/ZupIT/horusec-operator/controllers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AdapterFactory struct {
	builder ResourceBuilder
	client  KubernetesClient
}

func NewAdapterFactory(builder ResourceBuilder, client KubernetesClient) *AdapterFactory {
	return &AdapterFactory{builder: builder, client: client}
}

func (a *AdapterFactory) CreateHorusecPlatformAdapter(
	ctx context.Context, key client.ObjectKey) (controllers.HorusecPlatformAdapter, error) {
	resource, err := a.client.GetHorus(ctx, key)
	if err != nil {
		return nil, err
	}
	return &Adapter{builder: a.builder, client: a.client, resource: resource}, err
}
