package horusec

import (
	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
)

type ResourceBuilder interface {
	AutoscalingFor(resource *v2alpha1.HorusecPlatform) ([]autoscalingv2beta2.HorizontalPodAutoscaler, error)
	DeploymentsFor(resource *v2alpha1.HorusecPlatform) ([]appsv1.Deployment, error)
	IngressFor(resource *v2alpha1.HorusecPlatform) ([]networkingv1.Ingress, error)
	JobsFor(resource *v2alpha1.HorusecPlatform) ([]batchv1.Job, error)
	ServiceAccountsFor(resource *v2alpha1.HorusecPlatform) ([]corev1.ServiceAccount, error)
	ServicesFor(resource *v2alpha1.HorusecPlatform) ([]corev1.Service, error)
}
