package v2alpha1

import (
	"fmt"
)

func (h *HorusecPlatform) GetMessagesComponent() Messages {
	return h.Spec.Components.Messages
}

func (h *HorusecPlatform) GetMessagesAutoscaling() Autoscaling {
	return h.GetMessagesComponent().Pod.Autoscaling
}

func (h *HorusecPlatform) GetMessagesName() string {
	name := h.GetMessagesComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-messages", h.GetName())
	}
	return name
}

func (h *HorusecPlatform) GetMessagesPath() string {
	path := h.GetMessagesComponent().Ingress.Path
	if path == "" {
		return "/messages"
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

func (h *HorusecPlatform) GetMessagesLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "messages",
		"app.kubernetes.io/managed-by": "horusec",
	}
}

func (h *HorusecPlatform) GetMessagesReplicaCount() *int32 {
	if !h.GetMessagesAutoscaling().Enabled {
		count := h.GetMessagesComponent().ReplicaCount
		return &count
	}
	return nil
}

func (h *HorusecPlatform) GetMessagesDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetMessagesName(), h.GetMessagesPortHTTP())
}

func (h *HorusecPlatform) GetMessagesRegistry() string {
	registry := h.GetMessagesComponent().Container.Image.Registry
	if registry == "" {
		return "docker.io/"
	}
	return registry
}

func (h *HorusecPlatform) GetMessagesRepository() string {
	repository := h.GetMessagesComponent().Container.Image.Repository
	if repository == "" {
		return "horuszup/horusec-messages"
	}
	return repository
}

func (h *HorusecPlatform) GetMessagesTag() string {
	tag := h.GetMessagesComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}

func (h *HorusecPlatform) GetMessagesImage() string {
	return fmt.Sprintf("%s/%s:%s", h.GetMessagesRegistry(), h.GetMessagesRepository(), h.GetMessagesTag())
}

func (h *HorusecPlatform) GetMessagesMailServer() MailServer {
	return h.GetMessagesComponent().MailServer
}

func (h *HorusecPlatform) GetMessagesHost() string {
	host := h.Spec.Components.Messages.Ingress.Host
	if host == "" {
		return "messages.local"
	}

	return host
}

func (h *HorusecPlatform) IsMessagesIngressEnabled() bool {
	component := h.GetMessagesComponent()
	if !component.Enabled {
		return false
	}

	enabled := component.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
