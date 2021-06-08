package resources

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/resources/ingress"
	networkingv1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) IngressFor(resource *v2alpha1.HorusecPlatform) ([]networkingv1.Ingress, error) {
	var desiredList []networkingv1.Ingress
	if !resource.GetAllIngressIsDisabled() {
		desired := ingress.NewIngress(resource)
		if err := controllerutil.SetControllerReference(resource, &desired, b.scheme); err != nil {
			return nil, fmt.Errorf("failed to set ingress %q owner reference: %v", desired.GetName(), err)
		}
		desiredList = append(desiredList, desired)
	}
	return desiredList, nil
}
