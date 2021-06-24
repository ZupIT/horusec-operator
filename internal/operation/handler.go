// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operation

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
	"github.com/ZupIT/horusec-operator/internal/requeue"
)

type Handler struct {
	operations []Func
}

func NewHandler(operations ...Func) *Handler {
	return &Handler{operations: operations}
}

func (h *Handler) Handle(ctx context.Context, resource *v2alpha1.HorusecPlatform) (reconcile.Result, error) {
	for _, op := range h.operations {
		result, err := op(ctx, resource)
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
