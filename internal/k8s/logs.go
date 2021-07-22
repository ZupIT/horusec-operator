package k8s

import (
	"context"

	"github.com/ZupIT/horusec-operator/internal/tracing"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ContainerClient struct {
	typedcorev1.CoreV1Interface
}

func NewContainerClient(coreV1Interface typedcorev1.CoreV1Interface) *ContainerClient {
	return &ContainerClient{CoreV1Interface: coreV1Interface}
}

func (l *ContainerClient) PreviousContainerLogs(ctx context.Context, pod types.NamespacedName, container string) ([]byte, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return l.Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Container: container,
		Follow:    false,
		Previous:  true,
	}).DoRaw(ctx)
}
