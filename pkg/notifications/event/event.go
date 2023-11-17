package event

import (
	v1 "github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/pkg/notifications/reason"
)

type Phase string

type StatusColor string

type LoggingLevel string

type Event struct {
	Ship   v1.Ship
	Phase  Phase
	Reason reason.Reason
	Level  v1.NotificationLevel
}

const (
	PhaseBase Phase = "base"
	PhaseUser Phase = "user"
)
