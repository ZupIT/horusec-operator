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

package controllers

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

type Conditions interface {
	SetCondition(conditions []v2alpha1.Condition, conditionType v2alpha1.ConditionType, status corev1.ConditionStatus, reason, message string)
	FindCondition(conditions []v2alpha1.Condition, conditionType v2alpha1.ConditionType) (*v2alpha1.Condition, bool)
	HasCondition(conditions []v2alpha1.Condition, conditionType v2alpha1.ConditionType) bool
}
