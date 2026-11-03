package commands

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"time"
)

type UpsertReceiptCommand struct {
	Name            string                          `json:"name"`
	Amount          decimal.Decimal                 `json:"amount"`
	Date            time.Time                       `json:"date"`
	GroupId         uint                            `json:"groupId"`
	PaidByUserID    uint                            `json:"paidByUserId"`
	Status          models.ReceiptStatus            `json:"status"`
	Categories      []UpsertCategoryCommand         `json:"categories"`
	Tags            []UpsertTagCommand              `json:"tags"`
	Items           []UpsertItemCommand             `json:"receiptItems"`
	Comments        []UpsertCommentCommand          `json:"comments"`
	CustomFields    []UpsertCustomFieldValueCommand `json:"customFields"`
	CreatedByString string                          `json:"createdByString"`
	// EmailInbox and EmailSender are server-set by the email task and never
	// read from a request (json:"-"); ToReceipt copies them explicitly.
	EmailInbox  string `json:"-"`
	EmailSender string `json:"-"`
}

func (receipt *UpsertReceiptCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &receipt)
	if err != nil {
		return err
	}

	return nil
}

func (receipt *UpsertReceiptCommand) Validate(tokenUserId uint, isCreate bool) structs.ValidatorError {
	errors := make(map[string]string)
	vErr := structs.ValidatorError{}

	if len(receipt.Name) == 0 {
		errors["name"] = "Name is required"
	}

	if receipt.Date.IsZero() {
		errors["date"] = "Date is required"
	}

	if receipt.GroupId == 0 {
		errors["groupId"] = "Group Id is required"
	}

	if receipt.PaidByUserID == 0 {
		errors["paidByUserId"] = "Paid By User Id is required"
	}

	if receipt.Status == "" {
		errors["status"] = "Status is required"
	}

	for i, category := range receipt.Categories {
		basePath := "categories." + fmt.Sprintf("%d", i)
		categoryErrors := category.Validate()
		for key, value := range categoryErrors.Errors {
			errors[basePath+"."+key] = value
		}
	}

	for i, tag := range receipt.Tags {
		basePath := "tags." + fmt.Sprintf("%d", i)
		tagErrors := tag.Validate()
		for key, value := range tagErrors.Errors {
			errors[basePath+"."+key] = value
		}
	}

	for i, item := range receipt.Items {
		basePath := "receiptItems." + fmt.Sprintf("%d", i)
		itemErrors := item.Validate(receipt.Amount, isCreate)
		for key, value := range itemErrors.Errors {
			errors[basePath+"."+key] = value
		}
	}

	for i, comment := range receipt.Comments {
		basePath := "comments." + fmt.Sprintf("%d", i)
		commentErrors := comment.Validate(tokenUserId, isCreate)
		for key, value := range commentErrors.Errors {
			errors[basePath+"."+key] = value
		}
	}

	vErr.Errors = errors
	return vErr
}

func (receipt *UpsertReceiptCommand) ToReceipt() (models.Receipt, error) {
	var result models.Receipt
	bytes, err := json.Marshal(receipt)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}

	// json:"-" keeps these off the round trip above, so they are copied here.
	result.EmailInbox = receipt.EmailInbox
	result.EmailSender = receipt.EmailSender

	return result, nil
}
