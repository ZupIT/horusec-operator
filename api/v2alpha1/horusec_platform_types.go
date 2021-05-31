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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HorusecPlatformSpec defines the desired state of HorusecPlatform
type HorusecPlatformSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Components Components `json:"components,omitempty"`
	Global     Global     `json:"global,omitempty"`
}

type Global struct {
	Administrator Administrator `json:"administrator,omitempty"`
	Broker        Broker        `json:"broker,omitempty"`
	Database      Database      `json:"database,omitempty"`
	Jwt           Jwt           `json:"jwt,omitempty"`
	Keycloak      Keycloak      `json:"keycloak,omitempty"`
}

type Keycloak struct {
	Clients     Clients `json:"clients,omitempty"`
	InternalURL string  `json:"internalURL,omitempty"`
	Otp         bool    `json:"otp,omitempty"`
	PublicURL   string  `json:"publicURL,omitempty"`
	Realm       string  `json:"realm,omitempty"`
}

type Clients struct {
	Confidential Confidential `json:"clients,omitempty"`
	Public       Public       `json:"public,omitempty"`
}

type Confidential struct {
	ID           string       `json:"id,omitempty"`
	SecretKeyRef SecretKeyRef `json:"secretKeyRef,omitempty"`
}

type Public struct {
	ID string `json:"id,omitempty"`
}

type Jwt struct {
	SecretKeyRef SecretKeyRef `json:"secretKeyRef,omitempty"`
}

type Broker struct {
	Host     string   `json:"host,omitempty"`
	Password Password `json:"password,omitempty"`
	Port     int      `json:"port,omitempty"`
	User     User     `json:"user,omitempty"`
}

type Administrator struct {
	Email    string   `json:"email,omitempty"`
	Enabled  bool     `json:"enabled,omitempty"`
	Password Password `json:"password,omitempty"`
	User     User     `json:"user,omitempty"`
}

//nolint:golint, stylecheck // no need to be API in uppercase
type Components struct {
	Analytic      Analytic      `json:"analytic,omitempty"`
	Api           Api           `json:"api,omitempty"`
	Auth          Auth          `json:"auth,omitempty"`
	Core          Core          `json:"core,omitempty"`
	Manager       Manager       `json:"manager,omitempty"`
	Messages      Messages      `json:"messages,omitempty"`
	Vulnerability Vulnerability `json:"vulnerability,omitempty"`
	Webhook       Webhook       `json:"webhook,omitempty"`
}

type Analytic struct {
	BaseComponent
	Database Database `json:"database,omitempty"`
}

//nolint:golint, stylecheck // no need to be API in uppercase
type Api struct {
	BaseComponent
	Database Database `json:"database,omitempty"`
}

type Auth struct {
	BaseComponent
	Type     string   `json:"type,omitempty"`
	Database Database `json:"database,omitempty"`
}

type Core struct {
	BaseComponent
	Database Database `json:"database,omitempty"`
}

type Manager struct {
	BaseComponent
}

type Messages struct {
	BaseComponent
	Enabled    bool       `json:"enabled,omitempty"`
	MailServer MailServer `json:"mailServer,omitempty"`
}

type Vulnerability struct {
	BaseComponent
	Database Database `json:"database,omitempty"`
}

type Webhook struct {
	BaseComponent
	Database Database `json:"database,omitempty"`
}

type Container struct {
	Image           Image           `json:"image,omitempty"`
	LivenessProbe   *string         `json:"livenessProbe,omitempty"`
	ReadinessProbe  *string         `json:"readinessProbe,omitempty"`
	Resources       *string         `json:"resources,omitempty"`
	SecurityContext SecurityContext `json:"securityContext,omitempty"`
}

type Image struct {
	PullPolicy  string   `json:"pullPolicy,omitempty"`
	PullSecrets []string `json:"pullSecrets,omitempty"`
	Registry    string   `json:"registry,omitempty"`
	Repository  string   `json:"repository,omitempty"`
	Tag         string   `json:"tag,omitempty"`
}

type SecurityContext struct {
	Enabled      bool `json:"enabled,omitempty"`
	RunAsNonRoot bool `json:"runAsNonRoot,omitempty"`
	RunAsUser    int  `json:"runAsUser,omitempty"`
	FsGroup      int  `json:"fsGroup,omitempty"`
}

type Database struct {
	Dialect  string   `json:"dialect,omitempty"`
	Host     string   `json:"host,omitempty"`
	LogMode  bool     `json:"logMode,omitempty"`
	Name     string   `json:"name,omitempty"`
	Password Password `json:"password,omitempty"`
	Port     int      `json:"port,omitempty"`
	SslMode  bool     `json:"sslMode,omitempty"`
	User     User     `json:"user,omitempty"`
}

type Password struct {
	SecretKeyRef SecretKeyRef `json:"secretKeyRef,omitempty"`
}

type User struct {
	SecretKeyRef SecretKeyRef `json:"secretKeyRef,omitempty"`
}

type SecretKeyRef struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type Ingress struct {
	Enabled bool    `json:"enabled,omitempty"`
	Host    string  `json:"host,omitempty"`
	Path    string  `json:"path,omitempty"`
	TLS     *string `json:"tls,omitempty"`
}

type Pod struct {
	Autoscaling     Autoscaling     `json:"autoscaling,omitempty"`
	SecurityContext SecurityContext `json:"securityContext,omitempty"`
}

type Autoscaling struct {
	Enabled      bool   `json:"enabled,omitempty"`
	MaxReplicas  int32  `json:"maxReplicas,omitempty"`
	MinReplicas  *int32 `json:"minReplicas,omitempty"`
	TargetCPU    *int32 `json:"targetCPU,omitempty"`
	TargetMemory *int32 `json:"targetMemory,omitempty"`
}

type Port struct {
	HTTP int `json:"http,omitempty"`
	Grpc int `json:"grpc,omitempty"`
}

type MailServer struct {
	Host     string   `json:"host,omitempty"`
	Password Password `json:"password,omitempty"`
	User     User     `json:"user,omitempty"`
	Port     int      `json:"port,omitempty"`
}

// HorusecPlatformStatus defines the observed state of HorusecPlatform
type HorusecPlatformStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HorusecPlatform is the Schema for the horusecs API
type HorusecPlatform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HorusecPlatformSpec   `json:"spec,omitempty"`
	Status HorusecPlatformStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HorusecPlatformList contains a list of HorusecPlatform
type HorusecPlatformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HorusecPlatform `json:"items"`
}

// nolint // autogenerated by operator-sdk
func init() {
	SchemeBuilder.Register(&HorusecPlatform{}, &HorusecPlatformList{})
}
