package resources

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/analytic"
	"github.com/ZupIT/horusec-operator/internal/horusec/api"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
	"github.com/ZupIT/horusec-operator/internal/horusec/core"
	"github.com/ZupIT/horusec-operator/internal/horusec/manager"
	"github.com/ZupIT/horusec-operator/internal/horusec/messages"
	"github.com/ZupIT/horusec-operator/internal/horusec/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/horusec/webhook"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) ServicesFor(resource *v2alpha1.HorusecPlatform) ([]coreV1.Service, error) {
	desired := b.listServices(resource)
	for index := range desired {
		if err := b.ensureServices(resource, &desired[index]); err != nil {
			return nil, err
		}
	}
	return desired, nil
}

func (b *Builder) ensureServices(resource *v2alpha1.HorusecPlatform, desired *coreV1.Service) error {
	if err := controllerutil.SetControllerReference(resource, desired, b.scheme); err != nil {
		return fmt.Errorf("failed to set service %q owner reference: %v", desired.GetName(), err)
	}

	return nil
}

func (b *Builder) listServices(resource *v2alpha1.HorusecPlatform) []coreV1.Service {
	services := []coreV1.Service{
		analytic.NewService(resource),
		api.NewService(resource),
		auth.NewService(resource),
		core.NewService(resource),
		manager.NewService(resource),
		vulnerability.NewService(resource),
		webhook.NewService(resource),
	}
	msg := resource.GetMessagesComponent()
	if msg.Enabled {
		services = append(services, messages.NewService(resource))
	}
	return services
}
