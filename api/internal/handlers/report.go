package handlers

import (
	"errors"
	"net/http"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// GenerateReport builds and streams a report over one or more groups' receipts.
//
// The request command is parsed and validated up front — before the handler is
// built — because the group ids it carries drive the permission gate: HandleRequest
// re-checks group.reports.read (and membership) in every group the report covers
// before the handler function runs. The generation itself is synchronous; the
// resulting file (or a zip of several formats) is streamed back as an attachment.
func GenerateReport(w http.ResponseWriter, r *http.Request) {
	command := commands.ReportRequestCommand{}
	if err := command.LoadDataFromRequest(w, r); err != nil {
		utils.WriteCustomErrorResponse(w, "Error reading report request", http.StatusInternalServerError)
		return
	}

	if vErr := command.Validate(); len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
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
