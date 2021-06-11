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

package analytic

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

func NewV1ToV2Job(resource *v2alpha1.HorusecPlatform) batchv1.Job {
	var terminationPeriod int64 = 30
	component := resource.Spec.Components.Analytic
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-analytic-v1-2-v2-", resource.GetName()),
			Namespace:    resource.GetNamespace(),
			Labels:       resource.GetAnalyticV1ToV2Labels(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:                 corev1.RestartPolicyOnFailure,
					TerminationGracePeriodSeconds: &terminationPeriod,
					Containers: []corev1.Container{
						{
							Name:            "horusec-analytic-v1-2-v2",
							Image:           resource.GetAnalyticImage(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"/horusec-analytic-v1-to-v2-migrate"},
							Env: []corev1.EnvVar{
								resource.NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", component.Database.User.KeyRef),
								resource.NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", component.Database.Password.KeyRef),
								{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetAnalyticDatabaseURI()},
								{Name: "HORUSEC_DATABASE_HORUSEC_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
							},
						},
					},
				},
			},
		},
	}
}
