package horusec

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/inventory"
	"github.com/ZupIT/horusec-operator/internal/tracing"

	appsv1 "k8s.io/api/apps/v1"
	autoScalingV2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct{ client k8s.Client }

func NewService(client k8s.Client) *Service {
	return &Service{client: client}
}

func (s *Service) LookupHorusecPlatform(ctx context.Context, key k8s.ObjectKey) (*v2alpha1.HorusecPlatform, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	r := new(v2alpha1.HorusecPlatform)
	err := s.client.Get(ctx, key, r)
	if err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to lookup resource: %w", err))
	}

	r, err = loadDefaults(r)
	if err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to merge default values: %w", err))
	}

	return r, nil
}

func loadDefaults(b *v2alpha1.HorusecPlatform) (*v2alpha1.HorusecPlatform, error) {
	data, err := ioutil.ReadFile("defaults.json")
	if err != nil {
		return nil, err
	}

	a := new(v2alpha1.HorusecPlatform)
	if err = json.Unmarshal(data, &a.Spec); err != nil {
		return nil, err
	}

	jb, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jb, &a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) UpdateHorusecPlatformStatus(ctx context.Context, resource *v2alpha1.HorusecPlatform) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	err := s.client.Status().Update(ctx, resource)
	if err != nil {
		return span.HandleError(err)
	}
	span.Logger().Info(fmt.Sprintf("%T %q status updated", resource, resource.GetName()))
	return nil
}

//nolint:funlen // to improve in the future
func (s *Service) Apply(ctx context.Context, inv inventory.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	deleteOptions := []k8s.DeleteOption{k8s.PropagationPolicy(metav1.DeletePropagationBackground)}
	for _, obj := range inv.Delete {
		if err := s.client.Delete(ctx, obj, deleteOptions...); err != nil {
			return span.HandleError(fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err))
		}
		log.Info(fmt.Sprintf("%T %q deleted", obj, obj.GetName()))
	}

	for _, obj := range inv.Update {
		if err := s.client.Update(ctx, obj); err != nil {
			return span.HandleError(fmt.Errorf("failed to update %T %q: %w", obj, obj.GetName(), err))
		}
		log.Info(fmt.Sprintf("%T %q updated", obj, obj.GetName()))
	}

	for _, obj := range inv.Create {
		if err := s.client.Create(ctx, obj); err != nil {
			return span.HandleError(fmt.Errorf("failed to create %T %q: %w", obj, obj.GetName(), err))
		}
		log.Info(fmt.Sprintf("%T %q created", obj, obj.GetName()))
	}

	return nil
}

func (s *Service) ListDeployments(ctx context.Context,
	namespace string, matchingLabels map[string]string) (*appsv1.DeploymentList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(matchingLabels),
	}
	list := &appsv1.DeploymentList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s deployments: %w", matchingLabels["app.kubernetes.io/name"], err))
	}
	return list, nil
}

func (s *Service) ListAutoscaling(ctx context.Context,
	namespace string, matchingLabels map[string]string) (*autoScalingV2beta2.HorizontalPodAutoscalerList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(matchingLabels),
	}
	list := &autoScalingV2beta2.HorizontalPodAutoscalerList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s Autoscaling: %w", matchingLabels["app.kubernetes.io/name"], err))
	}
	return list, nil
}

func (s *Service) ListServiceAccounts(
	ctx context.Context, namespace, name string, labels map[string]string) (*corev1.ServiceAccountList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(labels),
	}

	list := &corev1.ServiceAccountList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s service accounts: %w", name, err))
	}

	return list, nil
}

func (s *Service) ListServices(
	ctx context.Context, namespace, name string, labels map[string]string) (*corev1.ServiceList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(labels),
	}
	list := &corev1.ServiceList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s services: %w", name, err))
	}
	return list, nil
}

func (s *Service) ListIngress(
	ctx context.Context, namespace, name string, labels map[string]string) (*networkingv1.IngressList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(labels),
	}

	list := &networkingv1.IngressList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list %s ingress: %w", name, err))
	}

	return list, nil
}

func (s *Service) ListJobs(ctx context.Context, namespace string, labels map[string]string) (*batchv1.JobList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []k8s.ListOption{
		k8s.InNamespace(namespace),
		k8s.MatchingLabels(labels),
	}

	list := &batchv1.JobList{}
	if err := s.client.List(ctx, list, opts...); err != nil {
		return nil, span.HandleError(fmt.Errorf("failed to list jobs: %w", err))
	}

	return list, nil
}
