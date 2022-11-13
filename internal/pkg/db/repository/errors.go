package repository

import (
	"errors"
	"fmt"
)

type DBError error
type ValidationError error

// ErrBindNoSQLRows used to check the error returned.
var ErrBindNoSQLRows = errors.New("bind failed to execute query: sql: no rows in result set")

func NewDBError(format string, args ...interface{}) DBError {
	return DBError(
		fmt.Errorf(format, args...),
	)
}

func NewValidationError(format string, args ...interface{}) ValidationError {
	return ValidationError(
		fmt.Errorf(format, args...),
	)
}
