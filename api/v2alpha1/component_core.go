package v2alpha1

func (h *HorusecPlatform) GetCoreComponent() Core {
	return h.Spec.Components.Core
}
func (h *HorusecPlatform) GetCoreAutoscaling() Autoscaling {
	return h.GetCoreComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetCoreName() string {
	name := h.GetCoreComponent().Name
	if name == "" {
		return "core"
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
