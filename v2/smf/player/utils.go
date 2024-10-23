package player

import (
	"context"
	"errors"
	"fmt"
)

// IgnoreError will call function f and ignore error return.
// Useful to explicitly ignore errors in deferred functions.
func IgnoreError(f func() error) {
	_ = f()
}

// WrapOnError will return nil if errInner is nil.
// Otherwise returns an error that wraps both errInner and errOuter.
func WrapOnError(errInner error, errOuter error) error {
	if errInner != nil {
		return fmt.Errorf("%w: %w", errOuter, errInner)
	}
	return nil
}

var ErrUnavailable = errors.New("unavailable context")

// UnavailableContext returns a context that is already canceled
func UnavailableContext() context.Context {
	// Create a context and cancel immediately
	var ctx, cancel = context.WithCancelCause(context.TODO())
	cancel(ErrUnavailable)
	return ctx
}
