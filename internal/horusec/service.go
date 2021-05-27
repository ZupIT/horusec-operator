package horusec

import (
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	client k8s.Client
	log    logr.Logger
}

func NewService(client k8s.Client) *Service {
	return &Service{
		client: client,
		log:    ctrl.Log.WithName("services").WithName("Horusec"),
	}
}
