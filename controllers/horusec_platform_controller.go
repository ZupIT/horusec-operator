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

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	installv2 "github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/internal/requeue"
	"github.com/ZupIT/horusec-operator/internal/tracing"
)

// HorusecPlatformReconciler reconciles a HorusecPlatform object
type HorusecPlatformReconciler struct {
	adapter HorusecPlatformAdapter
	client  HorusecPlatformClient
}

func NewHorusecPlatformReconciler(adapter HorusecPlatformAdapter, client HorusecPlatformClient) *HorusecPlatformReconciler {
	return &HorusecPlatformReconciler{adapter: adapter, client: client}
}

//+kubebuilder:rbac:groups=install.horusec.io,resources=horusecplatforms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=install.horusec.io,resources=horusecplatforms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extensions,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

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
	span, ctx := tracing.StartSpanFromContext(ctx, tracing.WithCustomResource(req.NamespacedName))
	defer span.Finish()

	log := span.Logger()
	log.Info("reconciling")

	resource, err := r.client.GetHorus(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return requeue.Not()
		}
		span.SetError(err)
		return requeue.OnErr(err)
	}

	result, err := operation.NewHandler(
		r.adapter.EnsureCurrentState,
		r.adapter.EnsureServiceAccounts,
		r.adapter.EnsureDatabaseMigrations,
		r.adapter.EnsureServices,
		r.adapter.EnsureDeployments,
		r.adapter.EnsureAutoscaling,
		r.adapter.EnsureIngressRules,
		r.adapter.EnsureDeploymentsAvailable,
		r.adapter.EnsureUnavailabilityReason,
	).Handle(ctx, resource)
	log.V(1).
		WithValues("error", err != nil, "requeing", result.Requeue, "delay", result.RequeueAfter).
		Info("finished reconcile")
	return result, span.HandleError(err)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HorusecPlatformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&installv2.HorusecPlatform{}).
		Complete(r)
}
