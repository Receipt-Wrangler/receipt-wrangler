package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func BuildWidgetRouter() *chi.Mux {
	widgetRouter := chi.NewRouter()

	widgetRouter.Use(middleware.UnifiedAuthMiddleware)
	widgetRouter.Post("/pieChart/{groupId}", handlers.GetPieChartData)

	return widgetRouter
}
