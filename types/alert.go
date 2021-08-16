package types

import "github.com/google/uuid"

type AlertID string

func NewAlertID() AlertID {
	return AlertID(uuid.New().String())
}

type AlertStatus string

const (
	StatusNew      AlertStatus = "new"
	StatusExecuted AlertStatus = "executed"
	StatusClosed   AlertStatus = "closed"
)

type Severity string

const (
	SevUnclassified Severity = "unclassified"
	SevSafe         Severity = "safe"
	SevAffected     Severity = "affected"
	SevUrgent       Severity = "urgent"
)
