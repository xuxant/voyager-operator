package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
type ShipSpec struct {
	Image string `json:"image,omitempty"`

	// +optional
	Domain string `json:"domain,omitempty"`

	// +optional
	Managed bool `json:"managed,omitempty"`

	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// +optional
	Network bool `json:"network,omitempty"`

	// +optional
	Pack bool `json:"pack,omitempty"`

	// +optional
	Keys corev1.SecretReference `json:"keys,omitempty"`

	// +optional
	Dock string `json:"dock,omitempty"`

	// +optional
	Resource corev1.ResourceRequirements `json:"resource,omitempty"`

	// +optional
	SnapTime int `json:"snapTime,omitempty"`

	// +optional
	Loom int `json:"loom,omitempty"`

	// +optional
	Demand bool `json:"demand,omitempty"`

	// +optional
	PioneerNodeSelector map[string]string `json:"pioneerNodeSelector"`

	// +optional
	Tolerations corev1.Toleration `json:"tolerations,omitempty"`

	// +optional
	Affinity corev1.Affinity `json:"affinity,omitempty"`

	// +optional
	Swap bool `json:"swap,omitempty"`
}

// +k8s:openapi-gen=true
type ShipStatus struct {

	//	+optional
	OperatorVersion string `json:"operatorVersion,omitempty"`

	// +optional
	ProvisionStartTime *metav1.Time `json:"provisionStartTime,omitempty"`

	//	+optional
	BaseConfigurationCompletedTime *metav1.Time `json:"BaseConfigurationCompletedTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
type Ship struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShipSpec   `json:"spec,omitempty"`
	Status ShipStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ShipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ship `json:"items"`
}

type NotificationLevel string

const (
	NotificationLevelWarning NotificationLevel = "warning"
	NotificationLevelInfo    NotificationLevel = "info"
)

type Service struct {

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	//	+optional
	Labels map[string]string `json:"labels,omitempty"`

	//+optional
	Type corev1.ServiceType `json:"type,omitempty"`

	Port int32 `json:"port,omitempty"`

	//+optional
	NodePort int32 `json:"nodePort,omitempty"`

	//+optional
	LoadBalancerIP string `json:"LoadBalancerIP,omitempty"`
}
