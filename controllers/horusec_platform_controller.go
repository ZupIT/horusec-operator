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
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	installv2 "github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/internal/requeue"
)

// HorusecPlatformReconciler reconciles a HorusecPlatform object
type HorusecPlatformReconciler struct {
	factory AdapterFactory
	log     logr.Logger
}

func NewHorusecPlatformReconciler(factory AdapterFactory) *HorusecPlatformReconciler {
	return &HorusecPlatformReconciler{
		factory: factory,
		log:     ctrl.Log.WithName("controllers").WithName("HorusecPlatform"),
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
//nolint:funlen // to improve in the future
func (r *HorusecPlatformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.log.WithValues("horusec", req.NamespacedName)
	log.Info("reconciling")

	adapter, err := r.factory.CreateHorusecPlatformAdapter(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return requeue.Not()
		}
		return requeue.OnErr(err)
	}

	result, err := operation.NewHandler(
		adapter.EnsureAuthDeployments,
		adapter.EnsureServiceAccounts,
		adapter.EnsureInitialization,
	).Handle(ctx)
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
