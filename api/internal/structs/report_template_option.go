package structs

// ReportTemplateOption is a lightweight report template for the role-form access
// matrix: just enough to list templates and show which groups each covers, without
// shipping the full configuration blob. GroupIds is never nil (empty slice instead).
type ReportTemplateOption struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	GroupIds []uint `json:"groupIds"`
}
