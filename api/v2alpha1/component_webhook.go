package v2alpha1

import (
	"fmt"
	"reflect"
)

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
func (h *HorusecPlatform) GetWebhookDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetWebhookName(), h.GetWebhookPortHTTP())
}
func (h *HorusecPlatform) GetWebhookImage() string {
	image := h.GetWebhookComponent().Container.Image
	if reflect.ValueOf(image).IsZero() {
		return fmt.Sprintf("docker.io/horuszup/horusec-webhook:%s", h.GetLatestVersion())
	}

	return fmt.Sprintf("%s:%s", image.Registry, image.Tag)
}

func (h *HorusecPlatform) GetWebhookHost() string {
	host := h.Spec.Components.Webhook.Ingress.Host
	if host == "" {
		return "webhook.local"
	}

	return host
}
