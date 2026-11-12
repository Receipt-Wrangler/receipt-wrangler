import { HttpErrorResponse } from "@angular/common/http";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { NgxsModule } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { ClaimLedgerEntry, ReimbursementClaimService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services/snackbar.service";
import { SystemSettingsState } from "../../store/system-settings.state";
import {
  ClaimLedgerEntryDialogComponent,
  ClaimLedgerEntryDialogData,
} from "./claim-ledger-entry-dialog.component";

describe("ClaimLedgerEntryDialogComponent", () => {
  let fixture: ComponentFixture<ClaimLedgerEntryDialogComponent>;
  let component: ClaimLedgerEntryDialogComponent;

  const claimService = {
    createClaimLedgerEntry: jest.fn(),
    updateClaimLedgerEntry: jest.fn(),
  };
  const matDialogRef = { close: jest.fn() };
  const snackbarService = { success: jest.fn(), error: jest.fn() };

  const recoupment: ClaimLedgerEntry = {
    id: 7,
    claimId: 12,
    type: "RECOUPMENT",
    date: "2026-10-05",
    amount: "40.00",
    note: "Per diem overpaid",
    createdAt: "2026-10-05T10:00:00.000-04:00",
    updatedAt: "2026-10-05T10:00:00.000-04:00",
  };

  const setUp = async (data: ClaimLedgerEntryDialogData) => {
    await TestBed.configureTestingModule({
      imports: [ClaimLedgerEntryDialogComponent, NgxsModule.forRoot([SystemSettingsState])],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ReimbursementClaimService, useValue: claimService },
        { provide: MatDialogRef, useValue: matDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
        { provide: SnackbarService, useValue: snackbarService },
      ],
    })
      .overrideComponent(ClaimLedgerEntryDialogComponent, {
        set: { imports: [ReactiveFormsModule, PipesModule], schemas: [CUSTOM_ELEMENTS_SCHEMA] },
      })
      .compileComponents();

    fixture = TestBed.createComponent(ClaimLedgerEntryDialogComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  };

  beforeEach(() => jest.resetAllMocks());

  describe("adding", () => {
    beforeEach(async () => {
      await setUp({ claimId: 12 });
    });

    it("defaults to a payment dated today", () => {
      const date = component.form.controls.date.value as Date;

      expect(component.headerText).toBe("Add ledger entry");
      expect(component.form.controls.type.value).toBe("PAYMENT");
      expect(date).toBeInstanceOf(Date);
      expect([date.getHours(), date.getMinutes()]).toEqual([0, 0]);
    });

    it("requires an amount for a payment, and not a note", () => {
      expect(component.form.controls.amount.hasError("required")).toBe(true);
      expect(component.form.controls.note.errors).toBeNull();
    });

    // A note's amount never counts; its text is the point.
    it("swaps the requirement to the note text for a note", () => {
      component.form.controls.type.setValue("NOTE");

      expect(component.isNote()).toBe(true);
      expect(component.form.controls.amount.errors).toBeNull();
      component.form.controls.note.setValue("  ");
      expect(component.form.controls.note.hasError("required")).toBe(true);
    });

    // app-input / app-textarea derive their messages from statusChanges, so a
    // validator swap that does not emit leaves the field red with no message.
    it("tells the inputs when the type swaps which field is required", () => {
      const noteStatuses: string[] = [];
      const amountStatuses: string[] = [];
      component.form.controls.note.statusChanges.subscribe((status) => noteStatuses.push(status));
      component.form.controls.amount.statusChanges.subscribe((status) => amountStatuses.push(status));

      component.form.controls.type.setValue("NOTE");

      expect(noteStatuses).toContain("INVALID");
      expect(amountStatuses).toContain("VALID");
    });

    it("creates the entry with a bare calendar day and closes with true", () => {
      claimService.createClaimLedgerEntry.mockReturnValue(of({}));
      component.form.setValue({
        type: "PAYMENT",
        date: new Date(2026, 9, 1),
        amount: "300.00",
        note: "EFT 4471",
      });

      component.submit();

      expect(claimService.createClaimLedgerEntry).toHaveBeenCalledWith(12, {
        type: "PAYMENT",
        date: "2026-10-01",
        amount: "300.00",
        note: "EFT 4471",
      });
      expect(snackbarService.success).toHaveBeenCalledWith("Ledger entry added");
      expect(matDialogRef.close).toHaveBeenCalledWith(true);
    });

    it("does not submit an invalid form", () => {
      component.submit();

      expect(claimService.createClaimLedgerEntry).not.toHaveBeenCalled();
      expect(component.form.controls.amount.touched).toBe(true);
    });

    it("puts the server's field errors on the form and stays open", () => {
      claimService.createClaimLedgerEntry.mockReturnValue(
        throwError(() => new HttpErrorResponse({ status: 400, error: { amount: "Amount must be greater than zero" } }))
      );
      component.form.patchValue({ amount: "0.001" });

      component.submit();

      expect(component.form.controls.amount.errors).toEqual({ error: "Amount must be greater than zero" });
      expect(matDialogRef.close).not.toHaveBeenCalled();
      expect(component.submitting()).toBe(false);
    });
  });

  describe("editing", () => {
    beforeEach(async () => {
      await setUp({ claimId: 12, entry: recoupment });
    });

    it("loads the entry, keeping its calendar day", () => {
      const date = component.form.controls.date.value as Date;

      expect(component.headerText).toBe("Edit ledger entry");
      expect(component.form.controls.type.value).toBe("RECOUPMENT");
      expect([date.getFullYear(), date.getMonth(), date.getDate()]).toEqual([2026, 9, 5]);
      expect(component.form.controls.amount.value).toBe("40.00");
    });

    it("updates the entry", () => {
      claimService.updateClaimLedgerEntry.mockReturnValue(of({}));

      component.submit();

      expect(claimService.updateClaimLedgerEntry).toHaveBeenCalledWith(12, 7, {
        type: "RECOUPMENT",
        date: "2026-10-05",
        amount: "40.00",
        note: "Per diem overpaid",
      });
      expect(snackbarService.success).toHaveBeenCalledWith("Ledger entry updated");
    });
  });
});
