package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

func CreateApiKey(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error creating API key",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppApiKeysCreate},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.UpsertApiKeyCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErrs := command.Validate()
			if len(vErrs.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErrs, http.StatusBadRequest)
				return 0, nil
			}

			token := structs.GetClaims(r)
			apiKeyService := services.NewApiKeyService(nil)

			generatedKey, err := apiKeyService.CreateApiKey(token.UserId, command)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			response := structs.ApiKeyResult{
				Key: generatedKey,
			}

			bytes, err := utils.MarshalResponseData(response)
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

func GetPagedApiKeys(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error retrieving API keys.",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppApiKeysRead},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.PagedApiKeyRequestCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErrs := command.Validate(r)
			if len(vErrs.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErrs, http.StatusBadRequest)
				return 0, nil
			}

			token := structs.GetClaims(r)

			// Viewing every user's API keys requires a dedicated permission,
			// resolved from the database (never trusted from the JWT).
			if command.ApiKeyFilter.AssociatedApiKeys == commands.ASSOCIATED_API_KEYS_ALL {
				permissionService := services.NewPermissionService(nil)
				canReadAny, err := permissionService.HasAppPermissions(token.UserId, permissions.AppApiKeysReadAny)
				if err != nil {
					return http.StatusInternalServerError, err
				}
				if !canReadAny {
					return http.StatusForbidden, errors.New("user is unauthorized to view all API keys")
				}
			}

			userIdString := utils.UintToString(token.UserId)
			apiKeyService := services.NewApiKeyService(nil)

			apiKeys, count, err := apiKeyService.GetPagedApiKeys(command, userIdString)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			anyData := make([]any, len(apiKeys))
			for i := 0; i < len(apiKeys); i++ {
				anyData[i] = apiKeys[i]
			}

			bytes, err := utils.MarshalResponseData(structs.PagedData{
				TotalCount: count,
				Data:       anyData,
			})
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

func UpdateApiKey(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error updating API key",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppApiKeysUpdate},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id := chi.URLParam(r, "id")
			// URL decode the ID parameter in case it was encoded by the frontend
			if decodedId, err := url.QueryUnescape(id); err == nil {
				id = decodedId
			}
			command := commands.UpsertApiKeyCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErrs := command.Validate()
			if len(vErrs.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErrs, http.StatusBadRequest)
				return 0, nil
			}

			token := structs.GetClaims(r)
			apiKeyService := services.NewApiKeyService(nil)

			err = apiKeyService.UpdateApiKey(id, token.UserId, command)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			return 0, nil
		},
	}

	HandleRequest(handler)
}

func DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error deleting API key",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppApiKeysDelete},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			id := chi.URLParam(r, "id")
			// URL decode the ID parameter in case it was encoded by the frontend
			if decodedId, err := url.QueryUnescape(id); err == nil {
				id = decodedId
			}

			token := structs.GetClaims(r)
			apiKeyService := services.NewApiKeyService(nil)

			// Whether the caller may delete API keys belonging to other users is
			// resolved from the database, never trusted from the JWT.
			permissionService := services.NewPermissionService(nil)
			canDeleteAny, err := permissionService.HasAppPermissions(token.UserId, permissions.AppApiKeysDeleteAny)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			err = apiKeyService.DeleteApiKey(id, token.UserId, canDeleteAny)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			return 0, nil
		},
	}

	HandleRequest(handler)
}
