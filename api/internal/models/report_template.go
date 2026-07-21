package models

import "encoding/json"

// ReportTemplate is a saved report configuration. Configuration holds the whole
// commands.ReportRequestCommand the builder submits, stored verbatim as a JSON
// blob so a template round-trips back into the builder unchanged. Name mirrors the
// report's own name; it is not unique (two templates may share a name).
//
// ConfigurationVersion stamps the schema the Configuration blob was written under
// (currently 1). It lets a future breaking change to the config shape upcast — or
// fail loud on — old blobs instead of silently misdeserializing them.
type ReportTemplate struct {
	BaseModel
	Name                 string          `gorm:"not null" json:"name"`
	Configuration        json.RawMessage `json:"configuration"`
	ConfigurationVersion int             `gorm:"not null;default:1" json:"configurationVersion"`

	// AllowedActions is not persisted (gorm:"-"): the list handler fills it per row
	// with the actions the requesting user may perform on this template
	// (read/generate/update/delete/duplicate), already factoring in the "*All"
	// permissions, the group-access ceiling, and the per-template grant matrix. It
	// lets the client gate each row's action buttons without re-deriving access.
	AllowedActions []string `gorm:"-" json:"allowedActions,omitempty"`
}
