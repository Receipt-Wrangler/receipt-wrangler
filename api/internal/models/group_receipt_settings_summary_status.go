package models

// GroupReceiptSettingsSummaryStatus is one receipt status a group has declared as a
// SUMMARY breakdown row: the figures under the receipts table are repeated for the
// filtered receipts carrying this status.
//
// A configured status that matches nothing still renders, as a zero row, so the
// block keeps its shape as the user narrows the filter. That is why this is stored
// configuration rather than derived from the data.
//
// Keyed on GroupId for the same reason as the sibling join tables (a lazily-created
// settings row still has ID 0 on the call that created it). Status is part of the
// composite key, so the same status cannot be configured twice for a group;
// replaceGroupSummaryStatuses dedupes the submitted list before insert rather than
// letting that surface as a constraint violation.
//
// `size:32` is load-bearing, not decoration. A Go string maps to an unbounded TEXT by
// default, and MySQL/MariaDB reject an unbounded column in a key — so without it this
// table creates fine on SQLite and fails AutoMigrate on MySQL, which is exactly the
// class of bug that never shows up in local development. 32 clears the longest status
// (NEEDS_ATTENTION, 15) with room to spare. GroupRolePermission bounds its keyed
// permission string at 128 for the same reason.
type GroupReceiptSettingsSummaryStatus struct {
	GroupId uint          `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	Status  ReceiptStatus `gorm:"primaryKey;autoIncrement:false;size:32" json:"status"`
}
