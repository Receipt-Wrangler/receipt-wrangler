package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"
	"receipt-wrangler/api/internal/oidc"

	"github.com/go-chi/chi/v5"
)

// BuildOidcRouter mounts the relying-party flow at /api/oidc.
//
// The static segments (link, exchange, connections) are declared as their own
// paths rather than under {name}, so chi's static-beats-param resolution stays
// unambiguous. The command validator additionally rejects those words as provider
// slugs, so no configured provider can shadow a route from the other direction.
func BuildOidcRouter() *chi.Mux {
	router := chi.NewRouter()

	// Unauthenticated: this IS the sign-in path.
	router.Get("/{name}/login", oidc.Login)
	router.Get("/{name}/callback", oidc.Callback)
	router.Post("/exchange", handlers.OidcExchange)

	// Authenticated: connecting a provider to an account that already exists. The
	// session is what proves identity here, so nothing has to be inferred from a
	// claim -- which is why linkByUsername can safely default to off.
	router.Group(func(authenticated chi.Router) {
		authenticated.Use(middleware.UnifiedAuthMiddleware)

		authenticated.Get("/link/{name}", oidc.LinkStart)
		authenticated.Get("/connections", handlers.GetOidcConnections)
		authenticated.Delete("/connections/{name}", handlers.DeleteOidcConnection)
	})

	return router
}

// BuildOidcProviderRouter mounts administrator CRUD at /api/oidcProvider.
func BuildOidcProviderRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.UnifiedAuthMiddleware)

	router.Get("/{oidcProviderId}", handlers.GetOidcProviderById)
	router.Put("/{oidcProviderId}", handlers.UpdateOidcProvider)
	router.Delete("/{oidcProviderId}", handlers.DeleteOidcProvider)
	router.Post("/getPagedOidcProviders", handlers.GetPagedOidcProviders)
	router.Post("/", handlers.CreateOidcProvider)

	return router
}
