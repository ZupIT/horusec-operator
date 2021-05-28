package api

import (
	"github.com/ZupIT/horusec-operator/internal/horusec"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

func NewServiceAccount(resource *v2alpha1.HorusecPlatform) *core.ServiceAccount {
	return &core.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    horusec.Labels,
		},
	}
}
