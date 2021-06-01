package v2alpha1

func (h *HorusecPlatform) GetAnalyticComponent() Analytic {
	return h.Spec.Components.Analytic
}
func (h *HorusecPlatform) GetAPIComponent() Api {
	return h.Spec.Components.Api
}
func (h *HorusecPlatform) GetAuthComponent() Auth {
	return h.Spec.Components.Auth
}
func (h *HorusecPlatform) GetCoreComponent() Core {
	return h.Spec.Components.Core
}
func (h *HorusecPlatform) GetManagerComponent() Manager {
	return h.Spec.Components.Manager
}
func (h *HorusecPlatform) GetMessagesComponent() Messages {
	return h.Spec.Components.Messages
}
func (h *HorusecPlatform) GetVulnerabilityComponent() Vulnerability {
	return h.Spec.Components.Vulnerability
}
func (h *HorusecPlatform) GetWebhookComponent() Webhook {
	return h.Spec.Components.Webhook
}
func (h *HorusecPlatform) GetAnalyticAutoscaling() Autoscaling {
	return h.GetAnalyticComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetAPIAutoscaling() Autoscaling {
	return h.GetAPIComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetAuthAutoscaling() Autoscaling {
	return h.GetAuthComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetCoreAutoscaling() Autoscaling {
	return h.GetCoreComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetManagerAutoscaling() Autoscaling {
	return h.GetManagerComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetMessagesAutoscaling() Autoscaling {
	return h.GetMessagesComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetVulnerabilityAutoscaling() Autoscaling {
	return h.GetVulnerabilityComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetWebhookAutoscaling() Autoscaling {
	return h.GetWebhookComponent().Pod.Autoscaling
}
