// Package v1alpha2 contains API Schema definitions for the jenkins.io v1alpha2 API group
// +k8s:deepcopy-gen=package,register
// +groupName=jenkins.io
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Kind = "Ship"
)

var (
	GroupVersion       = schema.GroupVersion{Group: "voyager.tlon.io", Version: "v1"}
	SchemeGroupVersion = schema.GroupVersion{Group: "voyager.tlon.io", Version: "v1"}
	SchemeBuilder      = &scheme.Builder{GroupVersion: GroupVersion}
	AddToScheme        = SchemeBuilder.AddToScheme
)

func (in *Ship) GetObjectKind() schema.ObjectKind { return in }

func (in *Ship) SetGroupVersionKind(kind schema.GroupVersionKind) {}

func (in *Ship) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   SchemeGroupVersion.Group,
		Version: SchemeGroupVersion.Version,
		Kind:    Kind,
	}
}

func ShipTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       Kind,
		APIVersion: SchemeGroupVersion.String(),
	}
}

func init() {
	SchemeBuilder.Register(&Ship{}, &ShipList{})
}
