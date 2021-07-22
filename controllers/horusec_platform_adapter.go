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

package controllers

import (
	"context"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/operation"
)

type HorusecPlatformAdapter interface {
	EnsureCurrentState(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureDatabaseMigrations(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureDeployments(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureAutoscaling(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureServices(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureIngressRules(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureServiceAccounts(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureDeploymentsAvailable(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
	EnsureUnavailabilityReason(context.Context, *v2alpha1.HorusecPlatform) (*operation.Result, error)
}
