package v2alpha1

import "fmt"

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
		return "/" + h.GetAPIName()
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
