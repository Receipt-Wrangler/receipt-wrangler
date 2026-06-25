package models

type SystemSettings struct {
	BaseModel
	EnableLocalSignUp                   bool                      `json:"enableLocalSignUp" gorm:"default:false"`
	DebugOcr                            bool                      `json:"debugOcr" gorm:"default:false"`
	NumWorkers                          int                       `json:"numWorkers"`
	EmailPollingInterval                int                       `json:"emailPollingInterval" gorm:"default:1800"`
	CurrencyDisplay                     string                    `json:"currencyDisplay" gorm:"default:$"`
	CurrencyThousandthsSeparator        CurrencySeparator         `json:"currencyThousandthsSeparator" gorm:"default:,"`
	CurrencyDecimalSeparator            CurrencySeparator         `json:"currencyDecimalSeparator" gorm:"default:."`
	CurrencySymbolPosition              CurrencySymbolPosition    `json:"currencySymbolPosition" gorm:"default:START"`
	CurrencyHideDecimalPlaces           bool                      `json:"currencyHideDecimalPlaces" gorm:"default:false"`
	ReceiptProcessingSettings           ReceiptProcessingSettings `json:"-"`
	ReceiptProcessingSettingsId         *uint                     `json:"receiptProcessingSettingsId"`
	FallbackReceiptProcessingSettings   ReceiptProcessingSettings `json:"-"`
	FallbackReceiptProcessingSettingsId *uint                     `json:"fallbackReceiptProcessingSettingsId"`
	TaskConcurrency                     int                       `json:"taskConcurrency" gorm:"default:10"`
	PdfDpi                              int                       `json:"pdfDpi" gorm:"default:300"`
	TaskQueueConfigurations             []TaskQueueConfiguration  `json:"taskQueueConfigurations"`
	// McpEnabled toggles the OAuth 2.1-protected MCP server live, without a
	// restart. The routes are always mounted; this gates them at request time.
	McpEnabled bool `json:"mcpEnabled" gorm:"default:false"`
	// McpPublicUrl is the externally reachable origin (scheme + host) used to
	// build the OAuth issuer/metadata/redirect URLs and the MCP token audience.
	McpPublicUrl string `json:"mcpPublicUrl"`
}
