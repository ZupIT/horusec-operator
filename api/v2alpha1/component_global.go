package v2alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"reflect"
)

func (h *HorusecPlatform) NewEnvFromSecret(variableName string, secretKeyRef *corev1.SecretKeySelector) corev1.EnvVar {
	return corev1.EnvVar{
		Name: variableName,
		ValueFrom: &corev1.EnvVarSource{SecretKeyRef: secretKeyRef},
	}
}

func (h *HorusecPlatform) GetDefaultLabel() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/managed-by": "horusec",
	}
}

func (h *HorusecPlatform) GetLatestVersion() string {
	return "v2.12.1"
}

func (h *HorusecPlatform) GetGlobalBrokerUsername() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Broker.User).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-broker"},
			Key:                  "user",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Broker.User.SecretKeyRef
}
func (h *HorusecPlatform) GetGlobalBrokerPassword() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Broker.Password).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-broker"},
			Key:                  "password",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Broker.Password.SecretKeyRef
}

func (h *HorusecPlatform) GetGlobalDatabaseUsername() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Database.User).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-database"},
			Key:                  "user",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Database.User.SecretKeyRef
}
func (h *HorusecPlatform) GetGlobalDatabasePassword() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Database.Password).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-broker"},
			Key:                  "password",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Database.Password.SecretKeyRef
}