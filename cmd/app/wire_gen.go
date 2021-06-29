// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/ZupIT/horusec-operator/controllers"
	"github.com/ZupIT/horusec-operator/internal/horusec"
	"github.com/ZupIT/horusec-operator/internal/k8s"
	"github.com/ZupIT/horusec-operator/internal/resources"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Injectors from wire.go:

func newHorusecPlatformReconciler(mgr manager.Manager) (*controllers.HorusecPlatformReconciler, error) {
	runtimeScheme := extractScheme(mgr)
	builder := resources.NewBuilder(runtimeScheme)
	client := extractClient(mgr)
	k8sClient := k8s.NewClient(client)
	config := extractRestConfig(mgr)
	coreV1Interface := k8s.NewTypedCoreClient(config)
	containerClient := k8s.NewContainerClient(coreV1Interface)
	adapter := horusec.NewAdapter(builder, k8sClient, containerClient)
	horusecPlatformReconciler := controllers.NewHorusecPlatformReconciler(adapter, k8sClient)
	return horusecPlatformReconciler, nil
}
