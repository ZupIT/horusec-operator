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

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/resources/analytic"
	"github.com/ZupIT/horusec-operator/internal/resources/api"
	"github.com/ZupIT/horusec-operator/internal/resources/auth"
	"github.com/ZupIT/horusec-operator/internal/resources/core"
	"github.com/ZupIT/horusec-operator/internal/resources/manager"
	"github.com/ZupIT/horusec-operator/internal/resources/messages"
	"github.com/ZupIT/horusec-operator/internal/resources/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/resources/webhook"
)

func (b *Builder) ServiceAccountsFor(resource *v2alpha1.HorusecPlatform) ([]corev1.ServiceAccount, error) {
	desired := b.listServiceAccounts(resource)
	for index := range desired {
		desiredItem := &desired[index]
		if err := controllerutil.SetControllerReference(resource, desiredItem, b.scheme); err != nil {
			return nil, fmt.Errorf("failed to set service account %q owner reference: %v", desiredItem.GetName(), err)
		}
	}
	return desired, nil
}

func (b *Builder) listServiceAccounts(resource *v2alpha1.HorusecPlatform) []corev1.ServiceAccount {
	serviceAccounts := []corev1.ServiceAccount{
		analytic.NewServiceAccount(resource),
		api.NewServiceAccount(resource),
		auth.NewServiceAccount(resource),
		core.NewServiceAccount(resource),
		manager.NewServiceAccount(resource),
		vulnerability.NewServiceAccount(resource),
		webhook.NewServiceAccount(resource),
	}
	msg := resource.GetMessagesComponent()
	if msg.Enabled {
		serviceAccounts = append(serviceAccounts, messages.NewServiceAccount(resource))
	}
	return serviceAccounts
}
