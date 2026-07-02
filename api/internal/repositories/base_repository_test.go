package repositories

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// TestGetCountBindsParameters is a regression test for GHSA-q6h3-4g3r-gg2x.
// GetCount forwards its args to GORM's Where so callers can bind values with a
// "?" placeholder instead of interpolating them into the clause. A SQL
// tautology passed as a bound value is therefore treated as a literal and
// matches nothing, rather than being parsed as SQL.
func TestGetCountBindsParameters(t *testing.T) {
	defer TruncateTestDb()
	CreateTestCategories() // seeds "test", "test2", "test3"

	repository := NewCategoryRepository(nil)

	// Bound value matches the literal row.
	legit, err := repository.GetCount("categories", "name = ?", "test")
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if legit != 1 {
		utils.PrintTestError(t, legit, 1)
	}

	// A tautology bound as a parameter is a literal name, not SQL, so it
	// matches nothing (pre-fix it interpolated to name = 'test' OR '1'='1'
	// and matched every row).
	injection, err := repository.GetCount("categories", "name = ?", "test' OR '1'='1")
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if injection != 0 {
		utils.PrintTestError(t, injection, 0)
	}

	// An empty clause counts every row (the system-email call site relies on
	// this and must keep working with the variadic signature).
	all, err := repository.GetCount("categories", "")
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if all != 3 {
		utils.PrintTestError(t, all, 3)
	}
}
