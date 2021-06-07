package inventory

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:gocritic, funlen // to improve in the future
func ForIngresses(existing, desired []networkingv1.Ingress) Object {
	update := []client.Object{}
	mcreate := ingressMap(desired)
	mdelete := ingressMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

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

//nolint:gocritic // to improve in the future
func ingressMap(deps []networkingv1.Ingress) map[string]networkingv1.Ingress {
	m := map[string]networkingv1.Ingress{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint // to improve in the future
func ingressList(m map[string]networkingv1.Ingress) []client.Object {
	l := []client.Object{}
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}
