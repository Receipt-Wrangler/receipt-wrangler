package handlers

import (
	"errors"
	"net/http"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/stacktrace"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Reimbursement claims (the swagger ReimbursementClaim tag). Reading is gated on
// group.receipts.read and every write on group.receipts.update, in the claim's own
// group: claims are a record about the group's receipts, so they reuse the receipt
// permissions rather than adding keys to the Permission enum, which released mobile
// builds parse as a closed set. Member isolation and receipt visibility are applied
// in services.ClaimService.

const invalidRequestBodyMessage = "The request body is not valid"

func GetClaimsForGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")

	HandleRequest(structs.Handler{
		ErrorMessage:     "Error getting claims",
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			uintGroupId, err := utils.StringToUint(groupId)
			if err != nil {
				return http.StatusBadRequest, stacktrace.WithStack(err)
			}

			list, err := services.NewClaimService(nil).GetClaimsForGroup(structs.GetClaims(r).UserId, uintGroupId)
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}
			return writeClaimResponse(w, list)
		},
	})
}

func CreateClaim(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")

	HandleRequest(structs.Handler{
		ErrorMessage:     "Error creating claim",
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupReceiptsUpdate},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			uintGroupId, err := utils.StringToUint(groupId)
			if err != nil {
				return http.StatusBadRequest, stacktrace.WithStack(err)
			}

			command := commands.UpsertClaimCommand{}
			if !loadAndValidate(w, r, &command, command.Validate) {
				return 0, nil
			}

			claim, err := services.NewClaimService(nil).CreateClaim(structs.GetClaims(r).UserId, uintGroupId, command)
			if errors.Is(err, services.ErrClaimAllGroupDestination) {
				utils.WriteCustomErrorResponse(w, "A claim cannot be created in the All group", http.StatusBadRequest)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}
			return writeClaimResponse(w, claim)
		},
	})
}

func GetClaimById(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error getting claim", permissions.GroupReceiptsRead,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			claim, err := services.NewClaimService(nil).GetClaimForUser(structs.GetClaims(r).UserId, claimId)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			return writeClaimResponse(w, claim)
		})
}

func UpdateClaim(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error updating claim", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			command := commands.UpsertClaimCommand{}
			if !loadAndValidate(w, r, &command, command.Validate) {
				return 0, nil
			}

			claim, err := services.NewClaimService(nil).UpdateClaim(structs.GetClaims(r).UserId, claimId, command)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			return writeClaimResponse(w, claim)
		})
}

func DeleteClaim(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error deleting claim", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			err := services.NewClaimService(nil).DeleteClaim(structs.GetClaims(r).UserId, claimId)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			w.WriteHeader(http.StatusOK)
			return 0, nil
		})
}

func AddReceiptsToClaim(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error adding receipts to claim", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			command := commands.ClaimReceiptsCommand{}
			if !loadAndValidate(w, r, &command, command.Validate) {
				return 0, nil
			}

			result, err := services.NewClaimService(nil).AddReceiptsToClaim(structs.GetClaims(r).UserId, claimId, command.ReceiptIds)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			return writeClaimResponse(w, result)
		})
}

func RemoveReceiptsFromClaim(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error removing receipts from claim", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			command := commands.ClaimReceiptsCommand{}
			if !loadAndValidate(w, r, &command, command.Validate) {
				return 0, nil
			}

			result, err := services.NewClaimService(nil).RemoveReceiptsFromClaim(structs.GetClaims(r).UserId, claimId, command.ReceiptIds)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			return writeClaimResponse(w, result)
		})
}

func CreateClaimLedgerEntry(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error adding ledger entry", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			command := commands.UpsertClaimLedgerEntryCommand{}
			if !loadAndValidate(w, r, &command, command.Validate) {
				return 0, nil
			}

			entry, err := services.NewClaimService(nil).CreateLedgerEntry(structs.GetClaims(r).UserId, claimId, command)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			return writeClaimResponse(w, entry)
		})
}

func UpdateClaimLedgerEntry(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error updating ledger entry", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			entryId, ok := ledgerEntryIdFromRequest(w, r)
			if !ok {
				return 0, nil
			}

			command := commands.UpsertClaimLedgerEntryCommand{}
			if !loadAndValidate(w, r, &command, command.Validate) {
				return 0, nil
			}

			entry, err := services.NewClaimService(nil).UpdateLedgerEntry(structs.GetClaims(r).UserId, claimId, entryId, command)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			return writeClaimResponse(w, entry)
		})
}

func DeleteClaimLedgerEntry(w http.ResponseWriter, r *http.Request) {
	handleClaimRequest(w, r, "Error deleting ledger entry", permissions.GroupReceiptsUpdate,
		func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error) {
			entryId, ok := ledgerEntryIdFromRequest(w, r)
			if !ok {
				return 0, nil
			}

			err := services.NewClaimService(nil).DeleteLedgerEntry(structs.GetClaims(r).UserId, claimId, entryId)
			if handled, status, err := claimServiceError(w, err); handled {
				return status, err
			}
			w.WriteHeader(http.StatusOK)
			return 0, nil
		})
}

// refuseDeletingReceiptInClaim answers 400, naming the claim, when the receipt is
// in one, and reports whether it did. A receipt in a claim is evidence of what was
// claimed, so the user removes it from the claim first. Called by DeleteReceipt
// after HandleRequest has authorized the delete.
func refuseDeletingReceiptInClaim(w http.ResponseWriter, userId uint, receiptId string) (bool, error) {
	uintReceiptId, err := utils.StringToUint(receiptId)
	if err != nil {
		return false, err
	}

	refusals, err := services.NewClaimService(nil).ReceiptDeletionRefusals(userId, []uint{uintReceiptId})
	if err != nil {
		return false, err
	}

	message, inClaim := refusals[uintReceiptId]
	if !inClaim {
		return false, nil
	}

	utils.WriteCustomErrorResponse(w, message, http.StatusBadRequest)
	return true, nil
}

// handleClaimRequest resolves the claim's group before HandleRequest, so the group
// permission is checked where the claim lives. A claim that does not exist gets the
// same 403 as one the caller may not reach.
func handleClaimRequest(
	w http.ResponseWriter,
	r *http.Request,
	errorMessage string,
	permission string,
	handlerFunction func(w http.ResponseWriter, r *http.Request, claimId uint) (int, error),
) {
	claimId, err := utils.StringToUint(chi.URLParam(r, "claimId"))
	if err != nil {
		utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
		return
	}

	claim, err := repositories.NewClaimRepository(nil).GetClaimById(claimId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
		return
	}
	if err != nil {
		logHandlerError(r, err)
		utils.WriteCustomErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	HandleRequest(structs.Handler{
		ErrorMessage:     errorMessage,
		Writer:           w,
		Request:          r,
		GroupId:          utils.UintToString(claim.GroupId),
		GroupPermissions: []string{permission},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			return handlerFunction(w, r, claimId)
		},
	})
}

// claimServiceError maps a service error: ErrClaimNotFound is the standard 403,
// anything else a 500. It reports whether there was an error to handle.
func claimServiceError(w http.ResponseWriter, err error) (bool, int, error) {
	if err == nil {
		return false, 0, nil
	}
	if errors.Is(err, services.ErrClaimNotFound) {
		utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
		return true, 0, nil
	}
	return true, http.StatusInternalServerError, stacktrace.WithStack(err)
}

// loadAndValidate decodes a command and runs its validation, answering 400 for an
// unreadable body or a failed check. It reports whether the handler may go on.
func loadAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	command interface {
		LoadDataFromRequest(http.ResponseWriter, *http.Request) error
	},
	validate func() structs.ValidatorError,
) bool {
	if err := command.LoadDataFromRequest(w, r); err != nil {
		logHandlerError(r, err)
		utils.WriteCustomErrorResponse(w, invalidRequestBodyMessage, http.StatusBadRequest)
		return false
	}

	vErr := validate()
	if len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
		return false
	}
	return true
}

func ledgerEntryIdFromRequest(w http.ResponseWriter, r *http.Request) (uint, bool) {
	entryId, err := utils.StringToUint(chi.URLParam(r, "ledgerEntryId"))
	if err != nil {
		utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
		return 0, false
	}
	return entryId, true
}

func writeClaimResponse(w http.ResponseWriter, data interface{}) (int, error) {
	bytes, err := utils.MarshalResponseData(data)
	if err != nil {
		return http.StatusInternalServerError, stacktrace.WithStack(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return 0, nil
}
