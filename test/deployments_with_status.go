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

package test

import (
	"os"
	"path"
	"runtime"

	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func DeploymentsWithStatus() ([]apps.Deployment, error) {
	_, filepath, _, _ := runtime.Caller(0)
	data, err := os.Open(path.Join(path.Dir(filepath), "data/deployments_with_status.yaml"))
	if err != nil {
		return nil, err
	}

	var deploys apps.DeploymentList
	if err = yaml.NewYAMLOrJSONDecoder(data, 100).Decode(&deploys); err != nil {
		return nil, err
	}

	return deploys.Items, nil
}
