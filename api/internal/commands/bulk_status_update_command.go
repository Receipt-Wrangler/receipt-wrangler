package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
)

// MaxBulkStatusUpdateIds bounds one request, for the reason MaxBulkReceiptDeleteIds
// does: the ids reach `id IN ?` queries, and each one now also costs two full
// receipt loads and a history row. The desktop sends at most one table page (100).
const MaxBulkStatusUpdateIds = MaxBulkReceiptDeleteIds

type BulkStatusUpdateCommand struct {
	Comment    string
	Status     models.ReceiptStatus
	ReceiptIds []uint
}

func (command *BulkStatusUpdateCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	return nil
}
