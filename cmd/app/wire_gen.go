// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/ZupIT/horusec-operator/controllers"
	"github.com/ZupIT/horusec-operator/internal/horusec"
	"github.com/ZupIT/horusec-operator/internal/k8s"
	"github.com/ZupIT/horusec-operator/internal/resources"
)

// Injectors from wire.go:

func newHorusecPlatformReconciler(mgr manager.Manager) (*controllers.HorusecPlatformReconciler, error) {
	client := extractClient(mgr)
	k8sClient := k8s.NewClient(client)
	runtimeScheme := extractScheme(mgr)
	builder := resources.NewBuilder(runtimeScheme)
	adapter := horusec.NewAdapter(k8sClient, builder)
	horusecPlatformReconciler := controllers.NewHorusecPlatformReconciler(adapter, k8sClient)
	return horusecPlatformReconciler, nil
}
