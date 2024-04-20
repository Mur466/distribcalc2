package errors

import "errors"

var (
	ErrInvalidCreds = errors.New("invalid credetials")
	ErrIvalidUsername = errors.New("username does not meet the requirements")
	ErrIvalidPassword = errors.New("password does not meet the requirements")
	ErrPasswordMismatch = errors.New("passwords mismatch")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDatabaseInternalError = errors.New("database internal error")
	ErrInvalidInput = errors.New("input data invalid")
)
