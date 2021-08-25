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
	_ "embed"
	"encoding/json"
)

//go:embed horusec_platform_defaults.json
var defaults []byte

func MergeWithDefaultValues(horus *HorusecPlatform) (*HorusecPlatform, error) {
	merged := new(HorusecPlatform)
	if err := json.Unmarshal(defaults, &merged.Spec); err != nil {
		return nil, err
	}

	jb, err := json.Marshal(horus)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jb, &merged)
	if err != nil {
		return nil, err
	}

	return merged, nil
}
