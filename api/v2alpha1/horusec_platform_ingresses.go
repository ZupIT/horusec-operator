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

import "net"

//+k8s:deepcopy-gen=false
type HorusecIngress interface {
	IsEnabled() bool
	GetPath() string
	GetHost() string
	GetName() string
	GetSecretName() string
}

func (in *HorusecPlatform) Ingresses() []HorusecIngress {
	return []HorusecIngress{
		&in.Spec.Components.Analytic,
		&in.Spec.Components.API,
		&in.Spec.Components.Auth,
		&in.Spec.Components.Core,
		&in.Spec.Components.Manager,
		&in.Spec.Components.Messages,
		&in.Spec.Components.Vulnerability,
		&in.Spec.Components.Webhook,
	}
}

func (in *Messages) IsEnabled() bool {
	if !in.Enabled {
		return false
	}

	return in.ExposableComponent.IsEnabled()
}

func (in *ExposableComponent) IsEnabled() bool {
	enabled := in.Ingress.Enabled
	if enabled == nil {
		return true
	}

	return *enabled
}

func (in *ExposableComponent) GetPath() string { return in.Ingress.Path }

func (in *ExposableComponent) GetHost() string {
	host := in.Ingress.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	return host
}

func (in *ExposableComponent) GetName() string {
	return in.Component.Name
}

func (in *ExposableComponent) GetSecretName() string {
	return in.Ingress.TLS.SecretName
}
