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

package api

import (
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

//nolint:funlen // improve in the future
func NewService(resource *v2alpha1.HorusecPlatform) coreV1.Service {
	return coreV1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetAPIName(),
			Namespace: resource.GetNamespace(),
			Labels:    resource.GetApiLabels(),
		},
		Spec: coreV1.ServiceSpec{
			Ports: []coreV1.ServicePort{
				{
					Name:       "http",
					Port:       int32(resource.GetAPIPortHTTP()),
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(resource.GetAPIPortHTTP()),
				},
			},
			Selector: resource.GetApiLabels(),
			Type:     "ClusterIP",
		},
	}
}
