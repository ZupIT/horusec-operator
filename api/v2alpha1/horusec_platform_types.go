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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1/state"
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
	Broker       Broker   `json:"broker,omitempty"`
	Database     Database `json:"database,omitempty"`
	JWT          JWT      `json:"jwt,omitempty"`
	Keycloak     Keycloak `json:"keycloak,omitempty"`
	Ldap         Ldap     `json:"ldap,omitempty"`
	GrpcUseCerts bool     `json:"grpcUseCerts,omitempty"`
}

type Keycloak struct {
	Clients     Clients `json:"clients,omitempty"`
	InternalURL string  `json:"internalURL,omitempty"`
	Otp         bool    `json:"otp,omitempty"`
	PublicURL   string  `json:"publicURL,omitempty"`
	Realm       string  `json:"realm,omitempty"`
}

type Ldap struct {
	Base               string           `json:"base,omitempty"`
	Host               string           `json:"host,omitempty"`
	Port               int              `json:"port,omitempty"`
	UseSSL             bool             `json:"useSsl,omitempty"`
	SkipTLS            bool             `json:"skipTls,omitempty"`
	InsecureSkipVerify bool             `json:"insecureSkipVerify,omitempty"`
	BindDN             string           `json:"bindDn,omitempty"`
	BindPassword       LdapBindPassword `json:"bindPassword,omitempty"`
	UserFilter         string           `json:"userFilter,omitempty"`
	AdminGroup         string           `json:"adminGroup,omitempty"`
}

type LdapBindPassword struct {
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

type Clients struct {
	Confidential Confidential `json:"confidential,omitempty"`
	Public       Public       `json:"public,omitempty"`
}

type Confidential struct {
	ID           string                    `json:"id,omitempty"`
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

type Public struct {
	ID string `json:"id,omitempty"`
}

type JWT struct {
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

type Broker struct {
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	Credentials `json:",inline,omitempty"`
}

type UserInfo struct {
	Email       string `json:"email,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	Credentials `json:",inline,omitempty"`
}

//nolint:golint, stylecheck // no need to be API in uppercase
type Components struct {
	Analytic      Analytic           `json:"analytic,omitempty"`
	API           ExposableComponent `json:"api,omitempty"`
	Auth          Auth               `json:"auth,omitempty"`
	Core          ExposableComponent `json:"core,omitempty"`
	Manager       ExposableComponent `json:"manager,omitempty"`
	Messages      Messages           `json:"messages,omitempty"`
	Vulnerability ExposableComponent `json:"vulnerability,omitempty"`
	Webhook       Webhook            `json:"webhook,omitempty"`
}

type Webhook struct {
	Timeout            int `json:"timeout,omitempty"`
	ExposableComponent `json:",inline,omitempty"`
}

type Analytic struct {
	ExposableComponent `json:",inline,omitempty"`
	Database           Database `json:"database,omitempty"`
}

type Auth struct {
	Type               AuthType `json:"type,omitempty"`
	User               User     `json:"user,omitempty"`
	ExposableComponent `json:",inline,omitempty"`
}

type User struct {
	Administrator UserInfo `json:"administrator,omitempty"`
	Default       UserInfo `json:"default,omitempty"`
}

type AuthType string

type Messages struct {
	Enabled            bool       `json:"enabled,omitempty"`
	MailServer         MailServer `json:"mailServer,omitempty"`
	EmailFrom          string     `json:"emailFrom,omitempty"`
	ExposableComponent `json:",inline,omitempty"`
}

type Container struct {
	Image           Image                       `json:"image,omitempty"`
	LivenessProbe   corev1.Probe                `json:"livenessProbe,omitempty"`
	ReadinessProbe  corev1.Probe                `json:"readinessProbe,omitempty"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty"`
	SecurityContext ContainerSecurityContext    `json:"securityContext,omitempty"`
}

type Image struct {
	PullPolicy  string   `json:"pullPolicy,omitempty"`
	PullSecrets []string `json:"pullSecrets,omitempty"`
	Registry    string   `json:"registry,omitempty"`
	Repository  string   `json:"repository,omitempty"`
	Tag         string   `json:"tag,omitempty"`
}

type ContainerSecurityContext struct {
	Enabled                bool `json:"enabled,omitempty"`
	corev1.SecurityContext `json:",inline,omitempty"`
}

type PodSecurityContext struct {
	Enabled                   bool `json:"enabled,omitempty"`
	corev1.PodSecurityContext `json:",inline,omitempty"`
}

type Database struct {
	Host        string    `json:"host,omitempty"`
	LogMode     bool      `json:"logMode,omitempty"`
	Name        string    `json:"name,omitempty"`
	Port        int       `json:"port,omitempty"`
	SslMode     *bool     `json:"sslMode,omitempty"`
	Migration   Migration `json:"migration,omitempty"`
	Credentials `json:",inline,omitempty"`
}

type Migration struct {
	Image Image `json:"image,omitempty"`
}

type Ingress struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Host    string `json:"host,omitempty"`
	Path    string `json:"path,omitempty"`
	TLS     TLS    `json:"tls,omitempty"`
}

type TLS struct {
	SecretName string `json:"secretName,omitempty"`
}

type Pod struct {
	Autoscaling     Autoscaling        `json:"autoscaling,omitempty"`
	SecurityContext PodSecurityContext `json:"securityContext,omitempty"`
}

type Autoscaling struct {
	Enabled      bool   `json:"enabled,omitempty"`
	MaxReplicas  int32  `json:"maxReplicas,omitempty"`
	MinReplicas  *int32 `json:"minReplicas,omitempty"`
	TargetCPU    *int32 `json:"targetCPU,omitempty"`
	TargetMemory *int32 `json:"targetMemory,omitempty"`
}

type Ports struct {
	HTTP int `json:"http,omitempty"`
	GRPC int `json:"grpc,omitempty"`
}

type MailServer struct {
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	Credentials `json:",inline,omitempty"`
}

type Credentials struct {
	User     SecretRef `json:"user,omitempty"`
	Password SecretRef `json:"password,omitempty"`
}

type SecretRef struct {
	KeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

type Component struct {
	Name         string          `json:"name,omitempty"`
	Port         Ports           `json:"port,omitempty"`
	ExtraEnv     []corev1.EnvVar `json:"extraEnv,omitempty"`
	ReplicaCount int32           `json:"replicaCount,omitempty"`
	Pod          Pod             `json:"pod,omitempty"`
	Container    Container       `json:"container,omitempty"`
}

type ExposableComponent struct {
	Component `json:",inline,omitempty"`
	Ingress   Ingress `json:"ingress,omitempty"`
}

// HorusecPlatformStatus defines the observed state of HorusecPlatform
type HorusecPlatformStatus struct {
	Conditions []metav1.Condition `json:"conditions"`
	State      state.Type         `json:"state"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=horus
// nolint:lll // kubebuilder tags could not be split into multiple lines.
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.state",description="The status of the platform installation"
//+kubebuilder:printcolumn:name="Reason",type="string",priority=1,JSONPath=".status.conditions[?(@.status=='False')].reason",description="It indicates the reason for a condition error"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

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
