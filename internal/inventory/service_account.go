package inventory

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:funlen, gocritic // to improve in the future
func ForServiceAccount(existing, desired []corev1.ServiceAccount) Object {
	update := []client.Object{}
	mcreate := serviceAccountMap(desired)
	mdelete := serviceAccountMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

			// we can't blindly DeepCopyInto, so, we select what we bring from the new to the old object
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
func serviceAccountMap(deps []corev1.ServiceAccount) map[string]corev1.ServiceAccount {
	m := map[string]corev1.ServiceAccount{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint // to improve in the future
func serviceAccountList(m map[string]corev1.ServiceAccount) []client.Object {
	l := []client.Object{}
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}
