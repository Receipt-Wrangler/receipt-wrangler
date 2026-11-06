package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/stacktrace"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

func GetFeatureConfig(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage: "Error getting feature config.",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			systemSettingsService := services.NewSystemSettingsService(nil)

			featureConfig, err := systemSettingsService.GetFeatureConfig()
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			bytes, err := utils.MarshalResponseData(featureConfig)
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			w.WriteHeader(200)
			w.Write(bytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}
