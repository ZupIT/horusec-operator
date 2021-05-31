package analytic

import (
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

func NewIngress(resource *v2alpha1.HorusecPlatform) *v1.Ingress {
	return &v1.Ingress{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.IngressSpec{},
		Status:     v1.IngressStatus{},
	}
}
