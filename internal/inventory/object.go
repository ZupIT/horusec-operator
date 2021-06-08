package inventory

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object struct {
	Create []client.Object
	Update []client.Object
	Delete []client.Object
}

func (o *Object) ToBeCreated() []client.Object {
	return o.Create
}

func (o *Object) ToBeUpdated() []client.Object {
	return o.Update
}

func (o *Object) ToBeDeleted() []client.Object {
	return o.Delete
}
