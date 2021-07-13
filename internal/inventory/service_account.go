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

//nolint:funlen, gocritic // to improve in the future
func ForServiceAccount(existing, desired []corev1.ServiceAccount) k8s.Objects {
	var update []client.Object
	mcreate := serviceAccountMap(desired)
	mdelete := serviceAccountMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(t, v, ignore(ignoredServiceAccountFields...))
			if diff != "" {
				tp := t.DeepCopy()

				// we can't blindly DeepCopyInto, so, we select what we bring from the new to the old object
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
		Create: serviceAccountList(mcreate),
		Update: update,
		Delete: serviceAccountList(mdelete),
	}
}

//nolint:gocritic // to improve in the future
func serviceAccountMap(deps []corev1.ServiceAccount) map[string]corev1.ServiceAccount {
	m := map[string]corev1.ServiceAccount{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint // to improve in the future
func serviceAccountList(m map[string]corev1.ServiceAccount) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}

var ignoredServiceAccountFields = []string{
	"TypeMeta",
	"ObjectMeta",
	"Secrets",
}
