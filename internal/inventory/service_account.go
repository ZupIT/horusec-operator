package inventory

import (
	"fmt"

	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:funlen, gocritic // to improve in the future
func ForServiceAccount(existing, desired []core.ServiceAccount) Object {
	var update []client.Object
	mcreate := serviceAccountMap(desired)
	mdelete := serviceAccountMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

			tp.ObjectMeta.OwnerReferences = v.ObjectMeta.OwnerReferences

			for k, v := range v.ObjectMeta.Annotations {
				tp.ObjectMeta.Annotations[k] = v
			}

			for k, v := range v.ObjectMeta.Labels {
				tp.ObjectMeta.Labels[k] = v
			}

			update = append(update, tp)
			delete(mcreate, k)
			delete(mdelete, k)
		}
	}

	return Object{
		Create: serviceAccountList(mcreate),
		Update: update,
		Delete: serviceAccountList(mdelete),
	}
}

//nolint:gocritic // to improve in the future
func serviceAccountMap(deps []core.ServiceAccount) map[string]core.ServiceAccount {
	m := map[string]core.ServiceAccount{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint:gocritic, gosec, exportloopref // to improve in the future
func serviceAccountList(m map[string]core.ServiceAccount) []client.Object {
	var l []client.Object
	for _, v := range m {
		l = append(l, &v)
	}
	return l
}
