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

package auth

import (
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:lll, funlen // to improve in the future
func NewDeployment(resource *v2alpha1.HorusecPlatform) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetAuthName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetAuthLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: resource.GetAuthReplicaCount(),
			Selector: &metav1.LabelSelector{MatchLabels: resource.GetAuthLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: resource.GetAuthLabels()},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:  resource.GetAuthName(),
					Image: resource.GetAuthImage(),
					Env:   getEnvVars(resource),
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: int32(resource.GetAuthPortHTTP())},
						{Name: "grpc", ContainerPort: int32(resource.GetAuthPortGRPC())},
					},
					LivenessProbe:  newLivenessProbe(resource),
					ReadinessProbe: newReadinessProbe(resource),
				}}},
			},
		},
	}
}

func getEnvVars(resource *v2alpha1.HorusecPlatform) []corev1.EnvVar {
	var envs []corev1.EnvVar

	global := resource.Spec.Global
	defaultEnvs := []corev1.EnvVar{
		resource.NewEnvFromSecret("HORUSEC_BROKER_USERNAME", global.Broker.User.KeyRef),
		resource.NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", global.Broker.Password.KeyRef),
		resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_USERNAME", global.Database.User.KeyRef),
		resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_PASSWORD", global.Database.Password.KeyRef),
	}
	switch resource.Spec.Components.Auth.Type {
	case "keycloak":
		defaultEnvs = append(defaultEnvs, resource.NewEnvFromSecret("HORUSEC_KEYCLOAK_CLIENT_SECRET", global.Keycloak.Clients.Confidential.SecretKeyRef))
	case "ldap":
		defaultEnvs = append(defaultEnvs, resource.NewEnvFromSecret("HORUSEC_LDAP_BINDPASSWORD", global.Ldap.BindPassword.SecretKeyRef))
	case "horusec":
		defaultEnvs = append(defaultEnvs, resource.NewEnvFromSecret("HORUSEC_JWT_SECRET_KEY", global.JWT.SecretKeyRef))
	}
	defaultEnvs = append(defaultEnvs, []corev1.EnvVar{
		{Name: "HORUSEC_PORT", Value: strconv.Itoa(resource.GetAuthPortHTTP())},
		{Name: "HORUSEC_GRPC_PORT", Value: strconv.Itoa(resource.GetAuthPortGRPC())},
		{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: resource.GetGlobalDatabaseLogMode()},
		{Name: "HORUSEC_DISABLED_EMAILS", Value: resource.IsEmailsEnabled()},
		{Name: "HORUSEC_GRPC_USE_CERTS", Value: strconv.FormatBool(global.GrpcUseCerts)},
		{Name: "HORUSEC_BROKER_HOST", Value: resource.GetGlobalBrokerHost()},
		{Name: "HORUSEC_BROKER_PORT", Value: resource.GetGlobalBrokerPort()},
		{Name: "HORUSEC_AUTH_TYPE", Value: string(resource.Spec.Components.Auth.Type)},
		{Name: "HORUSEC_ENABLE_APPLICATION_ADMIN", Value: strconv.FormatBool(resource.Spec.Components.Auth.User.Administrator.Enabled)},
		{Name: "HORUSEC_ENABLE_DEFAULT_USER", Value: strconv.FormatBool(resource.Spec.Components.Auth.User.Default.Enabled)},
		{Name: "HORUSEC_MANAGER_URL", Value: resource.GetManagerDefaultURL()},
		{Name: "HORUSEC_AUTH_URL", Value: resource.GetAuthEndpoint()},
		{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
		{Name: "HORUSEC_KEYCLOAK_BASE_PATH", Value: global.Keycloak.PublicURL},
		{Name: "HORUSEC_KEYCLOAK_CLIENT_ID", Value: global.Keycloak.Clients.Public.ID},
		{Name: "HORUSEC_KEYCLOAK_REALM", Value: global.Keycloak.Realm},
		{Name: "HORUSEC_LDAP_BASE", Value: global.Ldap.Base},
		{Name: "HORUSEC_LDAP_HOST", Value: global.Ldap.Host},
		{Name: "HORUSEC_LDAP_PORT", Value: strconv.Itoa(global.Ldap.Port)},
		{Name: "HORUSEC_LDAP_USESSL", Value: strconv.FormatBool(global.Ldap.UseSSL)},
		{Name: "HORUSEC_LDAP_SKIP_TLS", Value: strconv.FormatBool(global.Ldap.SkipTLS)},
		{Name: "HORUSEC_LDAP_INSECURE_SKIP_VERIFY", Value: strconv.FormatBool(global.Ldap.InsecureSkipVerify)},
		{Name: "HORUSEC_LDAP_BINDDN", Value: global.Ldap.BindDN},
		{Name: "HORUSEC_LDAP_USERFILTER", Value: global.Ldap.UserFilter},
		{Name: "HORUSEC_LDAP_ADMIN_GROUP", Value: global.Ldap.AdminGroup},
		{Name: "HORUSEC_APPLICATION_ADMIN_DATA", Value: resource.GetAuthAdminData()},
		{Name: "HORUSEC_DEFAULT_USER_DATA", Value: resource.GetAuthDefaultUserData()},
	}...)

	for _, envVar := range resource.GetAuthOptionalEnvs() {
		if envVar.Value == "" {
			continue
		}

		envs = append(envs, envVar)
	}

	return append(envs, defaultEnvs...)
}

func newLivenessProbe(resource *v2alpha1.HorusecPlatform) *corev1.Probe {
	p := resource.Spec.Components.Auth.Container.LivenessProbe
	p.Handler = corev1.Handler{HTTPGet: &corev1.HTTPGetAction{
		Path: "/auth/health",
		Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
	}}
	return &p
}

func newReadinessProbe(resource *v2alpha1.HorusecPlatform) *corev1.Probe {
	p := resource.Spec.Components.Auth.Container.ReadinessProbe
	p.Handler = corev1.Handler{HTTPGet: &corev1.HTTPGetAction{
		Path: "/auth/health",
		Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
	}}
	return &p
}
