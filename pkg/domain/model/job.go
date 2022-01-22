package model

import (
	"time"
)

type Job struct {
	Name string `json:"name"`

	Timeout   time.Duration `json:"-"`
	ExitOnErr bool          `json:"-"`
	Actions   []Action      `json:"-"`
}

type JobDefinition struct {
	Job

	Actions []string `json:"actions"`
}
