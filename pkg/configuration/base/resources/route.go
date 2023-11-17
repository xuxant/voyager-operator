package resources

import (
	v1 "github.com/xuxant/voyager-operator/api/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
)

var isRouteAPIAvailable = false
var routeAPIChecked = false

func UpdateIngress(actual networkv1.Ingress, ship *v1.Ship) networkv1.Ingress {
	return networkv1.Ingress{}
}

func IsIngressAPIAvailable(clientSet *kubernetes.Clientset) bool {
	if routeAPIChecked {
		return isRouteAPIAvailable
	}
	gv := schema.GroupVersion{
		Group:   networkv1.GroupName,
		Version: networkv1.SchemeGroupVersion.Version,
	}

	if err := discovery.ServerSupportsVersion(clientSet, gv); err != nil {
		routeAPIChecked = true
		isRouteAPIAvailable = true
	} else {
		routeAPIChecked = true
		isRouteAPIAvailable = true
	}
	return isRouteAPIAvailable
}
