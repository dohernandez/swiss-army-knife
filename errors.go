package swissarmyknife

import "errors"

var (
	// ErrTypeMismatch is returned when the casting is not ok.
	ErrTypeMismatch = errors.New("casting type is not ok, type mismatch")

	// ErrDoNotEmit is returned when the operation don't want to emit the current value to the next operation.
	ErrDoNotEmit = errors.New("do not emit")
)
