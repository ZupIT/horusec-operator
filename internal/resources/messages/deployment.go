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

package messages

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
	component := resource.Spec.Components.Messages
	global := resource.Spec.Global
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetMessagesName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetMessagesLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: resource.GetMessagesReplicaCount(),
			Selector: &metav1.LabelSelector{MatchLabels: resource.GetMessagesLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: resource.GetMessagesLabels()},
				Spec: corev1.PodSpec{
					ImagePullSecrets: component.Container.Image.PullSecrets,
					Containers: []corev1.Container{{
						Name:            resource.GetMessagesName(),
						Image:           resource.GetMessagesImage(),
						ImagePullPolicy: component.Container.Image.PullPolicy,
						Env: []corev1.EnvVar{
							{Name: "HORUSEC_PORT", Value: strconv.Itoa(resource.GetMessagesPortHTTP())},
							{Name: "HORUSEC_GRPC_USE_CERTS", Value: strconv.FormatBool(resource.Spec.Global.GrpcUseCerts)},
							{Name: "HORUSEC_GRPC_AUTH_URL", Value: resource.GetAuthDefaultGRPCURL()},
							{Name: "HORUSEC_BROKER_HOST", Value: resource.GetGlobalBrokerHost()},
							{Name: "HORUSEC_BROKER_PORT", Value: resource.GetGlobalBrokerPort()},
							{Name: "HORUSEC_SMTP_HOST", Value: resource.GetMessagesMailServer().Host},
							{Name: "HORUSEC_SMTP_PORT", Value: strconv.Itoa(resource.GetMessagesMailServer().Port)},
							{Name: "HORUSEC_EMAIL_FROM", Value: resource.Spec.Components.Messages.EmailFrom},
							resource.NewEnvFromSecret("HORUSEC_BROKER_USERNAME", global.Broker.User.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", global.Broker.Password.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_SMTP_USERNAME", component.MailServer.User.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_SMTP_PASSWORD", component.MailServer.Password.KeyRef),
							resource.NewEnvFromSecret("HORUSEC_JWT_SECRET_KEY", global.JWT.SecretKeyRef),
						},
						Ports: []corev1.ContainerPort{
							{Name: "http", ContainerPort: int32(resource.GetMessagesPortHTTP())},
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
	p := resource.Spec.Components.Messages.Container.LivenessProbe
	p.Handler = corev1.Handler{HTTPGet: &corev1.HTTPGetAction{
		Path: "/messages/health",
		Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
	}}
	return &p
}

func newReadinessProbe(resource *v2alpha1.HorusecPlatform) *corev1.Probe {
	p := resource.Spec.Components.Messages.Container.ReadinessProbe
	p.Handler = corev1.Handler{HTTPGet: &corev1.HTTPGetAction{
		Path: "/messages/health",
		Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
	}}
	return &p
}
