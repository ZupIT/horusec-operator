package resources

import (
	"fmt"
	"reflect"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/analytic"
	"github.com/ZupIT/horusec-operator/internal/horusec/api"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
	"github.com/ZupIT/horusec-operator/internal/horusec/core"
	"github.com/ZupIT/horusec-operator/internal/horusec/manager"
	"github.com/ZupIT/horusec-operator/internal/horusec/messages"
	"github.com/ZupIT/horusec-operator/internal/horusec/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/horusec/webhook"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) AutoscalingFor(resource *v2alpha1.HorusecPlatform) ([]autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	result := []autoscalingv2beta2.HorizontalPodAutoscaler{}

	desired := b.listOfAutoscaling(resource)
	for index := range desired {
		if reflect.ValueOf(desired[index]).IsZero() {
			continue
		}

		if err := b.ensureAutoscaling(resource, &desired[index]); err != nil {
			return nil, err
		}

		result = append(result, desired[index])
	}

	return result, nil
}

func (b *Builder) listOfAutoscaling(resource *v2alpha1.HorusecPlatform) []autoscalingv2beta2.HorizontalPodAutoscaler {
	autoscalers := []autoscalingv2beta2.HorizontalPodAutoscaler{
		auth.NewAutoscaling(resource),
		core.NewAutoscaling(resource),
		api.NewAutoscaling(resource),
		analytic.NewAutoscaling(resource),
		manager.NewAutoscaling(resource),
		vulnerability.NewAutoscaling(resource),
		webhook.NewAutoscaling(resource),
	}
	msg := resource.GetMessagesComponent()
	if msg.Enabled {
		autoscalers = append(autoscalers, messages.NewAutoscaling(resource))
	}
	return autoscalers
}

func (b *Builder) ensureAutoscaling(resource *v2alpha1.HorusecPlatform, hpa *autoscalingv2beta2.HorizontalPodAutoscaler) error {
	if err := controllerutil.SetControllerReference(resource, hpa, b.scheme); err != nil {
		return fmt.Errorf("failed to set autoscaling %q owner reference: %v", hpa.GetName(), err)
	}

	return nil
}
