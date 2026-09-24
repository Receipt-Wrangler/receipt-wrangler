package services

import (
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
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
	customFieldRepository := repositories.NewCustomFieldRepository(service.TX)
	permissionService := NewPermissionService(service.TX)

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

	// The extra preloads are what receiptsource reads beyond the always-loaded
	// Categories/Tags.
	receipts, _, err := service.fetchReceipts(userId, groupId, filter, []string{"PaidByUser", "Group", "CustomFields"}, -1)
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

// Receipts fetches the same receipts Rows turns into report rows, for a caller
// that lists them rather than reporting on them (the Report Builder's drill-in):
// at most limit of them, newest first, plus the count of every receipt that
// matched. Sharing fetchReceipts is what keeps the list and the report's count in
// step. Where the two differ is presentation, and there it follows the receipts
// list: each custom field value carries its definition, categories/tags the caller
// may not see are stripped rather than marked (Restricted), and user references
// are masked for member visibility.
func (service ReportDataService) Receipts(
	userId uint,
	groupId string,
	filter commands.ReceiptPagedRequestFilter,
	limit int,
) ([]models.Receipt, int64, error) {
	permissionService := NewPermissionService(service.TX)

	receipts, count, err := service.fetchReceipts(userId, groupId, filter, constants.CUSTOM_FIELD_ASSOCIATIONS, limit)
	if err != nil {
		return nil, 0, err
	}
	if err := permissionService.FilterReceiptCategoriesTags(userId, receipts); err != nil {
		return nil, 0, err
	}
	if err := permissionService.MaskReceiptsForMemberVisibility(userId, receipts); err != nil {
		return nil, 0, err
	}
	return receipts, count, nil
}

// fetchReceipts loads a group's receipts matching the filter, newest first, with
// the two controls that decide which receipts a caller may see at all: the filter
// is narrowed to the caller's category/tag grants, so a restricted caller cannot
// probe for hidden ones through it, and paid-by visibility is enforced in the
// query, so a receipt the caller may not see is never fetched. A limit above zero
// loads only that many; the returned count is every receipt the query matched,
// taken after both controls, so it holds either way. A limit of -1 loads them all.
func (service ReportDataService) fetchReceipts(
	userId uint,
	groupId string,
	filter commands.ReceiptPagedRequestFilter,
	associations []string,
	limit int,
) ([]models.Receipt, int64, error) {
	receiptRepository := repositories.NewReceiptRepository(service.TX)
	permissionService := NewPermissionService(service.TX)

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return nil, 0, err
	}

	page := -1
	if limit > 0 {
		page = 1
	}
	pagedRequest := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          page,
			PageSize:      limit,
			OrderBy:       "date",
			SortDirection: commands.DESCENDING,
		},
		Filter: filter,
	}

	if err := permissionService.IntersectReceiptFilterWithGrants(userId, uintGroupId, &pagedRequest.Filter); err != nil {
		return nil, 0, err
	}

	return receiptRepository.GetPagedReceiptsByGroupId(
		userId,
		groupId,
		pagedRequest,
		associations,
		permissionService.PaidByListResolver(userId),
		nil,
	)
}
