package core

import (
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

func NewIngressRule(resource *v2alpha1.HorusecPlatform, pathType v1beta1.PathType) v1beta1.IngressRule {
	if !resource.Spec.Components.Core.Ingress.Enabled {
		return v1beta1.IngressRule{}
	}

	return v1beta1.IngressRule{
		Host: resource.Spec.Components.Core.Ingress.Host,
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: []v1beta1.HTTPIngressPath{
					{
						Path:     resource.Spec.Components.Core.Ingress.Path,
						PathType: &pathType,
						Backend: v1beta1.IngressBackend{
							ServiceName: resource.Spec.Components.Core.Name,
							ServicePort: intstr.IntOrString{
								Type:   0,
								IntVal: int32(resource.Spec.Components.Core.Port.HTTP),
							},
						},
					},
				},
			},
		},
	}
}

func NewIngressTLS(resource *v2alpha1.HorusecPlatform) v1beta1.IngressTLS {
	if !resource.Spec.Components.Core.Ingress.Enabled {
		return v1beta1.IngressTLS{}
	}

	return v1beta1.IngressTLS{
		Hosts:      []string{resource.Spec.Components.Core.Ingress.Host},
		SecretName: resource.Spec.Components.Core.Ingress.TLS.SecretName,
	}
}
