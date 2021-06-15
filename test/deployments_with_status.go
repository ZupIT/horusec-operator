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
