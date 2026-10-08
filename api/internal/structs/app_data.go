package structs

import "receipt-wrangler/api/internal/models"

type AppData struct {
	TokenPair
	About                        About                         `json:"about"`
	Claims                       Claims                        `json:"claims"`
	Groups                       []models.Group                `json:"groups"`
	CurrencyDisplay              string                        `json:"currencyDisplay"`
	CurrencyThousandthsSeparator models.CurrencySeparator      `json:"currencyThousandthsSeparator"`
	CurrencyDecimalSeparator     models.CurrencySeparator      `json:"currencyDecimalSeparator"`
	CurrencySymbolPosition       models.CurrencySymbolPosition `json:"currencySymbolPosition"`
	CurrencyHideDecimalPlaces    bool                          `json:"currencyHideDecimalPlaces"`
	// TimeZone is the resolved app time zone (an IANA name, "UTC" when unset or
	// invalid) clients work out "today" and display instants in.
	TimeZone         string                     `json:"timeZone"`
	Categories       []models.Category          `json:"categories"`
	Tags             []models.Tag               `json:"tags"`
	GroupCategories  map[uint][]models.Category `json:"groupCategories"`
	GroupTags        map[uint][]models.Tag      `json:"groupTags"`
	UserPreferences  models.UserPrefernces      `json:"userPreferences"`
	Users            []UserView                 `json:"users"`
	FeatureConfig    FeatureConfig              `json:"featureConfig"`
	Icons            []Icon                     `json:"icons"`
	AppPermissions   []string                   `json:"appPermissions"`
	GroupPermissions map[uint][]string          `json:"groupPermissions"`
	// GroupReceiptRequirements is what the caller must supply on each group's
	// receipts, resolved server-side from their group role. Only groups where
	// something is required are present (absent ⇒ nothing required); never null.
	GroupReceiptRequirements map[uint]ReceiptRequirements `json:"groupReceiptRequirements"`
}
