package ingress

import (
	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/analytic"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewIngress(resource *v2alpha1.HorusecPlatform) *v1beta1.Ingress {
	pathType := v1beta1.PathTypePrefix

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.Labels,
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{},
			Rules: []v1beta1.IngressRule{
				analytic.NewIngress(resource, pathType),
			},
		},
	}
}
