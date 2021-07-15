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
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:gocritic, funlen // improve in the future
func ForService(existing, desired []corev1.Service) k8s.Objects {
	var update []client.Object
	mdelete := serviceMap(existing)
	mcreate := serviceMap(desired)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(t, v, ignore(ignoredServiceFields...))
			if diff != "" {
				tp := t.DeepCopy()

				if v.Spec.ClusterIP == "" && len(tp.Spec.ClusterIP) > 0 {
					v.Spec.ClusterIP = tp.Spec.ClusterIP
				}

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
		Create: serviceList(mcreate),
		Update: update,
		Delete: serviceList(mdelete),
	}
}

//nolint:gocritic // improve in the future
func serviceMap(deps []corev1.Service) map[string]corev1.Service {
	m := map[string]corev1.Service{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint // improve in the future
func serviceList(m map[string]corev1.Service) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}

var ignoredServiceFields = []string{
	"TypeMeta",
	"ObjectMeta",
	"Spec.ClusterIP",
	"Spec.ClusterIPs",
	"Spec.SessionAffinity",
	"Status",
}
