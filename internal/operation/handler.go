package operation

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/ZupIT/horusec-operator/internal/requeue"
)

type Handler struct {
	operations []Func
}

func NewHandler(operations ...Func) *Handler {
	return &Handler{operations: operations}
}

func (h *Handler) Handle(ctx context.Context) (reconcile.Result, error) {
	for _, op := range h.operations {
		result, err := op(ctx)
		if err != nil {
			return requeue.OnErr(err)
		}
		if result == nil || result.CancelRequest {
			return requeue.Not()
		}
		if result.RequeueRequest {
			return requeue.After(result.RequeueDelay, err)
		}
	}
	return requeue.Not()
}
