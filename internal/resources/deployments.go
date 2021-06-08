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
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) DeploymentsFor(resource *v2alpha1.HorusecPlatform) ([]appsv1.Deployment, error) {
	desired := b.listOfDeployments(resource)
	for index := range desired {
		if err := b.ensureDeployments(resource, &desired[index]); err != nil {
			return nil, err
		}
	}
	return desired, nil
}

func (b *Builder) ensureDeployments(resource *v2alpha1.HorusecPlatform, desired *appsv1.Deployment) error {
	if err := controllerutil.SetControllerReference(resource, desired, b.scheme); err != nil {
		return fmt.Errorf("failed to set service %q owner reference: %v", desired.GetName(), err)
	}

	return nil
}

func (b *Builder) listOfDeployments(resource *v2alpha1.HorusecPlatform) []appsv1.Deployment {
	deployments := []appsv1.Deployment{
		auth.NewDeployment(resource),
		core.NewDeployment(resource),
		api.NewDeployment(resource),
		analytic.NewDeployment(resource),
		manager.NewDeployment(resource),
		vulnerability.NewDeployment(resource),
		webhook.NewDeployment(resource),
	}
	msg := resource.GetMessagesComponent()
	if msg.Enabled {
		deployments = append(deployments, messages.NewDeployment(resource))
	}
	return deployments
}
