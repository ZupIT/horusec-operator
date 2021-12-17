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
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/horusec-operator/internal/k8s"
)

//nolint:gocritic, funlen // improve in the future
func ForJobs(existing, desired []batchv1.Job) k8s.Objects {
	var update []client.Object
	mdelete := jobMap(existing)
	mcreate := jobMap(desired)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(t.Spec.Template, v.Spec.Template, ignore(ignoredJobFields...))
			if diff == "" {
				delete(mcreate, k)
				delete(mdelete, k)
			}
		}
	}

	return &Object{
		Create: jobList(mcreate),
		Update: update,
		Delete: jobList(mdelete),
	}
}

//nolint:gocritic // improve in the future
func jobMap(jobs []batchv1.Job) map[string]batchv1.Job {
	m := map[string]batchv1.Job{}
	for _, j := range jobs {
		labels := j.GetLabels()
		component := labels["app.kubernetes.io/component"]
		name := labels["app.kubernetes.io/name"]
		m[fmt.Sprintf("%s.%s.%s", j.Namespace, component, name)] = j
	}
	return m
}

//nolint:gosec, exportloopref, gocritic // improve in the future
func jobList(m map[string]batchv1.Job) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}

var ignoredJobFields = []string{
	"ObjectMeta.Labels",
	"Spec.Containers.TerminationMessagePath",
	"Spec.Containers.TerminationMessagePolicy",
	"Spec.DNSPolicy",
	"Spec.SchedulerName",
	"Spec.SecurityContext",
}
