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

package resources

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/resources/ingress"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) IngressFor(resource *v2alpha1.HorusecPlatform) ([]networkingv1beta1.Ingress, error) {
	var desiredList []networkingv1beta1.Ingress
	if !resource.GetAllIngressIsDisabled() {
		desired := ingress.NewIngress(resource)
		if err := controllerutil.SetControllerReference(resource, &desired, b.scheme); err != nil {
			return nil, fmt.Errorf("failed to set ingress %q owner reference: %v", desired.GetName(), err)
		}
		desiredList = append(desiredList, desired)
	}
	return desiredList, nil
}
