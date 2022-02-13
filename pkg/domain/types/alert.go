package types

import "github.com/google/uuid"

type (
	AlertID      string
	AttributeID  string
	AnnotationID string
	ReferenceID  string
)

func NewAlertID() AlertID {
	return AlertID(uuid.New().String())
}

type Severity string

const (
	SevUnclassified Severity = "unclassified"
	SevSafe         Severity = "safe"
	SevAffected     Severity = "affected"
	SevUrgent       Severity = "urgent"
)
