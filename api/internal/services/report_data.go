package services

import (
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/reporting/receiptsource"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
)

// ReportDataService resolves a group's receipts into the reporting engine's two
// inputs — a field catalog and the rows — with the reporting access controls
// applied. It is the only place that reads the database for a report; the engine
// itself fetches nothing and enforces nothing.
//
// It stops at the engine's inputs: it does not build a ReportSpec or call
// reporting.Run. A caller assembles the spec and runs the report.
type ReportDataService struct {
	BaseService
}

func NewReportDataService(tx *gorm.DB) ReportDataService {
	service := ReportDataService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

// Rows fetches a group's receipts, applies the three reporting access controls,
// and maps the survivors to engine rows. The controls run in the order they must:
//
//   - the request filter is narrowed to what the caller may see, so a restricted
//     caller cannot probe for hidden categories/tags through the filter;
//   - paid-by visibility is enforced in the query, so a receipt the caller may
//     not see is never fetched (whole-receipt hiding);
//   - categories/tags the caller may not see are replaced with a (Restricted)
//     marker, so a hidden category still counts toward the totals in its own
//     bucket rather than vanishing.
//
// The returned catalog carries every built-in field plus one per custom field.
func (service ReportDataService) Rows(userId uint, groupId string, filter commands.ReceiptPagedRequestFilter) (reporting.FieldCatalog, []reporting.Row, error) {
	receiptRepository := repositories.NewReceiptRepository(service.TX)
	customFieldRepository := repositories.NewCustomFieldRepository(service.TX)
	permissionService := NewPermissionService(service.TX)

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return reporting.FieldCatalog{}, nil, err
	}

	// The catalog spans every custom field (a global pool), with their options
	// loaded so a select value resolves to its text rather than a bare option id.
	customFields, _, err := customFieldRepository.GetPagedCustomFields(commands.PagedRequestCommand{
		Page:          -1,
		PageSize:      -1,
		OrderBy:       "name",
		SortDirection: commands.ASCENDING,
	})
	if err != nil {
		return reporting.FieldCatalog{}, nil, err
	}

	source, err := receiptsource.New(customFields)
	if err != nil {
		return reporting.FieldCatalog{}, nil, err
	}

	pagedRequest := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          -1,
			PageSize:      -1,
			OrderBy:       "date",
			SortDirection: commands.DESCENDING,
		},
		Filter: filter,
	}

	// Narrow any category/tag filter to what the caller may see (anti-probing).
	if err := permissionService.IntersectReceiptFilterWithGrants(userId, uintGroupId, &pagedRequest.Filter); err != nil {
		return reporting.FieldCatalog{}, nil, err
	}

	// The paid-by resolver hides whole receipts in the query; the extra preloads
	// are what receiptsource reads beyond the always-loaded Categories/Tags.
	receipts, _, err := receiptRepository.GetPagedReceiptsByGroupId(
		userId,
		groupId,
		pagedRequest,
		[]string{"PaidByUser", "Group", "CustomFields"},
		permissionService.PaidByListResolver(userId),
	)
	if err != nil {
		return reporting.FieldCatalog{}, nil, err
	}

	// Replace categories/tags the caller cannot see with a (Restricted) marker so
	// they aggregate into their own bucket rather than disappearing.
	if err := permissionService.SubstituteRestrictedCategoriesTags(userId, receipts); err != nil {
		return reporting.FieldCatalog{}, nil, err
	}

	return source.Catalog(), source.Rows(receipts), nil
}
