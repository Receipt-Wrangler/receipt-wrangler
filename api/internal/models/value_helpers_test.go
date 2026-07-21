package models

import (
	"database/sql/driver"
	"errors"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// assertValuerValid asserts that v.Value() returns (expected, nil). It is shared
// by the enum tests, all of which implement database/sql/driver.Valuer via a
// value-receiver Value() method.
func assertValuerValid(t *testing.T, name string, v driver.Valuer, expected string) {
	t.Run(name, func(t *testing.T) {
		got, err := v.Value()
		if err != nil {
			utils.PrintTestError(t, err, nil)
		}
		if got != expected {
			utils.PrintTestError(t, got, expected)
		}
	})
}

// assertValuerInvalid asserts that v.Value() returns a non-nil error.
func assertValuerInvalid(t *testing.T, name string, v driver.Valuer) {
	t.Run(name, func(t *testing.T) {
		_, err := v.Value()
		if err == nil {
			utils.PrintTestError(t, err, "an error")
		}
	})
}

// errReader is an io.Reader whose Read always fails. It is used to exercise the
// request-body read-error branch of the LoadDataFromRequest handlers.
type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("simulated read error")
}
