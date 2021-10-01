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

package main

import (
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ZupIT/horusec-operator/controllers"
	"github.com/ZupIT/horusec-operator/internal/horusec"
	"github.com/ZupIT/horusec-operator/internal/horusec/usecase"
	"github.com/ZupIT/horusec-operator/internal/k8s"
	"github.com/ZupIT/horusec-operator/internal/resources"
)

// nolint // used for dependency injection container generation
var providers = wire.NewSet(
	extractClient,
	extractRestConfig,
	extractScheme,

	controllers.NewHorusecPlatformReconciler,
	horusec.NewAdapter,
	k8s.NewClient,
	k8s.NewContainerClient,
	k8s.NewTypedCoreClient,
	resources.NewBuilder,
	wire.Bind(new(controllers.HorusecPlatformAdapter), new(*horusec.Adapter)),
	wire.Bind(new(controllers.HorusecPlatformClient), new(*k8s.Client)),
	wire.Bind(new(usecase.KubernetesClient), new(*k8s.Client)),
	wire.Bind(new(usecase.KubernetesLogs), new(*k8s.ContainerClient)),
	wire.Bind(new(usecase.ResourceBuilder), new(*resources.Builder)),
)

func extractScheme(mgr manager.Manager) *runtime.Scheme {
	return mgr.GetScheme()
}

func extractClient(mgr manager.Manager) client.Client {
	return mgr.GetClient()
}

func extractRestConfig(mgr manager.Manager) *rest.Config {
	return mgr.GetConfig()
}
