package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func BuildReportRouter() *chi.Mux {
	reportRouter := chi.NewRouter()

	reportRouter.Use(middleware.UnifiedAuthMiddleware)
	reportRouter.Post("/generate", handlers.GenerateReport)
	reportRouter.Post("/preview", handlers.PreviewReport)

	return reportRouter
}
