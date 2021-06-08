package v2alpha1

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
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
	port := h.GetAuthComponent().Port.GRPC
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
		count := h.GetAuthComponent().ReplicaCount
		return &count
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
		return "docker.io/"
	}
	return registry
}

func (h *HorusecPlatform) GetAuthRepository() string {
	repository := h.GetAuthComponent().Container.Image.Repository
	if repository == "" {
		return "horuszup/horusec-auth"
	}
	return repository
}

func (h *HorusecPlatform) GetAuthTag() string {
	tag := h.GetAuthComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}

func (h *HorusecPlatform) GetAuthImage() string {
	return fmt.Sprintf("%s/%s:%s", h.GetAuthRegistry(), h.GetAuthRepository(), h.GetAuthTag())
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

func (h *HorusecPlatform) GetAuthAdminData() string {
	email := h.Spec.Global.Administrator.Email
	if !h.Spec.Global.Administrator.Enabled || email == "" {
		return ""
	}

	return fmt.Sprintf(
		"{\"username\": \"$(HORUSEC_ADMIN_USERNAME)\", \"email\":\"%s\", \"password\":\"$(HORUSEC_ADMIN_PASSWORD)\"}",
		email)
}

func (h *HorusecPlatform) GetAuthDefaultUserData() string {
	email := h.Spec.Components.Auth.DefaultUser.Email
	if !h.Spec.Global.Administrator.Enabled || email == "" {
		return ""
	}

	return fmt.Sprintf(
		"{\"username\": \"$(HORUSEC_DEFAULT_USER_USERNAME)\", \"email\":\"%s\", \"password\":\"$(HORUSEC_DEFAULT_USER_PASSWORD)\"}",
		email)
}

func (h *HorusecPlatform) GetAuthAdminUsernameEnv() v1.EnvVar {
	if h.Spec.Global.Administrator.Enabled {
		return v1.EnvVar{}
	}

	return h.NewEnvFromSecret("HORUSEC_ADMIN_USERNAME", h.Spec.Global.Administrator.Credentials.User.KeyRef)
}

func (h *HorusecPlatform) GetAuthAdminPasswordEnv() v1.EnvVar {
	if h.Spec.Global.Administrator.Enabled {
		return v1.EnvVar{}
	}

	return h.NewEnvFromSecret("HORUSEC_ADMIN_PASSWORD", h.Spec.Global.Administrator.Credentials.Password.KeyRef)
}

func (h *HorusecPlatform) GetAuthDefaultUserUsername() v1.EnvVar {
	if h.Spec.Global.Administrator.Enabled {
		return v1.EnvVar{}
	}

	return h.NewEnvFromSecret("HORUSEC_DEFAULT_USER_USERNAME", h.Spec.Components.Auth.DefaultUser.Credentials.User.KeyRef)
}

func (h *HorusecPlatform) GetAuthDefaultUserPassword() v1.EnvVar {
	if h.Spec.Global.Administrator.Enabled {
		return v1.EnvVar{}
	}

	return h.NewEnvFromSecret("HORUSEC_DEFAULT_USER_PASSWORD", h.Spec.Components.Auth.DefaultUser.Credentials.Password.KeyRef)
}

func (h *HorusecPlatform) GetAuthKeycloakClientSecret() v1.EnvVar {
	if h.Spec.Global.Keycloak.Realm == "" || h.Spec.Global.Keycloak.PublicURL == "" ||
		h.Spec.Global.Keycloak.InternalURL == "" {
		return v1.EnvVar{}
	}

	return h.NewEnvFromSecret("HORUSEC_KEYCLOAK_CLIENT_SECRET", h.Spec.Global.Keycloak.Clients.Confidential.SecretKeyRef)
}

func (h *HorusecPlatform) GetAuthOptionalEnvs() []v1.EnvVar {
	return []v1.EnvVar{
		h.GetAuthAdminUsernameEnv(),
		h.GetAuthAdminPasswordEnv(),
		h.GetAuthDefaultUserUsername(),
		h.GetAuthDefaultUserPassword(),
		h.GetAuthKeycloakClientSecret(),
	}
}
