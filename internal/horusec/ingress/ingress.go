package ingress

import (
	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//nolint:funlen // improve in the future
func NewIngress(resource *v2alpha1.HorusecPlatform) networkingv1.Ingress {
	return networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetDefaultLabel(),
		},
		Spec: networkingv1.IngressSpec{
			Rules: newIngressRules(resource),
			TLS:   newIngressTLS(resource),
		},
	}
}

func newIngressRules(resource *v2alpha1.HorusecPlatform) []networkingv1.IngressRule {
	hosts := mapHosts(resource)
	tls := make([]networkingv1.IngressRule, 0, len(hosts))
	for host, backends := range hosts {
		tls = append(tls, networkingv1.IngressRule{
			Host: host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{Paths: backends},
			},
		})
	}
	return tls
}

func newIngressTLS(resource *v2alpha1.HorusecPlatform) []networkingv1.IngressTLS {
	secrets := mapTLSSecrets(resource)
	tls := make([]networkingv1.IngressTLS, 0, len(secrets))
	for secret, hosts := range secrets {
		tls = append(tls, networkingv1.IngressTLS{
			Hosts:      hosts,
			SecretName: secret,
		})
	}
	return tls
}

func newHTTPIngressPath(path, service string) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
		Path: path,
		Backend: networkingv1.IngressBackend{
			ServiceName: service,
			ServicePort: intstr.FromString("http"),
		},
	}
}

func mapHosts(r *v2alpha1.HorusecPlatform) map[string][]networkingv1.HTTPIngressPath {
	hosts := make(map[string][]networkingv1.HTTPIngressPath, 0)
	if r.IsAnalyticIngressEnabled() {
		component := r.GetAnalyticComponent()
		path := r.GetAnalyticPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsAPIIngressEnabled() {
		component := r.GetAPIComponent()
		path := r.GetAPIPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsAuthIngressEnabled() {
		component := r.GetAuthComponent()
		path := r.GetAuthPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsCoreIngressEnabled() {
		component := r.GetCoreComponent()
		path := r.GetCorePath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsMessagesIngressEnabled() {
		component := r.GetMessagesComponent()
		path := r.GetMessagesPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsVulnerabilityIngressEnabled() {
		component := r.GetVulnerabilityComponent()
		path := r.GetVulnerabilityPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsWebhookIngressEnabled() {
		component := r.GetWebhookComponent()
		path := r.GetWebhookPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	if r.IsManagerIngressEnabled() {
		component := r.GetManagerComponent()
		path := r.GetManagerPath()
		host := component.Ingress.Host
		if host != "" {
			hosts[host] = append(hosts[host], newHTTPIngressPath(path, component.Name))
		}
	}
	return hosts
}

func mapTLSSecrets(r *v2alpha1.HorusecPlatform) map[string][]string {
	tlsSecrets := make(map[string][]string, 0)
	if r.IsAnalyticIngressEnabled() {
		component := r.GetAnalyticComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetAnalyticHost())
		}
	}
	if r.IsAPIIngressEnabled() {
		component := r.GetAPIComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetAPIHost())
		}
	}
	if r.IsAuthIngressEnabled() {
		component := r.GetAuthComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetAuthHost())
		}
	}
	if r.IsCoreIngressEnabled() {
		component := r.GetCoreComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetCoreHost())
		}
	}
	if r.IsManagerIngressEnabled() {
		component := r.GetManagerComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetManagerHost())
		}
	}
	if r.IsMessagesIngressEnabled() {
		component := r.GetMessagesComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetMessagesHost())
		}
	}
	if r.IsVulnerabilityIngressEnabled() {
		component := r.GetVulnerabilityComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetVulnerabilityHost())
		}
	}
	if r.IsWebhookIngressEnabled() {
		component := r.GetWebhookComponent()
		secretName := component.Ingress.TLS.SecretName
		if secretName != "" {
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.GetWebhookHost())
		}
	}
	return tlsSecrets
}

func dedupe(a []string, b ...string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}
	for letter := range check {
		res = append(res, letter)
	}
	return res
}
