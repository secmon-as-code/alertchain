package types

type Error struct {
	msg string
}

func (x *Error) Error() string { return x.msg }

func newError(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

var (
	ErrConfigNoActionID   = newError("no action ID")
	ErrConfigNoActionName = newError("no action name")
	ErrConfigNoProbeID    = newError("no probe ID")
	ErrConfigNoProbeName  = newError("no probe name")

	ErrConfigConflictActionID = newError("conflict action ID")
	ErrConfigConflictProbeID  = newError("conflict probe ID")
	ErrConfigNoPolicyPath     = newError("no policy path")

	ErrActionInvalidConfig = newError("invalid action config")

	ErrNoSuchActionName = newError("no such action name")
	ErrNoSuchActionID   = newError("no such action ID")
	ErrNoSuchProbeName  = newError("no such probe name")
	ErrNoSuchProbeID    = newError("no such probe ID")

	ErrInvalidHTTPRequest = newError("invalid HTTP request")

	// runtime errors
	ErrMaxStackDepth = newError("max stack depth")
)
