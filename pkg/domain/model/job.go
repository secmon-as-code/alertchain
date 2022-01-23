package model

import (
	"time"
)

type Job struct {
	Name      string `json:"name"`
	ExitOnErr bool   `json:"exit_on_error"`

	Timeout time.Duration `json:"-"`
	Actions []Action      `json:"-"`
}

type JobDefinition struct {
	Job

	Actions []string `json:"actions"`
	Timeout string   `json:"timeout"`
}
