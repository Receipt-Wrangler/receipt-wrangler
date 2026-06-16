package handlers

import (
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error retrieving roles",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppRolesRead},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			roleService := services.NewRoleService(nil)
			roles, err := roleService.GetRoles()
			if err != nil {
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(&roles)
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

func CreateRole(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error creating role",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppRolesCreate},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.UpsertRoleCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErrs := command.Validate()
			if len(vErrs.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErrs, http.StatusBadRequest)
				return 0, nil
			}

			roleService := services.NewRoleService(nil)
			createdRole, err := roleService.CreateRole(command)
			if err != nil {
				if errors.Is(err, services.ErrInvalidGrant) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"grants": "One or more category or tag grants do not exist"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(&createdRole)
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

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error updating role",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppRolesUpdate},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id, err := utils.StringToUint(chi.URLParam(r, "roleId"))
			if err != nil {
				structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
					Errors: map[string]string{"roleId": "Invalid role id"},
				}, http.StatusBadRequest)
				return 0, nil
			}

			command := commands.UpsertRoleCommand{}
			err = command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErrs := command.Validate()
			if len(vErrs.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErrs, http.StatusBadRequest)
				return 0, nil
			}

			roleService := services.NewRoleService(nil)
			updatedRole, err := roleService.UpdateRole(id, command)
			if err != nil {
				if errors.Is(err, services.ErrRoleTypeMismatch) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"scope": "Role type cannot be changed"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrSystemRoleImmutable) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"role": "System roles cannot be modified"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrRoleNotFound) {
					utils.WriteCustomErrorResponse(w, "Role not found", http.StatusNotFound)
					return 0, nil
				}
				if errors.Is(err, services.ErrInvalidGrant) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"grants": "One or more category or tag grants do not exist"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(&updatedRole)
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

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error deleting role",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppRolesDelete},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id, err := utils.StringToUint(chi.URLParam(r, "roleId"))
			if err != nil {
				structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
					Errors: map[string]string{"roleId": "Invalid role id"},
				}, http.StatusBadRequest)
				return 0, nil
			}

			scope := permissions.Scope(r.URL.Query().Get("scope"))
			if scope != permissions.ScopeApp && scope != permissions.ScopeGroup {
				structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
					Errors: map[string]string{"scope": "Scope must be either APP or GROUP"},
				}, http.StatusBadRequest)
				return 0, nil
			}

			roleService := services.NewRoleService(nil)
			err = roleService.DeleteRole(id, scope)
			if err != nil {
				if errors.Is(err, services.ErrRoleTypeMismatch) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"scope": "Role type cannot be changed"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrSystemRoleUndeletable) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"role": "System roles cannot be deleted"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrRoleAssigned) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"role": "Role is assigned and cannot be deleted"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrRoleIsDefault) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"role": "The default role cannot be deleted; choose another default first"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrRoleNotFound) {
					utils.WriteCustomErrorResponse(w, "Role not found", http.StatusNotFound)
					return 0, nil
				}
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

func SetDefaultRole(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error setting default role",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppRolesUpdate},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id, err := utils.StringToUint(chi.URLParam(r, "roleId"))
			if err != nil {
				structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
					Errors: map[string]string{"roleId": "Invalid role id"},
				}, http.StatusBadRequest)
				return 0, nil
			}

			scope := permissions.Scope(r.URL.Query().Get("scope"))
			if scope != permissions.ScopeApp && scope != permissions.ScopeGroup {
				structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
					Errors: map[string]string{"scope": "Scope must be either APP or GROUP"},
				}, http.StatusBadRequest)
				return 0, nil
			}

			roleService := services.NewRoleService(nil)
			updatedRole, err := roleService.SetDefaultRole(id, scope)
			if err != nil {
				if errors.Is(err, services.ErrRoleTypeMismatch) {
					structs.WriteValidatorErrorResponse(w, structs.ValidatorError{
						Errors: map[string]string{"scope": "Role type cannot be changed"},
					}, http.StatusBadRequest)
					return 0, nil
				}
				if errors.Is(err, services.ErrRoleNotFound) {
					utils.WriteCustomErrorResponse(w, "Role not found", http.StatusNotFound)
					return 0, nil
				}
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(&updatedRole)
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
