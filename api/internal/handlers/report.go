package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

// loadReportCommand parses and validates the report request body, writing the
// appropriate error response and returning ok=false on any failure. Both the
// generate and preview handlers parse up front — before building the
// structs.Handler — because the groupIds the command carries drive the permission
// gate (HandleRequest re-checks group.reports.read in every covered group).
func loadReportCommand(w http.ResponseWriter, r *http.Request) (commands.ReportRequestCommand, bool) {
	command := commands.ReportRequestCommand{}
	if err := command.LoadDataFromRequest(w, r); err != nil {
		// A decode failure is a malformed client payload (400); only a genuine
		// body-read failure is a server error (500).
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			utils.WriteCustomErrorResponse(w, "Malformed report request", http.StatusBadRequest)
			return command, false
		}
		utils.WriteCustomErrorResponse(w, "Error reading report request", http.StatusInternalServerError)
		return command, false
	}

	if vErr := command.Validate(); len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
		return command, false
	}

	return command, true
}

// GenerateReport builds and streams a report over one or more groups' receipts.
//
// The request command is parsed and validated up front — before the handler is
// built — because the group ids it carries drive the permission gate: HandleRequest
// re-checks group.reports.read (and membership) in every group the report covers
// before the handler function runs. The generation itself is synchronous; the
// resulting file (or a zip of several formats) is streamed back as an attachment.
func GenerateReport(w http.ResponseWriter, r *http.Request) {
	command, ok := loadReportCommand(w, r)
	if !ok {
		return
	}

	handler := structs.Handler{
		ErrorMessage:     "Error generating report",
		Writer:           w,
		Request:          r,
		GroupIds:         command.GroupIds,
		GroupPermissions: []string{permissions.GroupReportsRead},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			reportService := services.NewReportService(nil)

			report, err := reportService.Generate(token.UserId, command)
			if err != nil {
				// A malformed spec is the caller's fault (bad field key, wrong-role
				// column, formula cycle) — surface it as a 400 rather than a 500.
				var specErr *services.ReportSpecError
				if errors.As(err, &specErr) {
					utils.WriteCustomErrorResponse(w, "Invalid report configuration: "+specErr.Error(), http.StatusBadRequest)
					return 0, nil
				}
				return http.StatusInternalServerError, err
			}

			w.Header().Set("Content-Type", report.ContentType)
			w.Header().Set("Content-Disposition", "attachment; filename=\""+report.Filename+"\"")
			w.WriteHeader(http.StatusOK)
			w.Write(report.Bytes)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// PreviewReport renders the current report configuration as HTML for the builder's
// live preview. It shares GenerateReport's front-loaded parse/validate and the
// same per-group gate (group.reports.read in every covered group), but returns a
// JSON { html, receiptCount } body instead of a downloadable file — the preview is
// the engine's own rendered HTML (row-capped), so the builder never re-implements
// the engine.
func PreviewReport(w http.ResponseWriter, r *http.Request) {
	command, ok := loadReportCommand(w, r)
	if !ok {
		return
	}

	handler := structs.Handler{
		ErrorMessage:     "Error generating report preview",
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupIds:         command.GroupIds,
		GroupPermissions: []string{permissions.GroupReportsRead},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			reportService := services.NewReportService(nil)

			preview, err := reportService.Preview(token.UserId, command)
			if err != nil {
				var specErr *services.ReportSpecError
				if errors.As(err, &specErr) {
					utils.WriteCustomErrorResponse(w, "Invalid report configuration: "+specErr.Error(), http.StatusBadRequest)
					return 0, nil
				}
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(preview)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// DeleteReportTemplate removes a saved report template by id. App-scoped behind
// app.reports.delete, matching CreateReportTemplate. Deleting a non-existent id
// still returns 200 (GORM treats it as a no-op), so it is idempotent.
func DeleteReportTemplate(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error deleting report template",
		Writer:         w,
		Request:        r,
		ResponseType:   constants.ApplicationJson,
		AppPermissions: []string{permissions.AppReportsDelete},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id := chi.URLParam(r, "id")

			err := repositories.NewReportTemplateRepository(nil).DeleteReportTemplateById(id)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// CreateReportTemplate saves a report configuration as a reusable template. Unlike
// generate/preview it is app-scoped (app.reports.create): it persists a
// configuration and touches no group's receipts, so it does not gate on per-group
// generation access. It reuses the shared loadReportCommand parse+validate, so a
// saved template is always a complete, buildable configuration.
func CreateReportTemplate(w http.ResponseWriter, r *http.Request) {
	command, ok := loadReportCommand(w, r)
	if !ok {
		return
	}

	// A template is identified by its name, so require a non-empty one. The shared
	// loadReportCommand validator doesn't check Name (generate/preview don't need
	// it), so enforce it here.
	if strings.TrimSpace(command.Name) == "" {
		utils.WriteCustomErrorResponse(w, "A report template name is required", http.StatusBadRequest)
		return
	}

	handler := structs.Handler{
		ErrorMessage:   "Error saving report template",
		Writer:         w,
		Request:        r,
		ResponseType:   constants.ApplicationJson,
		AppPermissions: []string{permissions.AppReportsCreate},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)

			template, err := repositories.NewReportTemplateRepository(nil).CreateReportTemplate(command, token.UserId)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(template)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
			return 0, nil
		},
	}

	HandleRequest(handler)
}
