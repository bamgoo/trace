package trace

import (
	"errors"
)

const NAME = "TRACE"

const (
	StatusOK    = "ok"
	StatusFail  = "fail"
	StatusError = "error"
)

var (
	errInvalidTraceDriver = errors.New("invalid trace driver")
)
