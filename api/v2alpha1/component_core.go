package v2alpha1

import (
	"fmt"
	"reflect"
)

func (h *HorusecPlatform) GetCoreComponent() Core {
	return h.Spec.Components.Core
}
func (h *HorusecPlatform) GetCoreAutoscaling() Autoscaling {
	return h.GetCoreComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetCoreName() string {
	name := h.GetCoreComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-core", h.GetName())
	}
	return name
}
func (h *HorusecPlatform) GetCorePath() string {
	path := h.GetCoreComponent().Ingress.Path
	if path == "" {
		return "/core"
	}
	return path
}
func (h *HorusecPlatform) GetCorePortHTTP() int {
	port := h.GetCoreComponent().Port.HTTP
	if port == 0 {
		return 8003
	}
	return port
}
func (h *HorusecPlatform) GetCoreLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "core",
		"app.kubernetes.io/managed-by": "horusec",
	}
}
func (h *HorusecPlatform) GetCoreReplicaCount() *int32 {
	if !h.GetCoreAutoscaling().Enabled {
		return h.GetCoreComponent().ReplicaCount
	}
	return nil
}
func (h *HorusecPlatform) GetCoreDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetCoreName(), h.GetCorePortHTTP())
}
func (h *HorusecPlatform) GetCoreImage() string {
	image := h.GetCoreComponent().Container.Image
	if reflect.ValueOf(image).IsZero() {
		return fmt.Sprintf("docker.io/horuszup/horusec-core:%s", h.GetLatestVersion())
	}

	return fmt.Sprintf("%s:%s", image.Registry, image.Tag)
}

func (h *HorusecPlatform) GetCoreHost() string {
	host := h.Spec.Components.Core.Ingress.Host
	if host == "" {
		return "core.local"
	}

	return host
}

func (h *HorusecPlatform) IsCoreIngressEnabled() bool {
	enabled := h.Spec.Components.Core.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
