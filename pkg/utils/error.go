package utils

import (
	"errors"
	"fmt"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func CopyErrorToExecLog(err error, execLog *ent.ExecLog) {
	execLog.Errmsg = err.Error()
	var goErr *goerr.Error
	if errors.As(err, &goErr) {
		for k, v := range goErr.Values() {
			execLog.ErrValues = append(execLog.ErrValues, fmt.Sprintf("%s=%v", k, v))
		}
		for _, st := range goErr.StackTrace() {
			execLog.StackTrace = append(execLog.StackTrace, fmt.Sprintf("%v", st))
		}
	}
	execLog.Status = types.ExecFailure
}

func HandleError(err error) {
	log := Logger.Log()

	var goErr *goerr.Error
	if errors.As(err, &goErr) {
		for k, v := range goErr.Values() {
			log = log.With(k, v)
		}
		for _, st := range goErr.Stacks() {
			Logger.With("stack", st).Debug("stack trace")
		}
	}
	log.Error(err.Error())
}
