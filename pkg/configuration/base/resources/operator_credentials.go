package resources

import (
	v1 "github.com/xuxant/voyager-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	UserKeyName = "username"
	PassKeyName = "password"
)

func buildSecretTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "Secret",
		APIVersion: "v1",
	}
}

func NewOperatorCredentialsSecret(meta metav1.ObjectMeta, ship *v1.Ship) *corev1.Secret {
	meta.Name = ship.Spec.Keys.Name
	return &corev1.Secret{
		TypeMeta:   buildSecretTypeMeta(),
		ObjectMeta: meta,
		Data: map[string][]byte{
			"user":     []byte("username"),
			"password": []byte("password"),
		},
	}
}
