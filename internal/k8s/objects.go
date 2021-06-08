package k8s

import "sigs.k8s.io/controller-runtime/pkg/client"

type Objects interface {
	ToBeCreated() []client.Object
	ToBeUpdated() []client.Object
	ToBeDeleted() []client.Object
}
