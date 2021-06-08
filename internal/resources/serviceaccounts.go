package resources

import (
	"fmt"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/resources/analytic"
	"github.com/ZupIT/horusec-operator/internal/resources/api"
	"github.com/ZupIT/horusec-operator/internal/resources/auth"
	"github.com/ZupIT/horusec-operator/internal/resources/core"
	"github.com/ZupIT/horusec-operator/internal/resources/manager"
	"github.com/ZupIT/horusec-operator/internal/resources/messages"
	"github.com/ZupIT/horusec-operator/internal/resources/vulnerability"
	"github.com/ZupIT/horusec-operator/internal/resources/webhook"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) ServiceAccountsFor(resource *v2alpha1.HorusecPlatform) ([]corev1.ServiceAccount, error) {
	desired := b.listServiceAccounts(resource)
	for index := range desired {
		desiredItem := &desired[index]
		if err := controllerutil.SetControllerReference(resource, desiredItem, b.scheme); err != nil {
			return nil, fmt.Errorf("failed to set service account %q owner reference: %v", desiredItem.GetName(), err)
		}
	}
	return desired, nil
}

func (b *Builder) listServiceAccounts(resource *v2alpha1.HorusecPlatform) []corev1.ServiceAccount {
	serviceAccounts := []corev1.ServiceAccount{
		analytic.NewServiceAccount(resource),
		api.NewServiceAccount(resource),
		auth.NewServiceAccount(resource),
		core.NewServiceAccount(resource),
		manager.NewServiceAccount(resource),
		vulnerability.NewServiceAccount(resource),
		webhook.NewServiceAccount(resource),
	}
	msg := resource.GetMessagesComponent()
	if msg.Enabled {
		serviceAccounts = append(serviceAccounts, messages.NewServiceAccount(resource))
	}
	return serviceAccounts
}
