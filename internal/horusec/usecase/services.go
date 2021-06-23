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
	"github.com/ZupIT/horusec-operator/internal/tracing"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/inventory"
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type Services struct {
	client  KubernetesClient
	builder ResourceBuilder
}

func NewServices(client KubernetesClient, builder ResourceBuilder) *Services {
	return &Services{client: client, builder: builder}
}

func (s *Services) EnsureServices(ctx context.Context, resource *v2alpha1.HorusecPlatform) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	existing, err := s.client.ListServicesByOwner(ctx, resource)
	if err != nil {
		return nil, err
	}

	desired, err := s.builder.ServicesFor(resource)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForService(existing, desired)
	if err := s.client.Apply(ctx, inv); err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}
