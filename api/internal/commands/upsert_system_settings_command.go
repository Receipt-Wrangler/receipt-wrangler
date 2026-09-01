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
	// Pointers so an omitted key is distinguishable from an explicit zero value.
	// The repository writes every column (Select("*")), so a plain field would
	// persist as its zero value and silently reset a configured setting whenever a
	// client PUTs a body without these keys. Same reasoning as the pointer fields
	// on UpdateGroupReceiptSettingsCommand.
	RefreshTokenValidForHours    *int `json:"refreshTokenValidForHours"`
	McpRefreshTokenValidForHours *int `json:"mcpRefreshTokenValidForHours"`
	// ServerPublicUrl is a pointer for the same reason, and the stakes are higher
	// than a reset default: it builds the OIDC redirect URI registered at each
	// identity provider, so clearing it makes every provider reject the callback
	// with a redirect_uri mismatch — a failure that surfaces at the IdP rather
	// than in our logs.
	ServerPublicUrl *string `json:"serverPublicUrl"`
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

	// A nil pointer means the key was omitted, which leaves the stored value alone
	// and is always valid. An explicit empty string clears it, which is allowed —
	// it is only required once an OIDC provider is enabled, and that cross-setting
	// rule is enforced where providers are saved (OidcProviderService), so an admin
	// cannot save a provider that has no redirect URI to register.
	if command.ServerPublicUrl != nil {
		trimmedServerPublicUrl := strings.TrimSpace(*command.ServerPublicUrl)
		if len(trimmedServerPublicUrl) > 0 && !isValidAbsoluteUrl(trimmedServerPublicUrl) {
			errorMap["serverPublicUrl"] = "Server public URL must be an absolute origin like https://receipts.example.com"
		}
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
// (see ApplyOmittedValues) and is always valid. An explicit 0 means "unset":
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

// OmittedColumns names the pointer-typed fields the request did not send, so the
// repository can leave those columns out of the UPDATE entirely.
//
// Skipping the column is what makes a concurrent update safe. Copying the stored
// value onto the row instead (see ApplyOmittedValues) would still write it, so
// two requests that each set one field and omit the other would clobber each
// other with the values they read before the write. A column that is never
// written cannot be clobbered, and unlike a row lock this works identically on
// SQLite, MySQL and Postgres.
func (command *UpsertSystemSettingsCommand) OmittedColumns() []string {
	columns := make([]string, 0, 3)

	if command.RefreshTokenValidForHours == nil {
		columns = append(columns, "RefreshTokenValidForHours")
	}

	if command.McpRefreshTokenValidForHours == nil {
		columns = append(columns, "McpRefreshTokenValidForHours")
	}

	if command.ServerPublicUrl == nil {
		columns = append(columns, "ServerPublicUrl")
	}

	return columns
}

// ApplyOmittedValues carries the stored value of every omitted pointer field onto
// the settings a PUT is about to write.
//
// ToSystemSettings round-trips the command through JSON, so a nil pointer lands
// as the zero value on the model. The columns themselves are excluded from the
// UPDATE by OmittedColumns, so this exists purely so the object echoed back in
// the response carries the stored value rather than a misleading zero.
func (command *UpsertSystemSettingsCommand) ApplyOmittedValues(existing models.SystemSettings, updated *models.SystemSettings) {
	if command.RefreshTokenValidForHours == nil {
		updated.RefreshTokenValidForHours = existing.RefreshTokenValidForHours
	}

	if command.McpRefreshTokenValidForHours == nil {
		updated.McpRefreshTokenValidForHours = existing.McpRefreshTokenValidForHours
	}

	if command.ServerPublicUrl == nil {
		updated.ServerPublicUrl = existing.ServerPublicUrl
	}
}
