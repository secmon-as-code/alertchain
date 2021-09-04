package types

type ExecStatus string

const (
	ExecStart   ExecStatus = "start"
	ExecSucceed ExecStatus = "succeed"
	ExecFailure ExecStatus = "failure"
)
