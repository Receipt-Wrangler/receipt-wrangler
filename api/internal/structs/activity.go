package structs

import (
	"receipt-wrangler/api/internal/models"
	"time"
)

type Activity struct {
	Id                     uint                    `json:"id"`
	Type                   models.SystemTaskType   `json:"type"`
	Status                 models.SystemTaskStatus `json:"status"`
	StartedAt              time.Time               `json:"startedAt"`
	EndedAt                *time.Time              `json:"endedAt"`
	RanByUserId            *uint                   `json:"ranByUserId"`
	ReceiptId              *uint                   `json:"receiptId"`
	GroupId                *uint                   `json:"groupId"`
	CanBeRestarted         bool                    `json:"canBeRestarted"`
	HasSourceFile          bool                    `json:"hasSourceFile"`
	AssociatedSystemTaskId *uint                   `json:"-"`
	// AsynqTaskId is selected so the flags above can be resolved without a
	// per-activity database round trip. It is an internal Redis key and is kept
	// off the wire.
	AsynqTaskId string `json:"-"`
}
