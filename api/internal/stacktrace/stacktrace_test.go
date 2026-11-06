package stacktrace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

var errBoom = errors.New("boom")

func callerLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// wrapHere attaches a stack and reports the line it did it on, so assertions
// do not hard-code line numbers.
func wrapHere(err error) (int, error) {
	return callerLine(), WithStack(err)
}

func TestWithStackKeepsNilNil(t *testing.T) {
	if err := WithStack(nil); err != nil {
		t.Fatalf("WithStack(nil) = %v, want nil", err)
	}
}

func TestWithStackLeavesTheMessageAndIdentityAlone(t *testing.T) {
	_, err := wrapHere(errBoom)

	if err.Error() != errBoom.Error() {
		t.Errorf("Error() = %q, want %q", err.Error(), errBoom.Error())
	}
	if !errors.Is(err, errBoom) {
		t.Error("errors.Is no longer finds the wrapped error")
	}
}

func TestFormatNamesTheFunctionAndLineTheErrorWasReturnedFrom(t *testing.T) {
	line, err := wrapHere(errBoom)

	got := Format(err)
	want := fmt.Sprintf("\n\tinternal/stacktrace.wrapHere (stacktrace_test.go:%d)", line)
	if !strings.HasPrefix(got, want) {
		t.Fatalf("Format() = %q, want it to start with %q", got, want)
	}
	if !strings.Contains(got, "internal/stacktrace.TestFormatNamesTheFunctionAndLineTheErrorWasReturnedFrom") {
		t.Errorf("Format() = %q, want the calling test in the stack", got)
	}
}

func TestFormatShowsOnlyThisModulesFrames(t *testing.T) {
	_, err := wrapHere(errBoom)

	got := Format(err)
	for _, foreign := range []string{"testing.tRunner", "runtime.goexit"} {
		if strings.Contains(got, foreign) {
			t.Errorf("Format() = %q, should not contain %q", got, foreign)
		}
	}
}

func TestWithStackKeepsTheInnermostStack(t *testing.T) {
	innerLine, inner := wrapHere(errBoom)
	outer := fmt.Errorf("loading receipt: %w", inner)

	_, rewrapped := wrapHere(outer)

	if rewrapped != outer {
		t.Fatal("an error that already carries a stack should be returned unchanged")
	}
	if rewrapped.Error() != "loading receipt: boom" {
		t.Errorf("Error() = %q, want the outer message kept", rewrapped.Error())
	}
	want := fmt.Sprintf("(stacktrace_test.go:%d)", innerLine)
	if !strings.HasPrefix(Format(rewrapped), "\n\tinternal/stacktrace.wrapHere "+want) {
		t.Errorf("Format() = %q, want the innermost frame %s first", Format(rewrapped), want)
	}
}

func TestFormatIsEmptyWithoutAStack(t *testing.T) {
	if got := Format(errBoom); got != "" {
		t.Errorf("Format() = %q, want empty for an error with no stack", got)
	}
	if got := Format(nil); got != "" {
		t.Errorf("Format(nil) = %q, want empty", got)
	}
}
