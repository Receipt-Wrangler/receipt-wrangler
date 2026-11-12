package services

import (
	"errors"
	"fmt"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ErrClaimNotFound covers a claim (or ledger entry) that does not exist and one
// the caller may not see. Handlers answer both with the same 403, so a response
// never says which.
var ErrClaimNotFound = errors.New("claim not found")

// ErrClaimAllGroupDestination is returned when a claim would be created in the
// synthetic All group, which only aggregates real groups and so is never a
// destination. Handlers answer 400.
var ErrClaimAllGroupDestination = errors.New("a claim cannot be created in the all group")

// The per-receipt messages addReceiptsToClaim and removeReceiptsFromClaim report.
// A receipt that does not exist, lives in another group, or is hidden from the
// caller always gets the same message, so a failure never reveals that a receipt
// exists.
const (
	ClaimReceiptNotAddableMessage   = "This receipt cannot be added: it is not a receipt of this claim's group that you can see"
	ClaimReceiptInOtherClaimMessage = "This receipt is already in another claim; remove it from that claim first"
	ClaimReceiptNotInClaimMessage   = "This receipt is not in this claim"
)

// claimDayFormat is how every calendar day goes out on the wire.
const claimDayFormat = time.DateOnly

// ClaimService holds the reimbursement-claim rules: who sees which claim, which
// receipts and ledger entries count for a viewer, and the figures.
//
// What a viewer sees, and so what every figure is computed from:
//   - Claims are gated on group.receipts.read (writes on group.receipts.update),
//     which the handlers check through HandleRequest.
//   - In a group with isolated members, a claim or ledger entry recorded by a
//     member the viewer cannot see is left out, as comments are. In any other
//     group every claim and entry is visible.
//   - A linked receipt counts only if the viewer may see it (paid-by visibility,
//     which already folds in member isolation).
//
// Nothing in a response says that something was left out.
type ClaimService struct {
	BaseService
}

func NewClaimService(tx *gorm.DB) ClaimService {
	return ClaimService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
}

// GetClaimsForGroup lists a group's claims, or for the synthetic All group the
// claims of every group where the caller holds group.receipts.read, mirroring how
// the receipts list expands it.
func (service ClaimService) GetClaimsForGroup(userId uint, groupId uint) (structs.ClaimListView, error) {
	groupIds, err := service.readableGroupIds(userId, groupId)
	if err != nil {
		return structs.ClaimListView{}, err
	}

	claims, err := repositories.NewClaimRepository(service.TX).GetClaimsByGroupIds(groupIds)
	if err != nil {
		return structs.ClaimListView{}, err
	}

	today := todayInAppZone()
	folds, err := service.foldClaims(userId, claims, today)
	if err != nil {
		return structs.ClaimListView{}, err
	}

	return buildClaimList(folds, today), nil
}

// GetClaimForUser returns one claim the caller may see, or ErrClaimNotFound.
func (service ClaimService) GetClaimForUser(userId uint, claimId uint) (structs.ClaimView, error) {
	claim, err := service.loadVisibleClaim(userId, claimId)
	if err != nil {
		return structs.ClaimView{}, err
	}

	folds, err := service.foldClaims(userId, []models.Claim{claim}, todayInAppZone())
	if err != nil {
		return structs.ClaimView{}, err
	}
	if len(folds) == 0 {
		return structs.ClaimView{}, ErrClaimNotFound
	}
	return folds[0].view, nil
}

func (service ClaimService) CreateClaim(userId uint, groupId uint, command commands.UpsertClaimCommand) (structs.ClaimView, error) {
	isAllGroup, err := repositories.NewGroupRepository(service.TX).IsAllGroup(groupId)
	if err != nil {
		return structs.ClaimView{}, err
	}
	if isAllGroup {
		return structs.ClaimView{}, ErrClaimAllGroupDestination
	}

	claim := command.ToClaim()
	claim.GroupId = groupId
	claim.CreatedBy = &userId

	created, err := repositories.NewClaimRepository(service.TX).CreateClaim(claim)
	if err != nil {
		return structs.ClaimView{}, err
	}

	return service.GetClaimForUser(userId, created.ID)
}

func (service ClaimService) UpdateClaim(userId uint, claimId uint, command commands.UpsertClaimCommand) (structs.ClaimView, error) {
	if _, err := service.loadVisibleClaim(userId, claimId); err != nil {
		return structs.ClaimView{}, err
	}

	if err := repositories.NewClaimRepository(service.TX).UpdateClaim(claimId, command.ToClaim()); err != nil {
		return structs.ClaimView{}, err
	}

	return service.GetClaimForUser(userId, claimId)
}

func (service ClaimService) DeleteClaim(userId uint, claimId uint) error {
	if _, err := service.loadVisibleClaim(userId, claimId); err != nil {
		return err
	}

	return repositories.NewClaimRepository(service.TX).DeleteClaim(claimId)
}

// AddReceiptsToClaim links each receipt the caller may add and reports each one it
// did not, in request order. A receipt may be added when it belongs to the claim's
// group and the caller may see it. One already in this claim counts as added; one
// in another claim is refused rather than moved, so a submitted claim's record is
// never changed by accident.
func (service ClaimService) AddReceiptsToClaim(userId uint, claimId uint, receiptIds []uint) (structs.ClaimReceiptsResult, error) {
	result := structs.NewClaimReceiptsResult()

	claim, err := service.loadVisibleClaim(userId, claimId)
	if err != nil {
		return result, err
	}

	visible, err := service.visibleReceiptIdsInGroup(userId, claim.GroupId, receiptIds)
	if err != nil {
		return result, err
	}

	claimRepository := repositories.NewClaimRepository(service.TX)
	links, err := claimRepository.GetClaimLinksForReceipts(receiptIds)
	if err != nil {
		return result, err
	}
	linkByReceipt := make(map[uint]repositories.ReceiptClaimLink, len(links))
	for _, link := range links {
		linkByReceipt[link.ReceiptId] = link
	}

	toLink := make([]uint, 0, len(receiptIds))
	for _, receiptId := range receiptIds {
		if _, ok := visible[receiptId]; !ok {
			result.Failures = append(result.Failures, structs.ClaimReceiptsFailure{ReceiptId: receiptId, Message: ClaimReceiptNotAddableMessage})
			continue
		}

		link, linked := linkByReceipt[receiptId]
		switch {
		case linked && link.ClaimId == claimId:
			// Already here: nothing to do.
		case linked && link.ClaimGroupId == claim.GroupId:
			result.Failures = append(result.Failures, structs.ClaimReceiptsFailure{ReceiptId: receiptId, Message: ClaimReceiptInOtherClaimMessage})
			continue
		default:
			// Unlinked, or a stale link to a claim of a group the receipt has left.
			toLink = append(toLink, receiptId)
		}

		result.ReceiptIds = append(result.ReceiptIds, receiptId)
	}

	if err := claimRepository.LinkReceipts(claimId, toLink); err != nil {
		return structs.NewClaimReceiptsResult(), err
	}

	return result, nil
}

// RemoveReceiptsFromClaim unlinks each listed receipt that is in the claim and
// visible to the caller, and reports each other one in request order. A hidden
// receipt reads exactly like one that is not in the claim.
func (service ClaimService) RemoveReceiptsFromClaim(userId uint, claimId uint, receiptIds []uint) (structs.ClaimReceiptsResult, error) {
	result := structs.NewClaimReceiptsResult()

	claim, err := service.loadVisibleClaim(userId, claimId)
	if err != nil {
		return result, err
	}

	claimRepository := repositories.NewClaimRepository(service.TX)
	linked, err := claimRepository.GetLinkedReceipts([]uint{claimId})
	if err != nil {
		return result, err
	}
	visible, err := service.visibleLinkedReceiptIds(userId, linked[claimId])
	if err != nil {
		return result, err
	}

	toUnlink := make([]uint, 0, len(receiptIds))
	for _, receiptId := range receiptIds {
		if _, ok := visible[receiptId]; !ok {
			result.Failures = append(result.Failures, structs.ClaimReceiptsFailure{ReceiptId: receiptId, Message: ClaimReceiptNotInClaimMessage})
			continue
		}
		toUnlink = append(toUnlink, receiptId)
		result.ReceiptIds = append(result.ReceiptIds, receiptId)
	}

	if err := claimRepository.UnlinkReceipts(claim.ID, toUnlink); err != nil {
		return structs.NewClaimReceiptsResult(), err
	}

	return result, nil
}

func (service ClaimService) CreateLedgerEntry(userId uint, claimId uint, command commands.UpsertClaimLedgerEntryCommand) (structs.ClaimLedgerEntryView, error) {
	if _, err := service.loadVisibleClaim(userId, claimId); err != nil {
		return structs.ClaimLedgerEntryView{}, err
	}

	entry := command.ToClaimLedgerEntry()
	entry.ClaimId = claimId
	entry.CreatedBy = &userId

	created, err := repositories.NewClaimRepository(service.TX).CreateLedgerEntry(entry)
	if err != nil {
		return structs.ClaimLedgerEntryView{}, err
	}
	return buildLedgerEntryView(created), nil
}

func (service ClaimService) UpdateLedgerEntry(userId uint, claimId uint, entryId uint, command commands.UpsertClaimLedgerEntryCommand) (structs.ClaimLedgerEntryView, error) {
	if _, err := service.loadVisibleLedgerEntry(userId, claimId, entryId); err != nil {
		return structs.ClaimLedgerEntryView{}, err
	}

	claimRepository := repositories.NewClaimRepository(service.TX)
	if err := claimRepository.UpdateLedgerEntry(entryId, command.ToClaimLedgerEntry()); err != nil {
		return structs.ClaimLedgerEntryView{}, err
	}

	updated, err := claimRepository.GetLedgerEntryById(entryId)
	if err != nil {
		return structs.ClaimLedgerEntryView{}, err
	}
	return buildLedgerEntryView(updated), nil
}

func (service ClaimService) DeleteLedgerEntry(userId uint, claimId uint, entryId uint) error {
	if _, err := service.loadVisibleLedgerEntry(userId, claimId, entryId); err != nil {
		return err
	}

	return repositories.NewClaimRepository(service.TX).DeleteLedgerEntry(entryId)
}

// ReceiptDeletionRefusals returns, for each of the given receipts that is in a
// claim, the message refusing its deletion. A receipt in a claim is evidence of
// what was claimed, so it is never dropped out of one silently: the user removes it
// from the claim first. This applies whatever the claim's status, since removing a
// receipt from a draft is one click.
//
// Only the receipt delete endpoints call this. Deleting a group or a user still
// deletes their receipts, and the links go with them through the cascading foreign
// key, because refusing there would make those deletes impossible to finish.
//
// The claim is named only when the caller may see it; otherwise the message says
// "a claim", so an isolated member is not told about another member's claim.
func (service ClaimService) ReceiptDeletionRefusals(userId uint, receiptIds []uint) (map[uint]string, error) {
	refusals := make(map[uint]string)

	links, err := repositories.NewClaimRepository(service.TX).GetClaimLinksForGroupReceipts(receiptIds)
	if err != nil {
		return nil, err
	}

	visibility := newClaimVisibility(NewPermissionService(service.TX), userId)
	for _, link := range links {
		visible, err := visibility.creatorVisible(link.ClaimGroupId, link.ClaimCreatedBy)
		if err != nil {
			return nil, err
		}
		if visible {
			refusals[link.ReceiptId] = fmt.Sprintf("This receipt is in a claim; remove it from claim '%s' first", link.ClaimName)
		} else {
			refusals[link.ReceiptId] = "This receipt is in a claim; remove it from the claim first"
		}
	}

	return refusals, nil
}

// readableGroupIds resolves the groups a list covers: the group itself, or for the
// All group every group of the caller's where they hold group.receipts.read. The
// All group's own membership grants nothing, since no claim lives in it.
func (service ClaimService) readableGroupIds(userId uint, groupId uint) ([]uint, error) {
	isAllGroup, err := repositories.NewGroupRepository(service.TX).IsAllGroup(groupId)
	if err != nil {
		return nil, err
	}
	if !isAllGroup {
		return []uint{groupId}, nil
	}

	memberGroupIds, err := repositories.NewGroupMemberRepository(service.TX).GetGroupIdsByUserId(utils.UintToString(userId))
	if err != nil {
		return nil, err
	}

	permissionService := NewPermissionService(service.TX)
	readable := make([]uint, 0, len(memberGroupIds))
	for _, memberGroupId := range memberGroupIds {
		if memberGroupId == groupId {
			continue
		}
		canRead, err := permissionService.HasGroupPermissions(userId, memberGroupId, permissions.GroupReceiptsRead)
		if err != nil {
			return nil, err
		}
		if canRead {
			readable = append(readable, memberGroupId)
		}
	}
	return readable, nil
}

// loadVisibleClaim loads a claim the caller may see. HandleRequest has already
// checked the group permission; this adds member isolation.
func (service ClaimService) loadVisibleClaim(userId uint, claimId uint) (models.Claim, error) {
	claim, err := repositories.NewClaimRepository(service.TX).GetClaimById(claimId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Claim{}, ErrClaimNotFound
	}
	if err != nil {
		return models.Claim{}, err
	}

	visible, err := newClaimVisibility(NewPermissionService(service.TX), userId).creatorVisible(claim.GroupId, claim.CreatedBy)
	if err != nil {
		return models.Claim{}, err
	}
	if !visible {
		return models.Claim{}, ErrClaimNotFound
	}
	return claim, nil
}

// loadVisibleLedgerEntry loads an entry of a claim the caller may see, which is
// itself visible to them and belongs to that claim.
func (service ClaimService) loadVisibleLedgerEntry(userId uint, claimId uint, entryId uint) (models.ClaimLedgerEntry, error) {
	claim, err := service.loadVisibleClaim(userId, claimId)
	if err != nil {
		return models.ClaimLedgerEntry{}, err
	}

	entry, err := repositories.NewClaimRepository(service.TX).GetLedgerEntryById(entryId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ClaimLedgerEntry{}, ErrClaimNotFound
	}
	if err != nil {
		return models.ClaimLedgerEntry{}, err
	}
	if entry.ClaimId != claimId {
		return models.ClaimLedgerEntry{}, ErrClaimNotFound
	}

	visible, err := newClaimVisibility(NewPermissionService(service.TX), userId).creatorVisible(claim.GroupId, entry.CreatedBy)
	if err != nil {
		return models.ClaimLedgerEntry{}, err
	}
	if !visible {
		return models.ClaimLedgerEntry{}, ErrClaimNotFound
	}
	return entry, nil
}

// visibleReceiptIdsInGroup returns which of the given receipts exist, live in
// groupId and are visible to the caller.
func (service ClaimService) visibleReceiptIdsInGroup(userId uint, groupId uint, receiptIds []uint) (map[uint]struct{}, error) {
	receipts, err := repositories.NewReceiptRepository(service.TX).GetReceiptsForAuthorization(receiptIds)
	if err != nil {
		return nil, err
	}

	inGroup := make([]models.Receipt, 0, len(receipts))
	for _, receipt := range receipts {
		if receipt.GroupId == groupId {
			inGroup = append(inGroup, receipt)
		}
	}

	visible, err := NewPermissionService(service.TX).FilterReceiptsByPaidBy(userId, inGroup)
	if err != nil {
		return nil, err
	}

	ids := make(map[uint]struct{}, len(visible))
	for _, receipt := range visible {
		ids[receipt.ID] = struct{}{}
	}
	return ids, nil
}

// visibleLinkedReceiptIds returns the ids of the linked receipts the caller may see.
func (service ClaimService) visibleLinkedReceiptIds(userId uint, linked []repositories.LinkedReceipt) (map[uint]struct{}, error) {
	receipts := make([]models.Receipt, 0, len(linked))
	for _, row := range linked {
		receipts = append(receipts, models.Receipt{
			BaseModel:    models.BaseModel{ID: row.ReceiptId},
			GroupId:      row.GroupId,
			PaidByUserID: row.PaidByUserId,
		})
	}

	visible, err := NewPermissionService(service.TX).FilterReceiptsByPaidBy(userId, receipts)
	if err != nil {
		return nil, err
	}

	ids := make(map[uint]struct{}, len(visible))
	for _, receipt := range visible {
		ids[receipt.ID] = struct{}{}
	}
	return ids, nil
}

// claimFold is one claim's view plus the exact figure the list summary adds up,
// so the summary never re-parses a formatted string.
type claimFold struct {
	view        structs.ClaimView
	outstanding decimal.Decimal
}

// foldClaims drops the claims the caller may not see and builds a view of each of
// the rest from what the caller may see of it. Links and entries are loaded in one
// query each for the whole batch.
func (service ClaimService) foldClaims(userId uint, claims []models.Claim, today time.Time) ([]claimFold, error) {
	visibility := newClaimVisibility(NewPermissionService(service.TX), userId)

	visibleClaims := make([]models.Claim, 0, len(claims))
	claimIds := make([]uint, 0, len(claims))
	for _, claim := range claims {
		visible, err := visibility.creatorVisible(claim.GroupId, claim.CreatedBy)
		if err != nil {
			return nil, err
		}
		if visible {
			visibleClaims = append(visibleClaims, claim)
			claimIds = append(claimIds, claim.ID)
		}
	}

	claimRepository := repositories.NewClaimRepository(service.TX)
	linked, err := claimRepository.GetLinkedReceipts(claimIds)
	if err != nil {
		return nil, err
	}
	entries, err := claimRepository.GetLedgerEntries(claimIds)
	if err != nil {
		return nil, err
	}

	allLinked := make([]repositories.LinkedReceipt, 0)
	for _, rows := range linked {
		allLinked = append(allLinked, rows...)
	}
	visibleReceipts, err := service.visibleLinkedReceiptIds(userId, allLinked)
	if err != nil {
		return nil, err
	}

	folds := make([]claimFold, 0, len(visibleClaims))
	for _, claim := range visibleClaims {
		receipts := make([]repositories.LinkedReceipt, 0, len(linked[claim.ID]))
		for _, row := range linked[claim.ID] {
			if _, ok := visibleReceipts[row.ReceiptId]; ok {
				receipts = append(receipts, row)
			}
		}

		claimEntries := make([]models.ClaimLedgerEntry, 0, len(entries[claim.ID]))
		for _, entry := range entries[claim.ID] {
			visible, err := visibility.creatorVisible(claim.GroupId, entry.CreatedBy)
			if err != nil {
				return nil, err
			}
			if visible {
				claimEntries = append(claimEntries, entry)
			}
		}

		folds = append(folds, foldClaim(claim, receipts, claimEntries, today))
	}

	return folds, nil
}

// foldClaim computes one claim's figures from the receipts and entries the viewer
// may see. Money is folded with shopspring/decimal, never a SQL SUM: SQLite has no
// decimal type, so a SUM would report different cents there than on
// Postgres/MySQL. Pure, so the arithmetic is testable without a database.
func foldClaim(claim models.Claim, receipts []repositories.LinkedReceipt, entries []models.ClaimLedgerEntry, today time.Time) claimFold {
	receiptIds := make([]uint, 0, len(receipts))
	receiptsTotal := decimal.Zero
	for _, receipt := range receipts {
		receiptIds = append(receiptIds, receipt.ReceiptId)
		receiptsTotal = receiptsTotal.Add(receipt.Amount)
	}

	effectiveClaimed := receiptsTotal
	if claim.ClaimedAmount != nil {
		effectiveClaimed = *claim.ClaimedAmount
	}

	payments := decimal.Zero
	recoupments := decimal.Zero
	entryViews := make([]structs.ClaimLedgerEntryView, 0, len(entries))
	for _, entry := range entries {
		entryViews = append(entryViews, buildLedgerEntryView(entry))
		if entry.Amount == nil {
			continue
		}
		switch entry.EntryType {
		case models.ClaimLedgerEntryPayment:
			payments = payments.Add(*entry.Amount)
		case models.ClaimLedgerEntryRecoupment:
			recoupments = recoupments.Add(*entry.Amount)
		}
	}

	paid := payments.Sub(recoupments)
	outstanding := effectiveClaimed.Sub(paid)
	awaiting := models.ClaimStatusAwaitsSettlement(claim.Status)

	view := structs.ClaimView{
		Id:                     claim.ID,
		GroupId:                claim.GroupId,
		Name:                   claim.Name,
		ReferenceNumber:        claim.ReferenceNumber,
		Type:                   claim.ClaimType,
		Status:                 claim.Status,
		AwaitingSettlement:     awaiting,
		SubmittedDate:          formatClaimDay(claim.SubmittedDate),
		DueDate:                formatClaimDay(claim.DueDate),
		Notes:                  claim.Notes,
		CreatedAt:              claim.CreatedAt,
		UpdatedAt:              claim.UpdatedAt,
		CreatedBy:              claim.CreatedBy,
		ReceiptIds:             receiptIds,
		ReceiptCount:           len(receiptIds),
		ReceiptsTotal:          formatClaimMoney(receiptsTotal),
		EffectiveClaimedAmount: formatClaimMoney(effectiveClaimed),
		PaymentsTotal:          formatClaimMoney(payments),
		RecoupmentsTotal:       formatClaimMoney(recoupments),
		PaidAmount:             formatClaimMoney(paid),
		OutstandingAmount:      formatClaimMoney(outstanding),
		LedgerEntries:          entryViews,
	}
	if claim.ClaimedAmount != nil {
		view.ClaimedAmount = formatClaimMoney(*claim.ClaimedAmount)
	}
	if awaiting && claim.SubmittedDate != nil {
		days := daysBetween(utils.CalendarDayUTC(claim.SubmittedDate.UTC()), today)
		view.DaysWaiting = &days
	}

	return claimFold{view: view, outstanding: outstanding}
}

// buildClaimList adds the header figures over the claims in the list.
func buildClaimList(folds []claimFold, today time.Time) structs.ClaimListView {
	claims := make([]structs.ClaimView, 0, len(folds))
	outstandingTotal := decimal.Zero
	awaitingCount := 0
	var oldest *int

	for _, fold := range folds {
		claims = append(claims, fold.view)
		if !fold.view.AwaitingSettlement {
			continue
		}
		awaitingCount++
		outstandingTotal = outstandingTotal.Add(fold.outstanding)
		if fold.view.DaysWaiting != nil && (oldest == nil || *fold.view.DaysWaiting > *oldest) {
			days := *fold.view.DaysWaiting
			oldest = &days
		}
	}

	return structs.ClaimListView{
		Claims: claims,
		Summary: structs.ClaimListSummaryView{
			AsOfDate:          today.Format(claimDayFormat),
			ClaimCount:        len(claims),
			AwaitingCount:     awaitingCount,
			OutstandingTotal:  formatClaimMoney(outstandingTotal),
			OldestDaysWaiting: oldest,
		},
	}
}

func buildLedgerEntryView(entry models.ClaimLedgerEntry) structs.ClaimLedgerEntryView {
	view := structs.ClaimLedgerEntryView{
		Id:        entry.ID,
		ClaimId:   entry.ClaimId,
		Type:      entry.EntryType,
		Date:      formatClaimDay(&entry.Date),
		Note:      entry.Note,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
		CreatedBy: entry.CreatedBy,
	}
	if entry.Amount != nil {
		view.Amount = formatClaimMoney(*entry.Amount)
	}
	return view
}

// todayInAppZone is today's calendar day in the app time zone, as midnight UTC so
// it compares directly with stored calendar days.
func todayInAppZone() time.Time {
	return utils.CalendarDayUTC(time.Now().In(repositories.GetAppLocation()))
}

// daysBetween counts calendar days from one midnight-UTC day to another, never
// below zero (a submitted date in the future has waited no time yet).
func daysBetween(from time.Time, to time.Time) int {
	days := int(to.Sub(from).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// formatClaimDay renders a stored calendar day by its UTC day, so no zone shifts it.
func formatClaimDay(day *time.Time) string {
	if day == nil || day.IsZero() {
		return ""
	}
	return day.UTC().Format(claimDayFormat)
}

func formatClaimMoney(value decimal.Decimal) string {
	return value.StringFixed(2)
}

// claimVisibility answers "may this viewer see what this user recorded in this
// group", memoized per group across a batch. A record with no creator is visible.
type claimVisibility struct {
	permissionService PermissionService
	viewerId          uint
	byGroup           map[uint]claimGroupVisibility
}

type claimGroupVisibility struct {
	visible      map[uint]struct{}
	unrestricted bool
}

func newClaimVisibility(permissionService PermissionService, viewerId uint) *claimVisibility {
	return &claimVisibility{
		permissionService: permissionService,
		viewerId:          viewerId,
		byGroup:           make(map[uint]claimGroupVisibility),
	}
}

func (visibility *claimVisibility) creatorVisible(groupId uint, createdBy *uint) (bool, error) {
	if createdBy == nil {
		return true, nil
	}

	group, ok := visibility.byGroup[groupId]
	if !ok {
		visible, unrestricted, err := visibility.permissionService.GetVisibleUserIdsForUserInGroup(visibility.viewerId, groupId)
		if err != nil {
			return false, err
		}
		group = claimGroupVisibility{visible: visible, unrestricted: unrestricted}
		visibility.byGroup[groupId] = group
	}

	if group.unrestricted {
		return true, nil
	}
	return isUserVisible(*createdBy, visibility.viewerId, group.visible), nil
}
