package migration

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewJob(resource *v2alpha1.HorusecPlatform) batchv1.Job {
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-platform-migration", resource.GetName()),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetDefaultLabel(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "migration",
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "horusec-database-migration",
							Image:   resource.GetDatabaseMigrationImage(),
							Command: []string{"migrate.sh"},
							Env: []corev1.EnvVar{
								{
									Name:      "HORUSEC_DATABASE_USERNAME",
									ValueFrom: &corev1.EnvVarSource{SecretKeyRef: resource.GetDatabaseUserSecretKeyRef()},
								},
								{
									Name:      "HORUSEC_DATABASE_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{SecretKeyRef: resource.GetDatabasePasswordSecretKeyRef()},
								},
								{Name: "MIGRATION_NAME", Value: "platform"},
								{Name: "HORUSEC_DATABASE_SQL_URI", Value: "postgresql://$(HORUSEC_DATABASE_USERNAME):$(HORUSEC_DATABASE_PASSWORD)@db.svc.cluster.local:5432/horusec_db?sslmode=disable"},
							},
						},
					},
				},
			},
		},
	}
}
