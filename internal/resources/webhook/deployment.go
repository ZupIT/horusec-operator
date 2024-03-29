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

package webhook

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
	global := resource.Spec.Global
	component := resource.Spec.Components.Webhook
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetWebhookName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetWebhookLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: resource.GetWebhookReplicaCount(),
			Selector: &metav1.LabelSelector{MatchLabels: resource.GetWebhookLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: resource.GetWebhookLabels()},
				Spec: corev1.PodSpec{
					ImagePullSecrets: component.Container.Image.PullSecrets,
					Containers: []corev1.Container{{
						Name:            resource.GetWebhookName(),
						Image:           resource.GetWebhookImage(),
						ImagePullPolicy: component.Container.Image.PullPolicy,
						Env: []corev1.EnvVar{
							{Name: "HORUSEC_PORT", Value: strconv.Itoa(resource.GetWebhookPortHTTP())},
							{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: resource.GetGlobalDatabaseLogMode()},
							{Name: "HORUSEC_GRPC_USE_CERTS", Value: strconv.FormatBool(global.GrpcUseCerts)},
							{Name: "HORUSEC_GRPC_AUTH_URL", Value: resource.GetAuthDefaultGRPCURL()},
							{Name: "HORUSEC_BROKER_HOST", Value: resource.GetGlobalBrokerHost()},
							{Name: "HORUSEC_HTTP_TIMEOUT", Value: strconv.Itoa(resource.Spec.Components.Webhook.Timeout)},
							{Name: "HORUSEC_BROKER_PORT", Value: resource.GetGlobalBrokerPort()},
							resource.NewEnvFromSecret("HORUSEC_BROKER_USERNAME", global.Broker.User.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", global.Broker.Password.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_USERNAME", global.Database.User.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_PASSWORD", global.Database.Password.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_JWT_SECRET_KEY", global.JWT.SecretKeyRef),
							{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
						},
						Ports: []corev1.ContainerPort{
							{Name: "http", ContainerPort: int32(resource.GetWebhookPortHTTP())},
						},
						LivenessProbe:  newLivenessProbe(resource),
						ReadinessProbe: newReadinessProbe(resource),
					}},
				},
			},
		},
	}
}

func newLivenessProbe(resource *v2alpha1.HorusecPlatform) *corev1.Probe {
	p := resource.Spec.Components.Webhook.Container.LivenessProbe
	p.Handler = corev1.Handler{HTTPGet: &corev1.HTTPGetAction{
		Path: "/webhook/health",
		Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
	}}
	return &p
}

func newReadinessProbe(resource *v2alpha1.HorusecPlatform) *corev1.Probe {
	p := resource.Spec.Components.Webhook.Container.ReadinessProbe
	p.Handler = corev1.Handler{HTTPGet: &corev1.HTTPGetAction{
		Path: "/webhook/health",
		Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
	}}
	return &p
}
