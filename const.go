package trace

import (
	"errors"
)

const NAME = "TRACE"

const (
	StatusOK    = "ok"
	StatusError = "error"
)

var (
	errInvalidTraceDriver = errors.New("invalid trace driver")
)
