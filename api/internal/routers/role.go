package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func BuildRoleRouter() *chi.Mux {
	roleRouter := chi.NewRouter()

	roleRouter.Use(middleware.UnifiedAuthMiddleware)

	roleRouter.Get("/", handlers.GetRoles)
	roleRouter.Post("/", handlers.CreateRole)
	roleRouter.Put("/{roleId}", handlers.UpdateRole)
	roleRouter.Delete("/{roleId}", handlers.DeleteRole)

	return roleRouter
}
