package handlers

import (
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage: "Error retrieving roles",
		Writer:       w,
		Request:      r,
		UserRole:     models.ADMIN,
		ResponseType: constants.ApplicationJson,
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
		ErrorMessage: "Error creating role",
		Writer:       w,
		Request:      r,
		UserRole:     models.ADMIN,
		ResponseType: constants.ApplicationJson,
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
		ErrorMessage: "Error updating role",
		Writer:       w,
		Request:      r,
		UserRole:     models.ADMIN,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id, err := utils.StringToUint(chi.URLParam(r, "roleId"))
			if err != nil {
				return http.StatusInternalServerError, err
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
