package horusec

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/ZupIT/horusec-operator/internal/horusec/analytic"
	"github.com/ZupIT/horusec-operator/internal/horusec/api"
	"github.com/ZupIT/horusec-operator/internal/horusec/core"
	"github.com/ZupIT/horusec-operator/internal/horusec/manager"
	"github.com/ZupIT/horusec-operator/internal/horusec/messages"
	"github.com/ZupIT/horusec-operator/internal/horusec/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/horusec/webhook"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
	"github.com/ZupIT/horusec-operator/internal/inventory"
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type Adapter struct {
	scheme *runtime.Scheme
	svc    *Service

	resource *v2alpha1.HorusecPlatform
}

//nolint:funlen // to improve in the future
func (a *Adapter) EnsureAuthDeployments(ctx context.Context) (*operation.Result, error) {
	r := a.resource
	desired := auth.NewDeployment(r)
	if err := controllerutil.SetControllerReference(r, desired, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Deployment %q owner reference: %v", desired.GetName(), err)
	}

	deps, err := a.svc.ListAuthDeployments(ctx, r.Namespace)
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

func (a *Adapter) EnsureDeployments(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureServices(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) ensureServiceAccounts(
	ctx context.Context, desired *corev1.ServiceAccount) (*operation.Result, error) {
	servicesAccounts, err := a.svc.ListServiceAccounts(ctx, a.resource.GetNamespace(), a.resource.GetName())
	if err != nil {
		return nil, err
	}

	if err := controllerutil.SetControllerReference(a.resource, desired, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Service Account %q owner reference: %v", desired.GetName(), err)
	}

	inv := inventory.ForServiceAccount(servicesAccounts.Items, []corev1.ServiceAccount{*desired})
	if err := a.svc.Apply(ctx, inv); err != nil {
		return nil, err
	}
	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureAnalyticServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, analytic.NewServiceAccount(a.resource))
}

//nolint:golint, stylecheck // no need to be API
func (a *Adapter) EnsureApiServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, api.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureAuthServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, auth.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureCoreServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, core.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureManagerServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, manager.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureMessagesServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, messages.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureVulnerabilityServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, vulnerability.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureWebhookServiceAccounts(ctx context.Context) (*operation.Result, error) {
	return a.ensureServiceAccounts(ctx, webhook.NewServiceAccount(a.resource))
}

func (a *Adapter) EnsureAutoscalers(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureHPA(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureIngressRules(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureEverythingIsRunning(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}
