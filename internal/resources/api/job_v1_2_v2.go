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

package api

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

func NewV1ToV2Job(resource *v2alpha1.HorusecPlatform) batchv1.Job {
	component := resource.Spec.Components.API
	var terminationPeriod int64 = 30
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-api-v1-2-v2-", resource.GetName()),
			Namespace:    resource.GetNamespace(),
			Labels:       resource.GetApiV1ToV2Labels(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:                 corev1.RestartPolicyOnFailure,
					TerminationGracePeriodSeconds: &terminationPeriod,
					ImagePullSecrets:              component.Container.Image.PullSecrets,
					Containers: []corev1.Container{
						{
							Name:            "horusec-api-v1-2-v2",
							Image:           resource.GetAPIImage(),
							ImagePullPolicy: component.Container.Image.PullPolicy,
							Command:         []string{"/horusec-api-v1-to-v2-migrate"},
							Env: []corev1.EnvVar{
								resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_USERNAME", resource.Spec.Global.Database.User.KeyRef),
								resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_PASSWORD", resource.Spec.Global.Database.Password.KeyRef),
								{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
							},
						},
					},
				},
			},
		},
	}
}
