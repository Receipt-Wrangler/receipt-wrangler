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
