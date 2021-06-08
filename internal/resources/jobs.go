package resources

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/resources/analytic"
	"github.com/ZupIT/horusec-operator/internal/resources/migration"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) JobsFor(resource *v2alpha1.HorusecPlatform) ([]batchv1.Job, error) {
	mdesired := migration.NewJob(resource)
	if err := controllerutil.SetControllerReference(resource, &mdesired, b.scheme); err != nil {
		return nil, fmt.Errorf("failed to set job %q owner reference: %v", mdesired.GetName(), err)
	}
	adesired := analytic.NewJob(resource)
	if err := controllerutil.SetControllerReference(resource, &adesired, b.scheme); err != nil {
		return nil, fmt.Errorf("failed to set job %q owner reference: %v", adesired.GetName(), err)
	}

	return []batchv1.Job{mdesired, adesired}, nil
}
