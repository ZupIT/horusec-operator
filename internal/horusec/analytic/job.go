package analytic

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewJob(resource *v2alpha1.HorusecPlatform) batchv1.Job {
	var terminationPeriod int64 = 30
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-analytic-migration-", resource.GetName()),
			Namespace:    resource.GetNamespace(),
			Labels:       resource.GetAnalyticLabels(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:                 corev1.RestartPolicyNever,
					TerminationGracePeriodSeconds: &terminationPeriod,
					Containers: []corev1.Container{
						{
							Name:            "horusec-database-migration",
							Image:           resource.GetDatabaseMigrationImage(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"migrate.sh"},
							Env: []corev1.EnvVar{
								resource.NewEnvFromSecret("HORUSEC_DATABASE_USERNAME", resource.GetAnalyticDatabaseUsername()),
								resource.NewEnvFromSecret("HORUSEC_DATABASE_PASSWORD", resource.GetAnalyticDatabasePassword()),
								{Name: "MIGRATION_NAME", Value: "analytic"},
								{Name: "HORUSEC_DATABASE_SQL_URI", Value: resource.GetAnalyticDatabaseURI()},
							},
						},
					},
				},
			},
		},
	}
}
