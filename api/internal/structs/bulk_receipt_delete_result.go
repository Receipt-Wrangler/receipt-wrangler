package structs

// BulkReceiptDeleteResult reports, per requested receipt, whether a bulk delete
// removed it. Build it with NewBulkReceiptDeleteResult: both lists must serialize
// as [], never null, because the generated Dart client has no null guard.
type BulkReceiptDeleteResult struct {
	DeletedReceiptIds []uint                     `json:"deletedReceiptIds"`
	Failures          []BulkReceiptDeleteFailure `json:"failures"`
}

type BulkReceiptDeleteFailure struct {
	ReceiptId uint   `json:"receiptId"`
	Message   string `json:"message"`
}

func NewBulkReceiptDeleteResult() BulkReceiptDeleteResult {
	return BulkReceiptDeleteResult{
		DeletedReceiptIds: []uint{},
		Failures:          []BulkReceiptDeleteFailure{},
	}
}
