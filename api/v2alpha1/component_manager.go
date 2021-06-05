package v2alpha1

import (
	"fmt"
)

func (h *HorusecPlatform) GetManagerComponent() ExposableComponent {
	return h.Spec.Components.Manager
}

func (h *HorusecPlatform) GetManagerAutoscaling() Autoscaling {
	return h.GetManagerComponent().Pod.Autoscaling
}

func (h *HorusecPlatform) GetManagerName() string {
	name := h.GetManagerComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-manager", h.GetName())
	}
	return name
}

func (h *HorusecPlatform) GetManagerPath() string {
	path := h.GetManagerComponent().Ingress.Path
	if path == "" {
		return "/"
	}
	return path
}

func (h *HorusecPlatform) GetManagerPortHTTP() int {
	port := h.GetManagerComponent().Port.HTTP
	if port == 0 {
		return 8080
	}
	return port
}

func (h *HorusecPlatform) GetManagerLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "manager",
		"app.kubernetes.io/managed-by": "horusec",
	}
}

func (h *HorusecPlatform) GetManagerReplicaCount() *int32 {
	if !h.GetManagerAutoscaling().Enabled {
		count := h.GetManagerComponent().ReplicaCount
		return &count
	}
	return nil
}

func (h *HorusecPlatform) GetManagerDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetManagerName(), h.GetManagerPortHTTP())
}

func (h *HorusecPlatform) GetManagerRegistry() string {
	registry := h.GetManagerComponent().Container.Image.Registry
	if registry == "" {
		return "docker.io/"
	}
	return registry
}

func (h *HorusecPlatform) GetManagerRepository() string {
	repository := h.GetManagerComponent().Container.Image.Repository
	if repository == "" {
		return "horuszup/horusec-manager"
	}
	return repository
}

func (h *HorusecPlatform) GetManagerTag() string {
	tag := h.GetManagerComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}

func (h *HorusecPlatform) GetManagerImage() string {
	return fmt.Sprintf("%s%s:%s", h.GetManagerRegistry(), h.GetManagerRepository(), h.GetManagerTag())
}

func (h *HorusecPlatform) GetAnalyticEndpoint() string {
	host := h.GetAnalyticComponent().Ingress.Host
	if host == "" {
		return h.GetAnalyticHost()
	}

	return host
}

func (h *HorusecPlatform) GetAPIEndpoint() string {
	host := h.GetAPIComponent().Ingress.Host
	if host == "" {
		return h.GetAPIHost()
	}

	return host
}

func (h *HorusecPlatform) GetAuthEndpoint() string {
	host := h.GetAuthComponent().Ingress.Host
	if host == "" {
		return h.GetAuthHost()
	}

	return host
}

func (h *HorusecPlatform) GetCoreEndpoint() string {
	host := h.GetCoreComponent().Ingress.Host
	if host == "" {
		return h.GetCoreHost()
	}

	return host
}

func (h *HorusecPlatform) GetWebhookEndpoint() string {
	host := h.GetWebhookComponent().Ingress.Host
	if host == "" {
		return h.GetWebhookHost()
	}

	return host
}

func (h *HorusecPlatform) GetManagerHost() string {
	host := h.Spec.Components.Manager.Ingress.Host
	if host == "" {
		return "manager.local"
	}

	return host
}

func (h *HorusecPlatform) IsManagerIngressEnabled() bool {
	enabled := h.Spec.Components.Manager.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
