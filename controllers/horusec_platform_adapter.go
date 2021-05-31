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

	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/internal/operation"
)

type AdapterFactory interface {
	CreateHorusecPlatformAdapter(ctx context.Context, key k8s.ObjectKey) (HorusecPlatformAdapter, error)
}

type HorusecPlatformAdapter interface {
	EnsureDatabaseConnectivity(ctx context.Context) (*operation.Result, error)
	EnsureBrokerConnectivity(ctx context.Context) (*operation.Result, error)
	EnsureSMTPConnectivity(ctx context.Context) (*operation.Result, error)
	EnsureDatabaseMigrations(ctx context.Context) (*operation.Result, error)
	EnsureAuthDeployments(ctx context.Context) (*operation.Result, error)
	EnsureDeployments(ctx context.Context) (*operation.Result, error)
	EnsureServices(ctx context.Context) (*operation.Result, error)
	EnsureAutoscalers(ctx context.Context) (*operation.Result, error)
	EnsureHPA(ctx context.Context) (*operation.Result, error)
	EnsureIngressRules(ctx context.Context) (*operation.Result, error)
	EnsureEverythingIsRunning(ctx context.Context) (*operation.Result, error)
	EnsureServiceAccounts(ctx context.Context) (*operation.Result, error)
}
