package v2alpha1

import (
	"fmt"
	"reflect"
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

func (h *HorusecPlatform) GetAnalyticComponent() Analytic {
	return h.Spec.Components.Analytic
}

func (h *HorusecPlatform) GetAnalyticAutoscaling() Autoscaling {
	return h.GetAnalyticComponent().Pod.Autoscaling
}

func (h *HorusecPlatform) GetAnalyticName() string {
	name := h.GetAnalyticComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-analytic", h.GetName())
	}
	return name
}

func (h *HorusecPlatform) GetAnalyticPath() string {
	path := h.GetAnalyticComponent().Ingress.Path
	if path == "" {
		return "/analytic"
	}
	return path
}

func (h *HorusecPlatform) GetAnalyticPortHTTP() int {
	port := h.GetAnalyticComponent().Port.HTTP
	if port == 0 {
		return 8005
	}
	return port
}

func (h *HorusecPlatform) GetAnalyticLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "analytic",
		"app.kubernetes.io/managed-by": "horusec",
	}
}

func (h *HorusecPlatform) GetAnalyticReplicaCount() *int32 {
	if !h.GetAnalyticAutoscaling().Enabled {
		count := h.GetAnalyticComponent().ReplicaCount
		return &count
	}
	return nil
}

func (h *HorusecPlatform) GetAnalyticDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetAnalyticName(), h.GetAnalyticPortHTTP())
}

func (h *HorusecPlatform) GetAnalyticRegistry() string {
	registry := h.GetAnalyticComponent().Container.Image.Registry
	if registry == "" {
		return "docker.io/"
	}
	return registry
}

func (h *HorusecPlatform) GetAnalyticRepository() string {
	repository := h.GetAnalyticComponent().Container.Image.Repository
	if repository == "" {
		return "horuszup/horusec-analytic"
	}
	return repository
}

func (h *HorusecPlatform) GetAnalyticTag() string {
	tag := h.GetAnalyticComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}

func (h *HorusecPlatform) GetAnalyticImage() string {
	return fmt.Sprintf("%s%s:%s", h.GetAnalyticRegistry(), h.GetAnalyticRepository(), h.GetAnalyticTag())
}

func (h *HorusecPlatform) GetAnalyticDatabaseUsername() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.GetAnalyticComponent().Database.User).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-database"},
			Key:                  "user",
			Optional:             nil,
		}
	}
	secret := h.GetAnalyticComponent().Database.User.KeyRef
	return &secret
}

func (h *HorusecPlatform) GetAnalyticDatabasePassword() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.GetAnalyticComponent().Database.Password).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-database"},
			Key:                  "password",
			Optional:             nil,
		}
	}
	secret := h.GetAnalyticComponent().Database.Password.KeyRef
	return &secret
}

func (h *HorusecPlatform) GetAnalyticDatabaseLogMode() string {
	if h.Spec.Components.Analytic.Database.LogMode {
		return "true"
	}

	return "false"
}

func (h *HorusecPlatform) GetAnalyticDatabaseHost() string {
	host := h.Spec.Components.Analytic.Database.Host
	if host == "" {
		return "postgresql"
	}

	return host
}

func (h *HorusecPlatform) GetAnalyticDatabasePort() string {
	port := h.Spec.Components.Analytic.Database.Port
	if port <= 0 {
		return "5432"
	}

	return strconv.Itoa(port)
}

func (h *HorusecPlatform) GetAnalyticDatabaseName() string {
	name := h.Spec.Components.Analytic.Database.Name
	if name == "" {
		return "horusec_analytic_db"
	}

	return name
}

func (h *HorusecPlatform) GetAnalyticSSLMode() string {
	mode := h.Spec.Components.Analytic.Database.SslMode
	if mode == nil || *mode {
		return ""
	}

	return "?sslmode=disable"
}

func (h *HorusecPlatform) GetAnalyticDatabaseURI() string {
	return fmt.Sprintf("postgresql://$(HORUSEC_DATABASE_USERNAME):$(HORUSEC_DATABASE_PASSWORD)@%s:%s/%s%s",
		h.GetAnalyticDatabaseHost(), h.GetAnalyticDatabasePort(), h.GetAnalyticDatabaseName(), h.GetAnalyticSSLMode())
}

func (h *HorusecPlatform) GetAnalyticHost() string {
	host := h.Spec.Components.Analytic.Ingress.Host
	if host == "" {
		return "analytic.local"
	}

	return host
}

func (h *HorusecPlatform) IsAnalyticIngressEnabled() bool {
	enabled := h.Spec.Components.Analytic.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
