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
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/internal/k8s"
)

//nolint:funlen, gocritic // to improve in the future
func ForDeployments(existing, desired []appsv1.Deployment) k8s.Objects {
	var update []client.Object
	mcreate := deploymentMap(desired)
	mdelete := deploymentMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(t, v, ignore(ignoredDeploymentFields...))
			if diff != "" {
				tp := t.DeepCopy()

				// if we have a nil value for the replicas in the desired deployment
				// but we have a specific value in the current deployment, we override the desired with the current
				// as this might have been written by an HPA
				if tp.Spec.Replicas != nil && v.Spec.Replicas == nil {
					v.Spec.Replicas = tp.Spec.Replicas
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
		Create: deploymentList(mcreate),
		Update: update,
		Delete: deploymentList(mdelete),
	}
}

//nolint:gocritic // to improve in the future
func deploymentMap(deps []appsv1.Deployment) map[string]appsv1.Deployment {
	m := map[string]appsv1.Deployment{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint // to improve in the future
func deploymentList(m map[string]appsv1.Deployment) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}

var ignoredDeploymentFields = []string{
	"ObjectMeta",
	"Spec.ProgressDeadlineSeconds",
	"Spec.RevisionHistoryLimit",
	"Spec.Strategy",
	"Spec.Template.Spec.Containers.ImagePullPolicy",
	"Spec.Template.Spec.Containers.LivenessProbe.Handler.HTTPGet.Scheme",
	"Spec.Template.Spec.Containers.Ports.Protocol",
	"Spec.Template.Spec.Containers.ReadinessProbe.Handler.HTTPGet.Scheme",
	"Spec.Template.Spec.Containers.TerminationMessagePath",
	"Spec.Template.Spec.Containers.TerminationMessagePolicy",
	"Spec.Template.Spec.DNSPolicy",
	"Spec.Template.Spec.RestartPolicy",
	"Spec.Template.Spec.SchedulerName",
	"Spec.Template.Spec.SecurityContext",
	"Spec.Template.Spec.TerminationGracePeriodSeconds",
	"Status",
	"TypeMeta",
}
