package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
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

// authorizeTemplateAction loads the {id} template, confirms it exists, and checks
// the caller may perform action on it (the "*All" bypass, the per-group ceiling,
// and the per-template matrix — all resolved by CanActOnTemplate). It writes the
// 404/403/500 response and returns ok=false on any failure; on success it returns
// the loaded template so the handler can reuse it. Because base-OR-"*All" is an OR
// the declarative AppPermissions gate can't express, the six template handlers do
// their app-permission check here rather than on the structs.Handler.
func authorizeTemplateAction(w http.ResponseWriter, r *http.Request, action string) (models.ReportTemplate, bool) {
	id := chi.URLParam(r, "id")

	template, err := repositories.NewReportTemplateRepository(nil).GetReportTemplateById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteCustomErrorResponse(w, "Report template not found", http.StatusNotFound)
			return models.ReportTemplate{}, false
		}
		utils.WriteCustomErrorResponse(w, "Error loading report template", http.StatusInternalServerError)
		return models.ReportTemplate{}, false
	}

	token := structs.GetClaims(r)
	allowed, err := services.NewPermissionService(nil).CanActOnTemplate(token.UserId, template.ID, action)
	if err != nil {
		utils.WriteCustomErrorResponse(w, "Error authorizing report template access", http.StatusInternalServerError)
		return models.ReportTemplate{}, false
	}
	if !allowed {
		utils.WriteCustomErrorResponse(w, "User is unauthorized to access this report template", http.StatusForbidden)
		return models.ReportTemplate{}, false
	}

	return template, true
}

// GenerateReport builds and streams a report over one or more groups' receipts.
//
// The request command is parsed and validated up front — before the handler is
// built — because the group ids it carries drive the permission gate: HandleRequest
// re-checks group.reports.read (and membership) in every group the report covers
// before the handler function runs. It also requires the app-level
// app.reports.generate (ANDed with the per-group check). The generation itself is
// synchronous; the resulting file (or a zip of several formats) is streamed back as
// an attachment.
func GenerateReport(w http.ResponseWriter, r *http.Request) {
	command, ok := loadReportCommand(w, r)
	if !ok {
		return
	}

	handler := structs.Handler{
		ErrorMessage:     "Error generating report",
		Writer:           w,
		Request:          r,
		AppPermissions:   []string{permissions.AppReportsGenerate},
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

// DeleteReportTemplate removes a saved report template by id. Access is resolved by
// CanActOnTemplate (delete): a missing id maps to 404, an unauthorized caller to 403.
// On success the grant cache is flushed, since the deleted template's grant rows
// cascade out of every role's matrix.
func DeleteReportTemplate(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage: "Error deleting report template",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			template, ok := authorizeTemplateAction(w, r, "delete")
			if !ok {
				return 0, nil
			}

			err := repositories.NewReportTemplateRepository(nil).DeleteReportTemplateById(utils.UintToString(template.ID))
			if err != nil {
				return http.StatusInternalServerError, err
			}

			services.EvictAllGroupRoleGrants()

			w.WriteHeader(http.StatusOK)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// GetPagedReportTemplates returns a paged, sorted list of saved report templates.
// Visibility is resolved per user (readAll sees all; otherwise only templates whose
// covered groups the caller can read AND whose read action the per-template matrix
// permits), and each row carries the caller's allowed actions so the client can gate
// its buttons. The app-permission gate lives in VisibleTemplateIds (read/readAll).
func GetPagedReportTemplates(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage: "Error getting report templates",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			permissionService := services.NewPermissionService(nil)

			// Report access at all is app.reports.read OR readAll (an OR the declarative
			// AppPermissions gate can't express); deny outright without either.
			canList, err := permissionService.HasAnyAppPermission(token.UserId, permissions.AppReportsRead, permissions.AppReportsReadAll)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			if !canList {
				utils.WriteCustomErrorResponse(w, "User is unauthorized to access report templates", http.StatusForbidden)
				return 0, nil
			}

			command := commands.PagedRequestCommand{}
			if err := command.LoadDataFromRequest(w, r); err != nil {
				return http.StatusInternalServerError, err
			}

			vErr := command.Validate()
			if len(vErr.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
				return 0, nil
			}

			visibleIds, unrestricted, err := permissionService.VisibleTemplateIds(token.UserId)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			var idFilter *[]uint
			if !unrestricted {
				idFilter = &visibleIds
			}

			templates, count, err := repositories.NewReportTemplateRepository(nil).GetPagedReportTemplates(command, idFilter)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			pagedData := structs.PagedData{}
			data := make([]interface{}, 0, len(templates))
			for i := 0; i < len(templates); i++ {
				actions, err := permissionService.AllowedActionsForTemplate(token.UserId, templates[i].ID)
				if err != nil {
					return http.StatusInternalServerError, err
				}
				templates[i].AllowedActions = actions
				data = append(data, templates[i])
			}
			pagedData.TotalCount = count
			pagedData.Data = data

			responseBytes, err := utils.MarshalResponseData(pagedData)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// GetReportTemplate returns a saved report template by id. Access is resolved by
// CanActOnTemplate (read): a missing id maps to 404, an unauthorized caller to 403.
func GetReportTemplate(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage: "Error getting report template",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			template, ok := authorizeTemplateAction(w, r, "read")
			if !ok {
				return 0, nil
			}

			responseBytes, err := utils.MarshalResponseData(template)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// DuplicateReportTemplate copies a saved report template and returns the new copy.
// Access to the source is resolved by CanActOnTemplate (duplicate); a missing source
// id maps to 404, an unauthorized caller to 403. The copy is owned by the caller,
// its name suffixed " duplicate", and it starts unrestricted (no grant rows).
func DuplicateReportTemplate(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage: "Error duplicating report template",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)

			source, ok := authorizeTemplateAction(w, r, "duplicate")
			if !ok {
				return 0, nil
			}

			template, err := repositories.NewReportTemplateRepository(nil).DuplicateReportTemplate(token.UserId, utils.UintToString(source.ID))
			if err != nil {
				return http.StatusInternalServerError, err
			}

			responseBytes, err := utils.MarshalResponseData(template)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)
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
		ErrorMessage: "Error saving report template",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			permissionService := services.NewPermissionService(nil)

			// createAll implies the ability to create, so accept either. Then the
			// caller must be able to report over every group the template attaches to
			// (createAll bypasses that ceiling).
			canCreate, err := permissionService.HasAnyAppPermission(token.UserId, permissions.AppReportsCreate, permissions.AppReportsCreateAll)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			if !canCreate {
				utils.WriteCustomErrorResponse(w, "User is unauthorized to create report templates", http.StatusForbidden)
				return 0, nil
			}

			canReport, err := permissionService.CanReportOverGroups(token.UserId, command.GroupIds, permissions.AppReportsCreateAll)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			if !canReport {
				utils.WriteCustomErrorResponse(w, "User is unauthorized to create a report over one or more of these groups", http.StatusForbidden)
				return 0, nil
			}

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

// UpdateReportTemplate overwrites a saved report template in place with a new
// configuration. It mirrors CreateReportTemplate (shared loadReportCommand
// parse+validate, same non-empty-name guard) but targets an existing {id} and is
// app-scoped behind app.reports.update; a missing target id maps to a 404. The
// template's id and owner are preserved — only its name/config/version change.
func UpdateReportTemplate(w http.ResponseWriter, r *http.Request) {
	command, ok := loadReportCommand(w, r)
	if !ok {
		return
	}

	// A template is identified by its name, so require a non-empty one (the shared
	// loadReportCommand validator doesn't check Name), matching CreateReportTemplate.
	if strings.TrimSpace(command.Name) == "" {
		utils.WriteCustomErrorResponse(w, "A report template name is required", http.StatusBadRequest)
		return
	}

	handler := structs.Handler{
		ErrorMessage: "Error updating report template",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			permissionService := services.NewPermissionService(nil)

			existing, ok := authorizeTemplateAction(w, r, "update")
			if !ok {
				return 0, nil
			}

			// Retargeting the template onto new groups requires reporting access over
			// each of them (updateAll bypasses that ceiling).
			canReport, err := permissionService.CanReportOverGroups(token.UserId, command.GroupIds, permissions.AppReportsUpdateAll)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			if !canReport {
				utils.WriteCustomErrorResponse(w, "User is unauthorized to report over one or more of these groups", http.StatusForbidden)
				return 0, nil
			}

			template, err := repositories.NewReportTemplateRepository(nil).UpdateReportTemplate(command, utils.UintToString(existing.ID))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					utils.WriteCustomErrorResponse(w, "Report template not found", http.StatusNotFound)
					return 0, nil
				}
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
