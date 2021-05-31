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

//nolint:funlen // to improve in the future
func (a *Adapter) ensureDeployments(
	ctx context.Context, desired *appsv1.Deployment) (*operation.Result, error) {
	if err := controllerutil.SetControllerReference(a.resource, desired, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Deployment %q owner reference: %v", desired.GetName(), err)
	}

	deps, err := a.svc.ListDeployments(ctx, a.resource.Namespace, desired.ObjectMeta.Labels)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForDeployments(deps.Items, []appsv1.Deployment{*desired})
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureAuthDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, auth.NewDeployment(a.resource))
}

func (a *Adapter) EnsureCoreDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, core.NewDeployment(a.resource))
}

func (a *Adapter) EnsureAPIDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, api.NewDeployment(a.resource))
}

func (a *Adapter) EnsureMessagesDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, messages.NewDeployment(a.resource))
}

func (a *Adapter) EnsureAnalyticDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, analytic.NewDeployment(a.resource))
}

func (a *Adapter) EnsureManagerDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, manager.NewDeployment(a.resource))
}

func (a *Adapter) EnsureVulnerabilityDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, vulnerability.NewDeployment(a.resource))
}

func (a *Adapter) EnsureWebhookDeployments(ctx context.Context) (*operation.Result, error) {
	return a.ensureDeployments(ctx, webhook.NewDeployment(a.resource))
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

// nolint:funlen
func (a *Adapter) ensureAutoscaling(ctx context.Context,
	desired *autoScalingV2beta2.HorizontalPodAutoscaler) (*operation.Result, error) {
	deps, err := a.svc.ListAutoscaling(ctx, a.resource.Namespace, desired.Labels)
	if err != nil {
		return nil, err
	}
	horizontalScaler := []autoScalingV2beta2.HorizontalPodAutoscaler{}
	if desired != nil {
		horizontalScaler = append(horizontalScaler, *desired)
	}
	inv := inventory.ForHorizontalPodAutoscaling(deps.Items, horizontalScaler)
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureAuthAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, auth.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureCoreAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, core.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureAPIAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, api.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureMessagesAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, messages.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureAnalyticAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, analytic.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureManagerAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, manager.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureVulnerabilityAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, vulnerability.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureWebhookAutoscaling(ctx context.Context) (*operation.Result, error) {
	return a.ensureAutoscaling(ctx, webhook.NewAutoscaling(a.resource))
}

func (a *Adapter) EnsureIngressRules(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureEverythingIsRunning(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}
