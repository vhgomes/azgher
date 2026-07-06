package errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrInternalServerError = errors.New("internal server error")
var ErrEnvVarNotFound = errors.New("environment variable not found")
var ErrUserNotFound = errors.New("user not found")
