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

//nolint:funlen // improve in the future
func NewIngress(resource *v2alpha1.HorusecPlatform) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetDefaultLabel(),
		},
		Spec: v1beta1.IngressSpec{
			Rules: NewIngressRules(resource),
			TLS:   NewIngressTLS(resource),
		},
	}
}

func NewIngressTLS(resource *v2alpha1.HorusecPlatform) []v1beta1.IngressTLS {
	var ingressList []v1beta1.IngressTLS
	for key, value := range mapSecretsAndHosts(resource) {
		ingress := v1beta1.IngressTLS{
			Hosts:      value,
			SecretName: key,
		}

		ingressList = append(ingressList, ingress)
	}

	return ingressList
}

func mapSecretsAndHosts(resource *v2alpha1.HorusecPlatform) map[string][]string {
	defaultEnable := true
	ingressConfig := newIngressConfigList(resource)

	tlsMap := map[string][]string{}
	for index := range ingressConfig {
		if ingressConfig[index].TLS.SecretName != "" {
			if ingressConfig[index].Enabled == nil {
				ingressConfig[index].Enabled = &defaultEnable
			}
			if value, ok := tlsMap[ingressConfig[index].TLS.SecretName]; ok && *ingressConfig[index].Enabled {
				tlsMap[ingressConfig[index].TLS.SecretName] = append(value, ingressConfig[index].Host)
			} else {
				tlsMap[ingressConfig[index].TLS.SecretName] = append(tlsMap[ingressConfig[index].TLS.SecretName], ingressConfig[index].Host)
			}
		}
	}

	return tlsMap
}

func newIngressConfigList(resource *v2alpha1.HorusecPlatform) []v2alpha1.Ingress {
	return []v2alpha1.Ingress{
		resource.Spec.Components.Analytic.Ingress,
		resource.Spec.Components.Api.Ingress,
		resource.Spec.Components.Auth.Ingress,
		resource.Spec.Components.Core.Ingress,
		resource.Spec.Components.Manager.Ingress,
		resource.Spec.Components.Messages.Ingress,
		resource.Spec.Components.Vulnerability.Ingress,
		resource.Spec.Components.Webhook.Ingress,
	}
}

func NewIngressRules(resource *v2alpha1.HorusecPlatform) []v1beta1.IngressRule {
	var ingressList []v1beta1.IngressRule
	for key, value := range mapRulesAndHosts(resource) {
		ingress := v1beta1.IngressRule{
			Host: key,
			IngressRuleValue: v1beta1.IngressRuleValue{
				HTTP: &v1beta1.HTTPIngressRuleValue{
					Paths: value,
				},
			},
		}

		ingressList = append(ingressList, ingress)
	}

	return ingressList
}

func mapRulesAndHosts(resource *v2alpha1.HorusecPlatform) map[string][]v1beta1.HTTPIngressPath {
	rulesMap := map[string][]v1beta1.HTTPIngressPath{}

	pathType := v1beta1.PathTypePrefix
	analyticIngress := analytic.NewIngressRule(resource, pathType)

	if ingressPath, ok := rulesMap[analyticIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[analyticIngress.Host] = append(rulesMap[analyticIngress.Host], path)
		}
	} else {
		for _, path := range analyticIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[analyticIngress.Host] = append(rulesMap[analyticIngress.Host], path)
		}
	}

	apiIngress := api.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[apiIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[apiIngress.Host] = append(rulesMap[apiIngress.Host], path)
		}
	} else {
		for _, path := range apiIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[apiIngress.Host] = append(rulesMap[apiIngress.Host], path)
		}
	}

	authIngress := auth.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[authIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[authIngress.Host] = append(rulesMap[authIngress.Host], path)
		}
	} else {
		for _, path := range authIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[authIngress.Host] = append(rulesMap[authIngress.Host], path)
		}
	}
	coreIngress := core.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[coreIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[coreIngress.Host] = append(rulesMap[coreIngress.Host], path)
		}
	} else {
		for _, path := range coreIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[coreIngress.Host] = append(rulesMap[coreIngress.Host], path)
		}
	}
	managerIngress := manager.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[managerIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[managerIngress.Host] = append(rulesMap[managerIngress.Host], path)
		}
	} else {
		for _, path := range managerIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[managerIngress.Host] = append(rulesMap[managerIngress.Host], path)
		}
	}
	messagesIngress := messages.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[messagesIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[messagesIngress.Host] = append(rulesMap[messagesIngress.Host], path)
		}
	} else {
		for _, path := range messagesIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[messagesIngress.Host] = append(rulesMap[messagesIngress.Host], path)
		}
	}
	vulnerabilityIngress := vulnerability.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[vulnerabilityIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[vulnerabilityIngress.Host] = append(rulesMap[vulnerabilityIngress.Host], path)
		}
	} else {
		for _, path := range vulnerabilityIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[vulnerabilityIngress.Host] = append(rulesMap[vulnerabilityIngress.Host], path)
		}
	}
	webhookIngress := webhook.NewIngressRule(resource, pathType)
	if ingressPath, ok := rulesMap[webhookIngress.Host]; ok {
		for _, path := range ingressPath {
			rulesMap[webhookIngress.Host] = append(rulesMap[webhookIngress.Host], path)
		}
	} else {
		for _, path := range webhookIngress.IngressRuleValue.HTTP.Paths {
			rulesMap[webhookIngress.Host] = append(rulesMap[webhookIngress.Host], path)
		}
	}

	return rulesMap
}
