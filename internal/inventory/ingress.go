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

package inventory

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/internal/k8s"
	"github.com/google/go-cmp/cmp"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:gocritic, funlen // to improve in the future
func ForIngresses(existing, desired []networkingv1beta1.Ingress) k8s.Objects {
	var update []client.Object
	mcreate := ingressMap(desired)
	mdelete := ingressMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(t, v, ignore(ignoredIngressFields...))
			if diff != "" {
				tp := t.DeepCopy()

				tp.Spec = v.Spec
				tp.ObjectMeta.OwnerReferences = v.ObjectMeta.OwnerReferences

				for k, v := range v.ObjectMeta.Annotations {
					tp.ObjectMeta.Annotations[k] = v
				}

				for k, v := range v.ObjectMeta.Labels {
					tp.ObjectMeta.Labels[k] = v
				}

				update = append(update, tp)
			}
			delete(mcreate, k)
			delete(mdelete, k)
		}
	}

	return &Object{
		Create: ingressList(mcreate),
		Update: update,
		Delete: ingressList(mdelete),
	}
}

//nolint:gocritic // to improve in the future
func ingressMap(deps []networkingv1beta1.Ingress) map[string]networkingv1beta1.Ingress {
	m := map[string]networkingv1beta1.Ingress{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint // to improve in the future
func ingressList(m map[string]networkingv1beta1.Ingress) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}

var ignoredIngressFields = []string{
	"TypeMeta",
	"ObjectMeta.UID",
	"ObjectMeta.ResourceVersion",
	"ObjectMeta.Generation",
	"ObjectMeta.CreationTimestamp",
	"ObjectMeta.ManagedFields",
	"ObjectMeta.SelfLink",
	"Spec.Rules.IngressRuleValue.HTTP.Paths.PathType",
	"Status",
}
