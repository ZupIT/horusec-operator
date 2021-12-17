// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"
	"reflect"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/resources/analytic"
	"github.com/ZupIT/horusec-operator/internal/resources/api"
	"github.com/ZupIT/horusec-operator/internal/resources/auth"
	"github.com/ZupIT/horusec-operator/internal/resources/core"
	"github.com/ZupIT/horusec-operator/internal/resources/manager"
	"github.com/ZupIT/horusec-operator/internal/resources/messages"
	"github.com/ZupIT/horusec-operator/internal/resources/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/resources/webhook"
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
