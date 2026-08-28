package commands

import (
	"encoding/json"
	"net/http"
	"net/url"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
)

// Bounds on the configurable refresh-token lifetimes. They live here rather than
// alongside the resolver in services/ because internal/commands cannot import
// internal/services, and the validation and the read-side clamp must agree.
const (
	MinRefreshTokenValidForHours = 1
	MaxRefreshTokenValidForHours = 720 // 30 days
)

type UpsertSystemSettingsCommand struct {
	EnableLocalSignUp                   bool                                  `json:"enableLocalSignUp"`
	DebugOcr                            bool                                  `json:"debugOcr"`
	CurrencyDisplay                     string                                `json:"currencyDisplay"`
	CurrencyThousandthsSeparator        models.CurrencySeparator              `json:"currencyThousandthsSeparator"`
	CurrencyDecimalSeparator            models.CurrencySeparator              `json:"currencyDecimalSeparator"`
	CurrencySymbolPosition              models.CurrencySymbolPosition         `json:"currencySymbolPosition"`
	CurrencyHideDecimalPlaces           bool                                  `json:"currencyHideDecimalPlaces"`
	NumWorkers                          int                                   `json:"numWorkers"`
	EmailPollingInterval                int                                   `json:"emailPollingInterval"`
	ReceiptProcessingSettingsId         *uint                                 `json:"receiptProcessingSettingsId"`
	FallbackReceiptProcessingSettingsId *uint                                 `json:"fallbackReceiptProcessingSettingsId"`
	TaskConcurrency                     int                                   `json:"taskConcurrency"`
	PdfDpi                              int                                   `json:"pdfDpi"`
	TaskQueueConfigurations             []UpsertTaskQueueConfigurationCommand `json:"taskQueueConfigurations"`
	McpEnabled                          bool                                  `json:"mcpEnabled"`
	McpPublicUrl                        string                                `json:"mcpPublicUrl"`
	ShowLoginQr                         bool                                  `json:"showLoginQr"`
	MobileServerUrl                     string                                `json:"mobileServerUrl"`
	// Pointers so an omitted key is distinguishable from an explicit 0. The
	// repository writes every column (Select("*")), so a plain int would persist
	// as 0 and silently reset a configured lifetime to the default whenever a
	// client PUTs a body without these keys. Same reasoning as the pointer
	// fields on UpdateGroupReceiptSettingsCommand.
	RefreshTokenValidForHours    *int `json:"refreshTokenValidForHours"`
	McpRefreshTokenValidForHours *int `json:"mcpRefreshTokenValidForHours"`
}

func (command *UpsertSystemSettingsCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
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

func (command *UpsertSystemSettingsCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{}
	errorMap := make(map[string]string)
	vErr.Errors = errorMap

	if command.EmailPollingInterval < 0 {
		errorMap["emailPollingInterval"] = "Email polling interval must be greater than 0"
	}

	if command.ReceiptProcessingSettingsId != nil && *command.ReceiptProcessingSettingsId <= 0 {
		errorMap["receiptProcessingSettingsId"] = "Invalid receipt processing settings ID"
	}

	if command.FallbackReceiptProcessingSettingsId != nil && *command.FallbackReceiptProcessingSettingsId <= 0 {
		errorMap["fallbackReceiptProcessingSettingsId"] = "Invalid fallback receipt processing settings ID"
	}

	if command.ReceiptProcessingSettingsId == nil && command.FallbackReceiptProcessingSettingsId != nil {
		errorMap["fallbackReceiptProcessingSettingsId"] = "Fallback receipt processing settings ID cannot be set without receipt processing settings ID"
	}

	if command.ReceiptProcessingSettingsId != nil &&
		command.FallbackReceiptProcessingSettingsId != nil &&
		*command.ReceiptProcessingSettingsId ==
			*command.FallbackReceiptProcessingSettingsId {
		errorMap["fallbackReceiptProcessingSettingsId"] = "Fallback receipt processing settings ID cannot be the same as receipt processing settings ID"
	}

	if len(command.CurrencySymbolPosition) == 0 {
		errorMap["currencySymbolPosition"] = "Currency symbol position is required"
	}

	if len(command.CurrencyThousandthsSeparator) == 0 {
		errorMap["currencyThousandthsSeparator"] = "Currency thousandths separator is required"
	}

	if len(command.CurrencyDecimalSeparator) == 0 {
		errorMap["currencyDecimalSeparator"] = "Currency decimal separator is required"
	}

	if command.PdfDpi != 0 && (command.PdfDpi < 72 || command.PdfDpi > 1200) {
		errorMap["pdfDpi"] = "PDF DPI must be between 72 and 1200"
	}

	if command.TaskConcurrency < 0 {
		errorMap["taskConcurrency"] = "Task concurrency must be greater than or equal to 0"
	}

	queueNames := models.GetQueueNames()
	if len(command.TaskQueueConfigurations) != len(queueNames) {
		errorMap["taskQueueConfigurations"] = "Task queue configurations must be provided for all queues"
	}

	trimmedMcpPublicUrl := strings.TrimSpace(command.McpPublicUrl)
	if command.McpEnabled && len(trimmedMcpPublicUrl) == 0 {
		errorMap["mcpPublicUrl"] = "A public URL is required to enable the MCP server"
	} else if len(trimmedMcpPublicUrl) > 0 && !isValidAbsoluteUrl(trimmedMcpPublicUrl) {
		errorMap["mcpPublicUrl"] = "MCP public URL must be an absolute origin like https://receipts.example.com"
	}

	trimmedMobileServerUrl := strings.TrimSpace(command.MobileServerUrl)
	if command.ShowLoginQr && len(trimmedMobileServerUrl) == 0 {
		errorMap["mobileServerUrl"] = "A server URL is required to show the login QR code"
	} else if len(trimmedMobileServerUrl) > 0 && !isValidAbsoluteUrl(trimmedMobileServerUrl) {
		errorMap["mobileServerUrl"] = "Mobile server URL must be an absolute URL like https://receipts.example.com/api"
	}

	if msg := validateRefreshTokenValidForHours(command.RefreshTokenValidForHours); len(msg) > 0 {
		errorMap["refreshTokenValidForHours"] = msg
	}

	if msg := validateRefreshTokenValidForHours(command.McpRefreshTokenValidForHours); len(msg) > 0 {
		errorMap["mcpRefreshTokenValidForHours"] = msg
	}

	return vErr
}

// validateRefreshTokenValidForHours bounds a refresh-token lifetime, returning an
// empty string when the value is acceptable.
//
// A nil pointer means the key was omitted, which leaves the stored value alone
// (see ApplyOmittedLifetimes) and is always valid. An explicit 0 means "unset":
// the read side falls back to the built-in default. Shared by the app and MCP
// settings so the two cannot drift.
func validateRefreshTokenValidForHours(hours *int) string {
	if hours == nil || *hours == 0 {
		return ""
	}

	if *hours < MinRefreshTokenValidForHours || *hours > MaxRefreshTokenValidForHours {
		return "Refresh token lifetime must be between 1 and 720 hours (30 days)"
	}

	return ""
}

// isValidAbsoluteUrl reports whether the value is an absolute http(s) URL.
// A scheme and host are required; paths/queries/fragments are tolerated.
//
// Embedded credentials (https://user:token@host) are rejected: both settings
// this guards are published verbatim to unauthenticated clients — the mobile
// server URL is encoded into the login QR served by the public /featureConfig,
// and the MCP public URL is echoed in the OAuth discovery metadata.
func isValidAbsoluteUrl(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || parsed.User != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func (command *UpsertSystemSettingsCommand) ToSystemSettings(id uint) (models.SystemSettings, error) {
	var systemSettings models.SystemSettings

	bytes, err := json.Marshal(command)
	if err != nil {
		return systemSettings, err
	}

	err = json.Unmarshal(bytes, &systemSettings)
	if err != nil {
		return systemSettings, err
	}

	systemSettings.ID = id
	for _, config := range systemSettings.TaskQueueConfigurations {
		config.SystemSettingsId = id
	}

	return systemSettings, nil
}

// ApplyOmittedLifetimes carries the stored refresh-token lifetimes onto the
// settings a PUT is about to write for any key the request omitted.
//
// ToSystemSettings round-trips the command through JSON, so a nil pointer lands
// as 0 on the model; the repository then writes every column with Select("*").
// Without this, a client that PUTs a body lacking these keys would silently
// reset an admin's configured session length to the default.
func (command *UpsertSystemSettingsCommand) ApplyOmittedLifetimes(existing models.SystemSettings, updated *models.SystemSettings) {
	if command.RefreshTokenValidForHours == nil {
		updated.RefreshTokenValidForHours = existing.RefreshTokenValidForHours
	}

	if command.McpRefreshTokenValidForHours == nil {
		updated.McpRefreshTokenValidForHours = existing.McpRefreshTokenValidForHours
	}
}
