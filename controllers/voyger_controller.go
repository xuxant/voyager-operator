package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	v1 "github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/pkg/configuration"
	"github.com/xuxant/voyager-operator/pkg/configuration/base"
	"github.com/xuxant/voyager-operator/pkg/log"
	"github.com/xuxant/voyager-operator/pkg/notifications/event"
	"github.com/xuxant/voyager-operator/pkg/notifications/reason"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type reconcileError struct {
	err     error
	counter uint64
}

const (
	APIVersion         = "core/v1"
	SecretKind         = "Secret"
	ConfigMapKind      = "ConfigMap"
	ServiceKind        = "Service"
	containerProbePort = "http"
)

var reconcileErrors = map[string]reconcileError{}
var logx = log.Log

type ShipReconciler struct {
	Client                  client.Client
	Scheme                  *runtime.Scheme
	ClientSet               kubernetes.Clientset
	Config                  rest.Config
	KubernetesClusterDomain string
	NotificationEvents      *chan event.Event
}

func (r *ShipReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Ship{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&networkv1.Ingress{}).
		Watches(
			&v1.Ship{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

func (r *ShipReconciler) newShipReconciler(ship *v1.Ship) configuration.Configuration {
	config := configuration.Configuration{
		Client:                  r.Client,
		ClientSet:               r.ClientSet,
		Ship:                    ship,
		Scheme:                  r.Scheme,
		Config:                  &r.Config,
		KubernetesClusterDomain: r.KubernetesClusterDomain,
		Notifications:           r.NotificationEvents,
	}
	return config
}

// +kubebuilder:rbac:groups=voyager.tlon.io,resources=ships,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=voyager.tlon.io,resources=ships/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=voyager.tlon.io,resources=ships/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services;configmaps;secrets,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=*
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups=core,resources=pods/portforward,verbs=create
// +kubebuilder:rbac:groups=core,resources=pods/log,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods;pods/exec,verbs=*
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;watch;list;create;patch
// +kubebuilder:rbac:groups=voyager.tlon.io,resources=*,verbs=*
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch

func (r *ShipReconciler) Reconcile(_ context.Context, request ctrl.Request) (ctrl.Result, error) {
	reconcileFailedLimit := uint64(10)

	logger := logx.WithValues("cr", request.Name)
	logger.V(log.VDebug).Info("Reconciling Ship")

	result, ship, err := r.reconcile(request)
	if err != nil && apierror.IsConflict(err) {
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		lastErrors, found := reconcileErrors[request.Name]
		if found {
			if err.Error() == lastErrors.err.Error() {
				lastErrors.counter++
			} else {
				lastErrors.counter = 1
				lastErrors.err = err
			}

		} else {
			lastErrors = reconcileError{
				err:     err,
				counter: 1,
			}
		}
		reconcileErrors[request.Name] = lastErrors
		if lastErrors.counter >= reconcileFailedLimit {
			if log.Debug {
				logger.V(log.VWarn).Info(fmt.Sprintf("Reconcile loop failed %d times with same errors, giving up : %+v", reconcileFailedLimit, err))
			} else {
				logger.V(log.VWarn).Info(fmt.Sprintf("Reconcile loop failed %d times with same errors, giving up : %v", reconcileFailedLimit, err))
			}

			*r.NotificationEvents <- event.Event{
				Ship:  *ship,
				Phase: event.PhaseBase,
				Level: v1.NotificationLevelWarning,
				Reason: reason.NewReconcileLoopFailed(reason.OperatorSource,
					[]string{fmt.Sprintf("Reconcile loop failed %d times with the same errors, giving up: %s", reconcileFailedLimit, err)},
				),
			}
			return reconcile.Result{Requeue: false}, nil
		}
		if log.Debug {
			logger.V(log.VWarn).Info(fmt.Sprintf("Reconcile loop failed: %+v", err))
		} else if err.Error() != fmt.Sprintf("Operation cannot be fulfilled on ship.voyager.tlon.io \"%s\": the object has been modified; please apply your changes to the latest version and try again", request.Name) {
			logger.V(log.VWarn).Info(fmt.Sprintf("Reconcile loop failed: %s", err))
		}
		return reconcile.Result{Requeue: true}, nil
	}
	if result.Requeue && result.RequeueAfter == 0 {
		result.RequeueAfter = time.Duration(rand.Intn(10)) * time.Millisecond
	}
	return result, err
}

func (r *ShipReconciler) reconcile(request reconcile.Request) (reconcile.Result, *v1.Ship, error) {
	logger := logx.WithValues("cr", request.Name)

	ship := &v1.Ship{}
	var err error
	err = r.Client.Get(context.TODO(), request.NamespacedName, ship)
	if err != nil {
		if apierror.IsNotFound(err) {
			return reconcile.Result{}, nil, nil
		}
		return reconcile.Result{}, nil, errors.WithStack(err)
	}

	config := r.newShipReconciler(ship)
	baseConfiguration := base.New(config)

	baseMessage, err := baseConfiguration.Validate(ship)
	if err != nil {
		return reconcile.Result{}, ship, err
	}

	if len(baseMessage) > 0 {
		message := "Validation of the base configuration failed, please correct Ship CR."
		//*r.NotificationEvents <- event.Event{
		//	Ship:   *ship,
		//	Phase:  event.PhaseBase,
		//	Level:  v1.NotificationLevelWarning,
		//	Reason: reason.NewBaseConfigurationFailed(reason.HumanSource, []string{message}, append([]string{message}, baseMessage...)...),
		//}
		logger.V(log.VWarn).Info(message)
		for _, msg := range baseMessage {
			logger.V(log.VWarn).Info(msg)
		}
		return reconcile.Result{}, ship, nil
	}

	var result reconcile.Result
	result, err = baseConfiguration.Reconcile()
	if err != nil {
		return reconcile.Result{}, ship, err
	}
	if result.Requeue {
		return result, ship, nil
	}

	if ship.Status.BaseConfigurationCompletedTime == nil {
		now := metav1.Now()
		ship.Status.BaseConfigurationCompletedTime = &now

		err = r.Client.Status().Update(context.TODO(), ship)
		if err != nil {
			return reconcile.Result{}, ship, errors.WithStack(err)
		}

		message := fmt.Sprintf("Base configuration phase complete, took %s",
			ship.Status.BaseConfigurationCompletedTime.Sub(ship.Status.ProvisionStartTime.Time),
		)
		//*r.NotificationEvents <- event.Event{
		//	Ship:   *ship,
		//	Phase:  event.PhaseBase,
		//	Level:  v1.NotificationLevelInfo,
		//	Reason: reason.NewBaseConfigurationFailed(reason.OperatorSource, []string{message}),
		//}
		logger.Info(message)
	}

	return result, ship, nil
}
