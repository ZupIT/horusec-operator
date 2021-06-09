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

package core

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
				Path: "/core/health",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
			},
		},
	}
	global := resource.Spec.Global
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetCoreName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetCoreLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: resource.GetCoreReplicaCount(),
			Selector: &metav1.LabelSelector{MatchLabels: resource.GetCoreLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: resource.GetCoreLabels()},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:  resource.GetAuthName(),
					Image: resource.GetCoreImage(),
					Env: []corev1.EnvVar{
						{Name: "HORUSEC_PORT", Value: strconv.Itoa(resource.GetCorePortHTTP())},
						{Name: "HORUSEC_DATABASE_SQL_LOG_MODE", Value: resource.GetGlobalDatabaseLogMode()},
						{Name: "HORUSEC_GRPC_USE_CERTS", Value: strconv.FormatBool(global.GrpcUseCerts)},
						{Name: "HORUSEC_GRPC_AUTH_URL", Value: resource.GetAuthDefaultGRPCURL()},
						{Name: "HORUSEC_BROKER_HOST", Value: resource.GetGlobalBrokerHost()},
						{Name: "HORUSEC_BROKER_PORT", Value: resource.GetGlobalBrokerPort()},
						resource.NewEnvFromSecret("HORUSEC_BROKER_USERNAME", global.Broker.User.KeyRef),
						resource.NewEnvFromSecret("HORUSEC_BROKER_PASSWORD", global.Broker.Password.KeyRef),
						resource.NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", global.Database.User.KeyRef),
						resource.NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", global.Database.Password.KeyRef),
						{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
					},
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: int32(resource.GetCorePortHTTP())},
					},
					LivenessProbe:  &probe,
					ReadinessProbe: &probe,
				}}},
			},
		},
	}
}
