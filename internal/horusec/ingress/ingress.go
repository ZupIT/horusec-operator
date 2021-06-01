package ingress

import (
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/horusec/analytic"
	"github.com/ZupIT/horusec-operator/internal/horusec/api"
	"github.com/ZupIT/horusec-operator/internal/horusec/auth"
	"github.com/ZupIT/horusec-operator/internal/horusec/core"
	"github.com/ZupIT/horusec-operator/internal/horusec/manager"
	"github.com/ZupIT/horusec-operator/internal/horusec/messages"
	"github.com/ZupIT/horusec-operator/internal/horusec/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/horusec/webhook"
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
			Rules: []v1beta1.IngressRule{
				analytic.NewIngressRule(resource, pathType),
				api.NewIngressRule(resource, pathType),
				auth.NewIngressRule(resource, pathType),
				core.NewIngressRule(resource, pathType),
				manager.NewIngressRule(resource, pathType),
				messages.NewIngressRule(resource, pathType),
				vulnerability.NewIngressRule(resource, pathType),
				webhook.NewIngressRule(resource, pathType),
			},
			TLS: []v1beta1.IngressTLS{
				analytic.NewIngressTLS(resource),
				api.NewIngressTLS(resource),
				auth.NewIngressTLS(resource),
				core.NewIngressTLS(resource),
				manager.NewIngressTLS(resource),
				messages.NewIngressTLS(resource),
				vulnerability.NewIngressTLS(resource),
				webhook.NewIngressTLS(resource),
			},
		},
	}
}
