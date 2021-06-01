package requeue

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Not() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func OnErr(err error) (reconcile.Result, error) {
	return reconcile.Result{}, err
}

func After(duration time.Duration, err error) (reconcile.Result, error) {
	return reconcile.Result{RequeueAfter: duration, Requeue: true}, err
}
