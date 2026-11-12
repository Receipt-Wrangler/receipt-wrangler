import { CurrencyPipe } from "@angular/common";
import { HttpErrorResponse } from "@angular/common/http";
import { Component, Input, NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { FormMode } from "../../enums/form-mode.enum";
import { Claim, Group, Permission, Receipt, ReceiptService, ReimbursementClaimService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { ReceiptExportService } from "../../services/receipt-export.service";
import { SnackbarService } from "../../services/snackbar.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { AuthState, GroupState, SetGroups, SetPermissions, SetSelectedGroupId, UserState } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { CLAIM_NOT_FOUND_MESSAGE } from "../claim.resolver";
import { ClaimLedgerEntryDialogComponent } from "../claim-ledger-entry-dialog/claim-ledger-entry-dialog.component";
import { ClaimReceiptPickerDialogComponent } from "../claim-receipt-picker-dialog/claim-receipt-picker-dialog.component";
import { ClaimFormComponent } from "./claim-form.component";

// `onlyIcon` starts with "on", which Angular refuses to bind as a property on an
// unknown element; role-form.component.spec.ts stubs it the same way.
@Component({ selector: "app-submit-button", template: "", standalone: false })
class StubSubmitButtonComponent {
  @Input() onlyIcon = false;
  @Input() buttonText = "";
  @Input() disabled = false;
}

@Component({ selector: "app-status-chip", template: "{{ customStatus }}", standalone: false })
class StubStatusChipComponent {
  @Input() customStatus = "";
  @Input() customStatusColor?: string;
  @Input() status = "";
}

describe("ClaimFormComponent", () => {
  let fixture: ComponentFixture<ClaimFormComponent>;
  let component: ClaimFormComponent;
  let store: Store;

  const claimService = {
    createClaim: jest.fn(),
    updateClaim: jest.fn(),
    getClaimById: jest.fn(),
    addReceiptsToClaim: jest.fn(),
    removeReceiptsFromClaim: jest.fn(),
    deleteClaim: jest.fn(),
    deleteClaimLedgerEntry: jest.fn(),
  };
  const receiptService = { getReceiptsForGroup: jest.fn() };
  const receiptExportService = { exportReceiptsById: jest.fn() };
  const matDialog = { open: jest.fn() };
  const router = { navigate: jest.fn() };
  const snackbarService = { success: jest.fn(), error: jest.fn() };

  const household = { id: 5, name: "Household", isAllGroup: false } as Group;
  const work = { id: 6, name: "Work", isAllGroup: false } as Group;
  const all = { id: 9, name: "All", isAllGroup: true } as Group;

  const taxi = { id: 1, name: "Taxi", groupId: 5, amount: "40.00", date: "2026-09-02T00:00:00Z" } as Receipt;
  const hotel = { id: 2, name: "Hotel", groupId: 5, amount: "372.50", date: "2026-09-03T00:00:00Z" } as Receipt;
  const groceries = { id: 3, name: "Groceries", groupId: 5, amount: "80.00", date: "2026-09-04T00:00:00Z" } as Receipt;

  const recoupment = {
    id: 8,
    claimId: 12,
    type: "RECOUPMENT",
    date: "2026-10-20",
    amount: "40.00",
    createdAt: "2026-10-20T09:00:00.000-04:00",
    updatedAt: "2026-10-20T09:00:00.000-04:00",
  };
  const payment = { ...recoupment, id: 7, type: "PAYMENT", date: "2026-10-01", amount: "300.00" };

  const claim = (overrides: Partial<Claim> = {}): Claim => ({
    id: 12,
    groupId: 5,
    name: "Voucher 14",
    referenceNumber: "TA-14",
    type: "Travel voucher",
    status: "PARTIALLY_PAID",
    awaitingSettlement: true,
    submittedDate: "2026-09-30",
    dueDate: "2026-10-31",
    createdAt: "2026-09-30T10:00:00.000-04:00",
    updatedAt: "2026-10-20T09:00:00.000-04:00",
    // 99 is on the claim but not in the group list the viewer got back.
    receiptIds: [1, 2, 99],
    receiptCount: 3,
    receiptsTotal: "412.50",
    effectiveClaimedAmount: "412.50",
    paymentsTotal: "300.00",
    recoupmentsTotal: "40.00",
    paidAmount: "260.00",
    outstandingAmount: "152.50",
    daysWaiting: 42,
    ledgerEntries: [payment, recoupment],
    ...overrides,
  });

  const setUp = async (
    mode: FormMode,
    options: { claim?: Claim; permissions?: { [groupId: number]: string[] }; groups?: Group[]; selected?: string } = {}
  ) => {
    await TestBed.configureTestingModule({
      declarations: [ClaimFormComponent, StubSubmitButtonComponent, StubStatusChipComponent],
      imports: [
        PipesModule,
        ReactiveFormsModule,
        NgxsModule.forRoot([AuthState, GroupState, SystemSettingsState, UserState]),
      ],
      providers: [
        provideZonelessChangeDetection(),
        CurrencyPipe,
        { provide: ActivatedRoute, useValue: { snapshot: { data: { mode, claim: options.claim } } } },
        { provide: ReimbursementClaimService, useValue: claimService },
        { provide: ReceiptService, useValue: receiptService },
        { provide: ReceiptExportService, useValue: receiptExportService },
        { provide: MatDialog, useValue: matDialog },
        { provide: Router, useValue: router },
        { provide: SnackbarService, useValue: snackbarService },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    store = TestBed.inject(Store);
    store.dispatch(new SetGroups(options.groups ?? [all, household, work]));
    store.dispatch(new SetSelectedGroupId(options.selected ?? "5"));
    store.dispatch(new SetPermissions([], options.permissions ?? { 5: [Permission.GroupReceiptsUpdate] }));

    fixture = TestBed.createComponent(ClaimFormComponent);
    component = fixture.componentInstance;
    fixture.autoDetectChanges();
    await fixture.whenStable();
  };

  const stubDialog = (result: unknown) =>
    matDialog.open.mockReturnValue({ componentInstance: {}, afterClosed: () => of(result) });

  const text = (testId: string): string =>
    fixture.nativeElement.querySelector(`[data-testid='${testId}']`)?.textContent?.replace(/\s+/g, " ").trim();

  beforeEach(() => {
    jest.resetAllMocks();
    receiptService.getReceiptsForGroup.mockReturnValue(of({ totalCount: 3, data: [taxi, hotel, groceries] }));
  });

  describe("adding", () => {
    it("creates the claim in the group being browsed, then opens it", async () => {
      await setUp(FormMode.add);
      claimService.createClaim.mockReturnValue(of(claim()));
      component.form.patchValue({
        name: "Voucher 14",
        status: "SUBMITTED",
        submittedDate: new Date(2026, 8, 30),
        claimedAmount: "",
      });

      component.submit();

      expect(claimService.createClaim).toHaveBeenCalledWith(5, {
        name: "Voucher 14",
        status: "SUBMITTED",
        submittedDate: "2026-09-30",
      });
      expect(snackbarService.success).toHaveBeenCalledWith("Claim added");
      expect(router.navigate).toHaveBeenCalledWith(["/claims", 12, "view"]);
    });

    it("defaults a new claim to Draft and requires a name", async () => {
      await setUp(FormMode.add);

      expect(component.form.get("status")?.value).toBe("DRAFT");
      component.submit();
      expect(claimService.createClaim).not.toHaveBeenCalled();
      expect(component.form.get("name")?.touched).toBe(true);
    });

    it("names the group it will be added to when the viewer has several", async () => {
      await setUp(FormMode.add);

      expect(text("claim-group")).toBe("Will be added to Household");
    });

    it("loads no receipts until the claim exists", async () => {
      await setUp(FormMode.add);

      expect(receiptService.getReceiptsForGroup).not.toHaveBeenCalled();
    });

    // Only reachable by URL: the list offers Add Claim only with a real group.
    it("leaves when there is no real group to add to", async () => {
      await setUp(FormMode.add, { selected: "9", permissions: { 9: [Permission.GroupReceiptsUpdate] } });

      expect(snackbarService.error).toHaveBeenCalledWith("Choose a group before adding a claim.");
      expect(router.navigate).toHaveBeenCalledWith(["/dashboard/group/9"]);
    });
  });

  describe("editing", () => {
    beforeEach(async () => {
      await setUp(FormMode.edit, { claim: claim() });
    });

    it("loads the claim's fields, with its calendar days as local midnights of the same days", () => {
      const submitted = component.form.get("submittedDate")?.value as Date;

      expect(component.form.get("name")?.value).toBe("Voucher 14");
      expect([submitted.getFullYear(), submitted.getMonth(), submitted.getDate()]).toEqual([2026, 8, 30]);
      expect(component.headerText()).toBe("Edit Voucher 14");
    });

    it("replaces the claim's fields, clearing what was emptied", () => {
      claimService.updateClaim.mockReturnValue(of(claim()));
      component.form.patchValue({ dueDate: null, referenceNumber: "" });

      component.submit();

      expect(claimService.updateClaim).toHaveBeenCalledWith(12, {
        name: "Voucher 14",
        type: "Travel voucher",
        status: "PARTIALLY_PAID",
        submittedDate: "2026-09-30",
      });
      expect(router.navigate).toHaveBeenCalledWith(["/claims", 12, "view"]);
    });

    it("puts the server's field errors on the form", () => {
      claimService.updateClaim.mockReturnValue(
        throwError(() => new HttpErrorResponse({ status: 400, error: { dueDate: "Due date is not a valid day" } }))
      );

      component.submit();

      expect(component.form.get("dueDate")?.errors).toEqual({ error: "Due date is not a valid day" });
      expect(router.navigate).not.toHaveBeenCalled();
    });
  });

  describe("viewing", () => {
    beforeEach(async () => {
      await setUp(FormMode.view, { claim: claim() });
    });

    it("shows the server's figures", () => {
      expect(text("claim-figure-claimed")).toContain("$412.50");
      expect(text("claim-figure-paid")).toContain("$260.00");
      expect(text("claim-figure-paid")).toContain("Payments $300.00 − recouped $40.00");
      expect(text("claim-figure-outstanding")).toContain("$152.50");
      expect(text("claim-figure-waiting")).toContain("42 days");
      expect(text("claim-status")).toBe("Partially Paid");
    });

    it("never submits", () => {
      component.submit();

      expect(claimService.updateClaim).not.toHaveBeenCalled();
    });

    it("loads the claim group's receipts once, unpaged", () => {
      expect(receiptService.getReceiptsForGroup).toHaveBeenCalledTimes(1);
      expect(receiptService.getReceiptsForGroup).toHaveBeenCalledWith(5, {
        page: 1,
        pageSize: -1,
        orderBy: "date",
        sortDirection: "desc",
      });
    });

    it("shows the claim's receipts, naming one the list did not return by its id", () => {
      expect(component.claimReceipts().map((receipt) => receipt.name)).toEqual(["Taxi", "Hotel", "Receipt #99"]);
    });

    it("offers edit, delete and the row actions to a viewer who can update", () => {
      expect(component.canUpdate()).toBe(true);
      expect(component.receiptDisplayedColumns()).toContain("actions");
      expect(component.ledgerDisplayedColumns()).toContain("actions");
    });

    it("refetches the claim after a change without reloading the group's receipts", async () => {
      claimService.getClaimById.mockReturnValue(of(claim({ outstandingAmount: "0.00" })));

      component.refresh();
      await fixture.whenStable();

      expect(component.claim()?.outstandingAmount).toBe("0.00");
      expect(text("claim-figure-outstanding")).toContain("$0.00");
      expect(receiptService.getReceiptsForGroup).toHaveBeenCalledTimes(1);
    });

    it("leaves when the claim has gone in the meantime", () => {
      claimService.getClaimById.mockReturnValue(throwError(() => new HttpErrorResponse({ status: 403 })));

      component.refresh();

      expect(snackbarService.error).toHaveBeenCalledWith(CLAIM_NOT_FOUND_MESSAGE);
      expect(router.navigate).toHaveBeenCalledWith(["/claims/group", "5"]);
    });

    it("offers only the group's receipts not already on the claim, then adds the picked ones", () => {
      claimService.addReceiptsToClaim.mockReturnValue(of({ receiptIds: [3], failures: [] }));
      claimService.getClaimById.mockReturnValue(of(claim()));
      stubDialog([groceries]);

      component.addReceipts();

      expect(matDialog.open).toHaveBeenCalledWith(ClaimReceiptPickerDialogComponent, {
        ...DEFAULT_DIALOG_CONFIG,
        data: { claimName: "Voucher 14", candidates: [groceries] },
      });
      expect(claimService.addReceiptsToClaim).toHaveBeenCalledWith(12, { receiptIds: [3] });
      expect(snackbarService.success).toHaveBeenCalledWith("Added 1 receipt to Voucher 14");
      expect(claimService.getClaimById).toHaveBeenCalledWith(12);
    });

    it("adds nothing when the picker is cancelled", () => {
      stubDialog(undefined);

      component.addReceipts();

      expect(claimService.addReceiptsToClaim).not.toHaveBeenCalled();
    });

    it("removes a receipt without asking, then refetches", () => {
      claimService.removeReceiptsFromClaim.mockReturnValue(of({ receiptIds: [1], failures: [] }));
      claimService.getClaimById.mockReturnValue(of(claim()));

      component.removeReceipt(taxi);

      expect(matDialog.open).not.toHaveBeenCalled();
      expect(claimService.removeReceiptsFromClaim).toHaveBeenCalledWith(12, { receiptIds: [1] });
      expect(snackbarService.success).toHaveBeenCalledWith("Removed 1 receipt from Voucher 14");
      expect(claimService.getClaimById).toHaveBeenCalledWith(12);
    });

    it("refetches after a ledger entry is saved, and not after a cancel", () => {
      claimService.getClaimById.mockReturnValue(of(claim()));
      stubDialog(false);

      component.editLedgerEntry(payment);
      expect(matDialog.open).toHaveBeenCalledWith(ClaimLedgerEntryDialogComponent, {
        ...DEFAULT_DIALOG_CONFIG,
        data: { claimId: 12, entry: payment },
      });
      expect(claimService.getClaimById).not.toHaveBeenCalled();

      stubDialog(true);
      component.addLedgerEntry();
      expect(claimService.getClaimById).toHaveBeenCalledWith(12);
    });

    it("deletes a ledger entry once confirmed", () => {
      claimService.deleteClaimLedgerEntry.mockReturnValue(of({}));
      claimService.getClaimById.mockReturnValue(of(claim()));
      stubDialog(true);

      component.deleteLedgerEntry(recoupment);

      expect(matDialog.open).toHaveBeenCalledWith(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);
      expect(claimService.deleteClaimLedgerEntry).toHaveBeenCalledWith(12, 8);
      expect(claimService.getClaimById).toHaveBeenCalled();
    });

    it("keeps the ledger entry when the confirmation is dismissed", () => {
      stubDialog(undefined);

      component.deleteLedgerEntry(recoupment);

      expect(claimService.deleteClaimLedgerEntry).not.toHaveBeenCalled();
    });

    it("deletes the claim once confirmed and returns to the list", () => {
      claimService.deleteClaim.mockReturnValue(of({}));
      stubDialog(true);

      component.deleteClaim();

      expect(claimService.deleteClaim).toHaveBeenCalledWith(12);
      expect(router.navigate).toHaveBeenCalledWith(["/claims/group", "5"]);
    });

    it("exports the claim's receipts", () => {
      component.exportReceipts();

      expect(receiptExportService.exportReceiptsById).toHaveBeenCalledWith([1, 2, 99]);
    });
  });

  // Judged in the claim's own group (5), never the group being browsed (All, 9).
  it("is read-only for a viewer who cannot update the claim's group", async () => {
    await setUp(FormMode.view, {
      claim: claim(),
      selected: "9",
      permissions: { 5: [Permission.GroupReceiptsRead], 9: [Permission.GroupReceiptsUpdate] },
    });

    expect(component.canUpdate()).toBe(false);
    expect(component.receiptDisplayedColumns()).not.toContain("actions");
    expect(component.ledgerDisplayedColumns()).not.toContain("actions");
  });

  it("does not name the group for a viewer with only one", async () => {
    await setUp(FormMode.view, { claim: claim(), groups: [all, household] });

    expect(component.groupName()).toBe("");
  });

  it("shows an unknown status as its raw code", async () => {
    await setUp(FormMode.view, { claim: claim({ status: "ESCALATED" }) });

    expect(text("claim-status")).toBe("ESCALATED");
    expect(component.statusOptions.at(-1)?.value).toBe("ESCALATED");
  });
});
