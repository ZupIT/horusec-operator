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

package auth

import (
	autoScalingV2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

// nolint:funlen // constructor is required all data
func NewAutoscaling(resource *v2alpha1.HorusecPlatform) autoScalingV2beta2.HorizontalPodAutoscaler {
	autoScaling := resource.GetAuthAutoscaling()
	if !autoScaling.Enabled {
		return autoScalingV2beta2.HorizontalPodAutoscaler{}
	}
	return autoScalingV2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetAuthName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetAuthLabels(),
		},
		Spec: autoScalingV2beta2.HorizontalPodAutoscalerSpec{
			MinReplicas: autoScaling.MinReplicas,
			MaxReplicas: autoScaling.MaxReplicas,
			ScaleTargetRef: autoScalingV2beta2.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       "auth",
				APIVersion: "apps/v1",
			},
			Metrics: []autoScalingV2beta2.MetricSpec{
				{
					Type: "Resource",
					Resource: &autoScalingV2beta2.ResourceMetricSource{
						Name: "cpu",
						Target: autoScalingV2beta2.MetricTarget{
							AverageUtilization: autoScaling.TargetCPU,
						},
					},
				},
				{
					Type: "Resource",
					Resource: &autoScalingV2beta2.ResourceMetricSource{
						Name: "memory",
						Target: autoScalingV2beta2.MetricTarget{
							AverageUtilization: autoScaling.TargetMemory,
						},
					},
				},
			},
		},
	}
}
