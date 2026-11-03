package structs

import "time"

// ReceiptHistory is a receipt's change history, newest entry first, as
// GET /receipt/{id}/history returns it. The server computes every change and
// redacts it to the caller; a client only formats.
//
// Every code below (action, actor type, field, kind, value type, before state)
// is a plain string on the contract, never an enum: a generated Dart client
// fails a whole payload on an enum value it does not know, so a code added
// later must not break an already-released build. Pinned by
// receipt_history_contract_test.go. Build one with NewReceiptHistory: entries
// and changes serialize as [], never null.
type ReceiptHistory struct {
	Entries []ReceiptHistoryEntry `json:"entries"`
	// ChangeCount counts the change rows of the successful entries other than
	// the creation.
	ChangeCount int `json:"changeCount"`
	// SaveCount counts the successful entries, the creation included.
	SaveCount int `json:"saveCount"`
	// ContributorCount counts the distinct actors across every entry.
	ContributorCount int `json:"contributorCount"`
}

func NewReceiptHistory() ReceiptHistory {
	return ReceiptHistory{Entries: []ReceiptHistoryEntry{}}
}

// ReceiptHistoryEntry is one recorded event: the receipt's creation, a change
// or a failed save.
type ReceiptHistoryEntry struct {
	SystemTaskId uint                `json:"systemTaskId"`
	Action       string              `json:"action"`
	RecordedAt   time.Time           `json:"recordedAt"`
	Actor        ReceiptHistoryActor `json:"actor"`
	HasSnapshots bool                `json:"hasSnapshots"`
	ErrorMessage string              `json:"errorMessage,omitempty"`
	// BeforeState says how far an update's "before" can be trusted; see the
	// ReceiptHistoryBeforeState* codes.
	BeforeState            string                 `json:"beforeState,omitempty"`
	BeforeSourceType       string                 `json:"beforeSourceType,omitempty"`
	BeforeSourceRecordedAt *time.Time             `json:"beforeSourceRecordedAt,omitempty"`
	Changes                []ReceiptHistoryChange `json:"changes"`
}

// ReceiptHistoryActor is who or what made an entry. UserId is set for USER,
// and for QUICK_SCAN when the scanning user is visible to the caller.
type ReceiptHistoryActor struct {
	Type   string `json:"type"`
	UserId *uint  `json:"userId,omitempty"`
}

// ReceiptHistoryChange is one field-level change. A CHANGED row carries From
// and To, an ADDED row only To, a REMOVED row only From.
type ReceiptHistoryChange struct {
	Field string `json:"field"`
	Kind  string `json:"kind"`
	// Label is the item's name (ITEM, SHARE, ITEM_CATEGORIES, ITEM_TAGS) or the
	// custom field's name (CUSTOM_FIELD).
	Label string `json:"label,omitempty"`
	// UserId is the share's charged-to user (SHARE) or the comment's author
	// (COMMENT).
	UserId *uint                `json:"userId,omitempty"`
	From   *ReceiptHistoryValue `json:"from,omitempty"`
	To     *ReceiptHistoryValue `json:"to,omitempty"`
}

// ReceiptHistoryValue is a typed value for the client to format; see the
// ReceiptHistoryValue* codes for what Value holds for each type.
type ReceiptHistoryValue struct {
	Type   string `json:"type"`
	Value  string `json:"value,omitempty"`
	UserId *uint  `json:"userId,omitempty"`
}

// ReceiptHistorySnapshots is the receipt before and after one entry, as JSON
// text, redacted exactly as the history is. A creation has no Before.
type ReceiptHistorySnapshots struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

// What an entry did.
const (
	ReceiptHistoryActionCreated               = "CREATED"
	ReceiptHistoryActionUpdated               = "UPDATED"
	ReceiptHistoryActionStatusChanged         = "STATUS_CHANGED"
	ReceiptHistoryActionResolved              = "RESOLVED"
	ReceiptHistoryActionImagesCommentsUpdated = "IMAGES_COMMENTS_UPDATED"
	ReceiptHistoryActionFailed                = "FAILED"
)

// Who or what made an entry.
const (
	ReceiptHistoryActorUser      = "USER"
	ReceiptHistoryActorQuickScan = "QUICK_SCAN"
	ReceiptHistoryActorEmail     = "EMAIL"
	ReceiptHistoryActorSystem    = "SYSTEM"
)

// How far an update's "before" can be trusted. A version 2 row's own "before"
// is COMPLETE. A version 1 row's was loaded one level deep, so it is compared
// against the nearest earlier complete copy instead (REBUILT), or, when there
// is none, against its own (INCOMPLETE).
const (
	ReceiptHistoryBeforeStateComplete   = "COMPLETE"
	ReceiptHistoryBeforeStateRebuilt    = "REBUILT"
	ReceiptHistoryBeforeStateIncomplete = "INCOMPLETE"
)

// The kinds of change.
const (
	ReceiptHistoryKindChanged = "CHANGED"
	ReceiptHistoryKindAdded   = "ADDED"
	ReceiptHistoryKindRemoved = "REMOVED"
)

// The fields a change can be about, in the order an entry lists them (the
// email fields come first, and only on the creation).
const (
	ReceiptHistoryFieldEmailInbox     = "EMAIL_INBOX"
	ReceiptHistoryFieldEmailSender    = "EMAIL_SENDER"
	ReceiptHistoryFieldName           = "NAME"
	ReceiptHistoryFieldAmount         = "AMOUNT"
	ReceiptHistoryFieldDate           = "DATE"
	ReceiptHistoryFieldGroup          = "GROUP"
	ReceiptHistoryFieldPaidBy         = "PAID_BY"
	ReceiptHistoryFieldStatus         = "STATUS"
	ReceiptHistoryFieldResolvedDate   = "RESOLVED_DATE"
	ReceiptHistoryFieldCategories     = "CATEGORIES"
	ReceiptHistoryFieldTags           = "TAGS"
	ReceiptHistoryFieldCustomField    = "CUSTOM_FIELD"
	ReceiptHistoryFieldItem           = "ITEM"
	ReceiptHistoryFieldShare          = "SHARE"
	ReceiptHistoryFieldItemCategories = "ITEM_CATEGORIES"
	ReceiptHistoryFieldItemTags       = "ITEM_TAGS"
	ReceiptHistoryFieldComment        = "COMMENT"
	ReceiptHistoryFieldImage          = "IMAGE"
)

// How a value is formatted: TEXT as is, CURRENCY a decimal string, DATE a
// calendar day (yyyy-MM-dd), INSTANT an RFC 3339 instant, USER the UserId,
// STATUS a status key, BOOLEAN "true" or "false".
const (
	ReceiptHistoryValueText     = "TEXT"
	ReceiptHistoryValueCurrency = "CURRENCY"
	ReceiptHistoryValueDate     = "DATE"
	ReceiptHistoryValueInstant  = "INSTANT"
	ReceiptHistoryValueUser     = "USER"
	ReceiptHistoryValueStatus   = "STATUS"
	ReceiptHistoryValueBoolean  = "BOOLEAN"
)
