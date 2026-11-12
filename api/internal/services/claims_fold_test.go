package services

import (
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

func claimAmount(value string) *decimal.Decimal {
	amount := decimal.RequireFromString(value)
	return &amount
}

func claimDay(year int, month time.Month, day int) *time.Time {
	value := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &value
}

func linkedReceipt(id uint, amount string) repositories.LinkedReceipt {
	return repositories.LinkedReceipt{ReceiptId: id, GroupId: 1, Amount: decimal.RequireFromString(amount)}
}

func TestFoldClaim_FiguresFromReceiptsAndLedger(t *testing.T) {
	claim := models.Claim{
		BaseModel:     models.BaseModel{ID: 3},
		GroupId:       1,
		Name:          "PCS Rome",
		Status:        models.ClaimStatusPartiallyPaid,
		SubmittedDate: claimDay(2026, 9, 1),
	}
	receipts := []repositories.LinkedReceipt{linkedReceipt(1, "10.10"), linkedReceipt(2, "5.25")}
	entries := []models.ClaimLedgerEntry{
		{EntryType: models.ClaimLedgerEntryPayment, Date: *claimDay(2026, 9, 20), Amount: claimAmount("10.00")},
		{EntryType: models.ClaimLedgerEntryRecoupment, Date: *claimDay(2026, 10, 1), Amount: claimAmount("2.50")},
		// A note's amount is context only and never counts.
		{EntryType: models.ClaimLedgerEntryNote, Date: *claimDay(2026, 10, 2), Amount: claimAmount("99.00"), Note: "disallowed taxi"},
		// A type this build does not know counts toward nothing.
		{EntryType: "ADJUSTMENT", Date: *claimDay(2026, 10, 3), Amount: claimAmount("1000.00")},
	}

	fold := foldClaim(claim, receipts, entries, *claimDay(2026, 10, 11))
	view := fold.view

	checks := map[string][2]string{
		"receiptsTotal":          {view.ReceiptsTotal, "15.35"},
		"effectiveClaimedAmount": {view.EffectiveClaimedAmount, "15.35"},
		"paymentsTotal":          {view.PaymentsTotal, "10.00"},
		"recoupmentsTotal":       {view.RecoupmentsTotal, "2.50"},
		"paidAmount":             {view.PaidAmount, "7.50"},
		"outstandingAmount":      {view.OutstandingAmount, "7.85"},
		"submittedDate":          {view.SubmittedDate, "2026-09-01"},
	}
	for name, pair := range checks {
		if pair[0] != pair[1] {
			t.Errorf("%s: got %q, want %q", name, pair[0], pair[1])
		}
	}
	if view.ClaimedAmount != "" {
		utils.PrintTestError(t, view.ClaimedAmount, "absent when not set")
	}
	if view.ReceiptCount != 2 || len(view.ReceiptIds) != 2 {
		utils.PrintTestError(t, view.ReceiptIds, []uint{1, 2})
	}
	if len(view.LedgerEntries) != 4 {
		utils.PrintTestError(t, len(view.LedgerEntries), 4)
	}
	if !view.AwaitingSettlement || view.DaysWaiting == nil || *view.DaysWaiting != 40 {
		utils.PrintTestError(t, view.DaysWaiting, 40)
	}
	if !fold.outstanding.Equal(decimal.RequireFromString("7.85")) {
		utils.PrintTestError(t, fold.outstanding, "7.85")
	}
}

func TestFoldClaim_ExplicitClaimedAmountWinsAndOverpaymentIsNegative(t *testing.T) {
	claim := models.Claim{Status: models.ClaimStatusPaid, ClaimedAmount: claimAmount("100.00"), SubmittedDate: claimDay(2026, 1, 1)}
	entries := []models.ClaimLedgerEntry{
		{EntryType: models.ClaimLedgerEntryPayment, Amount: claimAmount("120.00")},
	}

	view := foldClaim(claim, []repositories.LinkedReceipt{linkedReceipt(1, "80.00")}, entries, *claimDay(2026, 2, 1)).view

	if view.EffectiveClaimedAmount != "100.00" || view.ClaimedAmount != "100.00" {
		utils.PrintTestError(t, view.EffectiveClaimedAmount, "100.00")
	}
	if view.OutstandingAmount != "-20.00" {
		utils.PrintTestError(t, view.OutstandingAmount, "-20.00")
	}
	// PAID does not await settlement, so it has no age.
	if view.AwaitingSettlement || view.DaysWaiting != nil {
		utils.PrintTestError(t, view.DaysWaiting, nil)
	}
}

// Money is folded in decimal, so amounts a float would smear add up to the cent.
func TestFoldClaim_DecimalPrecision(t *testing.T) {
	receipts := make([]repositories.LinkedReceipt, 0, 10)
	for i := uint(1); i <= 10; i++ {
		receipts = append(receipts, linkedReceipt(i, "0.10"))
	}

	view := foldClaim(models.Claim{Status: models.ClaimStatusDraft}, receipts, nil, *claimDay(2026, 1, 1)).view

	if view.ReceiptsTotal != "1.00" {
		utils.PrintTestError(t, view.ReceiptsTotal, "1.00")
	}
}

func TestFoldClaim_EmptyClaimSerializesArrays(t *testing.T) {
	view := foldClaim(models.Claim{Status: models.ClaimStatusDraft}, nil, nil, *claimDay(2026, 1, 1)).view

	if view.ReceiptIds == nil || view.LedgerEntries == nil {
		t.Errorf("receiptIds and ledgerEntries must be [] rather than null")
	}
	if view.ReceiptsTotal != "0.00" || view.OutstandingAmount != "0.00" {
		utils.PrintTestError(t, view.ReceiptsTotal, "0.00")
	}
}

func TestFoldClaim_FutureSubmittedDateWaitsZeroDays(t *testing.T) {
	claim := models.Claim{Status: models.ClaimStatusSubmitted, SubmittedDate: claimDay(2026, 12, 1)}

	view := foldClaim(claim, nil, nil, *claimDay(2026, 11, 1)).view

	if view.DaysWaiting == nil || *view.DaysWaiting != 0 {
		utils.PrintTestError(t, view.DaysWaiting, 0)
	}
}

func TestBuildClaimList_SummaryCoversAwaitingClaimsOnly(t *testing.T) {
	today := *claimDay(2026, 11, 11)
	submitted := foldClaim(models.Claim{Status: models.ClaimStatusSubmitted, SubmittedDate: claimDay(2026, 11, 1), ClaimedAmount: claimAmount("50.00")}, nil, nil, today)
	disputed := foldClaim(models.Claim{Status: models.ClaimStatusDisputed, SubmittedDate: claimDay(2026, 1, 1), ClaimedAmount: claimAmount("25.00")}, nil, nil, today)
	closed := foldClaim(models.Claim{Status: models.ClaimStatusClosed, ClaimedAmount: claimAmount("1000.00")}, nil, nil, today)
	unknown := foldClaim(models.Claim{Status: "SOMETHING_NEW", ClaimedAmount: claimAmount("1000.00")}, nil, nil, today)

	list := buildClaimList([]claimFold{submitted, disputed, closed, unknown}, today)

	if list.Summary.ClaimCount != 4 || list.Summary.AwaitingCount != 2 {
		utils.PrintTestError(t, list.Summary, "4 claims, 2 awaiting")
	}
	if list.Summary.OutstandingTotal != "75.00" {
		utils.PrintTestError(t, list.Summary.OutstandingTotal, "75.00")
	}
	if list.Summary.OldestDaysWaiting == nil || *list.Summary.OldestDaysWaiting != 314 {
		utils.PrintTestError(t, list.Summary.OldestDaysWaiting, 314)
	}
	if list.Summary.AsOfDate != "2026-11-11" {
		utils.PrintTestError(t, list.Summary.AsOfDate, "2026-11-11")
	}
}

func TestBuildClaimList_EmptyList(t *testing.T) {
	list := buildClaimList(nil, *claimDay(2026, 11, 11))

	if list.Claims == nil {
		t.Errorf("claims must be [] rather than null")
	}
	if list.Summary.OldestDaysWaiting != nil || list.Summary.OutstandingTotal != "0.00" {
		utils.PrintTestError(t, list.Summary, "zeroed summary")
	}
}

// "Today" is the app zone's calendar day, not the server's. UTC+14 is ahead of UTC
// for 14 hours of every day, so most runs of this test see a different day there.
func TestTodayInAppZone_UsesTheAppTimeZone(t *testing.T) {
	defer repositories.TruncateTestDb()
	if err := repositories.SetAppTimeZoneForTests("Pacific/Kiritimati"); err != nil {
		t.Fatalf("set zone: %v", err)
	}

	got := todayInAppZone()
	want := utils.CalendarDayUTC(time.Now().In(time.FixedZone("UTC+14", 14*60*60)))

	if !got.Equal(want) {
		utils.PrintTestError(t, got, want)
	}
}
