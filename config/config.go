package config

import (
	"time"
)

// Error messages
const (
	ErrEmptyUserID     = "Bad Request"
	ErrInvalidUserID   = "Bad Request"
	ErrUserNotFound    = "User Not Found"
	ErrFetchingUser    = "Internal Server Error"
	ErrMarshallingUser = "Internal Server Error"
)

// Default values
const (
	DefaultTimeout = 2 * time.Second
	DefaultPort    = "8080"
)
