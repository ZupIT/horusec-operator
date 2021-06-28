package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ZupIT/horusec-operator/internal/tracing"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/operation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UnavailabilityReason struct {
	client KubernetesClient
}

func NewUnavailabilityReason(client KubernetesClient) *UnavailabilityReason {
	return &UnavailabilityReason{client: client}
}

func (e *DeploymentsAvailability) EnsureUnavailabilityReason(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := span.Logger()
	defer span.Finish()

	for _, condition := range resource.Status.Conditions {
		if condition.Status == metav1.ConditionUnknown {
			log.Info(fmt.Sprintf("Condition of %q is Unknown", condition.Type))
		}
	}
	return operation.RequeueAfter(10*time.Second, nil)
}
