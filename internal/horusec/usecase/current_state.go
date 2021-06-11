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
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/api/v2alpha1/condition"
	"github.com/ZupIT/horusec-operator/api/v2alpha1/state"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/internal/tracing"
)

type CurrentState struct {
	client KubernetesClient
}

func NewCurrentState(client KubernetesClient) *CurrentState {
	return &CurrentState{client: client}
}

func (i *CurrentState) EnsureCurrentState(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*operation.Result, error) {
	span := tracing.SpanFromContext(ctx)
	log := span.Logger()

	if resource.Status.State == "" && resource.SetState(state.Pending) {
		log.Info(fmt.Sprintf("Updating status to %q", state.Pending))
		return operation.RequeueWithError(i.client.UpdateHorusStatus(ctx, resource))
	}

	if resource.IsStatusConditionTrue(condition.DeploymentsAvailable) {
		if resource.SetState(state.Ready) {
			log.Info(fmt.Sprintf("Updating status to %q", state.Ready))
			return operation.RequeueWithError(i.client.UpdateHorusStatus(ctx, resource))
		}
		return operation.ContinueProcessing()
	}

	if resource.IsStatusConditionFalse(condition.DeploymentsAvailable) {
		if resource.SetState(state.Pending) {
			log.Info(fmt.Sprintf("Updating status to %q", state.Pending))
			return operation.RequeueWithError(i.client.UpdateHorusStatus(ctx, resource))
		}
		return operation.ContinueProcessing()
	}

	return operation.ContinueProcessing()
}
