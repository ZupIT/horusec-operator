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

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	installv2 "github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/internal/requeue"
)

// HorusecPlatformReconciler reconciles a HorusecPlatform object
type HorusecPlatformReconciler struct {
	svc *horusec.Service
	log logr.Logger
}

func NewHorusecPlatformReconciler(svc *horusec.Service) *HorusecPlatformReconciler {
	return &HorusecPlatformReconciler{
		svc: svc,
		log: ctrl.Log.WithName("controllers").WithName("HorusecPlatform"),
	}
}

//+kubebuilder:rbac:groups=install.horusec.io,resources=horusecs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=install.horusec.io,resources=horusecs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=install.horusec.io,resources=horusecs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HorusecPlatform object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *HorusecPlatformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.log.WithValues("horusec", req.NamespacedName)
	log.Info("reconciling")

	adapter, err := r.svc.LookupResourceAdapter(ctx, req.NamespacedName)
	if err != nil {
		return requeue.OnErr(err)
	} else if adapter == nil {
		return requeue.Not()
	}

	result, err := operation.NewHandler(adapter.EnsureAuthDeployments).Handle(ctx)
	log.V(1).
		WithValues("error", err != nil, "requeing", result.Requeue, "delay", result.RequeueAfter).
		Info("finished reconcile")
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *HorusecPlatformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&installv2.HorusecPlatform{}).
		Complete(r)
}
