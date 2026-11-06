package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
)

// captureAppLog redirects what LogStd writes to the app log into a buffer for
// the rest of the test.
func captureAppLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	logger := logging.GetLogger()
	previous := logger.Writer()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	t.Cleanup(func() { logger.SetOutput(previous) })
	return &buf
}

// The case that prompted this: a client sends "" for a decimal, and
// encoding/json reports it without naming the field. The log used to say only
// "[generic_handler.go:71] ERROR: error decoding string ...", the same line for
// every endpoint. It has to say which request failed, and where in the handler.
func TestHandlerErrorLog_NamesTheRequestAndTheLineThatFailed(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedGroupWithRole(t, "error-log", []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate})
	receiptId := seedReceipt(t, groupId, userId, 0)
	logged := captureAppLog(t)

	path := "/api/receipt/" + utils.UintToString(receiptId)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", path+"?token=secret-value", strings.NewReader(`{"name":"x","amount":""}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", utils.UintToString(receiptId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	UpdateReceipt(w, r)

	assertStatus(t, w, http.StatusInternalServerError)
	got := logged.String()
	// The innermost frame comes straight after the message: the handler line
	// that returned the decode error, not HandleRequest.
	want := "PUT " + path + `: error decoding string '""': can't convert "" to decimal` +
		"\n\tinternal/handlers.UpdateReceipt.func1 (receipts.go:"
	if !strings.Contains(got, want) {
		t.Errorf("log = %q\nwant it to contain %q", got, want)
	}
	if !strings.Contains(got, "\n\tinternal/handlers.UpdateReceipt (receipts.go:") {
		t.Errorf("log = %q\nwant the endpoint's handler in the stack", got)
	}
	if strings.Contains(got, "secret-value") {
		t.Errorf("log = %q\nmust not contain the query string", got)
	}
}

// A handler that returns an error without stacktrace.WithStack still gets a
// stack from HandleRequest, which names the handler that called it.
func TestHandlerErrorLog_AnErrorWithoutAStackStillNamesTheCaller(t *testing.T) {
	defer tearDownGenericHandlerTest()
	logged := captureAppLog(t)
	w, r := requestForUser(1)

	HandleRequest(structs.Handler{
		ErrorMessage: "failed",
		Writer:       w,
		Request:      r,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			return http.StatusInternalServerError, errors.New("boom")
		},
	})

	assertStatus(t, w, http.StatusInternalServerError)
	got := logged.String()
	for _, want := range []string{
		"GET /api: boom\n",
		"\n\tinternal/handlers.HandleRequest (generic_handler.go:",
		"\n\tinternal/handlers.TestHandlerErrorLog_AnErrorWithoutAStackStillNamesTheCaller (handler_error_logging_test.go:",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("log = %q\nwant it to contain %q", got, want)
		}
	}
}
