package usecase

import (
	"context"
	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	"testing"
)

func TestDeploymentsAvailability_EnsureDeploymentsAvailable(t *testing.T) {
	// arrange
	usecase, ctrl := setupToEnsureDeploymentsAvailable(t)

	// act
	result, err := usecase.EnsureDeploymentsAvailable(context.TODO(), &v2alpha1.HorusecPlatform{})

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, operation.ContinueResult(), result)
	ctrl.Finish()
}

func setupToEnsureDeploymentsAvailable(t *testing.T) (*DeploymentsAvailability, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	client := test.NewMockKubernetesClient(ctrl)

	data, err := ioutil.ReadFile("/home/tiagoangelo/wrkspc/github.com/ZupIT/horusec-operator/test/data/deployments_availability_test.yaml")
	assert.NoError(t, err)

	var deploys appsv1.DeploymentList
	err = yaml.Unmarshal(data, &deploys)
	assert.NoError(t, err)

	client.EXPECT().
		ListDeploymentsByOwner(gomock.Any(), gomock.Any()).
		Times(1).
		Return(deploys.Items, nil)

	return NewDeploymentsAvailability(client), ctrl
}
