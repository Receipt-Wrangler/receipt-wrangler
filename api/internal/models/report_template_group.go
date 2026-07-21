package models

// ReportTemplateGroup is a queryable mirror of the group ids buried inside a
// ReportTemplate's Configuration JSON blob. It exists so the templates list can
// filter by "all of a template's groups are within the caller's readable set"
// without unmarshaling every config blob, and so the role-form matrix can list a
// template's groups cheaply.
//
// ReportTemplateID cascades with its template. GroupID is a plain indexed column
// with NO foreign key on purpose: the table is a denormalized copy of the blob,
// rebuilt on every template write, so keeping "table == blob" is the invariant. A
// group→CASCADE FK would let a group deletion silently drift these rows away from
// the blob (which still lists the group); instead a template covering a deleted
// group simply stays invisible, because its group-access ceiling can never pass.
type ReportTemplateGroup struct {
	ReportTemplateID uint           `gorm:"primaryKey;autoIncrement:false" json:"reportTemplateId"`
	GroupID          uint           `gorm:"primaryKey;autoIncrement:false;index" json:"groupId"`
	ReportTemplate   ReportTemplate `gorm:"foreignKey:ReportTemplateID;constraint:OnDelete:CASCADE" json:"-"`
}
