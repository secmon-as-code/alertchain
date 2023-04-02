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
	ErrConfigNoActionID       = newError("no action ID")
	ErrConfigConflictActionID = newError("conflict action ID")
	ErrConfigNoPolicyPath     = newError("no policy path")
	ErrConfigNoActionName     = newError("no action name")

	ErrNoSuchActionName = newError("no such action name")

	ErrNoSuchActionID = newError("no such action ID")

	ErrInvalidHTTPRequest = newError("invalid HTTP request")
)
