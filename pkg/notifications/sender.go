package notifications

import (
	v1 "github.com/xuxant/voyager-operator/api/v1"
	k8sevent "github.com/xuxant/voyager-operator/pkg/event"
)

//func Listen(events chan event.Event, k8sEvent k8sevent.Recorder, k8sClient k8sclient.Client) {
//	httpClient := http.Client{}
//	for e := range events {
//		logger := log.Log.WithValues("cr", e.Ship.Name)
//		if !e.Reason.HasMessages() {
//			logger.V(log.VWarn).Info("Reason has no message, this should not happen")
//			continue
//		}
//		k8sEvent.Emit(&e.Ship,
//			eventLevelToKubernetesEventType(e.Level),
//			k8sevent.Reason(reflect.TypeOf(e.Reason).Name()),
//			strings.Join(e.Reason.Short(), "; "),
//		)
//		for _, notificationConfig := range
//	}
//}

func eventLevelToKubernetesEventType(level v1.NotificationLevel) k8sevent.Type {
	switch level {
	case v1.NotificationLevelWarning:
		return k8sevent.TypeWarning
	case v1.NotificationLevelInfo:
		return k8sevent.TypeNormal
	default:
		return k8sevent.TypeNormal
	}
}
