package v1

import corev1 "k8s.io/api/core/v1"

// Container defines Kubernetes container attributes
type Container struct {
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy"`

	// +optional
	Resources corev1.ResourceRequirements `json:"resources"`

	// +optional
	Command []string `json:"command,omitempty"`

	// +optional
	Args []string `json:"args,omitempty"`

	// +optional
	Ports []corev1.ContainerPort `json:"ports,omitempty"`

	// +optional
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`

	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// +optional
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`

	// +optional
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`
}
