// Package stacktrace attaches a call stack to an error at the point it is
// returned, so whoever finally logs it can say where it came from.
//
// A Go error carries no location of its own. By the time one reaches a central
// logger such as handlers.HandleRequest, the frames that produced it have already
// returned, so a stack taken at the logger only ever shows the logger. The stack
// has to be captured while those frames are still live, which is what WithStack
// is for.
package stacktrace

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// modulePrefix marks the frames Format keeps. Library and runtime frames (chi,
// net/http, goroutine start-up) are the same on every request and bury the ones
// that differ.
const modulePrefix = "receipt-wrangler/api/"

// maxDepth comfortably covers a handler closure, HandleRequest, the handler and
// the router's middleware chain.
const maxDepth = 32

type withStack struct {
	err error
	pcs []uintptr
}

// Error is the wrapped error's message, unchanged, so nothing that shows or
// compares the text sees a difference.
func (e *withStack) Error() string { return e.err.Error() }

// Unwrap keeps errors.Is and errors.As working through the wrapper.
func (e *withStack) Unwrap() error { return e.err }

// WithStack returns err with its caller's stack attached. A nil error stays nil,
// so `return status, stacktrace.WithStack(err)` is safe on any path. An error
// that already carries a stack is returned unchanged: the innermost stack is the
// one closest to the failure, so it is the one worth keeping.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	var existing *withStack
	if errors.As(err, &existing) {
		return err
	}

	pcs := make([]uintptr, maxDepth)
	// Skip runtime.Callers and WithStack itself.
	n := runtime.Callers(2, pcs)
	return &withStack{err: err, pcs: pcs[:n]}
}

// Format renders the stack attached to err, innermost frame first, one frame per
// line, each starting with a newline and a tab so the result can be appended
// straight to a log message. It returns "" when err carries no stack. Only this
// module's frames are shown unless there are none, in which case every frame is.
func Format(err error) string {
	var ws *withStack
	if !errors.As(err, &ws) {
		return ""
	}

	var all, app []runtime.Frame
	frames := runtime.CallersFrames(ws.pcs)
	for {
		frame, more := frames.Next()
		all = append(all, frame)
		if strings.HasPrefix(frame.Function, modulePrefix) {
			app = append(app, frame)
		}
		if !more {
			break
		}
	}

	shown := app
	if len(shown) == 0 {
		shown = all
	}

	var b strings.Builder
	for _, frame := range shown {
		fmt.Fprintf(&b, "\n\t%s (%s:%d)", strings.TrimPrefix(frame.Function, modulePrefix), filepath.Base(frame.File), frame.Line)
	}
	return b.String()
}
