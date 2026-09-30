package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
)

const (
	receiptCommentRequiredMessage = "At least one comment is required"
	receiptImageRequiredMessage   = "At least one image is required"
	// legacyCreateImageRequiredMessage is returned by the deprecated JSON create,
	// which can never carry an image: only a client too old to know the multipart
	// create still calls it, so the fix is on the user's side.
	legacyCreateImageRequiredMessage = "This group requires an image on every receipt. Please update your app to create receipts with images."
	receiptCommentDeleteMessage      = "This receipt must keep at least one comment"
	receiptImageDeleteMessage        = "This receipt must keep at least one image"
)

// enforceReceiptRequirements rejects a receipt write that would leave a receipt in
// groupId without a comment or an image while the caller's group role requires
// one there (ReceiptService.ResolveReceiptRequirements, waivers included). It
// writes a 400 validator error keyed `comments` / `files` itself and returns
// false when it rejected the write. imageRequiredMessage lets a caller word the
// image error for its context.
func enforceReceiptRequirements(
	w http.ResponseWriter,
	userId uint,
	groupId uint,
	hasComment bool,
	hasImage bool,
	imageRequiredMessage string,
) (bool, error) {
	requirements, err := services.NewReceiptService(nil).ResolveReceiptRequirements(userId, groupId)
	if err != nil {
		return false, err
	}

	vErr := structs.ValidatorError{Errors: make(map[string]string)}
	if requirements.CommentRequired && !hasComment {
		vErr.Errors["comments"] = receiptCommentRequiredMessage
	}
	if requirements.ImageRequired && !hasImage {
		vErr.Errors["files"] = imageRequiredMessage
	}

	if len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
		return false, nil
	}

	return true, nil
}

// receiptCommandHasComment reports whether a receipt command carries at least one
// comment with non-blank text.
func receiptCommandHasComment(command commands.UpsertReceiptCommand) bool {
	for _, comment := range command.Comments {
		if len(strings.TrimSpace(comment.Comment)) > 0 {
			return true
		}
	}
	return false
}

// receiptHasComment reports whether a stored receipt has at least one comment with
// non-blank text.
func receiptHasComment(receipt models.Receipt) bool {
	for _, comment := range receipt.Comments {
		if len(strings.TrimSpace(comment.Comment)) > 0 {
			return true
		}
	}
	return false
}

// enforceCommentDeleteKeepsRequired refuses deleting a receipt's last comment
// while the caller's role requires one in the receipt's group. Replacing the only
// comment therefore means adding the new one first. Writes the 400 itself and
// returns false when it refused.
func enforceCommentDeleteKeepsRequired(w http.ResponseWriter, userId uint, receiptId uint, commentId string) (bool, error) {
	requirements, err := resolveRequirementsForReceipt(userId, receiptId)
	if err != nil {
		return false, err
	}
	if !requirements.CommentRequired {
		return true, nil
	}

	var remaining []string
	err = repositories.GetDB().Model(&models.Comment{}).
		Where("receipt_id = ? AND id <> ?", receiptId, commentId).
		Pluck("comment", &remaining).Error
	if err != nil {
		return false, err
	}
	for _, comment := range remaining {
		if len(strings.TrimSpace(comment)) > 0 {
			return true, nil
		}
	}

	writeRequirementError(w, "comments", receiptCommentDeleteMessage)
	return false, nil
}

// enforceImageDeleteKeepsRequired refuses deleting a receipt's last image while
// the caller's role requires one in the receipt's group. Swapping the only image
// therefore means uploading the new one first. Writes the 400 itself and returns
// false when it refused.
func enforceImageDeleteKeepsRequired(w http.ResponseWriter, userId uint, receiptId uint, fileDataId uint) (bool, error) {
	requirements, err := resolveRequirementsForReceipt(userId, receiptId)
	if err != nil {
		return false, err
	}
	if !requirements.ImageRequired {
		return true, nil
	}

	var remaining int64
	err = repositories.GetDB().Model(&models.FileData{}).
		Where("receipt_id = ? AND id <> ?", receiptId, fileDataId).
		Count(&remaining).Error
	if err != nil {
		return false, err
	}
	if remaining > 0 {
		return true, nil
	}

	writeRequirementError(w, "files", receiptImageDeleteMessage)
	return false, nil
}

func resolveRequirementsForReceipt(userId uint, receiptId uint) (structs.ReceiptRequirements, error) {
	groupId, err := repositories.NewReceiptRepository(nil).GetReceiptGroupIdByReceiptId(utils.UintToString(receiptId))
	if err != nil {
		return structs.ReceiptRequirements{}, err
	}

	return services.NewReceiptService(nil).ResolveReceiptRequirements(userId, groupId)
}

func writeRequirementError(w http.ResponseWriter, key string, message string) {
	structs.WriteValidatorErrorResponse(w, structs.ValidatorError{Errors: map[string]string{key: message}}, http.StatusBadRequest)
}
