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
	probe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/auth/health",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
			},
		},
	}
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
					Env: []corev1.EnvVar{
						{Name: "HORUSEC_PORT", Value: strconv.Itoa(resource.GetAuthPortHTTP())},
						{Name: "HORUSEC_GRPC_PORT", Value: strconv.Itoa(resource.GetAuthPortGRPC())},
						{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: resource.GetGlobalDatabaseLogMode()},
						{Name: "HORUSEC_DISABLED_EMAILS", Value: resource.IsEmailsEnabled()},
						{Name: "HORUSEC_GRPC_USE_CERTS", Value: strconv.FormatBool(resource.Spec.Global.GrpcUseCerts)},
						{Name: "HORUSEC_BROKER_HOST", Value: resource.GetGlobalBrokerHost()},
						{Name: "HORUSEC_BROKER_PORT", Value: resource.GetGlobalBrokerPort()},
						{Name: "HORUSEC_AUTH_TYPE", Value: string(resource.Spec.Components.Auth.Type)},
						{Name: "HORUSEC_ENABLE_APPLICATION_ADMIN", Value: strconv.FormatBool(resource.Spec.Components.Auth.EnableApplicationAdmin)},
						{Name: "HORUSEC_ENABLE_DEFAULT_USER", Value: strconv.FormatBool(resource.Spec.Components.Auth.EnableDefaultUser)},
						{Name: "HORUSEC_DEFAULT_USER_DATA", Value: resource.Spec.Components.Auth.DefaultUserData},
						{Name: "HORUSEC_MANAGER_URL", Value: resource.Spec.Components.Auth.ManagerUrl},
						{Name: "HORUSEC_AUTH_URL", Value: resource.Spec.Components.Auth.AuthUrl},
						resource.NewEnvFromSecret("HORUSEC_BROKER_USERNAME", resource.GetGlobalBrokerUsername()),
						resource.NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", resource.GetGlobalBrokerPassword()),
						resource.NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", resource.GetGlobalDatabaseUsername()),
						resource.NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", resource.GetGlobalDatabasePassword()),
						{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
						{Name: "HORUSEC_KEYCLOAK_BASE_PATH", Value: resource.Spec.Global.Keycloak.PublicURL},
						{Name: "HORUSEC_KEYCLOAK_CLIENT_ID", Value: resource.Spec.Global.Keycloak.Clients.Public.ID},
						resource.NewEnvFromSecret("HORUSEC_KEYCLOAK_CLIENT_SECRET", &resource.Spec.Global.Keycloak.Clients.Confidential.SecretKeyRef),
						{Name: "HORUSEC_KEYCLOAK_REALM", Value: resource.Spec.Global.Keycloak.Realm},
						{Name: "HORUSEC_LDAP_BASE", Value: resource.Spec.Components.Auth.Ldap.Base},
						{Name: "HORUSEC_LDAP_HOST", Value: resource.Spec.Components.Auth.Ldap.Host},
						{Name: "HORUSEC_LDAP_PORT", Value: strconv.Itoa(resource.Spec.Components.Auth.Ldap.Port)},
						{Name: "HORUSEC_LDAP_USESSL", Value: strconv.FormatBool(resource.Spec.Components.Auth.Ldap.UseSSL)},
						{Name: "HORUSEC_LDAP_SKIP_TLS", Value: strconv.FormatBool(resource.Spec.Components.Auth.Ldap.SkipTLS)},
						{Name: "HORUSEC_LDAP_INSECURE_SKIP_VERIFY", Value: resource.GetGlobalDatabaseURI()},
						{Name: "HORUSEC_LDAP_BINDDN", Value: resource.GetGlobalDatabaseURI()},
						{Name: "HORUSEC_LDAP_BINDPASSWORD", Value: resource.GetGlobalDatabaseURI()},
						{Name: "HORUSEC_LDAP_USERFILTER", Value: resource.GetGlobalDatabaseURI()},
						{Name: "HORUSEC_LDAP_ADMIN_GROUP", Value: resource.GetGlobalDatabaseURI()},
					},
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: int32(resource.GetAuthPortHTTP())},
						{Name: "grpc", ContainerPort: int32(resource.GetAuthPortGRPC())},
					},
					LivenessProbe:  &probe,
					ReadinessProbe: &probe,
				}}},
			},
		},
	}
}
