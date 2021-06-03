package v2alpha1

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
)

func (h *HorusecPlatform) GetAnalyticComponent() Analytic {
	return h.Spec.Components.Analytic
}
func (h *HorusecPlatform) GetAnalyticAutoscaling() Autoscaling {
	return h.GetAnalyticComponent().Pod.Autoscaling
}
func (h *HorusecPlatform) GetAnalyticName() string {
	name := h.GetAnalyticComponent().Name
	if name == "" {
		return fmt.Sprintf("%s-analytic", h.GetName())
	}
	return name
}
func (h *HorusecPlatform) GetAnalyticPath() string {
	path := h.GetAnalyticComponent().Ingress.Path
	if path == "" {
		return "/analytic"
	}
	return path
}
func (h *HorusecPlatform) GetAnalyticPortHTTP() int {
	port := h.GetAnalyticComponent().Port.HTTP
	if port == 0 {
		return 8005
	}
	return port
}
func (h *HorusecPlatform) GetAnalyticLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       h.GetName(),
		"app.kubernetes.io/component":  "analytic",
		"app.kubernetes.io/managed-by": "horusec",
	}
}
func (h *HorusecPlatform) GetAnalyticReplicaCount() *int32 {
	if !h.GetAnalyticAutoscaling().Enabled {
		return h.GetAnalyticComponent().ReplicaCount
	}
	return nil
}
func (h *HorusecPlatform) GetAnalyticDefaultURL() string {
	return fmt.Sprintf("http://%s:%v", h.GetAnalyticName(), h.GetAnalyticPortHTTP())
}
func (h *HorusecPlatform) GetAnalyticImage() string {
	image := h.GetAnalyticComponent().Container.Image
	if reflect.ValueOf(image).IsZero() {
		return fmt.Sprintf("docker.io/horuszup/horusec-analytic:%s", h.GetLatestVersion())
	}

	return fmt.Sprintf("%s:%s", image.Registry, image.Tag)
}

func (h *HorusecPlatform) GetAnalyticDatabaseUsername() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.GetAnalyticComponent().Database.User).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-database"},
			Key:                  "user",
			Optional:             nil,
		}
	}
	secret:= h.GetAnalyticComponent().Database.User.SecretKeyRef
	return &secret
}
func (h *HorusecPlatform) GetAnalyticDatabasePassword() *corev1.SecretKeySelector {
	if reflect.ValueOf(h.GetAnalyticComponent().Database.Password).IsZero() {
		return &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "horusec-broker"},
			Key:                  "password",
			Optional:             nil,
		}
	}
	secret:= h.GetAnalyticComponent().Database.Password.SecretKeyRef
	return &secret
}
