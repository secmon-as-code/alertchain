package types

import "github.com/m-mizutani/goerr"

var (
	ErrDatabaseUnexpected   = goerr.New("database failure")
	ErrDatabaseInvalidInput = goerr.New("invalid input for database")
	ErrInvalidChain         = goerr.New("invalid chain plugin")
	ErrInvalidInput         = goerr.New("invalid input data")
)
