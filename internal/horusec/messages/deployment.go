package messages

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
	component := resource.GetMessagesComponent()
	probe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/messages/health",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
			},
		},
	}
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: component.GetReplicaCount(),
			Selector: &metav1.LabelSelector{MatchLabels: Labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: Labels},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:  "horusec-messages",
					Image: "docker.io/horuszup/horusec-messages:v2.12.1",
					Env: []corev1.EnvVar{
						{Name: "HORUSEC_PORT", Value: strconv.Itoa(component.Port.HTTP)},
						{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: "false"},
						{Name: "HORUSEC_GRPC_USE_CERTS", Value: "false"},
						{Name: "HORUSEC_GRPC_AUTH_URL", Value: "horusec-auth:8007"},
						{Name: "HORUSEC_BROKER_HOST", Value: "localhost"},
						{Name: "HORUSEC_BROKER_PORT", Value: "5672"},
						{Name: "HORUSEC_SMTP_HOST", Value: ""},
						{Name: "HORUSEC_SMTP_PORT", Value: ""},
						{Name: "HORUSEC_EMAIL_FROM", Value: "horusec@zup.com.br"},
						{Name: "HORUSEC_DATABASE_SQL_URI", Value: "postgresql://$(HORUSEC_DATABASE_USERNAME):$(HORUSEC_DATABASE_PASSWORD)@db.svc.cluster.local:5432/horusec_db?sslmode=disable"},
						NewEnvFromSecret("HORUSEC_BROKER_USERNAME", "horusec-broker", "username"),
						NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", "horusec-broker", "password"),
						NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", "horusec-database", "username"),
						NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", "horusec-database", "password"),
						NewEnvFromSecret("HORUSEC_SMTP_USERNAME", "horusec-smtp", "username"),
						NewEnvFromSecret("HORUSEC_SMTP_PASSWORD", "horusec-smtp", "password"),
					},
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: int32(component.Port.HTTP)},
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
