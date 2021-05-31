package condition

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ZupIT/horusec-operator/api/v2alpha1"
)

type Manager struct{}

func (m *Manager) SetCondition(conditions []v2alpha1.Condition, conditionType v2alpha1.ConditionType, status corev1.ConditionStatus, reason, message string) {
	now := metav1.Now()
	condition, _ := m.FindCondition(conditions, conditionType)
	if message != condition.Message ||
		status != condition.Status ||
		reason != condition.Reason ||
		conditionType != condition.Type {

		condition.LastTransitionTime = now
	}
	if message != "" {
		condition.Message = message
	}
	condition.LastProbeTime = now
	condition.Reason = reason
	condition.Status = status
}

func (m *Manager) FindCondition(conditions []v2alpha1.Condition, conditionType v2alpha1.ConditionType) (*v2alpha1.Condition, bool) {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return &conditions[i], true
		}
	}
	conditions = append(conditions, v2alpha1.Condition{Type: conditionType})
	return &conditions[len(conditions)-1], false
}

func (m *Manager) HasCondition(conditions []v2alpha1.Condition, conditionType v2alpha1.ConditionType) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}
