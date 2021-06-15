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
	d, err := e.client.ListDeploymentsByOwner(ctx, resource)
	if err != nil {
		return nil, err
	}

	status := statusOfDeployments(d).UpdateConditions(resource)
	if status.HasChanges() {
		return operation.RequeueOnErrorOrStop(e.client.UpdateHorusStatus(ctx, resource))
	}

	return operation.ContinueProcessing()
}

type deployStatus struct{ item appsv1.DeploymentStatus }

func (ps *deployStatus) HasUnavailableReplicas() bool {
	return ps.item.UnavailableReplicas > 0
}

type deployStatuses struct {
	items      map[string]*deployStatus
	conditions map[string]condition.Type
	changed    bool
}

func statusOfDeployments(deployments []appsv1.Deployment) *deployStatuses {
	items := make(map[string]*deployStatus, len(deployments))
	for _, deploy := range deployments {
		if component, ok := deploy.Labels["app.kubernetes.io/component"]; ok {
			items[component] = &deployStatus{item: deploy.Status}
		}
	}
	return &deployStatuses{
		items: items,
		conditions: map[string]condition.Type{
			"analytic": condition.AnalyticAvailable, "api": condition.APIAvailable, "auth": condition.AuthAvailable,
			"core": condition.CoreAvailable, "manager": condition.ManagerAvailable,
			"vulnerability": condition.VulnerabilityAvailable, "webhook": condition.WebhookAvailable,
		},
	}
}

func (ds *deployStatuses) UpdateConditions(resource *v2alpha1.HorusecPlatform) *deployStatuses {
	reason := condition.Reason{
		Type:    "UnavailableReplicas",
		Message: "Deployment is unavailable but we could not discover the cause.",
	}

	for component, conditionType := range ds.conditions {
		if ds.checkAvailabilityOf(component) && resource.SetStatusCondition(condition.True(conditionType)) {
			ds.changed = true
			continue
		}

		if !resource.IsStatusConditionFalse(conditionType) && resource.SetStatusCondition(condition.Unknown(conditionType, reason)) {
			ds.changed = true
		}
	}

	return ds
}

func (ds *deployStatuses) HasChanges() bool {
	return ds.changed
}

func (ds *deployStatuses) checkAvailabilityOf(component string) bool {
	if status, ok := ds.items[component]; ok {
		return !status.HasUnavailableReplicas()
	}
	return false
}
