package analytic

import (
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:lll, funlen // to improve in the future
func NewDeployment(resource *v2alpha1.HorusecPlatform) appsv1.Deployment {
	probe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/analytic/health",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
			},
		},
	}
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetAnalyticName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetAnalyticLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: resource.GetAnalyticReplicaCount(),
			Selector: &metav1.LabelSelector{MatchLabels: resource.GetAnalyticLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: resource.GetAnalyticLabels()},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:  resource.GetAnalyticName(),
					Image: resource.GetAnalyticImage(),
					Env: []corev1.EnvVar{
						{Name: "HORUSEC_PORT", Value: strconv.Itoa(resource.GetAnalyticPortHTTP())},
						{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: "false"},
						{Name: "HORUSEC_GRPC_USE_CERTS", Value: "false"},
						{Name: "HORUSEC_GRPC_AUTH_URL", Value: resource.GetAuthDefaultGRPCURL()},
						{Name: "HORUSEC_BROKER_HOST", Value: "rabbitmq"},
						{Name: "HORUSEC_BROKER_PORT", Value: "5672"},
						{Name: "HORUSEC_DATABASE_SQL_URI", Value: "postgresql://$(HORUSEC_DATABASE_USERNAME):$(HORUSEC_DATABASE_PASSWORD)@postgresql:5432/horusec_analytic_db?sslmode=disable"},
						NewEnvFromSecret("HORUSEC_BROKER_USERNAME", resource.GetGlobalBrokerUsername()),
						NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", resource.GetGlobalBrokerPassword()),
						NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", "horusec-database", "username"),
						NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", "horusec-database", "password"),
					},
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: int32(resource.GetAnalyticPortHTTP())},
					},
					LivenessProbe:  &probe,
					ReadinessProbe: &probe,
				}}},
			},
		},
	}
}

func NewEnvFromSecret(variableName string, secretKeyRef *corev1.SecretKeySelector) corev1.EnvVar {
	return corev1.EnvVar{
		Name: variableName,
		ValueFrom: &corev1.EnvVarSource{SecretKeyRef: secretKeyRef},
	}
}
