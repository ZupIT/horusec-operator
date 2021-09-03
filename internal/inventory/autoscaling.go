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

	"github.com/google/go-cmp/cmp"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/internal/k8s"
)

//nolint:funlen,gocritic // to improve in the future
func ForHorizontalPodAutoscaling(existing []autoscalingv2beta2.HorizontalPodAutoscaler,
	desired []autoscalingv2beta2.HorizontalPodAutoscaler) k8s.Objects {
	var update []client.Object
	mcreate := hpaMap(desired)
	mdelete := hpaMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(t, v, ignore(ignoredAutoscalingFields...))
			if diff != "" {
				tp := t.DeepCopy()
				if tp.GetLabels() == nil {
					tp.SetLabels(map[string]string{})
				}
				if tp.GetAnnotations() == nil {
					tp.SetAnnotations(map[string]string{})
				}

				// we can't blindly DeepCopyInto, so, we select what we bring from the new to the old object
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
		Create: hpaList(mcreate),
		Update: update,
		Delete: hpaList(mdelete),
	}
}

// nolint:gocritic
func hpaMap(hpas []autoscalingv2beta2.HorizontalPodAutoscaler) map[string]autoscalingv2beta2.HorizontalPodAutoscaler {
	m := map[string]autoscalingv2beta2.HorizontalPodAutoscaler{}
	for _, d := range hpas {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

// nolint
func hpaList(m map[string]autoscalingv2beta2.HorizontalPodAutoscaler) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}

var ignoredAutoscalingFields = []string{
	"TypeMeta",
	"ObjectMeta",
}
