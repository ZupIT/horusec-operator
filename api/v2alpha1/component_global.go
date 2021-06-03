package v2alpha1

func (h *HorusecPlatform) GetDefaultLabel() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/managed-by": "horusec",
	}
}

func (h *HorusecPlatform) GetLatestVersion() string {
	return "v2.12.1"
}
