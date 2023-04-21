package errors

import (
	"runtime"

	"github.com/pkg/errors"
)

const depth = 32

type stack []uintptr

type errWrapper struct {
	error
	stack *stack
}

func callers(pos int) *stack {
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[pos:n]
	return &st
}

func wrapErr(err error, caller int) *errWrapper {
	wrapper := &errWrapper{}
	wrapper.error = err
	wrapper.stack = callers(caller)

	return wrapper
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	wrapper, ok := err.(*errWrapper)
	if !ok || wrapper == nil {
		wrapper = wrapErr(err, 1)
	}

	return wrapper
}

// New retrun new error annotated with stack trace
func New(message string) error {
	return wrapErr(errors.New(message), 1)
}
