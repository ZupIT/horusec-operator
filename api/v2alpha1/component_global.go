package v2alpha1

import (
	"fmt"
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

func (h *HorusecPlatform) IsEmailsEnabled() string {
	enabled := h.Spec.Components.Messages.Enabled
	if enabled {
		return "true"
	}

	return "false"
}
