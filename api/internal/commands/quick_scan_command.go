package commands

import (
	"mime/multipart"
	"net/http"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
)

type QuickScanCommand struct {
	Files         []multipart.File        `json:"file"`
	FileHeaders   []*multipart.FileHeader `json:"fileHeader"`
	PaidByUserIds []uint                  `json:"paidByUserId"`
	GroupIds      []uint                  `json:"groupId"`
	Statuses      []models.ReceiptStatus  `json:"status"`
	// CategoryIds and TagIds carry the per-file category/tag selections. Each outer element maps to
	// a file (by index); the multipart form sends one comma-joined id string per file (empty for
	// none), keeping the payload a flat string array the generated client can encode.
	CategoryIds [][]uint `json:"categoryIds"`
	TagIds      [][]uint `json:"tagIds"`
	// Comments carries the per-file comment, one entry per file (empty for none). Unlike CategoryIds
	// and TagIds this is free text and is deliberately NOT comma-split — a comment may contain commas
	// and newlines.
	Comments []string `json:"comments"`
}

func (command *QuickScanCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(constants.MultipartFormMaxSize)

	var form = r.Form

	var files = make([]multipart.File, 0)
	var fileHeaders = make([]*multipart.FileHeader, 0)
	var paidByUserIds = make([]uint, 0)
	var groupIds = make([]uint, 0)
	var statuses = make([]models.ReceiptStatus, 0)
	var categoryIds = make([][]uint, 0)
	var tagIds = make([][]uint, 0)
	var comments = make([]string, 0)

	var formPaidByUserIds = form["paidByUserIds"]
	var formGroupIds = form["groupIds"]
	var formStatuses = form["statuses"]
	var formCategoryIds = form["categoryIds"]
	var formTagIds = form["tagIds"]
	var formComments = form["comments"]

	if err != nil {
		return err
	}

	for _, fileHeader := range r.MultipartForm.File["files"] {
		fileHeaders = append(fileHeaders, fileHeader)
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		files = append(files, file)
	}

	for _, userId := range formPaidByUserIds {
		// An empty paid-by (optional/hidden field) is sent as a blank string; treat it as 0 (unset)
		// so the handler can apply the group's configured default.
		trimmedUserId := strings.TrimSpace(userId)
		if len(trimmedUserId) == 0 {
			paidByUserIds = append(paidByUserIds, 0)
			continue
		}

		formattedUserId, err := utils.StringToUint(trimmedUserId)
		if err != nil {
			return err
		}

		paidByUserIds = append(paidByUserIds, formattedUserId)
	}

	for _, groupId := range formGroupIds {
		formattedGroupId, err := utils.StringToUint(groupId)
		if err != nil {
			return err
		}

		groupIds = append(groupIds, formattedGroupId)
	}

	for _, status := range formStatuses {
		var formattedStatus models.ReceiptStatus
		err = formattedStatus.Scan(strings.TrimSpace(status))
		if err != nil {
			return err
		}

		statuses = append(statuses, formattedStatus)
	}

	for _, idList := range formCategoryIds {
		parsedIds, err := parseCommaSeparatedUints(idList)
		if err != nil {
			return err
		}

		categoryIds = append(categoryIds, parsedIds)
	}

	for _, idList := range formTagIds {
		parsedIds, err := parseCommaSeparatedUints(idList)
		if err != nil {
			return err
		}

		tagIds = append(tagIds, parsedIds)
	}

	for _, comment := range formComments {
		// Trim so a whitespace-only comment counts as empty: it then fails a required check instead
		// of persisting a blank comment.
		comments = append(comments, strings.TrimSpace(comment))
	}

	command.Files = files
	command.FileHeaders = fileHeaders
	command.PaidByUserIds = paidByUserIds
	command.GroupIds = groupIds
	command.Statuses = statuses
	command.CategoryIds = categoryIds
	command.TagIds = tagIds
	command.Comments = comments

	return nil
}

// parseCommaSeparatedUints converts a comma-joined id string (e.g. "3,7") into a []uint, skipping
// blank segments so an empty string yields an empty slice.
func parseCommaSeparatedUints(value string) ([]uint, error) {
	result := make([]uint, 0)
	for _, part := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(part)
		if len(trimmed) == 0 {
			continue
		}

		parsed, err := utils.StringToUint(trimmed)
		if err != nil {
			return nil, err
		}

		result = append(result, parsed)
	}

	return result, nil
}

func (command QuickScanCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{
		Errors: make(map[string]string),
	}

	var filesLength = len(command.Files)

	if filesLength == 0 {
		vErr.Errors["files"] = "At least one file is required."
	}

	if len(command.PaidByUserIds) != filesLength {
		vErr.Errors["paidByUserId"] = "Paid By User Ids must match the number of files."
	}

	if len(command.GroupIds) != filesLength {
		vErr.Errors["groupIds"] = "Group Ids must match the number of files."
	}

	if len(command.Statuses) != filesLength {
		vErr.Errors["statuses"] = "Statuses must match the number of files."
	}

	if len(command.PaidByUserIds) == 0 {
		vErr.Errors["paidByUserId"] = "Paid By User Id is required."
	}

	if len(command.GroupIds) == 0 {
		vErr.Errors["groupId"] = "Group Id is required."
	}

	if len(command.Statuses) == 0 {
		vErr.Errors["status"] = "Status is required."
	}

	// Category/tag/comment selections are optional (a client may omit them entirely — notably any
	// client released before the field existed), but when supplied they must carry one entry per file
	// so the handler can align them by index.
	if len(command.CategoryIds) > 0 && len(command.CategoryIds) != filesLength {
		vErr.Errors["categoryIds"] = "Category Ids must match the number of files."
	}

	if len(command.TagIds) > 0 && len(command.TagIds) != filesLength {
		vErr.Errors["tagIds"] = "Tag Ids must match the number of files."
	}

	if len(command.Comments) > 0 && len(command.Comments) != filesLength {
		vErr.Errors["comments"] = "Comments must match the number of files."
	}

	return vErr
}

// CategoryIdsForFile returns the category id selections for the file at index i, or an empty slice
// when the client omitted them (or supplied none for that file).
func (command QuickScanCommand) CategoryIdsForFile(i int) []uint {
	if i < len(command.CategoryIds) {
		return command.CategoryIds[i]
	}
	return []uint{}
}

// TagIdsForFile returns the tag id selections for the file at index i, or an empty slice when the
// client omitted them (or supplied none for that file).
func (command QuickScanCommand) TagIdsForFile(i int) []uint {
	if i < len(command.TagIds) {
		return command.TagIds[i]
	}
	return []uint{}
}

// CommentForFile returns the comment for the file at index i, or an empty string when the client
// omitted comments entirely (or supplied none for that file).
func (command QuickScanCommand) CommentForFile(i int) string {
	if i < len(command.Comments) {
		return command.Comments[i]
	}
	return ""
}

func (command *QuickScanCommand) LoadDataFromRequestAndValidate(w http.ResponseWriter, r *http.Request) (structs.ValidatorError, error) {
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		return structs.ValidatorError{}, err
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		return vErr, nil
	}

	return structs.ValidatorError{}, nil
}
