package v2alpha1

import (
	"fmt"
	"reflect"
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

func (h *HorusecPlatform) NewEnvFromSecret(variableName string, secretKeyRef *corev1.SecretKeySelector) corev1.EnvVar {
	return corev1.EnvVar{
		Name:      variableName,
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
	return &h.Spec.Global.Broker.User.KeyRef
}

func (h *HorusecPlatform) GetGlobalBrokerPassword() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Broker.Password).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-broker"},
			Key:                  "password",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Broker.Password.KeyRef
}

func (h *HorusecPlatform) GetGlobalDatabaseUsername() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Database.User).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-database"},
			Key:                  "user",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Database.User.KeyRef
}

func (h *HorusecPlatform) GetGlobalDatabasePassword() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.Spec.Global.Database.Password).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-database"},
			Key:                  "password",
			Optional:             nil,
		}
	}
	return &h.Spec.Global.Database.Password.KeyRef
}

func (h *HorusecPlatform) GetGlobalDatabaseLogMode() string {
	if h.Spec.Global.Database.LogMode {
		return "true"
	}

	return "false"
}

func (h *HorusecPlatform) GetGlobalDatabaseHost() string {
	host := h.Spec.Global.Database.Host
	if host == "" {
		return "postgresql"
	}

	return host
}

func (h *HorusecPlatform) GetGlobalDatabasePort() string {
	port := h.Spec.Global.Database.Port
	if port <= 0 {
		return "5432"
	}

	return strconv.Itoa(port)
}

func (h *HorusecPlatform) GetGlobalDatabaseName() string {
	name := h.Spec.Global.Database.Name
	if name == "" {
		return "horusec_db"
	}

	return name
}

func (h *HorusecPlatform) GetGlobalSSLMode() string {
	mode := h.Spec.Global.Database.SslMode
	if mode == nil || *mode {
		return ""
	}

	return "?sslmode=disable"
}

func (h *HorusecPlatform) GetGlobalBrokerHost() string {
	host := h.Spec.Global.Broker.Host
	if host == "" {
		return "rabbitmq"
	}

	return host
}

func (h *HorusecPlatform) GetGlobalBrokerPort() string {
	port := h.Spec.Global.Broker.Port
	if port <= 0 {
		return "5672"
	}

	return strconv.Itoa(port)
}

func (h *HorusecPlatform) GetGlobalDatabaseURI() string {
	return fmt.Sprintf(
		"postgresql://$(HORUSEC_DATABASE_USERNAME):$(HORUSEC_DATABASE_PASSWORD)@%s:%s/%s%s",
		h.GetGlobalDatabaseHost(), h.GetGlobalDatabasePort(), h.GetGlobalDatabaseName(), h.GetGlobalSSLMode())
}

func (h *HorusecPlatform) GetAllIngressIsDisabled() bool {
	return !h.IsAnalyticIngressEnabled() &&
		!h.IsAPIIngressEnabled() &&
		!h.IsAuthIngressEnabled() &&
		!h.IsCoreIngressEnabled() &&
		!h.IsManagerIngressEnabled() &&
		!h.IsMessagesIngressEnabled() &&
		!h.IsVulnerabilityIngressEnabled() &&
		!h.IsWebhookIngressEnabled()
}
