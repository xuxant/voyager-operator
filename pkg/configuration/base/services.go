package base

import (
	"context"
	"fmt"
	stackerr "github.com/pkg/errors"
	"github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/pkg/configuration/base/resources"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"time"
)

const (
	StartingPort = int32(30000)
	EndingPort   = int32(32000)
	AmisAppPort  = int32(8080)
)

func (r *ShipBaseConfigurationReconciler) createService(meta metav1.ObjectMeta, name string, config v1.Service, targetPort int32) error {
	service := corev1.Service{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: meta.Namespace}, &service)
	if err != nil && apierror.IsNotFound(err) {
		service = resources.UpdateService(corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: meta.Namespace,
				Labels:    meta.Labels,
			},
			Spec: corev1.ServiceSpec{
				Selector: meta.Labels,
			},
		}, config, targetPort)
		if err = r.CreateResource(&service); err != nil {
			return stackerr.WithStack(err)
		}
	} else if err != nil {
		return stackerr.WithStack(err)
	}
	service.Spec.Selector = meta.Labels
	service = resources.UpdateService(service, config, targetPort)
	return stackerr.WithStack(r.UpdateResource(&service))
}

func (r *ShipBaseConfigurationReconciler) createServiceForAmes(meta metav1.ObjectMeta, targetPort int32) error {
	service := corev1.Service{}
	config := v1.Service{}
	serviceName := fmt.Sprintf("%s-%s", r.Configuration.Ship.Name, "amis")
	config.Type = corev1.ServiceTypeNodePort

	port, err := findPortForNodePort(&r.ClientSet)
	if err != nil {
		return err
	}
	config.NodePort = port
	config.Port = AmisAppPort

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: serviceName, Namespace: meta.Namespace}, &service)
	if err != nil && apierror.IsNotFound(err) {
		service = resources.UpdateService(corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: meta.Namespace,
				Labels:    meta.Labels,
			},
			Spec: corev1.ServiceSpec{
				Selector: meta.Labels,
			},
		}, config, targetPort)
		if err = r.CreateResource(&service); err != nil {
			return stackerr.WithStack(err)
		}
	} else if err != nil {
		return stackerr.WithStack(err)
	}
	return nil
}

func isNodePortUsed(clientSet *kubernetes.Clientset, port int32) (bool, error) {
	services, err := clientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, service := range services.Items {
		if service.Spec.Type == corev1.ServiceTypeNodePort {
			for _, servicePort := range service.Spec.Ports {
				if servicePort.NodePort == int32(port) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func findPortForNodePort(clientSet *kubernetes.Clientset) (int32, error) {
	for i := StartingPort; i <= EndingPort; i++ {
		used, err := isNodePortUsed(clientSet, i)
		if err != nil {
			return 0, err
		}
		if !used {
			return i, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return 0, fmt.Errorf("no available NodePort for Ames")
}
