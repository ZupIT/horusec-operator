package operation

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
			return h.requeueOnErr(err)
		}
		if result == nil || result.CancelRequest {
			return h.doNotRequeue()
		}
		if result.RequeueRequest {
			return h.requeueAfter(result.RequeueDelay, err)
		}
	}
	return h.doNotRequeue()
}

func (h *Handler) doNotRequeue() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (h *Handler) requeueOnErr(err error) (reconcile.Result, error) {
	return reconcile.Result{}, err
}

func (h *Handler) requeueAfter(duration time.Duration, err error) (reconcile.Result, error) {
	return reconcile.Result{RequeueAfter: duration}, err
}
