package usecase

import (
	"context"
	"testing"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/operation"
	"github.com/ZupIT/horusec-operator/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	deps, err := test.DeploymentsWithStatus()
	assert.NoError(t, err)

	client.EXPECT().
		ListDeploymentsByOwner(gomock.Any(), gomock.Any()).
		Times(1).
		Return(deps, nil)

	return NewDeploymentsAvailability(client), ctrl
}
