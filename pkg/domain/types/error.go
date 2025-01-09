package types

import "github.com/m-mizutani/goerr/v2"

var (
	ErrNoPolicyResult        = goerr.New("no policy result")
	ErrActionInvalidArgument = goerr.New("invalid action argument", goerr.T(ErrTagAction))
	/*
		ErrInvalidOption = AsConfigErr(goerr.New("invalid option"))

		ErrActionInvalidArgument = AsPolicyErr(goerr.New("invalid action argument"))
		ErrActionNotFound        = AsPolicyErr(goerr.New("action not found"))
		ErrActionFailed          = AsPolicyErr(goerr.New("action failed"))

		ErrInvalidScenario = goerr.New("invalid play scenario")

		ErrInvalidHTTPRequest   = AsBadRequestErr(goerr.New("invalid HTTP request"))
		ErrInvalidLambdaRequest = AsBadRequestErr(goerr.New("invalid Lambda request"))
	*/
)

var (
	// ErrTagConfig is a tag for configuration and startup option error.
	ErrTagConfig = goerr.NewTag("config")

	// ErrTagPolicy is a tag for policy error. It is used for failure of policy evaluation or invalid policy result.
	ErrTagPolicy = goerr.NewTag("policy")

	// ErrTagAction is a tag for action error. It is used for failure of action execution or invalid action argument.
	ErrTagAction = goerr.NewTag("action")

	// ErrTagBadRequest is a tag for bad request to AlertChain server or runtime.
	ErrTagBadRequest = goerr.NewTag("bad_request")

	// ErrTagSystem is a tag for unexpected system behavior. E.g. I/O error, system call failure, database error, error from integrated system, connection error, etc.
	ErrTagSystem = goerr.NewTag("system")
)
