package main

import (
	"github.com/ZupIT/horusec-operator/internal/resources"
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ZupIT/horusec-operator/controllers"
	"github.com/ZupIT/horusec-operator/internal/horusec"
)

// nolint:deadcode,unused,varcheck // used for dependency injection container generation
var providers = wire.NewSet(
	extractClient,
	extractRestConfig,
	extractScheme,

	controllers.NewHorusecPlatformReconciler,
	horusec.NewAdapterFactory,
	horusec.NewService,
	resources.NewBuilder,
	wire.Bind(new(controllers.AdapterFactory), new(*horusec.AdapterFactory)),
	wire.Bind(new(horusec.ResourceBuilder), new(*resources.Builder)),
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
