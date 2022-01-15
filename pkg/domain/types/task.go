package types

type ExecStatus string

const (
	ExecRunning ExecStatus = "running"
	ExecSucceed ExecStatus = "succeed"
	ExecFailure ExecStatus = "failure"
)
