//+build wireinject

package main

import (
	"github.com/google/wire"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ZupIT/horusec-operator/controllers"
)

func newHorusecPlatformReconciler(mgr manager.Manager) (*controllers.HorusecPlatformReconciler, error) {
	wire.Build(providers)
	return nil, nil
}
