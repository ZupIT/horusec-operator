// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2alpha1

import (
	"fmt"
)

func (h *HorusecPlatform) GetManagerComponent() ExposableComponent {
	return h.Spec.Components.Manager
}

func (h *HorusecPlatform) GetManagerAutoscaling() Autoscaling {
	return h.GetManagerComponent().Pod.Autoscaling
}

func (h *HorusecPlatform) GetManagerName() string {
	name := h.GetManagerComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-manager", h.GetName())
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

func (h *HorusecPlatform) GetManagerLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "manager",
		"app.kubernetes.io/managed-by": "horusec",
	}
}

func (h *HorusecPlatform) GetManagerReplicaCount() *int32 {
	if !h.GetManagerAutoscaling().Enabled {
		count := h.GetManagerComponent().ReplicaCount
		return &count
	}
	return nil
}

func (h *HorusecPlatform) GetManagerDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetManagerName(), h.GetManagerPortHTTP())
}

func (h *HorusecPlatform) GetManagerRegistry() string {
	registry := h.GetManagerComponent().Container.Image.Registry
	if registry == "" {
		return "docker.io/"
	}
	return registry
}

func (h *HorusecPlatform) GetManagerRepository() string {
	repository := h.GetManagerComponent().Container.Image.Repository
	if repository == "" {
		return "horuszup/horusec-manager"
	}
	return repository
}

func (h *HorusecPlatform) GetManagerTag() string {
	tag := h.GetManagerComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}

func (h *HorusecPlatform) GetManagerImage() string {
	return fmt.Sprintf("%s/%s:%s", h.GetManagerRegistry(), h.GetManagerRepository(), h.GetManagerTag())
}

func (h *HorusecPlatform) GetAnalyticEndpoint() string {
	host := h.GetAnalyticHost()
	schema := h.GetAnalyticSchema()
	return fmt.Sprintf("%s:\\/\\/%s", schema, host)
}

func (h *HorusecPlatform) GetAPIEndpoint() string {
	host := h.GetAPIHost()
	schema := h.GetAPISchema()
	return fmt.Sprintf("%s:\\/\\/%s", schema, host)
}

func (h *HorusecPlatform) GetAuthEndpoint() string {
	host := h.GetAuthHost()
	schema := h.GetAuthSchema()
	return fmt.Sprintf("%s:\\/\\/%s", schema, host)
}

func (h *HorusecPlatform) GetVulnerabilityEndpoint() string {
	host := h.GetVulnerabilityHost()
	schema := h.GetVulnerabilitySchema()
	return fmt.Sprintf("%s:\\/\\/%s", schema, host)
}

func (h *HorusecPlatform) GetCoreEndpoint() string {
	host := h.GetCoreHost()
	schema := h.GetCoreSchema()
	return fmt.Sprintf("%s:\\/\\/%s", schema, host)
}

func (h *HorusecPlatform) GetWebhookEndpoint() string {
	host := h.GetWebhookHost()
	schema := h.GetWebhookSchema()
	return fmt.Sprintf("%s:\\/\\/%s", schema, host)
}

func (h *HorusecPlatform) GetManagerHost() string {
	host := h.Spec.Components.Manager.Ingress.Host
	if host == "" {
		return "manager.local"
	}

	return host
}

func (h *HorusecPlatform) IsManagerIngressEnabled() bool {
	enabled := h.Spec.Components.Manager.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}

func (h *HorusecPlatform) GetAnalyticSchema() string {
	component := h.Spec.Components.Analytic
	if component.Ingress.TLS.SecretName != "" {
		return "https"
	}
	return "http"
}

func (h *HorusecPlatform) GetAPISchema() string {
	component := h.Spec.Components.API
	if component.Ingress.TLS.SecretName != "" {
		return "https"
	}
	return "http"
}

func (h *HorusecPlatform) GetAuthSchema() string {
	component := h.Spec.Components.Auth
	if component.Ingress.TLS.SecretName != "" {
		return "https"
	}
	return "http"
}

func (h *HorusecPlatform) GetVulnerabilitySchema() string {
	component := h.Spec.Components.Vulnerability
	if component.Ingress.TLS.SecretName != "" {
		return "https"
	}
	return "http"
}

func (h *HorusecPlatform) GetCoreSchema() string {
	component := h.Spec.Components.Core
	if component.Ingress.TLS.SecretName != "" {
		return "https"
	}
	return "http"
}

func (h *HorusecPlatform) GetWebhookSchema() string {
	component := h.Spec.Components.Webhook
	if component.Ingress.TLS.SecretName != "" {
		return "https"
	}
	return "http"
}
