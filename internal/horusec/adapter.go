package horusec

import (
	"context"
	"fmt"

	autoScalingV2beta2 "k8s.io/api/autoscaling/v2beta2"

	"github.com/ZupIT/horusec-operator/internal/horusec/vulnerability"

	"github.com/ZupIT/horusec-operator/internal/horusec/analytic"
	"github.com/ZupIT/horusec-operator/internal/horusec/api"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
	"github.com/ZupIT/horusec-operator/internal/horusec/core"
	"github.com/ZupIT/horusec-operator/internal/horusec/manager"
	"github.com/ZupIT/horusec-operator/internal/horusec/messages"
	"github.com/ZupIT/horusec-operator/internal/horusec/webhook"

	coreV1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/inventory"
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type Adapter struct {
	scheme *runtime.Scheme
	svc    *Service

	resource *v2alpha1.HorusecPlatform
}

//nolint:funlen
func (a *Adapter) EnsureDeployments(ctx context.Context) (*operation.Result, error) {
	desired := a.listOfDeployments()
	for index := range desired {
		deps, err := a.svc.ListDeployments(ctx, a.resource.Namespace, desired[index].ObjectMeta.Labels)
		if err != nil {
			return nil, err
		}
		if err = controllerutil.SetControllerReference(a.resource, &desired[index], a.scheme); err != nil {
			return nil, fmt.Errorf("failed to set Deployment %q owner reference: %v", desired[index].GetName(), err)
		}
		inv := inventory.ForDeployments(deps.Items, desired)
		err = a.svc.Apply(ctx, inv)
		if err != nil {
			return nil, err
		}
	}
	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureDatabaseConnectivity(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureBrokerConnectivity(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureSMTPConnectivity(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureDatabaseMigrations(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureServices(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureServicesAccounts(ctx context.Context) (*operation.Result, error) {
	servicesAccounts, err := a.svc.ListAuthServiceAccounts(ctx, a.resource.GetNamespace())
	if err != nil {
		return nil, err
	}

	desired := auth.NewServiceAccount(a.resource)
	if err = controllerutil.SetControllerReference(a.resource, desired, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Service Account %q owner reference: %v", desired.GetName(), err)
	}

	inv := inventory.ForServiceAccount(servicesAccounts.Items, []coreV1.ServiceAccount{*desired})
	if err := a.svc.Apply(ctx, inv); err != nil {
		return nil, err
	}
	return operation.ContinueProcessing()
}

//nolint:funlen
func (a *Adapter) EnsureAutoscaling(ctx context.Context) (*operation.Result, error) {
	desired := a.listOfAutoscaling()
	for index := range desired {
		deps, err := a.svc.ListAutoscaling(ctx, a.resource.Namespace, desired[index].Labels)
		if err != nil {
			return nil, err
		}
		horizontalScaler := []autoScalingV2beta2.HorizontalPodAutoscaler{}
		if desired[index] != nil {
			horizontalScaler = append(horizontalScaler, *desired[index])
		}
		inv := inventory.ForHorizontalPodAutoscaling(deps.Items, horizontalScaler)
		err = a.svc.Apply(ctx, inv)
		if err != nil {
			return nil, err
		}
	}
	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureIngressRules(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureEverythingIsRunning(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) listOfDeployments() []appsv1.Deployment {
	return []appsv1.Deployment{
		auth.NewDeployment(a.resource),
		core.NewDeployment(a.resource),
		api.NewDeployment(a.resource),
		messages.NewDeployment(a.resource),
		analytic.NewDeployment(a.resource),
		manager.NewDeployment(a.resource),
		vulnerability.NewDeployment(a.resource),
		webhook.NewDeployment(a.resource),
	}
}

func (a *Adapter) listOfAutoscaling() []*autoScalingV2beta2.HorizontalPodAutoscaler {
	return []*autoScalingV2beta2.HorizontalPodAutoscaler{
		auth.NewAutoscaling(a.resource),
		core.NewAutoscaling(a.resource),
		api.NewAutoscaling(a.resource),
		messages.NewAutoscaling(a.resource),
		analytic.NewAutoscaling(a.resource),
		manager.NewAutoscaling(a.resource),
		vulnerability.NewAutoscaling(a.resource),
		webhook.NewAutoscaling(a.resource),
	}
}
