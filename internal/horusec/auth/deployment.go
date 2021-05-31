package auth

import (
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:lll, funlen // to improve in the future
func NewDeployment(resource *v2alpha1.HorusecPlatform) *appsv1.Deployment {
	var replicas int32 = 1
	probe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/auth/health",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
			},
		},
	}
	var httpPort int32 = 8006
	var grpcPort int32 = 8007
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: Labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: Labels},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:  "horusec-auth",
					Image: "docker.io/horuszup/horusec-auth:v2.12.0",
					Env: []corev1.EnvVar{
						{Name: "HORUSEC_PORT", Value: strconv.Itoa(int(httpPort))},
						{Name: "HORUSEC_GRPC_PORT", Value: strconv.Itoa(int(grpcPort))},
						{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: "false"},
						{Name: "HORUSEC_GRPC_USE_CERTS", Value: "false"},
						{Name: "HORUSEC_BROKER_HOST", Value: ""},
						{Name: "HORUSEC_BROKER_PORT", Value: "5672"},
						{Name: "HORUSEC_AUTH_TYPE", Value: "horusec"},
						{Name: "HORUSEC_ENABLE_APPLICATION_ADMIN", Value: "false"},
						{Name: "HORUSEC_DATABASE_SQL_URI", Value: "postgresql://$(HORUSEC_DATABASE_USERNAME):$(HORUSEC_DATABASE_PASSWORD)@db.svc.cluster.local:5432/horusec_db?sslmode=disable"},
						NewEnvFromSecret("HORUSEC_BROKER_USERNAME", "horusec-broker", "username"),
						NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", "horusec-broker", "password"),
						NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", "horusec-database", "username"),
						NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", "horusec-database", "password"),
						NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", "horusec-jwt", "jwt-token"),
					},
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: httpPort},
						{Name: "grpc", ContainerPort: grpcPort},
					},
					LivenessProbe:  &probe,
					ReadinessProbe: &probe,
				}}},
			},
		},
	}
}

func NewEnvFromSecret(variableName, secretName, secretKey string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: variableName,
		ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
			Key:                  secretKey,
		}},
	}
}