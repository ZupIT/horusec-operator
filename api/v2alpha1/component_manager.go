package v2alpha1

import (
	"fmt"
	"reflect"
)

func (h *HorusecPlatform) GetManagerComponent() Manager {
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
		return h.GetManagerComponent().ReplicaCount
	}
	return nil
}
func (h *HorusecPlatform) GetManagerDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetManagerName(), h.GetManagerPortHTTP())
}
func (h *HorusecPlatform) GetManagerImage() string {
	image := h.GetManagerComponent().Container.Image
	if reflect.ValueOf(image).IsZero() {
		return fmt.Sprintf("docker.io/horuszup/horusec-manager:%s", h.GetLatestVersion())
	}

	return fmt.Sprintf("%s:%s", image.Registry, image.Tag)
}
func (h *HorusecPlatform) GetAnalyticEndpoint() string {
	host := h.GetAnalyticComponent().Ingress.Host
	if host == "" {
		return h.GetAnalyticDefaultURL()
	}

	return host
}
func (h *HorusecPlatform) GetAPIEndpoint() string {
	host := h.GetAPIComponent().Ingress.Host
	if host == "" {
		return h.GetAPIDefaultURL()
	}

	return host
}
func (h *HorusecPlatform) GetAuthEndpoint() string {
	host := h.GetAuthComponent().Ingress.Host
	if host == "" {
		return h.GetAuthDefaultHTTPURL()
	}

	return host
}
func (h *HorusecPlatform) GetCoreEndpoint() string {
	host := h.GetCoreComponent().Ingress.Host
	if host == "" {
		return h.GetCoreDefaultURL()
	}

	return host
}
func (h *HorusecPlatform) GetWebhookEndpoint() string {
	host := h.GetWebhookComponent().Ingress.Host
	if host == "" {
		return h.GetWebhookDefaultURL()
	}

	return host
}
