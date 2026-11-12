import { HttpErrorResponse } from "@angular/common/http";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from "@angular/material/dialog";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { Claim, ClaimList, Group, Permission, ReimbursementClaimService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services/snackbar.service";
import { AuthState, GroupState, SetGroups, SetPermissions } from "../../store";
import { ClaimReceiptsResultDialogComponent } from "../claim-receipts-result-dialog/claim-receipts-result-dialog.component";
import { AddToClaimDialogComponent, AddToClaimDialogData, NEW_CLAIM_OPTION } from "./add-to-claim-dialog.component";

describe("AddToClaimDialogComponent", () => {
  let fixture: ComponentFixture<AddToClaimDialogComponent>;
  let component: AddToClaimDialogComponent;
  let store: Store;

  const claimService = {
    getClaimsForGroup: jest.fn(),
    createClaim: jest.fn(),
    addReceiptsToClaim: jest.fn(),
  };
  const matDialog = { open: jest.fn() };
  const matDialogRef = { close: jest.fn() };
  const snackbarService = { success: jest.fn(), error: jest.fn() };

  const existing = { id: 40, groupId: 5, name: "Voucher 13", status: "SUBMITTED" } as Claim;
  const created = { id: 41, groupId: 5, name: "Voucher 14", status: "DRAFT" } as Claim;

  const taxi = { id: 1, name: "Taxi", groupId: 5 };
  const hotel = { id: 2, name: "Hotel", groupId: 5 };
  const elsewhere = { id: 3, name: "Groceries", groupId: 6 };

  const setUp = async (data: AddToClaimDialogData) => {
    await TestBed.configureTestingModule({
      imports: [AddToClaimDialogComponent, NgxsModule.forRoot([AuthState, GroupState])],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ReimbursementClaimService, useValue: claimService },
        { provide: MatDialog, useValue: matDialog },
        { provide: MatDialogRef, useValue: matDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
        { provide: SnackbarService, useValue: snackbarService },
      ],
    })
      .overrideComponent(AddToClaimDialogComponent, {
        set: { imports: [ReactiveFormsModule, PipesModule], schemas: [CUSTOM_ELEMENTS_SCHEMA] },
      })
      .compileComponents();

    store = TestBed.inject(Store);
    store.dispatch(
      new SetGroups([
        { id: 5, name: "Household", isAllGroup: false },
        { id: 6, name: "Work", isAllGroup: false },
      ] as Group[])
    );
    store.dispatch(
      new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate], 6: [Permission.GroupReceiptsRead] })
    );

    fixture = TestBed.createComponent(AddToClaimDialogComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  };

  const text = (testId: string): string =>
    fixture.nativeElement.querySelector(`[data-testid='${testId}']`)?.textContent?.replace(/\s+/g, " ").trim();

  beforeEach(() => {
    jest.resetAllMocks();
    claimService.getClaimsForGroup.mockReturnValue(of({ claims: [existing], summary: {} } as ClaimList));
  });

  // The simplest honest behaviour for a selection made on the All view: a claim
  // is one group's, so say so instead of guessing.
  it("refuses a selection that spans groups, naming them", async () => {
    await setUp({ receipts: [taxi, elsewhere] });

    expect(text("add-to-claim-mixed-groups")).toContain("are in 2 groups (Household, Work)");
    expect(claimService.getClaimsForGroup).not.toHaveBeenCalled();
    expect(fixture.nativeElement.querySelector("form")).toBeNull();
  });

  it("says so when the selection's group does not allow changing claims", async () => {
    await setUp({ receipts: [elsewhere] });

    expect(text("add-to-claim-forbidden")).toContain("in Work");
    expect(claimService.getClaimsForGroup).not.toHaveBeenCalled();
  });

  describe("with one group", () => {
    beforeEach(async () => {
      await setUp({ receipts: [taxi, hotel] });
    });

    it("offers a new claim first, then the group's claims with their status", () => {
      expect(claimService.getClaimsForGroup).toHaveBeenCalledWith(5);
      expect(component.targetOptions()).toEqual([
        { value: NEW_CLAIM_OPTION, displayValue: "New claim from these receipts" },
        { value: 40, displayValue: "Voucher 13 (Submitted)" },
      ]);
    });

    // Adding to the wrong claim is the mistake to prevent, so nothing is pre-picked.
    it("does not submit until a claim is picked", () => {
      component.submit();

      expect(claimService.addReceiptsToClaim).not.toHaveBeenCalled();
      expect(component.form.controls.claimId.touched).toBe(true);
    });

    it("adds the selection to an existing claim and reports it", () => {
      const result = { receiptIds: [1, 2], failures: [] };
      claimService.addReceiptsToClaim.mockReturnValue(of(result));
      component.form.controls.claimId.setValue(40);

      component.submit();

      expect(claimService.createClaim).not.toHaveBeenCalled();
      expect(claimService.addReceiptsToClaim).toHaveBeenCalledWith(40, { receiptIds: [1, 2] });
      expect(matDialogRef.close).toHaveBeenCalledWith(existing);
      expect(snackbarService.success).toHaveBeenCalledWith("Added 2 receipts to Voucher 13");
    });

    it("lists each refused receipt by name", () => {
      claimService.addReceiptsToClaim.mockReturnValue(
        of({ receiptIds: [1], failures: [{ receiptId: 2, message: "Already in claim 'Voucher 12'" }] })
      );
      component.form.controls.claimId.setValue(40);

      component.submit();

      expect(matDialog.open).toHaveBeenCalledWith(
        ClaimReceiptsResultDialogComponent,
        expect.objectContaining({
          data: expect.objectContaining({
            failures: [{ name: "Hotel", message: "Already in claim 'Voucher 12'" }],
          }),
        })
      );
    });

    it("names no receipt as already claimed when none is", () => {
      expect(component.alreadyClaimed()).toEqual([]);
      expect(fixture.nativeElement.querySelector("[data-testid='add-to-claim-already-claimed']")).toBeNull();
    });

    it("requires a name for a new claim", () => {
      component.form.controls.claimId.setValue(NEW_CLAIM_OPTION);
      component.form.controls.name.setValue("   ");

      component.submit();

      expect(component.form.controls.name.hasError("required")).toBe(true);
      expect(claimService.createClaim).not.toHaveBeenCalled();
    });

    it("creates the claim, then adds the receipts to it", () => {
      claimService.createClaim.mockReturnValue(of(created));
      claimService.addReceiptsToClaim.mockReturnValue(of({ receiptIds: [1, 2], failures: [] }));
      component.form.controls.claimId.setValue(NEW_CLAIM_OPTION);
      component.form.controls.name.setValue(" Voucher 14 ");

      component.submit();

      expect(claimService.createClaim).toHaveBeenCalledWith(5, { name: "Voucher 14", status: "DRAFT" });
      expect(claimService.addReceiptsToClaim).toHaveBeenCalledWith(41, { receiptIds: [1, 2] });
      expect(matDialogRef.close).toHaveBeenCalledWith(created);
    });

    it("reuses the claim it created when retrying a failed add", () => {
      claimService.createClaim.mockReturnValue(of(created));
      claimService.addReceiptsToClaim.mockReturnValueOnce(
        throwError(() => new HttpErrorResponse({ status: 500 }))
      );
      component.form.controls.claimId.setValue(NEW_CLAIM_OPTION);
      component.form.controls.name.setValue("Voucher 14");

      component.submit();
      expect(matDialogRef.close).not.toHaveBeenCalled();

      claimService.addReceiptsToClaim.mockReturnValue(of({ receiptIds: [1, 2], failures: [] }));
      component.submit();

      expect(claimService.createClaim).toHaveBeenCalledTimes(1);
      expect(claimService.addReceiptsToClaim).toHaveBeenLastCalledWith(41, { receiptIds: [1, 2] });
      expect(matDialogRef.close).toHaveBeenCalledWith(created);
    });

    it("puts the server's field errors on the form", () => {
      claimService.createClaim.mockReturnValue(
        throwError(() => new HttpErrorResponse({ status: 400, error: { name: "Name is too long" } }))
      );
      component.form.controls.claimId.setValue(NEW_CLAIM_OPTION);
      component.form.controls.name.setValue("Voucher");

      component.submit();

      expect(component.form.controls.name.errors).toEqual({ error: "Name is too long" });
      expect(matDialogRef.close).not.toHaveBeenCalled();
    });
  });
  // Named up front, so the result dialog's failure list is the exception.
  describe("with receipts already in a claim", () => {
    const fuel = { id: 3, name: "Fuel", groupId: 5, claimId: 40, claimName: "Voucher 13" };
    const parking = { id: 4, name: "Parking", groupId: 5, claimId: 39, claimName: "Voucher 12" };
    const toll = { id: 5, name: "Toll", groupId: 5, claimId: 39, claimName: "Voucher 12" };

    beforeEach(async () => {
      await setUp({ receipts: [taxi, fuel, parking, toll] });
    });

    it("names them by claim before a claim is picked", () => {
      expect(text("add-to-claim-already-claimed")).toContain("Voucher 13: Fuel");
      expect(text("add-to-claim-already-claimed")).toContain("Voucher 12: Parking, Toll");
      expect(text("add-to-claim-already-claimed")).not.toContain("Taxi");
    });

    it("leaves out the ones already in the claim picked, which the server counts as added", async () => {
      component.form.controls.claimId.setValue(40);
      await fixture.whenStable();

      expect(component.alreadyClaimed()).toEqual([
        { claimId: 39, claimName: "Voucher 12", receiptNames: ["Parking", "Toll"] },
      ]);
      expect(text("add-to-claim-already-claimed")).not.toContain("Fuel");
    });

    it("still sends them, leaving the refusal to the server", () => {
      claimService.addReceiptsToClaim.mockReturnValue(of({ receiptIds: [1, 3], failures: [] }));
      component.form.controls.claimId.setValue(40);

      component.submit();

      expect(claimService.addReceiptsToClaim).toHaveBeenCalledWith(40, { receiptIds: [1, 3, 4, 5] });
    });
  });
});
