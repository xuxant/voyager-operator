package controllers

import (
	"fmt"
	"github.com/xuxant/voyager-operator/pkg/constants"
	"github.com/xuxant/voyager-operator/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type enqueueRequestForShip struct{}

func (e *enqueueRequestForShip) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if req := e.getOwnerReconcileRequests(evt.Object); req != nil {
		q.Add(*req)
	}
}

func (e *enqueueRequestForShip) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	req1 := e.getOwnerReconcileRequests(evt.ObjectOld)
	req2 := e.getOwnerReconcileRequests(evt.ObjectNew)

	if req1 != nil || req2 != nil {
		shipName := "unknown"
		if req1 != nil {
			shipName = req1.Name
		}
		if req2 != nil {
			shipName = req2.Name
		}

		log.Log.WithValues("cr", shipName).Info(
			fmt.Sprintf("%T/%s has been updated", evt.ObjectNew, evt.ObjectNew.GetName()),
		)
	}
	if req1 != nil {
		q.Add(*req1)
	}
	if req2 != nil {
		q.Add(*req2)
	}
}

func (e *enqueueRequestForShip) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if req := e.getOwnerReconcileRequests(evt.Object); req != nil {
		q.Add(*req)
	}
}

func (e *enqueueRequestForShip) getOwnerReconcileRequests(object metav1.Object) *reconcile.Request {
	if object.GetLabels()[constants.LabelAppKey] == constants.LabelAppValue &&
		object.GetLabels()[constants.LabelWatchKey] == constants.LabelWatchValue &&
		len(object.GetLabels()[constants.LabelShipCRKey]) > 0 {
		return &reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: object.GetNamespace(),
				Name:      object.GetName(),
			},
		}
	}
	return nil
}
