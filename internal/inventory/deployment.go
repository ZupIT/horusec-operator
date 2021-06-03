package inventory

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
)

//nolint:funlen, gocritic // to improve in the future
func ForDeployments(existing, desired []appsv1.Deployment) Object {
	update := []client.Object{}
	mcreate := deploymentMap(desired)
	mdelete := deploymentMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

			// if we have a nil value for the replicas in the desired deployment
			// but we have a specific value in the current deployment, we override the desired with the current
			// as this might have been written by an HPA
			if tp.Spec.Replicas != nil && v.Spec.Replicas == nil {
				v.Spec.Replicas = tp.Spec.Replicas
			}

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
		Create: deploymentList(mcreate),
		Update: update,
		Delete: deploymentList(mdelete),
	}
}

//nolint:gocritic // to improve in the future
func deploymentMap(deps []appsv1.Deployment) map[string]appsv1.Deployment {
	m := map[string]appsv1.Deployment{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

//nolint:gocritic, gosec, exportloopref // to improve in the future
func deploymentList(m map[string]appsv1.Deployment) []client.Object {
	l := []client.Object{}
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}
