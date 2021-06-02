package v2alpha1

func (h *HorusecPlatform) GetMessagesComponent() Messages {
	return h.Spec.Components.Messages
}
func (h *HorusecPlatform) GetMessagesAutoscaling() Autoscaling {
	return h.GetMessagesComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetMessagesName() string {
	name := h.GetMessagesComponent().Name
	if name == "" {
		return "messages"
	}
	return name
}
func (h *HorusecPlatform) GetMessagesPath() string {
	path := h.GetMessagesComponent().Ingress.Path
	if path == "" {
		return "/" + h.GetMessagesName()
	}
	return path
}
func (h *HorusecPlatform) GetMessagesPortHTTP() int {
	port := h.GetMessagesComponent().Port.HTTP
	if port == 0 {
		return 8004
	}
	return port
}
