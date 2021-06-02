package v2alpha1

import "fmt"

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
		return "/" + h.GetCoreName()
	}
	return path
}
func (h *HorusecPlatform) GetCorePortHTTP() int {
	port := h.GetCoreComponent().Port.HTTP
	if port == 0 {
		return 8008
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
