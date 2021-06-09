package inventory

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/internal/k8s"

	autoScalingV2beta2 "k8s.io/api/autoscaling/v2beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//nolint:funlen,gocritic // to improve in the future
func ForHorizontalPodAutoscaling(existing []autoScalingV2beta2.HorizontalPodAutoscaler,
	desired []autoScalingV2beta2.HorizontalPodAutoscaler) k8s.Objects {
	update := []client.Object{}
	mcreate := hpaMap(desired)
	mdelete := hpaMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()
			if tp.GetLabels() == nil {
				tp.SetLabels(map[string]string{})
			}
			if tp.GetAnnotations() == nil {
				tp.SetAnnotations(map[string]string{})
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

	return &Object{
		Create: hpaList(mcreate),
		Update: update,
		Delete: hpaList(mdelete),
	}
}

// nolint:gocritic
func hpaMap(hpas []autoScalingV2beta2.HorizontalPodAutoscaler) map[string]autoScalingV2beta2.HorizontalPodAutoscaler {
	m := map[string]autoScalingV2beta2.HorizontalPodAutoscaler{}
	for _, d := range hpas {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

// nolint
func hpaList(m map[string]autoScalingV2beta2.HorizontalPodAutoscaler) []client.Object {
	l := []client.Object{}
	for _, v := range m {
		obj := v
		l = append(l, &obj)
	}
	return l
}
