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
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type Initialization struct {
	client KubernetesClient
}

func NewInitialization(client KubernetesClient) *Initialization {
	return &Initialization{client: client}
}

func (i *Initialization) EnsureInitialization(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*operation.Result, error) {
	if resource.Status.Conditions != nil {
		return operation.ContinueProcessing()
	}
	resource.Status.Conditions = []v2alpha1.Condition{}
	resource.Status.State = v2alpha1.StatusPending
	err := i.client.UpdateHorusStatus(ctx, resource)
	if err != nil {
		return operation.RequeueWithError(err)
	}
	return operation.StopProcessing()
}
