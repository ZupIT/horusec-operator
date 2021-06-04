package auth

import (
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:funlen // improve in the future
func NewIngressRule(resource *v2alpha1.HorusecPlatform, pathType v1beta1.PathType) v1beta1.IngressRule {
	if resource.Spec.Components.Auth.Ingress.Enabled {
		return v1beta1.IngressRule{}
	}

	return v1beta1.IngressRule{
		Host: resource.GetAuthHost(),
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: []v1beta1.HTTPIngressPath{
					{
						Path:     resource.GetAuthPath(),
						PathType: &pathType,
						Backend: v1beta1.IngressBackend{
							ServiceName: resource.GetAuthName(),
							ServicePort: intstr.FromInt(resource.GetAuthPortHTTP()),
						},
					},
				},
			},
		},
	}
}
