package inventory

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:gocritic, funlen // improve in the future
func ForJobs(existing, desired []batchv1.Job) Object {
	var update []client.Object
	mdelete := jobMap(existing)
	mcreate := jobMap(desired)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()
			//
			//if reflect.ValueOf(v.Spec.Template).IsZero() && !reflect.ValueOf(tp.Spec.Template).IsZero() {
			//	v.Spec.Template = tp.Spec.Template
			//}
			//if v.Spec.Selector == nil && tp.Spec.Selector != nil {
			//	v.Spec.Selector = tp.Spec.Selector
			//}

			//tp.Spec = v.Spec
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
		Create: jobList(mcreate),
		Update: update,
		Delete: jobList(mdelete),
	}
}

//nolint:gocritic // improve in the future
func jobMap(deps []batchv1.Job) map[string]batchv1.Job {
	m := map[string]batchv1.Job{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint:gosec, exportloopref, gocritic // improve in the future
func jobList(m map[string]batchv1.Job) []client.Object {
	var l []client.Object
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}
