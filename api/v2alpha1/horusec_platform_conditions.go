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

package v2alpha1

import (
	"github.com/ZupIT/horusec-operator/api/v2alpha1/condition"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *HorusecPlatform) IsStatusConditionFalse(types ...condition.Type) bool {
	for _, conditionType := range types {
		if !meta.IsStatusConditionFalse(in.Status.Conditions, string(conditionType)) {
			return false
		}
	}
	return true
}

func (in *HorusecPlatform) IsStatusConditionTrue(types ...condition.Type) bool {
	for _, conditionType := range types {
		if !meta.IsStatusConditionTrue(in.Status.Conditions, string(conditionType)) {
			return false
		}
	}
	return true
}

func (in *HorusecPlatform) AnyStatusConditionFalse(types ...condition.Type) bool {
	for _, conditionType := range types {
		if !meta.IsStatusConditionFalse(in.Status.Conditions, string(conditionType)) { // TODO: if any condition is false, than return true
			return false
		}
	}
	return true
}

func (in *HorusecPlatform) SetStatusCondition(newCondition metav1.Condition) bool {
	conditionType := newCondition.Type
	status := newCondition.Status
	if meta.IsStatusConditionPresentAndEqual(in.Status.Conditions, conditionType, status) {
		return false
	}

	meta.SetStatusCondition(&in.Status.Conditions, newCondition)
	_ = in.UpdateState()
	return true
}
