package v2alpha1

func (h *HorusecPlatform) GetWebhookComponent() Webhook {
	return h.Spec.Components.Webhook
}
func (h *HorusecPlatform) GetWebhookAutoscaling() Autoscaling {
	return h.GetWebhookComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetWebhookName() string {
	name := h.GetWebhookComponent().Name
	if name == "" {
		return "webhook"
	}
	return name
}
func (h *HorusecPlatform) GetWebhookPath() string {
	path := h.GetWebhookComponent().Ingress.Path
	if path == "" {
		return "/" + h.GetWebhookName()
	}
	return path
}
func (h *HorusecPlatform) GetWebhookPortHTTP() int {
	port := h.GetWebhookComponent().Port.HTTP
	if port == 0 {
		return 8005
	}
	return port
}
