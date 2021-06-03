package v2alpha1

import (
	"fmt"
	"reflect"
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
		return h.GetMessagesComponent().ReplicaCount
	}
	return nil
}
func (h *HorusecPlatform) GetMessagesDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetMessagesName(), h.GetMessagesPortHTTP())
}
func (h *HorusecPlatform) GetMessagesImage() string {
	image := h.GetMessagesComponent().Container.Image
	if reflect.ValueOf(image).IsZero() {
		return fmt.Sprintf("docker.io/horuszup/horusec-messages:%s", h.GetLatestVersion())
	}

	return fmt.Sprintf("%s:%s", image.Registry, image.Tag)
}
