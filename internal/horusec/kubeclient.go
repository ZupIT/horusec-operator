package horusec

import (
	"context"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/k8s"
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

type KubernetesClient interface {
	Apply(ctx context.Context, objects k8s.Objects) error
	GetHorus(ctx context.Context, namespacedName types.NamespacedName) (*v2alpha1.HorusecPlatform, error)
	UpdateHorusStatus(ctx context.Context, horus *v2alpha1.HorusecPlatform) error
	ListAutoscalingByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]autoscaling.HorizontalPodAutoscaler, error)
	ListDeploymentsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]apps.Deployment, error)
	ListIngressByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]networking.Ingress, error)
	ListJobsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]batch.Job, error)
	ListServiceAccountsByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.ServiceAccount, error)
	ListServicesByOwner(ctx context.Context, owner *v2alpha1.HorusecPlatform) ([]core.Service, error)
}
