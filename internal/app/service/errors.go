package service

import "errors"

var (
	// Service errors:
	ErrNotImplemented  = errors.New("method is not implemented")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrNotFound        = errors.New("object has not been found")
	ErrTimeout         = errors.New("request timed out")
	ErrAlreadyExists   = errors.New("an object already exists")
	ErrDBError         = errors.New("database error")
	ErrCannotCreate    = errors.New("object cannot be created")
	ErrObjAccessDenied = errors.New("access to the object is denied")
	// this error should be returned if only it is related to some algoritmic or broken processing logic
	ErrInternalServerError = errors.New("internal server error")
)
