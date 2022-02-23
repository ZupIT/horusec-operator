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
	"context"

	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/k8s"
)

type KubernetesClient interface {
	Apply(ctx context.Context, objects k8s.Objects) error
	UpdateHorusStatus(ctx context.Context, horus *v2alpha1.HorusecPlatform) error
	ListAutoscalingByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]autoscaling.HorizontalPodAutoscaler, error)
	ListDeploymentsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]apps.Deployment, error)
	ListIngressByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]networkingv1.Ingress, error)
	ListJobsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]batch.Job, error)
	ListPodsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.Pod, error)
	ListServiceAccountsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.ServiceAccount, error)
	ListServicesByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.Service, error)
}

type KubernetesLogs interface {
	PreviousContainerLogs(ctx context.Context, pod types.NamespacedName, container string) ([]byte, error)
}
