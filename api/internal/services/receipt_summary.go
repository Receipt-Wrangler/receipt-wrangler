package services

import (
	"errors"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

// ErrConfigurationGroupForbidden is returned when the caller names a configuration
// group they may not read. The handler maps it to a 403; it is a sentinel rather than
// a bare error so a genuine failure is never mistaken for a denial.
var ErrConfigurationGroupForbidden = errors.New("no access to the requested configuration group")

// ErrConfigurationGroupNotAllGroup is returned when the caller names a configuration
// group other than the real group being viewed. Only the synthetic All group may borrow
// another group's configuration — it has none of its own. The handler maps this to a
// 400: it is a malformed request, not an access failure.
var ErrConfigurationGroupNotAllGroup = errors.New("a configuration group may only be named for the all group")

// ReceiptSummaryService aggregates a group's filtered receipts into the block of
// totals under the receipts table. It follows PieChartService: fetch the filtered set
// unpaged through the repository that already enforces the access controls, then fold
// in Go with shopspring/decimal.
//
// The fold is deliberately NOT a SQL SUM. Three reasons, strongest first:
//
//   - Exactness is not portable. amount is decimal(10,2), but SQLite has no decimal
//     type, so SUM(amount) comes back an IEEE double there and an exact decimal on
//     Postgres/MySQL. The same endpoint would report different cents on the three
//     supported engines, on exactly the money this feature exists to report.
//   - The filter, the "All" group's group_id IN (...) branch and the paid-by
//     visibility disjunction are reachable only through GetPagedReceiptsByGroupId. A
//     SUM query would re-derive all three, and the failure mode is a totals row that
//     silently disagrees with the rows above it — or worse, one that counts receipts
//     the viewer may not see.
//   - The per-field totals would need a dynamically generated
//     SUM(CASE WHEN custom_field_id = ? THEN currency_value END) per configured
//     field across three dialects, and would still need the lowest-id tie-break below.
//
// The cost is loading the filtered set. That is the same bargain PieChartService and
// ReportDataService already make, and it is bounded on three sides: a group that has
// not enabled the summary never reaches a query at all (below), the desktop does not
// re-request on paging or sorting, and CustomFields is preloaded only when a field is
// actually configured.
type ReceiptSummaryService struct {
	BaseService
}

func NewReceiptSummaryService(tx *gorm.DB) ReceiptSummaryService {
	service := ReceiptSummaryService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

func (service ReceiptSummaryService) GetReceiptSummary(
	userId uint,
	groupId string,
	command commands.ReceiptSummaryCommand,
) (structs.ReceiptSummary, error) {
	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return structs.ReceiptSummary{}, err
	}

	permissionService := NewPermissionService(service.TX)

	configurationGroupId, err := service.resolveConfigurationGroupId(userId, uintGroupId, command, permissionService)
	if err != nil {
		return structs.ReceiptSummary{}, err
	}

	settings, err := service.loadSettings(configurationGroupId)
	if err != nil {
		return structs.ReceiptSummary{}, err
	}

	// The off state costs one settings read and no receipt query. Every install that
	// has not opted in pays exactly this.
	if !settings.ReceiptSummaryEnabled {
		return emptyReceiptSummary(configurationGroupId), nil
	}

	fieldNames, fieldIds, err := service.resolveSummaryFields(settings)
	if err != nil {
		return structs.ReceiptSummary{}, err
	}

	receipts, err := service.fetchReceipts(userId, groupId, uintGroupId, command, fieldIds, permissionService)
	if err != nil {
		return structs.ReceiptSummary{}, err
	}

	return service.fold(receipts, settings.ReceiptSummaryStatuses, fieldIds, fieldNames, configurationGroupId), nil
}

// resolveConfigurationGroupId picks the group whose settings shape the breakdown.
//
// Borrowing another group's configuration is a privilege of the synthetic All group
// alone, which spans several groups and has no settings row of its own. A real group
// must use its own: otherwise a member could render group A's receipts under group B's
// statuses and currency fields, overriding what A's admin configured — and opting into
// a summary A has switched off. That would undo the invariant the feature is built on,
// that the server owns the configuration and the client cannot add a column or drop one.
//
// The All-group test comes BEFORE the permission check on purpose. Rejecting on the
// shape of the request first means the answer cannot depend on whether the caller can
// read the named group, so this endpoint can never be used to probe for another group's
// existence.
func (service ReceiptSummaryService) resolveConfigurationGroupId(
	userId uint,
	uintGroupId uint,
	command commands.ReceiptSummaryCommand,
	permissionService PermissionService,
) (uint, error) {
	if command.ConfigurationGroupId == nil || *command.ConfigurationGroupId == uintGroupId {
		return uintGroupId, nil
	}

	isAllGroup, err := repositories.NewGroupRepository(service.TX).IsAllGroup(uintGroupId)
	if err != nil {
		return 0, err
	}
	if !isAllGroup {
		return 0, ErrConfigurationGroupNotAllGroup
	}

	configurationGroupId := *command.ConfigurationGroupId
	canRead, err := permissionService.HasGroupPermissions(userId, configurationGroupId, permissions.GroupReceiptsRead)
	if err != nil {
		return 0, err
	}
	if !canRead {
		return 0, ErrConfigurationGroupForbidden
	}

	return configurationGroupId, nil
}

// loadSettings reads the configuration group's receipt settings. A missing row is not
// an error: settings rows are created lazily, so a group nobody has opened the settings
// page for legitimately has none — which reads as "summary off" (the same treatment
// ApplyGroupDefaultCustomFields gives it).
func (service ReceiptSummaryService) loadSettings(configurationGroupId uint) (models.GroupReceiptSettings, error) {
	settings, err := repositories.NewGroupReceiptSettingsRepository(service.TX).
		GetGroupReceiptSettingsByGroupId(configurationGroupId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.GroupReceiptSettings{}, nil
	}
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	return settings, nil
}

// resolveSummaryFields turns the configured ids into (name lookup, surviving id order).
// An id with no catalog row is dropped rather than rendered nameless — the delete
// cascade should make that impossible, but a column labelled "" would be worse than a
// column that quietly disappears.
//
// Names are read here rather than through a CustomFields.CustomField preload on every
// receipt: the preload would pull a catalog row per receipt for a handful of names.
func (service ReceiptSummaryService) resolveSummaryFields(
	settings models.GroupReceiptSettings,
) (map[uint]string, []uint, error) {
	if len(settings.ReceiptSummaryCustomFieldIds) == 0 {
		return map[uint]string{}, []uint{}, nil
	}

	customFields, err := repositories.NewCustomFieldRepository(service.TX).
		GetCustomFieldsByIds(settings.ReceiptSummaryCustomFieldIds)
	if err != nil {
		return nil, nil, err
	}

	names := make(map[uint]string, len(customFields))
	for _, customField := range customFields {
		names[customField.ID] = customField.Name
	}

	// Preserve the configured order so the columns do not reshuffle between requests.
	fieldIds := make([]uint, 0, len(settings.ReceiptSummaryCustomFieldIds))
	for _, customFieldId := range settings.ReceiptSummaryCustomFieldIds {
		if _, ok := names[customFieldId]; ok {
			fieldIds = append(fieldIds, customFieldId)
		}
	}

	return names, fieldIds, nil
}

// fetchReceipts pulls the whole filtered set. Going through GetPagedReceiptsByGroupId
// is what guarantees the summary describes exactly the rows the table shows: the same
// filter builder, the same "All" group fan-out, and the same paid-by visibility.
func (service ReceiptSummaryService) fetchReceipts(
	userId uint,
	groupId string,
	uintGroupId uint,
	command commands.ReceiptSummaryCommand,
	fieldIds []uint,
	permissionService PermissionService,
) ([]models.Receipt, error) {
	pagedRequest := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          -1,
			PageSize:      -1,
			OrderBy:       "date",
			SortDirection: commands.DESCENDING,
		},
		Filter: command.Filter,
	}

	// Narrow any category/tag filter to what the caller may see (anti-probing). Against
	// the group being VIEWED, not the configuration group: this is about the data scope.
	err := permissionService.IntersectReceiptFilterWithGrants(userId, uintGroupId, &pagedRequest.Filter)
	if err != nil {
		return nil, err
	}

	// Only pay for the custom field values when a field is actually configured.
	associations := []string{}
	if len(fieldIds) > 0 {
		associations = append(associations, "CustomFields")
	}

	receipts, _, err := repositories.NewReceiptRepository(service.TX).GetPagedReceiptsByGroupId(
		userId,
		groupId,
		pagedRequest,
		associations,
		permissionService.PaidByListResolver(userId),
	)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

// summaryAccumulator is one row under construction.
type summaryAccumulator struct {
	receiptCount int64
	total        decimal.Decimal
	fieldTotals  map[uint]decimal.Decimal
}

func newSummaryAccumulator(fieldIds []uint) *summaryAccumulator {
	fieldTotals := make(map[uint]decimal.Decimal, len(fieldIds))
	for _, fieldId := range fieldIds {
		fieldTotals[fieldId] = decimal.Zero
	}

	return &summaryAccumulator{total: decimal.Zero, fieldTotals: fieldTotals}
}

func (accumulator *summaryAccumulator) add(receipt models.Receipt, fieldValues map[uint]decimal.Decimal) {
	accumulator.receiptCount++
	accumulator.total = accumulator.total.Add(receipt.Amount)

	for fieldId, value := range fieldValues {
		if _, configured := accumulator.fieldTotals[fieldId]; !configured {
			continue
		}
		accumulator.fieldTotals[fieldId] = accumulator.fieldTotals[fieldId].Add(value)
	}
}

// fold walks the receipts once, accumulating the overall row and one row per configured
// status.
//
// Two rules that are easy to get backwards. Every configured status is seeded BEFORE the
// walk, so one matching no receipt still renders as a zero row rather than vanishing —
// that is what keeps the block's shape steady as the filter narrows. And a receipt whose
// status is not configured still counts toward the overall row: it is in the filter
// result, so excluding it would make the total disagree with the table's count.
func (service ReceiptSummaryService) fold(
	receipts []models.Receipt,
	statuses []models.ReceiptStatus,
	fieldIds []uint,
	fieldNames map[uint]string,
	configurationGroupId uint,
) structs.ReceiptSummary {
	overall := newSummaryAccumulator(fieldIds)

	buckets := make(map[models.ReceiptStatus]*summaryAccumulator, len(statuses))
	for _, status := range statuses {
		buckets[status] = newSummaryAccumulator(fieldIds)
	}

	for _, receipt := range receipts {
		fieldValues := summaryFieldValues(receipt, fieldIds)

		overall.add(receipt, fieldValues)
		if bucket, configured := buckets[receipt.Status]; configured {
			bucket.add(receipt, fieldValues)
		}
	}

	rows := make([]structs.ReceiptSummaryRow, 0, len(statuses))
	for _, status := range statuses {
		rows = append(rows, buckets[status].toRow(status, fieldIds, fieldNames))
	}

	return structs.ReceiptSummary{
		Enabled:              true,
		ConfigurationGroupId: configurationGroupId,
		Overall:              overall.toRow("", fieldIds, fieldNames),
		Statuses:             rows,
	}
}

// summaryFieldValues resolves a receipt's value for each configured field.
//
// Where a receipt holds SEVERAL values for one field, the lowest id wins — mirroring
// reporting/receiptsource.addCustomFields. Nothing stops the duplicate:
// custom_field_values carries no unique index on (receipt_id, custom_field_id), and the
// association is loaded without an ORDER BY, so preferring whichever came back first
// would hand the answer to the database and let two identical requests disagree.
//
// A row with no CurrencyValue never wins, so an empty low-id row cannot hide a real
// one. A field the receipt carries no value for is simply absent, which contributes
// nothing to the total.
func summaryFieldValues(receipt models.Receipt, fieldIds []uint) map[uint]decimal.Decimal {
	if len(fieldIds) == 0 || len(receipt.CustomFields) == 0 {
		return nil
	}

	values := make(map[uint]decimal.Decimal, len(fieldIds))
	winners := make(map[uint]uint, len(fieldIds))

	for _, customFieldValue := range receipt.CustomFields {
		if customFieldValue.CurrencyValue == nil {
			continue
		}
		if incumbent, held := winners[customFieldValue.CustomFieldId]; held && incumbent <= customFieldValue.ID {
			continue
		}

		winners[customFieldValue.CustomFieldId] = customFieldValue.ID
		values[customFieldValue.CustomFieldId] = *customFieldValue.CurrencyValue
	}

	return values
}

func (accumulator *summaryAccumulator) toRow(
	status models.ReceiptStatus,
	fieldIds []uint,
	fieldNames map[uint]string,
) structs.ReceiptSummaryRow {
	// Built from fieldIds rather than by ranging the map: a Go map has no order, and a
	// summary whose columns reshuffle between requests is unreadable.
	fieldTotals := make([]structs.ReceiptSummaryCustomFieldTotal, 0, len(fieldIds))
	for _, fieldId := range fieldIds {
		fieldTotals = append(fieldTotals, structs.ReceiptSummaryCustomFieldTotal{
			CustomFieldId: fieldId,
			Name:          fieldNames[fieldId],
			Total:         accumulator.fieldTotals[fieldId],
		})
	}

	return structs.ReceiptSummaryRow{
		Status:            status,
		ReceiptCount:      accumulator.receiptCount,
		Total:             accumulator.total,
		CustomFieldTotals: fieldTotals,
	}
}

// emptyReceiptSummary is the off state: a well-formed 200 with nothing in it. Every
// slice is non-nil so it serializes as [] rather than null.
func emptyReceiptSummary(configurationGroupId uint) structs.ReceiptSummary {
	return structs.ReceiptSummary{
		Enabled:              false,
		ConfigurationGroupId: configurationGroupId,
		Overall: structs.ReceiptSummaryRow{
			Total:             decimal.Zero,
			CustomFieldTotals: []structs.ReceiptSummaryCustomFieldTotal{},
		},
		Statuses: []structs.ReceiptSummaryRow{},
	}
}
