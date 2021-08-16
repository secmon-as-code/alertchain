package types

import "github.com/m-mizutani/goerr"

var (
	ErrDatabaseUnexpected = goerr.New("database failure")
	ErrInvalidChain       = goerr.New("invalid chain plugin")
	ErrInvalidInput       = goerr.New("invalid input data")
)
