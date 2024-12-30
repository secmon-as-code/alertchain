package types

import "github.com/m-mizutani/goerr"

var (
	ErrInvalidOption = AsConfigErr(goerr.New("invalid option"))
	ErrNoPolicyData  = AsConfigErr(goerr.New("no policy data"))

	ErrActionInvalidArgument = AsPolicyErr(goerr.New("invalid action argument"))
	ErrActionNotFound        = AsPolicyErr(goerr.New("action not found"))
	ErrActionFailed          = AsPolicyErr(goerr.New("action failed"))

	ErrInvalidScenario = goerr.New("invalid play scenario")

	ErrInvalidHTTPRequest   = AsBadRequestErr(goerr.New("invalid HTTP request"))
	ErrInvalidLambdaRequest = AsBadRequestErr(goerr.New("invalid Lambda request"))

	ErrNoPolicyResult = goerr.New("no policy result")
)

type ErrorType int

const (
	ErrTypeUnknown ErrorType = iota
	ErrTypeConfig
	ErrTypePolicy
	ErrTypeAction
	ErrTypeRuntime
	ErrTypeBadRequest
)

const (
	errorTypeKey = "errorType"
)

func AsConfigErr(err *goerr.Error) *goerr.Error {
	return err.With(errorTypeKey, ErrTypeConfig)
}
func AsPolicyErr(err *goerr.Error) *goerr.Error {
	return err.With(errorTypeKey, ErrTypePolicy)
}
func AsRuntimeErr(err *goerr.Error) *goerr.Error {
	return err.With(errorTypeKey, ErrTypeRuntime)
}
func AsActionErr(err *goerr.Error) *goerr.Error {
	return err.With(errorTypeKey, ErrTypeAction)
}
func AsBadRequestErr(err *goerr.Error) *goerr.Error {
	return err.With(errorTypeKey, ErrTypeBadRequest)
}

func GetErrorType(err error) ErrorType {
	if err == nil {
		return ErrTypeUnknown
	}

	if e, ok := err.(*goerr.Error); ok {
		if v, ok := e.Values()[errorTypeKey]; ok {
			return v.(ErrorType)
		}
	}

	return ErrTypeUnknown
}
