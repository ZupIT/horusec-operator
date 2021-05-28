package horusec

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
	"github.com/ZupIT/horusec-operator/internal/inventory"
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type Adapter struct {
	resource *v2alpha1.HorusecPlatform
	svc      *Service
}

func (a *Adapter) EnsureAuthDeployments(ctx context.Context) (*operation.Result, error) {
	r := a.resource
	desired := auth.NewDeployment(r)
	if err := controllerutil.SetControllerReference(r, desired, a.svc.scheme); err != nil {
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
