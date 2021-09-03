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

package usecase

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

type ResourceBuilder interface {
	AutoscalingFor(resource *v2alpha1.HorusecPlatform) ([]autoscalingv2beta2.HorizontalPodAutoscaler, error)
	DeploymentsFor(resource *v2alpha1.HorusecPlatform) ([]appsv1.Deployment, error)
	IngressFor(resource *v2alpha1.HorusecPlatform) ([]networkingv1beta1.Ingress, error)
	JobsFor(resource *v2alpha1.HorusecPlatform) ([]batchv1.Job, error)
	ServiceAccountsFor(resource *v2alpha1.HorusecPlatform) ([]corev1.ServiceAccount, error)
	ServicesFor(resource *v2alpha1.HorusecPlatform) ([]corev1.Service, error)
}
