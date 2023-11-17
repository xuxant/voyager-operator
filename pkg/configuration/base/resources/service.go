package resources

import (
	v1 "github.com/xuxant/voyager-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
)

func UpdateService(actual corev1.Service, config v1.Service, targetPort int32) corev1.Service {
	actual.ObjectMeta.Annotations = config.Annotations
	for key, value := range config.Labels {
		actual.ObjectMeta.Labels[key] = value
	}
	actual.Spec.Type = config.Type
	if len(actual.Spec.Ports) == 0 {
		actual.Spec.Ports = []corev1.ServicePort{{}}
	}
	actual.Spec.Ports[0].Port = config.Port
	if config.NodePort != 0 {
		actual.Spec.Ports[0].NodePort = config.NodePort
	}
	actual.Spec.Ports[0].TargetPort = intstr.IntOrString{IntVal: targetPort, Type: intstr.Int}
	return actual
}

func IsRunningInCluster() (bool, error) {
	const inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	_, err := os.Stat(inClusterNamespacePath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return false, err
}
