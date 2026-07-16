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
	reportRouter.Post("/template", handlers.CreateReportTemplate)
	reportRouter.Post("/template/list", handlers.GetPagedReportTemplates)
	reportRouter.Get("/template/{id}", handlers.GetReportTemplate)
	reportRouter.Put("/template/{id}", handlers.UpdateReportTemplate)
	reportRouter.Post("/template/{id}/duplicate", handlers.DuplicateReportTemplate)
	reportRouter.Delete("/template/{id}", handlers.DeleteReportTemplate)

	return reportRouter
}
