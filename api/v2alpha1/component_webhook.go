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
		count := h.GetWebhookComponent().ReplicaCount
		return &count
	}
	return nil
}

func (h *HorusecPlatform) GetWebhookDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetWebhookName(), h.GetWebhookPortHTTP())
}

func (h *HorusecPlatform) GetWebhookRegistry() string {
	registry := h.GetWebhookComponent().Container.Image.Registry
	if registry == "" {
		return "docker.io/"
	}
	return registry
}

func (h *HorusecPlatform) GetWebhookRepository() string {
	repository := h.GetWebhookComponent().Container.Image.Repository
	if repository == "" {
		return "horuszup/horusec-webhook"
	}
	return repository
}

func (h *HorusecPlatform) GetWebhookTag() string {
	tag := h.GetWebhookComponent().Container.Image.Tag
	if tag == "" {
		return h.GetLatestVersion()
	}
	return tag
}

func (h *HorusecPlatform) GetWebhookImage() string {
	return fmt.Sprintf("%s/%s:%s", h.GetWebhookRegistry(), h.GetWebhookRepository(), h.GetWebhookTag())
}

func (h *HorusecPlatform) GetWebhookHost() string {
	host := h.Spec.Components.Webhook.Ingress.Host
	if host == "" {
		return "webhook.local"
	}

	return host
}

func (h *HorusecPlatform) IsWebhookIngressEnabled() bool {
	enabled := h.Spec.Components.Webhook.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}
