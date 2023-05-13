package types

type Error struct {
	msg string
}

func (x *Error) Error() string { return x.msg }

func (x *Error) Is(err error) bool {
	if e, ok := err.(*Error); ok {
		return x.msg == e.msg
	}
	return false
}

func newError(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

var (
	ErrConfigConflictActionID = newError("conflict action ID")
	ErrConfigConflictProbeID  = newError("conflict probe ID")
	ErrConfigNoPolicyPath     = newError("no policy path")

	ErrActionInvalidArgument = newError("invalid action argument")

	ErrNoSuchActionName = newError("no such action name")
	ErrNoSuchActionID   = newError("no such action ID")

	ErrInvalidScenario = newError("invalid play scenario")

	ErrInvalidHTTPRequest = newError("invalid HTTP request")

	ErrPolicyClientFailed = newError("policy client failed")
	ErrNoPolicyData       = newError("no policy data")
	ErrNoPolicyResult     = newError("no policy result")

	// runtime errors
	ErrMaxStackDepth  = newError("max stack depth")
	ErrActionNotFound = newError("action not found")
	ErrActionFailed   = newError("action failed")
)
