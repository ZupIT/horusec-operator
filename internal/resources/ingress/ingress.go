// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ingress

import (
	"sort"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//nolint:funlen // improve in the future
func NewIngress(resource *v2alpha1.HorusecPlatform) networkingv1beta1.Ingress {
	return networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetDefaultLabel(),
		},
		Spec: networkingv1beta1.IngressSpec{
			Rules: newIngressRules(resource),
			TLS:   newIngressTLS(resource),
		},
	}
}

func newIngressRules(resource *v2alpha1.HorusecPlatform) []networkingv1beta1.IngressRule {
	hosts := mapHosts(resource)
	rules := make([]networkingv1beta1.IngressRule, 0, len(hosts))
	for host, backends := range hosts {
		rules = append(rules, networkingv1beta1.IngressRule{
			Host: host,
			IngressRuleValue: networkingv1beta1.IngressRuleValue{
				HTTP: &networkingv1beta1.HTTPIngressRuleValue{Paths: backends},
			},
		})
	}
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Host < rules[j].Host
	})
	return rules
}

func newIngressTLS(resource *v2alpha1.HorusecPlatform) []networkingv1beta1.IngressTLS {
	secrets := mapTLSSecrets(resource)
	tls := make([]networkingv1beta1.IngressTLS, 0, len(secrets))
	for secret, hosts := range secrets {
		tls = append(tls, networkingv1beta1.IngressTLS{
			Hosts:      hosts,
			SecretName: secret,
		})
	}
	if len(tls) == 0 {
		return nil
	}
	sort.Slice(tls, func(i, j int) bool {
		return tls[i].SecretName < tls[j].SecretName
	})
	return tls
}

func newHTTPIngressPath(path, service string) networkingv1beta1.HTTPIngressPath {
	prefix := networkingv1beta1.PathTypePrefix
	return networkingv1beta1.HTTPIngressPath{
		Path:     path,
		PathType: &prefix,
		Backend: networkingv1beta1.IngressBackend{
			ServiceName: service,
			ServicePort: intstr.FromString("http"),
		},
	}
}

func mapHosts(r *v2alpha1.HorusecPlatform) map[string][]networkingv1beta1.HTTPIngressPath {
	hosts := make(map[string][]networkingv1beta1.HTTPIngressPath, 0)
	for _, ingress := range r.Ingresses() {
		if ingress.IsEnabled() {
			path := ingress.GetPath()
			if host := ingress.GetHost(); host != "" {
				hosts[host] = append(hosts[host], newHTTPIngressPath(path, ingress.GetName()))
			}
		}
	}
	return hosts
}

func mapTLSSecrets(r *v2alpha1.HorusecPlatform) map[string][]string {
	tlsSecrets := make(map[string][]string, 0)
	for _, ingress := range r.Ingresses() {
		if ingress.IsEnabled() {
			secretName := ingress.GetSecretName()
			if secretName != "" {
				host := ingress.GetHost()
				tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], host)
			}
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
	sort.Strings(res)
	return res
}
