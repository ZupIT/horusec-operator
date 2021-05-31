// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ZupIT/horusec-operator/controllers"
	"github.com/ZupIT/horusec-operator/internal/horusec"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Injectors from wire.go:

func newHorusecPlatformReconciler(mgr manager.Manager) (*controllers.HorusecPlatformReconciler, error) {
	runtimeScheme := extractScheme(mgr)
	client := extractClient(mgr)
	service := horusec.NewService(client)
	adapterFactory := horusec.NewAdapterFactory(runtimeScheme, service)
	horusecPlatformReconciler := controllers.NewHorusecPlatformReconciler(adapterFactory)
	return horusecPlatformReconciler, nil
}