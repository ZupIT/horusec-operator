package horusec

import (
	"context"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/inventory"
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type Adapter struct {
	builder ResourceBuilder
	client  KubernetesClient

	resource *v2alpha1.HorusecPlatform
}

func (a *Adapter) EnsureInitialization(ctx context.Context) (*operation.Result, error) {
	if a.resource.Status.Conditions != nil {
		return operation.ContinueProcessing()
	}
	a.resource.Status.Conditions = []v2alpha1.Condition{}
	a.resource.Status.State = v2alpha1.StatusPending
	err := a.client.UpdateHorusStatus(ctx, a.resource)
	if err != nil {
		return operation.RequeueWithError(err)
	}
	return operation.StopProcessing()
}

func (a *Adapter) EnsureDatabaseMigrations(ctx context.Context) (*operation.Result, error) {
	existing, err := a.client.ListJobsByOwner(ctx, a.resource)
	if err != nil {
		return nil, err
	}

	desired, err := a.builder.JobsFor(a.resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForJobs(existing, desired)
	if err := a.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

//nolint:funlen
func (a *Adapter) EnsureDeployments(ctx context.Context) (*operation.Result, error) {
	existing, err := a.client.ListDeploymentsByOwner(ctx, a.resource)
	if err != nil {
		return nil, err
	}

	desired, err := a.builder.DeploymentsFor(a.resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForDeployments(existing, desired)
	if err := a.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

//nolint
func (a *Adapter) EnsureAutoscaling(ctx context.Context) (*operation.Result, error) {
	existing, err := a.client.ListAutoscalingByOwner(ctx, a.resource)
	if err != nil {
		return nil, err
	}

	desired, err := a.builder.AutoscalingFor(a.resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForHorizontalPodAutoscaling(existing, desired)
	if err := a.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

//nolint:funlen // improve in the future
func (a *Adapter) EnsureServices(ctx context.Context) (*operation.Result, error) {
	existing, err := a.client.ListServicesByOwner(ctx, a.resource)
	if err != nil {
		return nil, err
	}

	desired, err := a.builder.ServicesFor(a.resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForService(existing, desired)
	if err := a.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

//nolint:funlen // to improve in the future
func (a *Adapter) EnsureIngressRules(ctx context.Context) (*operation.Result, error) {
	existing, err := a.client.ListIngressByOwner(ctx, a.resource)
	if err != nil {
		return nil, err
	}

	desiredList, err := a.builder.IngressFor(a.resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForIngresses(existing, desiredList)
	if err := a.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) EnsureEverythingIsRunning(ctx context.Context) (*operation.Result, error) {
	panic("implement me") // TODO
}

//nolint // to improve in the future
func (a *Adapter) EnsureServiceAccounts(ctx context.Context) (*operation.Result, error) {
	existing, err := a.client.ListServiceAccountsByOwner(ctx, a.resource)
	if err != nil {
		return nil, err
	}

	desired, err := a.builder.ServiceAccountsFor(a.resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForServiceAccount(existing, desired)
	if err := a.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}
