package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// BuildClaimRouter serves the reimbursement claims endpoints (swagger tag
// ReimbursementClaim), mounted at /api/claim.
func BuildClaimRouter() *chi.Mux {
	claimRouter := chi.NewRouter()

	claimRouter.Use(middleware.UnifiedAuthMiddleware)
	claimRouter.Get("/group/{groupId}", handlers.GetClaimsForGroup)
	claimRouter.Post("/group/{groupId}", handlers.CreateClaim)
	claimRouter.Get("/{claimId}", handlers.GetClaimById)
	claimRouter.Put("/{claimId}", handlers.UpdateClaim)
	claimRouter.Delete("/{claimId}", handlers.DeleteClaim)
	claimRouter.Post("/{claimId}/addReceipts", handlers.AddReceiptsToClaim)
	claimRouter.Post("/{claimId}/removeReceipts", handlers.RemoveReceiptsFromClaim)
	claimRouter.Post("/{claimId}/ledgerEntry", handlers.CreateClaimLedgerEntry)
	claimRouter.Put("/{claimId}/ledgerEntry/{ledgerEntryId}", handlers.UpdateClaimLedgerEntry)
	claimRouter.Delete("/{claimId}/ledgerEntry/{ledgerEntryId}", handlers.DeleteClaimLedgerEntry)

	return claimRouter
}
