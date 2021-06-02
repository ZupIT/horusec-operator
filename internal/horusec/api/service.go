package api

import (
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:funlen // improve in the future
func NewService(resource *v2alpha1.HorusecPlatform) coreV1.Service {
	return coreV1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetAPIName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetApiLabels(),
		},
		Spec: coreV1.ServiceSpec{
			Ports: []coreV1.ServicePort{
				{
					Name:       "http",
					Port:       int32(resource.GetAPIPortHTTP()),
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(resource.GetAPIPortHTTP()),
				},
			},
			Selector: map[string]string{
				"app": "horusec-analytic",
			},
			Type: "ClusterIP",
		},
	}
}
