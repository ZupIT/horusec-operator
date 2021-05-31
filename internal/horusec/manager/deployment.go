package manager

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:funlen // to improve in the future
func NewDeployment(resource *v2alpha1.HorusecPlatform) *appsv1.Deployment {
	component := resource.GetManagerComponent()
	return &appsv1.Deployment{
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
					Name:  "horusec-manager",
					Image: "docker.io/horuszup/horusec-manager:v2.12.1",
					Env: []corev1.EnvVar{
						{Name: "REACT_APP_HORUSEC_ENDPOINT_API", Value: ""},
						{Name: "REACT_APP_HORUSEC_ENDPOINT_ANALYTIC", Value: ""},
						{Name: "REACT_APP_HORUSEC_ENDPOINT_CORE", Value: ""},
						{Name: "REACT_APP_HORUSEC_ENDPOINT_WEBHOOK", Value: ""},
						{Name: "REACT_APP_HORUSEC_ENDPOINT_AUTH", Value: ""},
						{Name: "REACT_APP_KEYCLOAK_BASE_PATH", Value: ""},
						{Name: "REACT_APP_KEYCLOAK_CLIENT_ID", Value: ""},
						{Name: "REACT_APP_KEYCLOAK_REALM", Value: ""},
						{Name: "REACT_APP_MICROFRONTEND_PUBLIC_PATH", Value: ""},
						{Name: "REACT_APP_HORUSEC_MANAGER_THEME", Value: ""},
					},
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: int32(component.Port.HTTP)},
					},
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
