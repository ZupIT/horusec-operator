package horusec

import (
	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

type Adapter struct {
	resource *v2alpha1.HorusecPlatform
	svc      *Service
}
