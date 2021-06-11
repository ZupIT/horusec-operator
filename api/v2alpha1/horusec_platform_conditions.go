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
	"github.com/ZupIT/horusec-operator/api/v2alpha1/state"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *HorusecPlatform) SetState(state state.Type) bool {
	if in.Status.State != state {
		in.Status.State = state
		return true
	}
	return false
}

func (in *HorusecPlatform) IsStatusConditionFalse(conditionType condition.Type) bool {
	return meta.IsStatusConditionFalse(in.Status.Conditions, string(conditionType))
}

func (in *HorusecPlatform) IsStatusConditionTrue(conditionType condition.Type) bool {
	return meta.IsStatusConditionTrue(in.Status.Conditions, string(conditionType))
}

func (in *HorusecPlatform) SetStatusConditionFalse(conditionType condition.Type) bool {
	if !in.IsStatusConditionFalse(conditionType) {
		switch conditionType {
		case condition.DeploymentsAvailable:
			meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
				Type:    string(conditionType),
				Status:  metav1.ConditionFalse,
				Reason:  "UnavailableReplicas",
				Message: "Deployment has unavailable replicas.",
			})
			return true
		}
	}
	return false
}

func (in *HorusecPlatform) SetStatusConditionTrue(conditionType condition.Type) bool {
	if !in.IsStatusConditionTrue(conditionType) {
		switch conditionType {
		case condition.DeploymentsAvailable:
			meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
				Type:    string(conditionType),
				Status:  metav1.ConditionTrue,
				Reason:  "AvailableReplicas",
				Message: "Deployment has minimum availability.",
			})
			return true
		}
	}
	return false
}
