package configuration

import (
	"context"
	v1 "github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/pkg/notifications/event"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	stackerr "github.com/pkg/errors"
)

type Configuration struct {
	Client                  client.Client
	ClientSet               kubernetes.Clientset
	Ship                    *v1.Ship
	Scheme                  *runtime.Scheme
	Config                  *rest.Config
	KubernetesClusterDomain string
	Notifications           *chan event.Event
}

func (c *Configuration) CreateResource(obj metav1.Object) error {
	clientObj, ok := obj.(client.Object)
	if !ok {
		return stackerr.Errorf("is not a %T a runtime.Object", obj)
	}

	if err := controllerutil.SetControllerReference(c.Ship, obj, c.Scheme); err != nil {
		return stackerr.WithStack(err)
	}

	return c.Client.Create(context.TODO(), clientObj)
}

func (c *Configuration) UpdateResource(obj metav1.Object) error {
	clientObj, ok := obj.(client.Object)
	if !ok {
		return stackerr.Errorf("is not a %T a runtime.Object", obj)
	}

	_ = controllerutil.SetControllerReference(c.Ship, obj, c.Scheme)

	return c.Client.Update(context.TODO(), clientObj)
}
