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

package migration

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

func NewJob(resource *v2alpha1.HorusecPlatform) batchv1.Job {
	var terminationPeriod int64 = 30
	global := resource.Spec.Global
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-platform-migration-", resource.GetName()),
			Namespace:    resource.GetNamespace(),
			Labels:       resource.GetDefaultLabel(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:                 corev1.RestartPolicyOnFailure,
					TerminationGracePeriodSeconds: &terminationPeriod,
					ImagePullSecrets:              global.Database.Migration.Image.PullSecrets,
					Containers: []corev1.Container{
						{
							Name:            "horusec-platform-database-migration",
							Image:           resource.GetDatabaseMigrationImage(),
							ImagePullPolicy: global.Database.Migration.Image.PullPolicy,
							Command:         []string{"migrate.sh"},
							Env: []corev1.EnvVar{
								resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_USERNAME", global.Database.User.KeyRef),
								resource.NewEnvFromSecret("HORUSEC_PLATFORM_DATABASE_PASSWORD", global.Database.Password.KeyRef),
								{Name: "MIGRATION_NAME", Value: "platform"},
								{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetGlobalDatabaseURI()},
							},
						},
					},
				},
			},
		},
	}
}
