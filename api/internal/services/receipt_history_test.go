package services

import (
	"encoding/json"
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// These tests write history the way the app does — CreateReceipt for the
// creation, the ReceiptChangeRecorder for every change — and hand-build rows
// only for what the app no longer writes (version 1 rows, older origins).

func historyUint(value uint) *uint {
	return &value
}

func resetHistoryTest() {
	repositories.TruncateTestDb()
	clearRolePermissionCacheAll()
	clearGroupRoleGrantCacheAll()
}

// createHistoryReceipt creates a receipt through CreateReceipt with its
// RECEIPT_UPLOADED row, as the receipt form does.
func createHistoryReceipt(t *testing.T, createdBy uint, command commands.UpsertReceiptCommand) models.Receipt {
	t.Helper()
	if command.Name == "" {
		command.Name = "Groceries"
	}
	if command.Amount.IsZero() {
		command.Amount = decimal.RequireFromString("20.00")
	}
	command.Date = time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	if command.Status == "" {
		command.Status = models.OPEN
	}

	receipt, err := repositories.NewReceiptRepository(nil).CreateReceipt(command, createdBy, true)
	if err != nil {
		t.Fatalf("create receipt: %v", err)
	}
	return receipt
}

// recordHistoryChange makes a change through the recorder, as every handler that
// changes a receipt does.
func recordHistoryChange(t *testing.T, userId uint, receiptId uint, change func(db *gorm.DB) error) {
	t.Helper()
	err := repositories.NewReceiptChangeRecorder(userId, receiptId).Record(func() error {
		return change(repositories.GetDB())
	})
	if err != nil {
		t.Fatalf("record change: %v", err)
	}
}

func renameReceipt(receiptId uint, name string) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		return db.Model(&models.Receipt{}).Where("id = ?", receiptId).Update("name", name).Error
	}
}

func getHistory(t *testing.T, viewerId uint, receiptId uint) structs.ReceiptHistory {
	t.Helper()
	history, err := NewReceiptHistoryService(nil).GetReceiptHistory(viewerId, receiptId)
	if err != nil {
		t.Fatalf("GetReceiptHistory: %v", err)
	}
	return history
}

func historyActions(history structs.ReceiptHistory) []string {
	actions := make([]string, len(history.Entries))
	for i, entry := range history.Entries {
		actions[i] = entry.Action
	}
	return actions
}

func changeFields(entry structs.ReceiptHistoryEntry) []string {
	fields := make([]string, len(entry.Changes))
	for i, change := range entry.Changes {
		fields[i] = change.Field + " " + change.Kind
	}
	return fields
}

func assertStrings(t *testing.T, what string, got []string, want []string) {
	t.Helper()
	gotJson, _ := json.Marshal(got)
	wantJson, _ := json.Marshal(want)
	if string(gotJson) != string(wantJson) {
		t.Errorf("%s = %s, want %s", what, gotJson, wantJson)
	}
}

func TestGetReceiptHistoryOrdersEntriesNewestFirstWithCounts(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := createHistoryReceipt(t, 1, commands.UpsertReceiptCommand{GroupId: 1, PaidByUserID: 1})

	recordHistoryChange(t, 2, receipt.ID, func(db *gorm.DB) error {
		return db.Model(&models.Receipt{}).Where("id = ?", receipt.ID).
			Updates(map[string]interface{}{"name": "Weekly shop", "amount": decimal.RequireFromString("25.00")}).Error
	})
	recordHistoryChange(t, 1, receipt.ID, func(db *gorm.DB) error {
		return db.Create(&models.Comment{Comment: "Paid in cash", ReceiptId: receipt.ID, UserId: historyUint(1)}).Error
	})
	recordHistoryChange(t, 2, receipt.ID, func(db *gorm.DB) error {
		resolvedAt := time.Now().UTC()
		return db.Model(&models.Receipt{}).Where("id = ?", receipt.ID).
			Updates(map[string]interface{}{"status": models.RESOLVED, "resolved_date": resolvedAt}).Error
	})
	repositories.NewReceiptChangeRecorder(3, receipt.ID).RecordFailure(errors.New("the save failed"))

	history := getHistory(t, 1, receipt.ID)

	assertStrings(t, "actions", historyActions(history), []string{
		structs.ReceiptHistoryActionFailed,
		structs.ReceiptHistoryActionResolved,
		structs.ReceiptHistoryActionImagesCommentsUpdated,
		structs.ReceiptHistoryActionUpdated,
		structs.ReceiptHistoryActionCreated,
	})
	for i := 1; i < len(history.Entries); i++ {
		if history.Entries[i-1].SystemTaskId <= history.Entries[i].SystemTaskId {
			t.Errorf("entries are not newest first: %d then %d", history.Entries[i-1].SystemTaskId, history.Entries[i].SystemTaskId)
		}
	}

	failed := history.Entries[0]
	if failed.ErrorMessage != "the save failed" || failed.HasSnapshots || len(failed.Changes) != 0 ||
		failed.Actor.Type != structs.ReceiptHistoryActorUser || *failed.Actor.UserId != 3 || failed.RecordedAt.IsZero() {
		t.Errorf("failed entry = %+v", failed)
	}

	update := history.Entries[3]
	assertStrings(t, "update changes", changeFields(update), []string{"NAME CHANGED", "AMOUNT CHANGED"})
	if update.BeforeState != structs.ReceiptHistoryBeforeStateComplete || !update.HasSnapshots {
		t.Errorf("update entry = %+v, want complete snapshots", update)
	}

	creation := history.Entries[4]
	if !creation.HasSnapshots || creation.BeforeState != "" || creation.Actor.Type != structs.ReceiptHistoryActorUser || *creation.Actor.UserId != 1 {
		t.Errorf("creation entry = %+v", creation)
	}

	// Changes leave out the creation: 2 (name, amount) + 1 (comment) + 2
	// (status, resolved date). Saves leave out the failure. Contributors are
	// users 1, 2 and 3.
	if history.ChangeCount != 5 || history.SaveCount != 4 || history.ContributorCount != 3 {
		t.Errorf("counts = %d changes, %d saves, %d contributors; want 5, 4, 3", history.ChangeCount, history.SaveCount, history.ContributorCount)
	}
}

func TestGetReceiptHistoryEmptyForAReceiptWithNoTasks(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := models.Receipt{Name: "Imported", GroupId: 1, PaidByUserID: 1, Date: time.Now()}
	repositories.GetDB().Create(&receipt)

	history := getHistory(t, 1, receipt.ID)
	if history.Entries == nil || len(history.Entries) != 0 || history.SaveCount != 0 || history.ContributorCount != 0 {
		t.Errorf("history = %+v, want empty", history)
	}
}

func TestGetReceiptHistoryMissingReceiptIsDenied(t *testing.T) {
	defer resetHistoryTest()

	_, err := NewReceiptHistoryService(nil).GetReceiptHistory(1, 404)
	if !errors.Is(err, ErrReceiptAccessDenied) {
		t.Errorf("err = %v, want ErrReceiptAccessDenied", err)
	}
}

func TestReceiptHistoryActor(t *testing.T) {
	emailReceipt := &models.Receipt{EmailInbox: "receipts@example.com"}
	integrationReceipt := &models.Receipt{BaseModel: models.BaseModel{CreatedByString: constants.EMAIL_INTEGRATION_CREATED_BY}}
	parents := map[uint]models.SystemTask{
		10: {BaseModel: models.BaseModel{ID: 10}, Type: models.QUICK_SCAN, RanByUserId: historyUint(5)},
		11: {BaseModel: models.BaseModel{ID: 11}, Type: models.QUICK_SCAN, RanByUserId: historyUint(0)},
		12: {BaseModel: models.BaseModel{ID: 12}, Type: models.EMAIL_UPLOAD},
	}

	tests := []struct {
		name   string
		record historyRecord
		want   structs.ReceiptHistoryActor
	}{
		{
			name:   "the user who ran it",
			record: historyRecord{task: models.SystemTask{RanByUserId: historyUint(3), AssociatedSystemTaskId: historyUint(10)}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorUser, UserId: historyUint(3)},
		},
		{
			name:   "a quick scan, with the scanning user",
			record: historyRecord{creation: true, task: models.SystemTask{AssociatedSystemTaskId: historyUint(10)}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorQuickScan, UserId: historyUint(5)},
		},
		{
			name:   "a ran-by of 0 is nobody",
			record: historyRecord{creation: true, task: models.SystemTask{RanByUserId: historyUint(0), AssociatedSystemTaskId: historyUint(11)}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorQuickScan},
		},
		{
			name:   "an email upload",
			record: historyRecord{creation: true, task: models.SystemTask{AssociatedSystemTaskId: historyUint(12)}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorEmail},
		},
		{
			name:   "a parent id of 0, on an email receipt",
			record: historyRecord{creation: true, after: integrationReceipt, task: models.SystemTask{AssociatedSystemTaskId: historyUint(0)}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorEmail},
		},
		{
			name:   "a missing parent, on a receipt with an inbox",
			record: historyRecord{creation: true, after: emailReceipt, task: models.SystemTask{AssociatedSystemTaskId: historyUint(99)}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorEmail},
		},
		{
			name:   "an update of an email receipt by nobody is the system",
			record: historyRecord{before: emailReceipt, after: emailReceipt},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorSystem},
		},
		{
			name:   "anything else is the system",
			record: historyRecord{creation: true, after: &models.Receipt{}},
			want:   structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorSystem},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := receiptHistoryActor(test.record, parents)
			gotJson, _ := json.Marshal(got)
			wantJson, _ := json.Marshal(test.want)
			if string(gotJson) != string(wantJson) {
				t.Errorf("actor = %s, want %s", gotJson, wantJson)
			}
		})
	}
}

func TestGetReceiptHistoryReadsTheUploadsParentTask(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := models.Receipt{Name: "Scanned", GroupId: 1, PaidByUserID: 1, Date: time.Now()}
	repositories.GetDB().Create(&receipt)
	snapshot, _ := receipt.ToString()

	quickScan := createHistoryTask(t, historyTaskOptions{taskType: models.QUICK_SCAN, receiptId: receipt.ID, ranByUserId: historyUint(2)})
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: receipt.ID, description: snapshot, associatedSystemTaskId: &quickScan.ID})

	history := getHistory(t, 1, receipt.ID)
	if len(history.Entries) != 1 {
		t.Fatalf("entries = %+v, want the upload alone (the QUICK_SCAN parent is not an entry)", history.Entries)
	}
	actor := history.Entries[0].Actor
	if actor.Type != structs.ReceiptHistoryActorQuickScan || actor.UserId == nil || *actor.UserId != 2 {
		t.Errorf("actor = %+v, want the quick scan by user 2", actor)
	}
}

func TestGetReceiptHistoryBeforeStates(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	db := repositories.GetDB()
	rebuiltReceipt := models.Receipt{Name: "Rebuilt", GroupId: 1, PaidByUserID: 1, Date: time.Now()}
	incompleteReceipt := models.Receipt{Name: "Incomplete", GroupId: 1, PaidByUserID: 1, Date: time.Now()}
	db.Create(&rebuiltReceipt)
	db.Create(&incompleteReceipt)

	snapshot := func(receipt models.Receipt, name string) string {
		receipt.Name = name
		value, _ := receipt.ToString()
		return value
	}

	created := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: rebuiltReceipt.ID, description: snapshot(rebuiltReceipt, "Created"), ranByUserId: historyUint(1)})
	// A version 1 row's own "before" is shallow; here it claims the name was
	// "Shallow", which the rebuild must not show.
	versionOne := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: rebuiltReceipt.ID, ranByUserId: historyUint(1),
		description: updateDescription(snapshot(rebuiltReceipt, "Shallow"), snapshot(rebuiltReceipt, "One"), 0)})
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: rebuiltReceipt.ID, ranByUserId: historyUint(1),
		description: updateDescription(snapshot(rebuiltReceipt, "One"), snapshot(rebuiltReceipt, "Two"), 2)})
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: incompleteReceipt.ID, ranByUserId: historyUint(1),
		description: updateDescription(snapshot(incompleteReceipt, "Shallow"), snapshot(incompleteReceipt, "One"), 0)})

	history := getHistory(t, 1, rebuiltReceipt.ID)
	if len(history.Entries) != 3 {
		t.Fatalf("entries = %+v, want 3", history.Entries)
	}

	complete, rebuilt := history.Entries[0], history.Entries[1]
	if complete.BeforeState != structs.ReceiptHistoryBeforeStateComplete || complete.BeforeSourceType != "" || complete.BeforeSourceRecordedAt != nil {
		t.Errorf("version 2 entry = %+v, want COMPLETE with no source", complete)
	}
	if rebuilt.SystemTaskId != versionOne.ID || rebuilt.BeforeState != structs.ReceiptHistoryBeforeStateRebuilt ||
		rebuilt.BeforeSourceType != string(models.RECEIPT_UPLOADED) || rebuilt.BeforeSourceRecordedAt == nil ||
		rebuilt.BeforeSourceRecordedAt.Sub(*created.EndedAt).Abs() > time.Millisecond {
		t.Errorf("version 1 entry = %+v, want REBUILT from the creation copy", rebuilt)
	}
	if len(rebuilt.Changes) != 1 || rebuilt.Changes[0].From.Value != "Created" || rebuilt.Changes[0].To.Value != "One" {
		t.Errorf("rebuilt changes = %+v, want the name Created -> One", rebuilt.Changes)
	}

	incomplete := getHistory(t, 1, incompleteReceipt.ID).Entries[0]
	if incomplete.BeforeState != structs.ReceiptHistoryBeforeStateIncomplete || incomplete.BeforeSourceType != "" ||
		len(incomplete.Changes) != 1 || incomplete.Changes[0].From.Value != "Shallow" {
		t.Errorf("entry with no earlier copy = %+v, want INCOMPLETE against its own before", incomplete)
	}
}

func TestGetReceiptHistoryIncludesFailedRowsLinkedOnlyByAssociatedEntity(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := createHistoryReceipt(t, 1, commands.UpsertReceiptCommand{GroupId: 1, PaidByUserID: 1})
	createHistoryTask(t, historyTaskOptions{
		taskType: models.RECEIPT_UPDATED, status: models.SYSTEM_TASK_FAILED, receiptId: receipt.ID,
		description: "record not found", associatedEntityOnly: true, ranByUserId: historyUint(2),
	})

	history := getHistory(t, 1, receipt.ID)
	assertStrings(t, "actions", historyActions(history), []string{structs.ReceiptHistoryActionFailed, structs.ReceiptHistoryActionCreated})
	if history.Entries[0].ErrorMessage != "record not found" {
		t.Errorf("error message = %q", history.Entries[0].ErrorMessage)
	}
}

func TestGetReceiptHistoryDropsAnEarlierReceiptWithTheSameId(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := models.Receipt{Name: "Current", GroupId: 1, PaidByUserID: 1, Date: time.Now()}
	repositories.GetDB().Create(&receipt)
	snapshot := func(name string) string {
		copied := receipt
		copied.Name = name
		value, _ := copied.ToString()
		return value
	}

	// A deleted receipt's rows survive it, and its id was handed out again.
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: receipt.ID, description: snapshot("Old"), ranByUserId: historyUint(1)})
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: receipt.ID, description: updateDescription(snapshot("Old"), snapshot("Older"), 2), ranByUserId: historyUint(1)})
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, status: models.SYSTEM_TASK_FAILED, receiptId: receipt.ID, description: "old failure", ranByUserId: historyUint(1)})
	created := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: receipt.ID, description: snapshot("New"), ranByUserId: historyUint(2)})
	updated := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: receipt.ID, description: updateDescription(snapshot("New"), snapshot("Current"), 2), ranByUserId: historyUint(2)})

	history := getHistory(t, 1, receipt.ID)
	if len(history.Entries) != 2 || history.Entries[0].SystemTaskId != updated.ID || history.Entries[1].SystemTaskId != created.ID {
		t.Errorf("entries = %+v, want only the current receipt's two", history.Entries)
	}
}

// An isolated plain member reads the history the way they read the receipt: a
// peer's changes vanish, a peer's share and comment are masked, and an email
// receipt's origin is hidden along with its creator.
func TestGetReceiptHistoryRedactsForAnIsolatedMember(t *testing.T) {
	defer resetHistoryTest()
	resetIsolationCaches()
	fx := seedIsolatedReceiptGroup(t)

	receipt := createHistoryReceipt(t, 0, commands.UpsertReceiptCommand{
		GroupId:         fx.groupId,
		PaidByUserID:    fx.memberAId,
		CreatedByString: constants.EMAIL_INTEGRATION_CREATED_BY,
		EmailInbox:      "receipts@example.com",
		EmailSender:     "store@example.com",
	})
	recordHistoryChange(t, fx.supervisorId, receipt.ID, func(db *gorm.DB) error {
		if err := db.Create(&models.Item{Name: "Snacks", Amount: decimal.RequireFromString("5"), ChargedToUserId: &fx.memberBId, ReceiptId: receipt.ID, Status: models.ITEM_OPEN}).Error; err != nil {
			return err
		}
		if err := db.Create(&models.Comment{Comment: "from B", ReceiptId: receipt.ID, UserId: &fx.memberBId}).Error; err != nil {
			return err
		}
		return db.Create(&models.Comment{Comment: "from the supervisor", ReceiptId: receipt.ID, UserId: &fx.supervisorId}).Error
	})
	recordHistoryChange(t, fx.memberBId, receipt.ID, renameReceipt(receipt.ID, "Renamed by B"))

	t.Run("the supervisor sees everything", func(t *testing.T) {
		history := getHistory(t, fx.supervisorId, receipt.ID)
		assertStrings(t, "actions", historyActions(history), []string{
			structs.ReceiptHistoryActionUpdated, structs.ReceiptHistoryActionUpdated, structs.ReceiptHistoryActionCreated,
		})
		assertStrings(t, "supervisor's change", changeFields(history.Entries[1]), []string{"SHARE ADDED", "COMMENT ADDED", "COMMENT ADDED"})
		creation := history.Entries[2]
		if creation.Actor.Type != structs.ReceiptHistoryActorEmail || creation.Changes[0].Field != structs.ReceiptHistoryFieldEmailInbox ||
			creation.Changes[1].Field != structs.ReceiptHistoryFieldEmailSender || creation.Changes[1].To.Value != "store@example.com" {
			t.Errorf("creation = %+v, want an email creation listing its inbox and sender first", creation)
		}
		if history.ContributorCount != 3 {
			t.Errorf("contributors = %d, want 3 (email, supervisor, B)", history.ContributorCount)
		}
	})

	t.Run("a plain member does not see a peer", func(t *testing.T) {
		history := getHistory(t, fx.memberAId, receipt.ID)
		assertStrings(t, "actions", historyActions(history), []string{structs.ReceiptHistoryActionUpdated, structs.ReceiptHistoryActionCreated})

		// B's share is masked to an item, and B's comment is gone.
		supervisorChange := history.Entries[0]
		assertStrings(t, "supervisor's change", changeFields(supervisorChange), []string{"ITEM ADDED", "COMMENT ADDED"})
		if supervisorChange.Changes[1].To.Value != "from the supervisor" {
			t.Errorf("comment = %+v, want the supervisor's", supervisorChange.Changes[1])
		}

		creation := history.Entries[1]
		if creation.Actor.Type != structs.ReceiptHistoryActorEmail {
			t.Errorf("creation actor = %+v, want EMAIL", creation.Actor)
		}
		for _, change := range creation.Changes {
			if change.Field == structs.ReceiptHistoryFieldEmailInbox || change.Field == structs.ReceiptHistoryFieldEmailSender {
				t.Errorf("creation shows %s to a member who cannot see its creator", change.Field)
			}
		}
		if history.ContributorCount != 2 {
			t.Errorf("contributors = %d, want 2 (email, supervisor)", history.ContributorCount)
		}
	})

	t.Run("the snapshots are redacted, never the stored description", func(t *testing.T) {
		history := getHistory(t, fx.memberAId, receipt.ID)
		service := NewReceiptHistoryService(nil)

		snapshots, err := service.GetReceiptHistorySnapshots(fx.memberAId, receipt.ID, history.Entries[0].SystemTaskId)
		if err != nil {
			t.Fatalf("GetReceiptHistorySnapshots: %v", err)
		}
		var after models.Receipt
		if err := json.Unmarshal([]byte(snapshots.After), &after); err != nil {
			t.Fatalf("read after: %v", err)
		}
		if len(after.Comments) != 1 || after.Comments[0].Comment != "from the supervisor" ||
			len(after.ReceiptItems) != 1 || after.ReceiptItems[0].ChargedToUserId != nil ||
			after.EmailInbox != "" || after.EmailSender != "" || after.GroupId != fx.groupId {
			t.Errorf("after = %s, want B's comment and charge masked and the email origin hidden", snapshots.After)
		}
		if snapshots.Before == "" {
			t.Error("an update's snapshots have no before")
		}

		creation, err := service.GetReceiptHistorySnapshots(fx.memberAId, receipt.ID, history.Entries[1].SystemTaskId)
		if err != nil || creation.Before != "" || creation.After == "" {
			t.Errorf("creation snapshots = %+v, %v; want an after alone", creation, err)
		}
	})

	t.Run("a hidden, unknown or foreign entry is denied", func(t *testing.T) {
		var bUpdate models.SystemTask
		repositories.GetDB().Where("type = ? AND ran_by_user_id = ?", models.RECEIPT_UPDATED, fx.memberBId).First(&bUpdate)
		otherReceipt := createHistoryReceipt(t, fx.memberAId, commands.UpsertReceiptCommand{GroupId: fx.groupId, PaidByUserID: fx.memberAId})
		var otherTask models.SystemTask
		repositories.GetDB().Where("receipt_id = ?", otherReceipt.ID).First(&otherTask)

		for name, taskId := range map[string]uint{"hidden": bUpdate.ID, "unknown": 9999, "zero": 0, "foreign": otherTask.ID} {
			_, err := NewReceiptHistoryService(nil).GetReceiptHistorySnapshots(fx.memberAId, receipt.ID, taskId)
			if !errors.Is(err, ErrReceiptAccessDenied) {
				t.Errorf("%s entry: err = %v, want ErrReceiptAccessDenied", name, err)
			}
		}
	})
}

func TestGetReceiptHistorySnapshotsDeniesAFailedEntry(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := createHistoryReceipt(t, 1, commands.UpsertReceiptCommand{GroupId: 1, PaidByUserID: 1})
	repositories.NewReceiptChangeRecorder(1, receipt.ID).RecordFailure(errors.New("failed"))
	history := getHistory(t, 1, receipt.ID)

	_, err := NewReceiptHistoryService(nil).GetReceiptHistorySnapshots(1, receipt.ID, history.Entries[0].SystemTaskId)
	if !errors.Is(err, ErrReceiptAccessDenied) {
		t.Errorf("err = %v, want ErrReceiptAccessDenied for an entry with no snapshots", err)
	}
}

func TestGetReceiptHistoryStripsRestrictedCategories(t *testing.T) {
	defer resetHistoryTest()

	visible := makeCategory(t, "Visible")
	hidden := makeCategory(t, "Hidden")
	viewerId, groupId, _ := seedMemberWithGroupRoleGrants(t, "history-viewer", []uint{visible}, nil)

	receipt := createHistoryReceipt(t, viewerId, commands.UpsertReceiptCommand{GroupId: groupId, PaidByUserID: viewerId})
	recordHistoryChange(t, viewerId, receipt.ID, func(db *gorm.DB) error {
		return db.Exec("INSERT INTO receipt_categories (receipt_id, category_id) VALUES (?, ?), (?, ?)",
			receipt.ID, visible, receipt.ID, hidden).Error
	})

	update := getHistory(t, viewerId, receipt.ID).Entries[0]
	if len(update.Changes) != 1 || update.Changes[0].To.Value != "Visible" {
		t.Errorf("changes = %+v, want the visible category alone", update.Changes)
	}
}

func TestGetReceiptHistoryDropsAPastStateWithAHiddenPayer(t *testing.T) {
	defer resetHistoryTest()

	// The viewer's role shows only receipts they paid for.
	viewerId, groupId, _ := seedMemberWithPaidByRole(t, "history-payer", nil, true)
	otherId := makeUser(t, "history-other-payer")
	seedIsoMember(t, groupId, otherId, nil)

	receipt := createHistoryReceipt(t, otherId, commands.UpsertReceiptCommand{GroupId: groupId, PaidByUserID: otherId})
	recordHistoryChange(t, otherId, receipt.ID, func(db *gorm.DB) error {
		return db.Model(&models.Receipt{}).Where("id = ?", receipt.ID).Update("paid_by_user_id", viewerId).Error
	})
	recordHistoryChange(t, otherId, receipt.ID, renameReceipt(receipt.ID, "Now mine"))

	history := getHistory(t, viewerId, receipt.ID)
	if len(history.Entries) != 1 || len(history.Entries[0].Changes) != 1 || history.Entries[0].Changes[0].Field != structs.ReceiptHistoryFieldName {
		t.Errorf("entries = %+v, want only the rename, made while the viewer paid", history.Entries)
	}
}

func TestGetReceiptHistoryNamesOnlyGroupsTheViewerMaySee(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()
	db := repositories.GetDB()

	// User 1 is in group 1 only; the receipt starts in group 2.
	receipt := createHistoryReceipt(t, 4, commands.UpsertReceiptCommand{GroupId: 2, PaidByUserID: 1})
	recordHistoryChange(t, 1, receipt.ID, func(db *gorm.DB) error {
		return db.Model(&models.Receipt{}).Where("id = ?", receipt.ID).Update("group_id", 1).Error
	})

	move := getHistory(t, 1, receipt.ID).Entries[0]
	if len(move.Changes) != 1 || move.Changes[0].Field != structs.ReceiptHistoryFieldGroup ||
		move.Changes[0].From.Value != restrictedGroupName || move.Changes[0].To.Value != "test" {
		t.Errorf("changes = %+v, want (Restricted) -> test", move.Changes)
	}

	role, err := repositories.NewRoleRepository(nil).CreateAppRole("Groups Reader", "", []string{permissions.AppGroupsRead}, false)
	if err != nil {
		t.Fatalf("create app role: %v", err)
	}
	db.Model(&models.User{}).Where("id = ?", 1).Update("app_role_id", role.ID)
	clearRolePermissionCacheAll()

	move = getHistory(t, 1, receipt.ID).Entries[0]
	if move.Changes[0].From.Value != "test2" {
		t.Errorf("from = %+v, want test2 for a holder of app.groups.read", move.Changes[0].From)
	}
}

// The email task records its RECEIPT_UPLOADED row through
// CreateReceiptUploadedSystemTask, under its EMAIL_UPLOAD. The description is
// the created receipt, so it carries the email origin, and the creation entry
// lists it.
func TestEmailUploadRowCarriesTheEmailOriginIntoTheHistory(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := models.Receipt{
		Name: "From email", GroupId: 1, PaidByUserID: 1, Date: time.Now(), Status: models.OPEN,
		EmailInbox: "receipts@example.com", EmailSender: "store@example.com",
	}
	repositories.GetDB().Create(&receipt)

	emailUpload := createHistoryTask(t, historyTaskOptions{taskType: models.EMAIL_UPLOAD, receiptId: receipt.ID})
	task, err := NewSystemTaskService(nil).CreateReceiptUploadedSystemTask(nil, receipt, structs.ReceiptProcessingSystemTasks{SystemTask: emailUpload}, time.Now())
	if err != nil {
		t.Fatalf("CreateReceiptUploadedSystemTask: %v", err)
	}
	var described models.Receipt
	if err := json.Unmarshal([]byte(task.ResultDescription), &described); err != nil ||
		described.EmailInbox != "receipts@example.com" || described.EmailSender != "store@example.com" {
		t.Errorf("description = %s, want the email origin", task.ResultDescription)
	}

	creation := getHistory(t, 1, receipt.ID).Entries[0]
	if creation.Actor.Type != structs.ReceiptHistoryActorEmail {
		t.Errorf("actor = %+v, want EMAIL", creation.Actor)
	}
	assertStrings(t, "first changes", changeFields(creation)[:2], []string{"EMAIL_INBOX ADDED", "EMAIL_SENDER ADDED"})
}

// A receipt moved into an isolated group is judged there for its whole
// history: a peer's comment made while it sat in an open group is hidden from
// an isolated member as it is on the receipt today.
func TestGetReceiptHistoryJudgesEverySnapshotInTheCurrentGroup(t *testing.T) {
	defer resetHistoryTest()
	resetIsolationCaches()
	fx := seedIsolatedReceiptGroup(t)

	openGroup := models.Group{Name: "open-group"}
	repositories.GetDB().Create(&openGroup)

	receipt := createHistoryReceipt(t, fx.supervisorId, commands.UpsertReceiptCommand{GroupId: openGroup.ID, PaidByUserID: fx.memberAId})
	recordHistoryChange(t, fx.supervisorId, receipt.ID, func(db *gorm.DB) error {
		return db.Create(&models.Comment{Comment: "from B", ReceiptId: receipt.ID, UserId: &fx.memberBId}).Error
	})
	recordHistoryChange(t, fx.supervisorId, receipt.ID, func(db *gorm.DB) error {
		return db.Model(&models.Receipt{}).Where("id = ?", receipt.ID).Update("group_id", fx.groupId).Error
	})

	history := getHistory(t, fx.memberAId, receipt.ID)
	if len(history.Entries) != 3 {
		t.Fatalf("entries = %+v, want 3", history.Entries)
	}
	for _, entry := range history.Entries {
		for _, change := range entry.Changes {
			if change.Field == structs.ReceiptHistoryFieldComment {
				t.Errorf("entry %d shows B's comment to an isolated member: %+v", entry.SystemTaskId, change)
			}
		}
	}

	// The move names the open group only to its members.
	move := history.Entries[0]
	if len(move.Changes) != 1 || move.Changes[0].From.Value != restrictedGroupName || move.Changes[0].To.Value != "iso-receipt-group" {
		t.Errorf("move = %+v, want (Restricted) -> iso-receipt-group", move.Changes)
	}
}

// A quick scan by someone the caller cannot see keeps its entry, since the
// receipt itself is visible, but no longer names the scanner.
func TestGetReceiptHistoryHidesAQuickScannerTheViewerCannotSee(t *testing.T) {
	defer resetHistoryTest()
	resetIsolationCaches()
	fx := seedIsolatedReceiptGroup(t)

	receipt := models.Receipt{Name: "Scanned", GroupId: fx.groupId, PaidByUserID: fx.memberAId, Date: time.Now()}
	repositories.GetDB().Create(&receipt)
	snapshot, _ := receipt.ToString()
	quickScan := createHistoryTask(t, historyTaskOptions{taskType: models.QUICK_SCAN, receiptId: receipt.ID, ranByUserId: &fx.memberBId})
	createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: receipt.ID, description: snapshot, associatedSystemTaskId: &quickScan.ID})

	for viewer, wantScanner := range map[uint]bool{fx.memberAId: false, fx.supervisorId: true} {
		history := getHistory(t, viewer, receipt.ID)
		if len(history.Entries) != 1 {
			t.Fatalf("viewer %d: entries = %+v, want the scan", viewer, history.Entries)
		}
		actor := history.Entries[0].Actor
		if actor.Type != structs.ReceiptHistoryActorQuickScan || (actor.UserId != nil) != wantScanner {
			t.Errorf("viewer %d: actor = %+v, want a quick scan that names the scanner: %v", viewer, actor, wantScanner)
		}
	}
}

// A row whose description cannot be read still shows that a save happened, with
// nothing to compare and no claim about its "before".
func TestGetReceiptHistoryKeepsAnUnreadableRowWithoutSnapshots(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	receipt := createHistoryReceipt(t, 1, commands.UpsertReceiptCommand{GroupId: 1, PaidByUserID: 1})
	createHistoryTask(t, historyTaskOptions{
		taskType: models.RECEIPT_UPDATED, receiptId: receipt.ID, ranByUserId: historyUint(1),
		description: updateDescription("not json", "not json either", 2),
	})

	history := getHistory(t, 1, receipt.ID)
	unreadable := history.Entries[0]
	if unreadable.Action != structs.ReceiptHistoryActionUpdated || unreadable.HasSnapshots ||
		unreadable.BeforeState != "" || unreadable.Changes == nil || len(unreadable.Changes) != 0 {
		t.Errorf("entry = %+v, want an UPDATED entry with no snapshots, no before state and no changes", unreadable)
	}
	if history.SaveCount != 2 {
		t.Errorf("saves = %d, want 2", history.SaveCount)
	}

	_, err := NewReceiptHistoryService(nil).GetReceiptHistorySnapshots(1, receipt.ID, unreadable.SystemTaskId)
	if !errors.Is(err, ErrReceiptAccessDenied) {
		t.Errorf("snapshots err = %v, want ErrReceiptAccessDenied", err)
	}
}

// A version 1 "before" stored custom field values without their definitions.
// When neither side carries one, the stored definition names the field.
func TestGetReceiptHistoryReadsAValueWithoutADefinitionAgainstTheStoredOne(t *testing.T) {
	defer resetHistoryTest()
	repositories.CreateTestGroupWithUsers()

	customField := models.CustomField{Name: "PO Number", Type: models.TEXT}
	repositories.GetDB().Create(&customField)
	receipt := models.Receipt{Name: "Old", GroupId: 1, PaidByUserID: 1, Date: time.Now()}
	repositories.GetDB().Create(&receipt)

	before := receipt
	poNumber := "A-1"
	before.CustomFields = []models.CustomFieldValue{{CustomFieldId: customField.ID, StringValue: &poNumber}}
	beforeJson, _ := before.ToString()
	afterJson, _ := receipt.ToString()
	createHistoryTask(t, historyTaskOptions{
		taskType: models.RECEIPT_UPDATED, receiptId: receipt.ID, ranByUserId: historyUint(1),
		description: updateDescription(beforeJson, afterJson, 0),
	})

	entry := getHistory(t, 1, receipt.ID).Entries[0]
	if entry.BeforeState != structs.ReceiptHistoryBeforeStateIncomplete || len(entry.Changes) != 1 ||
		entry.Changes[0].Label != "PO Number" || entry.Changes[0].Kind != structs.ReceiptHistoryKindRemoved ||
		entry.Changes[0].From.Value != "A-1" {
		t.Errorf("entry = %+v, changes %+v; want PO Number A-1 removed", entry, entry.Changes)
	}
}
