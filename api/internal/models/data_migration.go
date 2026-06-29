package models

import "time"

// DataMigration records a one-time data migration that has been applied, so it
// is not re-run on subsequent startups. It is distinct from the GORM schema
// AutoMigrate handled by repositories.MakeMigrations: this ledger tracks
// one-off transformations of existing rows.
//
// It deliberately does not embed BaseModel — a soft-deleted row would be hidden
// by GORM's default scope and cause the migration to re-run.
type DataMigration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	AppliedAt time.Time `gorm:"not null" json:"appliedAt"`
}
