package v2alpha1

func (h *HorusecPlatform) GetManagerComponent() Manager {
	return h.Spec.Components.Manager
}
func (h *HorusecPlatform) GetManagerAutoscaling() Autoscaling {
	return h.GetManagerComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetManagerName() string {
	name := h.GetManagerComponent().Name
	if name == "" {
		return "manager"
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
