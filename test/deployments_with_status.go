package test

import (
	"io/ioutil"
	"path"
	"runtime"

	"gopkg.in/yaml.v3"

	apps "k8s.io/api/apps/v1"
)

func DeploymentsWithStatus() ([]apps.Deployment, error) {
	_, filepath, _, _ := runtime.Caller(0)
	data, err := ioutil.ReadFile(path.Join(path.Dir(filepath), "data/deployments_with_status.yaml"))
	if err != nil {
		return nil, err
	}

	var deploys apps.DeploymentList
	if err = yaml.Unmarshal(data, &deploys); err != nil {
		return nil, err
	}

	return deploys.Items, nil
}
