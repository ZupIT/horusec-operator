package inventory

import (
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForIngresses(existing, desired []v1beta1.Ingress) Object {
	var update []client.Object
	mcreate := ingressMap(desired)
	mdelete := ingressMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

			// we can't blindly DeepCopyInto, so, we select what we bring from the new to the old object
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
		Create: ingressList(mcreate),
		Update: update,
		Delete: ingressList(mdelete),
	}
}

func ingressMap(deps []v1beta1.Ingress) map[string]v1beta1.Ingress {
	m := map[string]v1beta1.Ingress{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

func ingressList(m map[string]v1beta1.Ingress) []client.Object {
	var l []client.Object
	for _, v := range m {
		l = append(l, &v)
	}
	return l
}
