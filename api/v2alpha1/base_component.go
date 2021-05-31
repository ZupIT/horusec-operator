package v2alpha1

type BaseComponent struct {
	Container    Container `json:"container,omitempty"`
	ExtraEnv     []string  `json:"extraEnv,omitempty"`
	Ingress      Ingress   `json:"ingress,omitempty"`
	Name         string    `json:"name,omitempty"`
	Pod          Pod       `json:"pod,omitempty"`
	Port         Port      `json:"port,omitempty"`
	ReplicaCount *int32    `json:"replicaCount,omitempty"`
}

//nolint:gocritic // [auto created]
func (b BaseComponent) GetReplicaCount() *int32 {
	if !b.Pod.Autoscaling.Enabled {
		return b.ReplicaCount
	}
	return nil
}
