package event

import (
	"fmt"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

const (
	TypeNormal  = Type("Normal")
	TypeWarning = Type("Warning")
)

type Recorder interface {
	Emit(object runtime.Object, eventType Type, reason Reason, message string)
	Emitf(object runtime.Object, eventType Type, reason Reason, format string, args ...interface{})
}

type Type string

type Reason string

type recorder struct {
	recorder record.EventRecorder
}

func New(config *rest.Config, component string) (Recorder, error) {
	eventRecorded, err := initializeEventRecorder(config, component)
	if err != nil {
		return nil, err
	}

	return &recorder{
		recorder: eventRecorded,
	}, nil
}

func initializeEventRecorder(config *rest.Config, component string) (record.EventRecorder, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: client.CoreV1().Events("")})
	eventRecorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		v1.EventSource{Component: component},
	)
	return eventRecorder, nil
}

func (r recorder) Emit(object runtime.Object, eventType Type, reason Reason, message string) {
	r.recorder.Event(object, string(eventType), string(reason), message)
}

func (r recorder) Emitf(object runtime.Object, eventType Type, reason Reason, format string, args ...interface{}) {
	r.recorder.Event(object, string(eventType), string(reason), fmt.Sprintf(format, args...))
}
