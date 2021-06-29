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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1/condition"
	"github.com/ZupIT/horusec-operator/api/v2alpha1/state"
)

func (in *HorusecPlatform) UpdateState() bool {
	desired := make([]condition.Type, 0, len(condition.ComponentMap))
	for _, c := range condition.ComponentMap {
		desired = append(desired, c)
	}

	if in.IsStatusConditionTrue(desired...) {
		return in.setState(state.Ready)
	}

	if in.AnyStatusConditionFalse(desired...) {
		return in.setState(state.Error)
	}

	return in.setState(state.Pending)
}

func (in *HorusecPlatform) setState(state state.Type) bool {
	if in.Status.Conditions == nil {
		in.Status.Conditions = make([]metav1.Condition, 0)
	}
	if in.Status.State != state {
		in.Status.State = state
		return true
	}
	return false
}
