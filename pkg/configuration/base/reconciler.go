package base

import (
	"context"
	"github.com/go-logr/logr"
	stackerr "github.com/pkg/errors"
	v1 "github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/pkg/configuration"
	"github.com/xuxant/voyager-operator/pkg/configuration/base/resources"
	"github.com/xuxant/voyager-operator/pkg/log"
	"github.com/xuxant/voyager-operator/version"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ShipBaseConfigurationReconciler struct {
	configuration.Configuration
	logger logr.Logger
}

func New(config configuration.Configuration) *ShipBaseConfigurationReconciler {
	return &ShipBaseConfigurationReconciler{
		Configuration: config,
		logger:        log.Log.WithValues("cr", config.Ship.Name),
	}
}

func (r *ShipBaseConfigurationReconciler) Reconcile() (reconcile.Result, error) {
	metaObject := resources.NewResourceObjectMeta(r.Configuration.Ship)

	err := r.ensureResourcesRequiredForShip(metaObject)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ShipBaseConfigurationReconciler) ensureResourcesRequiredForShip(metaObject metav1.ObjectMeta) error {
	now := metav1.Now()
	r.Configuration.Ship.Status = v1.ShipStatus{
		OperatorVersion:    version.Version,
		ProvisionStartTime: &now,
	}
	if err := r.createOperatorCredentialsSecret(metaObject); err != nil {
		return err
	}
	if err := r.createServiceForAmes(metaObject, int32(8080)); err != nil {
		return err
	}
	return nil
}

func (r *ShipBaseConfigurationReconciler) createOperatorCredentialsSecret(meta metav1.ObjectMeta) error {
	found := &corev1.Secret{}

	err := r.Configuration.Client.Get(context.TODO(), types.NamespacedName{Name: r.Configuration.Ship.Spec.Keys.Name, Namespace: r.Configuration.Ship.ObjectMeta.Namespace}, found)

	if err != nil && apierrors.IsNotFound(err) {
		return stackerr.WithStack(r.CreateResource(resources.NewOperatorCredentialsSecret(meta, r.Configuration.Ship)))
	} else if err != nil && !apierrors.IsNotFound(err) {
		return stackerr.WithStack(err)
	}
	if found.Data[resources.UserKeyName] != nil && found.Data[resources.PassKeyName] != nil {
		return nil
	}
	return stackerr.WithStack(r.UpdateResource(resources.NewOperatorCredentialsSecret(meta, r.Configuration.Ship)))
}
