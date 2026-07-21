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
}
