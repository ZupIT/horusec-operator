package v2alpha1

import (
	"fmt"
)

func (h *HorusecPlatform) GetAuthComponent() Auth {
	return h.Spec.Components.Auth
}
func (h *HorusecPlatform) GetAuthAutoscaling() Autoscaling {
	return h.GetAuthComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetAuthName() string {
	name := h.GetAuthComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-auth", h.GetName())
	}
	return name
}
func (h *HorusecPlatform) GetAuthPath() string {
	path := h.GetAuthComponent().Ingress.Path
	if path == "" {
		return "/auth"
	}
	return path
}
func (h *HorusecPlatform) GetAuthPortHTTP() int {
	port := h.GetAuthComponent().Port.HTTP
	if port == 0 {
		return 8006
	}
	return port
}
func (h *HorusecPlatform) GetAuthPortGRPC() int {
	port := h.GetAuthComponent().Port.Grpc
	if port == 0 {
		return 8007
	}
	return port
}
func (h *HorusecPlatform) GetAuthLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "auth",
		"app.kubernetes.io/managed-by": "horusec",
	}
}
func (h *HorusecPlatform) GetAuthReplicaCount() *int32 {
	if !h.GetAuthAutoscaling().Enabled {
		return h.GetAuthComponent().ReplicaCount
	}
	return nil
}
func (h *HorusecPlatform) GetAuthDefaultHTTPURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetAuthName(), h.GetAuthPortHTTP())
}
func (h *HorusecPlatform) GetAuthDefaultGRPCURL() string {
	return fmt.Sprintf("%s:%v", h.GetAuthName(), h.GetAuthPortGRPC())
}
func (h *HorusecPlatform) GetAuthRegistry() string {
	registry := h.GetAuthComponent().Container.Image.Registry
	if registry == "" {
		return "docker.io/horuszup/horusec-auth"
	}
	return registry
}
func (h *HorusecPlatform) GetAuthTag() string {
	tag := h.GetAuthComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}
func (h *HorusecPlatform) GetAuthImage() string {
	return fmt.Sprintf("%s:%s", h.GetAuthRegistry(), h.GetAuthTag())
}

func (h *HorusecPlatform) GetAuthHost() string {
	host := h.Spec.Components.Auth.Ingress.Host
	if host == "" {
		return "auth.local"
	}

	return host
}

func (h *HorusecPlatform) IsAuthIngressEnabled() bool {
	enabled := h.Spec.Components.Auth.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
