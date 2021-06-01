package horusec

import (
	"context"
	"fmt"

	autoScalingV2beta2 "k8s.io/api/autoscaling/v2beta2"

	"k8s.io/api/networking/v1beta1"

	"github.com/ZupIT/horusec-operator/internal/horusec/vulnerability"

	"github.com/ZupIT/horusec-operator/internal/horusec/ingress"

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

func (a *Adapter) EnsureInitialization(ctx context.Context) (*operation.Result, error) {
	if a.resource.Status.Conditions != nil {
		return operation.ContinueProcessing()
	}
	a.resource.Status.Conditions = []v2alpha1.Condition{}
	a.resource.Status.State = v2alpha1.StatusPending
	err := a.svc.UpdateHorusecPlatformStatus(ctx, a.resource)
	if err != nil {
		return operation.RequeueWithError(err)
	}
	return operation.StopProcessing()
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

//nolint:funlen // improve in the future
func (a *Adapter) EnsureServices(ctx context.Context) (*operation.Result, error) {
	existing, err := a.svc.ListServices(ctx, a.resource.GetNamespace(),
		a.resource.GetName(), map[string]string{"app.kubernetes.io/managed-by": "horusec"})
	if err != nil {
		return nil, err
	}

	desired := a.listServices()
	for index := range desired {
		if err := a.ensureServices(&desired[index]); err != nil {
			return nil, err
		}
	}

	inv := inventory.ForService(existing.Items, desired)
	if err := a.svc.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) listServices() []coreV1.Service {
	return []coreV1.Service{
		analytic.NewService(a.resource),
		api.NewService(a.resource),
		auth.NewService(a.resource),
		core.NewService(a.resource),
		manager.NewService(a.resource),
		messages.NewService(a.resource),
		vulnerability.NewService(a.resource),
		webhook.NewService(a.resource),
	}
}

func (a *Adapter) ensureServices(desired *coreV1.Service) error {
	if err := controllerutil.SetControllerReference(a.resource, desired, a.scheme); err != nil {
		return fmt.Errorf("failed to set service %q owner reference: %v", desired.GetName(), err)
	}

	return nil
}

//nolint:funlen // to improve in the future
func (a *Adapter) EnsureIngressRules(ctx context.Context) (*operation.Result, error) {
	existing, err := a.svc.ListIngress(ctx, a.resource.GetNamespace(),
		a.resource.GetName(), map[string]string{"app.kubernetes.io/managed-by": "horusec"})
	if err != nil {
		return nil, err
	}

	desired := ingress.NewIngress(a.resource)
	if err := controllerutil.SetControllerReference(a.resource, desired, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set ingress %q owner reference: %v", desired.GetName(), err)
	}

	inv := inventory.ForIngresses(existing.Items, []v1beta1.Ingress{*desired})
	if err := a.svc.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureEverythingIsRunning(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

//nolint // to improve in the future
func (a *Adapter) EnsureServiceAccounts(ctx context.Context) (*operation.Result, error) {
	existing, err := a.svc.ListServiceAccounts(ctx, a.resource.GetNamespace(),
		a.resource.GetName(), map[string]string{"app.kubernetes.io/managed-by": "horusec"})
	if err != nil {
		return nil, err
	}

	desired := a.listServiceAccounts()
	for index := range desired {
		desiredItem := &desired[index]
		if err := controllerutil.SetControllerReference(a.resource, desiredItem, a.scheme); err != nil {
			return nil, fmt.Errorf("failed to set service account %q owner reference: %v", desiredItem.GetName(), err)
		}
	}

	inv := inventory.ForServiceAccount(existing.Items, desired)
	if err := a.svc.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) listServiceAccounts() []coreV1.ServiceAccount {
	return []coreV1.ServiceAccount{
		analytic.NewServiceAccount(a.resource),
		api.NewServiceAccount(a.resource),
		auth.NewServiceAccount(a.resource),
		core.NewServiceAccount(a.resource),
		manager.NewServiceAccount(a.resource),
		messages.NewServiceAccount(a.resource),
		vulnerability.NewServiceAccount(a.resource),
		webhook.NewServiceAccount(a.resource),
	}
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
