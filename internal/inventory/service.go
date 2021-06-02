package inventory

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:gocritic, funlen // improve in the future
func ForService(existing, desired []v1.Service) Object {
	var update []client.Object
	mdelete := serviceMap(existing)
	mcreate := serviceMap(desired)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

			if v.Spec.ClusterIP == "" && len(tp.Spec.ClusterIP) > 0 {
				v.Spec.ClusterIP = tp.Spec.ClusterIP
			}

			tp.Spec = v.Spec
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
		Create: serviceList(mcreate),
		Update: update,
		Delete: serviceList(mdelete),
	}
}

//nolint:gocritic // improve in the future
func serviceMap(deps []v1.Service) map[string]v1.Service {
	m := map[string]v1.Service{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint:gosec, exportloopref, gocritic // improve in the future
func serviceList(m map[string]v1.Service) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}
