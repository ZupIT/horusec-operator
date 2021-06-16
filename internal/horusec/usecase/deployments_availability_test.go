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
	"testing"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/api/v2alpha1/condition"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentsAvailability_EnsureDeploymentsAvailable(t *testing.T) {
	// arrange
	usecase, ctrl := setupToEnsureDeploymentsAvailable(t)

	// act
	resource := v2alpha1.HorusecPlatform{}
	result, err := usecase.EnsureDeploymentsAvailable(context.TODO(), &resource)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, operation.StopResult(), result)
	assert.True(t, resource.IsStatusConditionTrue(condition.AnalyticAvailable), `"AnalyticAvailable" condition should be true`)
	assert.False(t, resource.IsStatusConditionTrue(condition.APIAvailable), `"APIAvailable" condition should be unknown`)
	assert.False(t, resource.IsStatusConditionTrue(condition.AuthAvailable), `"AuthAvailable" condition should be unknown`)
	assert.False(t, resource.IsStatusConditionTrue(condition.CoreAvailable), `"CoreAvailable" condition should be unknown`)
	assert.True(t, resource.IsStatusConditionTrue(condition.ManagerAvailable), `"ManagerAvailable" condition should be true`)
	assert.False(t, resource.IsStatusConditionTrue(condition.VulnerabilityAvailable), `"VulnerabilityAvailable" condition should be unknown`)
	assert.False(t, resource.IsStatusConditionTrue(condition.WebhookAvailable), `"WebhookAvailable" condition should be unknown`)
	ctrl.Finish()
}

func setupToEnsureDeploymentsAvailable(t *testing.T) (*DeploymentsAvailability, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	client := test.NewMockKubernetesClient(ctrl)

	deps, err := test.DeploymentsWithStatus()
	assert.NoError(t, err)

	client.EXPECT().
		ListDeploymentsByOwner(gomock.Any(), gomock.Any()).
		Times(1).
		Return(deps, nil)
	client.EXPECT().
		UpdateHorusStatus(gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil)

	return NewDeploymentsAvailability(client), ctrl
}
