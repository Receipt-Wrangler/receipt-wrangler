package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func BuildPermissionRouter() *chi.Mux {
	permissionRouter := chi.NewRouter()

	permissionRouter.Use(middleware.UnifiedAuthMiddleware)
	permissionRouter.Get("/", handlers.GetPermissions)

	return permissionRouter
}
