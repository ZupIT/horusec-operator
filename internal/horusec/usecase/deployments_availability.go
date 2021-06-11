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
	"time"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/api/v2alpha1/condition"
	"github.com/ZupIT/horusec-operator/internal/operation"
	appsv1 "k8s.io/api/apps/v1"
)

type DeploymentsAvailability struct {
	client KubernetesClient
}

func NewDeploymentsAvailability(client KubernetesClient) *DeploymentsAvailability {
	return &DeploymentsAvailability{client: client}
}

func (e *DeploymentsAvailability) EnsureDeploymentsAvailable(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*operation.Result, error) {
	deployments, err := e.client.ListDeploymentsByOwner(ctx, resource)
	if err != nil {
		return nil, err
	}

	if statusOf(deployments).IsAvailable() {
		if resource.SetStatusConditionTrue(condition.DeploymentsAvailable) {
			return operation.RequeueOnErrorOrStop(e.client.UpdateHorusStatus(ctx, resource))
		}
		return operation.ContinueProcessing()
	}

	if resource.SetStatusConditionFalse(condition.DeploymentsAvailable) {
		return operation.RequeueOnErrorOrStop(e.client.UpdateHorusStatus(ctx, resource))
	}
	return operation.RequeueAfter(10*time.Second, nil)
}

type (
	deployStatuses struct{ items []*deployStatus }
	deployStatus   struct{ item *appsv1.DeploymentStatus }
)

func statusOf(deployments []appsv1.Deployment) *deployStatuses {
	items := make([]*deployStatus, 0, len(deployments))
	for _, pod := range deployments {
		items = append(items, &deployStatus{item: &pod.Status})
	}
	return &deployStatuses{items: items}
}

func (ds *deployStatuses) IsAvailable() bool {
	for _, item := range ds.items {
		if item.HasUnavailableReplicas() {
			return false
		}
	}
	return true
}

func (ps *deployStatus) HasUnavailableReplicas() bool {
	return ps.item.UnavailableReplicas > 0
}
