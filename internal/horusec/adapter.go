package horusec

import (
	"context"
	"fmt"

	"k8s.io/api/networking/v1beta1"

	"github.com/ZupIT/horusec-operator/internal/horusec/ingress"

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

	deps, err := a.svc.ListAuthDeployments(ctx, r.Namespace, auth.Labels)
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

func (a *Adapter) ensureServiceAccounts(desired *corev1.ServiceAccount) error {
	if err := controllerutil.SetControllerReference(a.resource, desired, a.scheme); err != nil {
		return fmt.Errorf("failed to set service account %q owner reference: %v", desired.GetName(), err)
	}

	return nil
}

//nolint:funlen // to improve in the future
func (a *Adapter) EnsureServiceAccounts(ctx context.Context) (*operation.Result, error) {
	existing, err := a.svc.ListServiceAccounts(ctx, a.resource.GetNamespace(),
		a.resource.GetName(), map[string]string{"app.kubernetes.io/managed-by": "horusec"})
	if err != nil {
		return nil, err
	}

	desired := a.listServiceAccounts()
	for index := range desired {
		if err := a.ensureServiceAccounts(&desired[index]); err != nil {
			return nil, err
		}
	}

	inv := inventory.ForServiceAccount(existing.Items, desired)
	if err := a.svc.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) listServiceAccounts() []corev1.ServiceAccount {
	return []corev1.ServiceAccount{
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

func (a *Adapter) EnsureAutoscalers(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

func (a *Adapter) EnsureHPA(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
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
