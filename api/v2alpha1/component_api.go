package v2alpha1

import (
	"fmt"
	"reflect"
)

func (h *HorusecPlatform) GetAPIComponent() Api {
	return h.Spec.Components.Api
}
func (h *HorusecPlatform) GetAPIAutoscaling() Autoscaling {
	return h.GetAPIComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetAPIName() string {
	name := h.GetAPIComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-api", h.GetName())
	}
	return name
}
func (h *HorusecPlatform) GetAPIPath() string {
	path := h.GetAPIComponent().Ingress.Path
	if path == "" {
		return "/api"
	}
	return path
}
func (h *HorusecPlatform) GetAPIPortHTTP() int {
	port := h.GetAPIComponent().Port.HTTP
	if port == 0 {
		return 8000
	}
	return port
}
func (h *HorusecPlatform) GetApiLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "api",
		"app.kubernetes.io/managed-by": "horusec",
	}
}
func (h *HorusecPlatform) GetAPIReplicaCount() *int32 {
	if !h.GetAPIAutoscaling().Enabled {
		return h.GetAPIComponent().ReplicaCount
	}
	return nil
}
func (h *HorusecPlatform) GetAPIDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetAPIName(), h.GetAPIPortHTTP())
}
func (h *HorusecPlatform) GetAPIImage() string {
	image := h.GetAPIComponent().Container.Image
	if reflect.ValueOf(image).IsZero() {
		return fmt.Sprintf("docker.io/horuszup/horusec-api:%s", h.GetLatestVersion())
	}

	return fmt.Sprintf("%s:%s", image.Registry, image.Tag)
}

func (h *HorusecPlatform) GetAPIHost() string {
	host := h.Spec.Components.Api.Ingress.Host
	if host == "" {
		return "api.local"
	}

	return host
}

func (h *HorusecPlatform) IsAPIIngressEnabled() bool {
	enabled := h.Spec.Components.Api.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
