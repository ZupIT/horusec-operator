package v2alpha1

import "fmt"

func (h *HorusecPlatform) GetWebhookComponent() Webhook {
	return h.Spec.Components.Webhook
}
func (h *HorusecPlatform) GetWebhookAutoscaling() Autoscaling {
	return h.GetWebhookComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetWebhookName() string {
	name := h.GetWebhookComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-webhook", h.GetName())
	}
	return name
}
func (h *HorusecPlatform) GetWebhookPath() string {
	path := h.GetWebhookComponent().Ingress.Path
	if path == "" {
		return "/webhook"
	}
	return path
}
func (h *HorusecPlatform) GetWebhookPortHTTP() int {
	port := h.GetWebhookComponent().Port.HTTP
	if port == 0 {
		return 8004
	}
	return port
}
func (h *HorusecPlatform) GetWebhookLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "webhook",
		"app.kubernetes.io/managed-by": "horusec",
	}
}
func (h *HorusecPlatform) GetWebhookReplicaCount() *int32 {
	if !h.GetWebhookAutoscaling().Enabled {
		return h.GetWebhookComponent().ReplicaCount
	}
	return nil
}
