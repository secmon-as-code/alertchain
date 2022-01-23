package types

import "github.com/m-mizutani/goerr"

var (
	ErrNoEvalResult         = goerr.New("no evaluation result")
	ErrDatabaseUnexpected   = goerr.New("database failure")
	ErrDatabaseInvalidInput = goerr.New("invalid input for database")
	ErrItemNotFound         = goerr.New("item not found")
	ErrInvalidInput         = goerr.New("invalid input data")

	// Configuration errors
	ErrInvalidChainConfig  = goerr.New("invalid chain config")
	ErrActionNotFound      = goerr.New("specified action is not found")
	ErrActionNotDefined    = goerr.New("specified action ID is not defined")
	ErrDuplicatedActionID  = goerr.New("duplicated action ID found")
	ErrInvalidActionConfig = goerr.New("invalid action config")
)
