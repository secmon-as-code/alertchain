package types

import "github.com/google/uuid"

type AlertID string
type AlertStatus string

const (
	StatusNew      AlertStatus = "new"
	StatusExecuted AlertStatus = "executed"
	StatusClosed   AlertStatus = "closed"
)

func NewAlertID() AlertID {
	return AlertID(uuid.New().String())
}
