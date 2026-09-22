package models

type GroupReceiptSettings struct {
	BaseModel
	GroupId               uint `gorm:"not null;unique" json:"groupId"`
	HideImages            bool `gorm:"not null;default:false" json:"hideImages"`
	HideReceiptCategories bool `gorm:"not null;default:false" json:"hideReceiptCategories"`
	HideReceiptTags       bool `gorm:"not null;default:false" json:"hideReceiptTags"`
	HideItemCategories    bool `gorm:"not null;default:false" json:"hideItemCategories"`
	HideItemTags          bool `gorm:"not null;default:false" json:"hideItemTags"`
	HideShareCategories   bool `gorm:"not null;default:false" json:"hideShareCategories"`
	HideShareTags         bool `gorm:"not null;default:false" json:"hideShareTags"`
	HideComments          bool `gorm:"not null;default:false" json:"hideComments"`

	// Quick scan field configuration. Controls which fields appear in the quick-scan dialog and
	// which the user must provide. Paid-by/status always resolve to a value (a configured default
	// backfills them when not shown+required), so receipts are never left without one.
	QuickScanPaidByEnabled     bool                       `gorm:"not null;default:true" json:"quickScanPaidByEnabled"`
	QuickScanPaidByRequired    bool                       `gorm:"not null;default:true" json:"quickScanPaidByRequired"`
	QuickScanDefaultPaidByType QuickScanDefaultPaidByType `json:"quickScanDefaultPaidByType"`
	QuickScanDefaultPaidById   *uint                      `json:"quickScanDefaultPaidById"`

	QuickScanStatusEnabled  bool          `gorm:"not null;default:true" json:"quickScanStatusEnabled"`
	QuickScanStatusRequired bool          `gorm:"not null;default:true" json:"quickScanStatusRequired"`
	QuickScanDefaultStatus  ReceiptStatus `json:"quickScanDefaultStatus"`

	QuickScanCategoriesEnabled  bool `gorm:"not null;default:false" json:"quickScanCategoriesEnabled"`
	QuickScanCategoriesRequired bool `gorm:"not null;default:false" json:"quickScanCategoriesRequired"`

	QuickScanTagsEnabled  bool `gorm:"not null;default:false" json:"quickScanTagsEnabled"`
	QuickScanTagsRequired bool `gorm:"not null;default:false" json:"quickScanTagsRequired"`

	QuickScanCommentEnabled  bool `gorm:"not null;default:false" json:"quickScanCommentEnabled"`
	QuickScanCommentRequired bool `gorm:"not null;default:false" json:"quickScanCommentRequired"`

	// Default custom fields. ApplyDefaultCustomFieldsOnIngest extends the group's default set to
	// receipts the SERVER creates (quick scan, email integration); it is off by default so existing
	// installs are unchanged until an admin opts in.
	//
	// DefaultCustomFieldIds is transient (`gorm:"-"`) and lives in GroupReceiptSettingsCustomField
	// rows. It is filled explicitly by GroupReceiptSettingsRepository.LoadDefaultCustomFieldIds at
	// the serialization boundaries — never by a GORM hook — so nothing loads it implicitly. It must
	// always serialize as `[]` when empty, never `null`: the generated Dart client has no null guard
	// and a null would fail the whole AppData payload on already-released Android builds.
	ApplyDefaultCustomFieldsOnIngest bool   `gorm:"not null;default:false" json:"applyDefaultCustomFieldsOnIngest"`
	DefaultCustomFieldIds            []uint `gorm:"-" json:"defaultCustomFieldIds"`

	// Receipt summary. A block of totals rendered under the receipts table, aggregated over the
	// WHOLE current filter result set rather than the visible page: a receipt count and amount total
	// overall, then the same figures per configured status. ReceiptSummaryEnabled is the master
	// switch and is off by default, so existing installs are unchanged until an admin opts in.
	//
	// Both sets are transient (`gorm:"-"`), stored in GroupReceiptSettingsSummaryCustomField and
	// GroupReceiptSettingsSummaryStatus rows, and carry the same rules as DefaultCustomFieldIds
	// above: loaded explicitly by the repository at the serialization boundaries, and always
	// serialized as `[]` when empty rather than `null`.
	//
	// ReceiptSummaryCustomFieldIds holds CURRENCY custom fields only — those are the only ones with
	// a value that can be summed. The handler rejects any other type.
	//
	// ReceiptSummaryStatuses is the set of statuses to break out. A configured status that matches
	// no receipt still renders as a zero row, so the block keeps its shape as the filter narrows.
	//
	// ReceiptSummaryPosition is where the block renders relative to the receipts list. A real
	// column with a DB default rather than a `gorm:"-"` projection, because it is a scalar and
	// AutoMigrate backfills existing rows with the default on all three engines. BOTTOM is where
	// the summary rendered before the setting existed, so an install that never touches it is
	// unchanged. Read it through OrDefault(): a settings row that predates the column reads "",
	// and an empty enum on the wire fails a closed Dart EnumClass -- and with it the whole payload.
	ReceiptSummaryEnabled        bool                   `gorm:"not null;default:false" json:"receiptSummaryEnabled"`
	ReceiptSummaryPosition       ReceiptSummaryPosition `gorm:"default:BOTTOM" json:"receiptSummaryPosition"`
	ReceiptSummaryCustomFieldIds []uint                 `gorm:"-" json:"receiptSummaryCustomFieldIds"`
	ReceiptSummaryStatuses       []ReceiptStatus        `gorm:"-" json:"receiptSummaryStatuses"`
}

// IsQuickScanCommentShown reports whether the quick-scan comment field should be shown. HideComments
// hides comments for the whole group, so it overrides the quick-scan toggle — without mutating it,
// which is what lets the configured value come back when HideComments is turned off again. Callers
// must additionally check the user's group.comments.create permission (see resolveQuickScanFields).
func (settings GroupReceiptSettings) IsQuickScanCommentShown() bool {
	return settings.QuickScanCommentEnabled && !settings.HideComments
}

// IsQuickScanCommentRequired reports whether a comment must be supplied. A hidden field is never
// required.
func (settings GroupReceiptSettings) IsQuickScanCommentRequired() bool {
	return settings.IsQuickScanCommentShown() && settings.QuickScanCommentRequired
}
