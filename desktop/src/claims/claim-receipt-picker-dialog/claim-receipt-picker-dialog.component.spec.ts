import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { Receipt } from "../../open-api";
import { PipesModule } from "../../pipes";
import {
  ClaimReceiptPickerDialogComponent,
  ClaimReceiptPickerDialogData,
} from "./claim-receipt-picker-dialog.component";

describe("ClaimReceiptPickerDialogComponent", () => {
  let fixture: ComponentFixture<ClaimReceiptPickerDialogComponent>;
  let component: ClaimReceiptPickerDialogComponent;
  const matDialogRef = { close: jest.fn() };

  const taxi = { id: 1, name: "Taxi" } as Receipt;

  const setUp = async (data: ClaimReceiptPickerDialogData) => {
    await TestBed.configureTestingModule({
      imports: [ClaimReceiptPickerDialogComponent],
      providers: [
        provideZonelessChangeDetection(),
        { provide: MatDialogRef, useValue: matDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    })
      .overrideComponent(ClaimReceiptPickerDialogComponent, {
        set: { imports: [ReactiveFormsModule, PipesModule], schemas: [CUSTOM_ELEMENTS_SCHEMA] },
      })
      .compileComponents();

    fixture = TestBed.createComponent(ClaimReceiptPickerDialogComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  };

  beforeEach(() => jest.resetAllMocks());

  it("closes with the picked receipts", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [taxi], inOtherClaimsCount: 0 });
    component.form.controls.receipts.push(new FormControl(taxi));

    component.submit();

    expect(matDialogRef.close).toHaveBeenCalledWith([taxi]);
  });

  it("stays open with nothing picked", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [taxi], inOtherClaimsCount: 0 });

    component.submit();

    expect(matDialogRef.close).not.toHaveBeenCalled();
    expect(component.form.controls.receipts.hasError("required")).toBe(true);
  });

  it("says when there is nothing left to add, and offers no submit", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [], inOtherClaimsCount: 0 });

    expect(
      fixture.nativeElement.querySelector("[data-testid='claim-receipt-picker-empty']").textContent
    ).toContain("already on Voucher 14");
    expect(fixture.nativeElement.querySelector("app-autocomlete")).toBeNull();
  });

  it("says how many receipts in other claims it left out", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [taxi], inOtherClaimsCount: 2 });

    expect(
      fixture.nativeElement.querySelector("[data-testid='claim-receipt-picker-other-claims']").textContent
    ).toContain("2 receipts already in another claim are not listed.");
  });

  it("stays quiet when nothing was left out", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [taxi], inOtherClaimsCount: 0 });

    expect(
      fixture.nativeElement.querySelector("[data-testid='claim-receipt-picker-other-claims']")
    ).toBeNull();
  });

  it("names the other claims when nothing is left to add", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [], inOtherClaimsCount: 1 });

    expect(
      fixture.nativeElement
        .querySelector("[data-testid='claim-receipt-picker-empty']")
        .textContent.replace(/\s+/g, " ")
    ).toContain("already on Voucher 14 or another claim.");
  });

  it("closes with nothing on cancel", async () => {
    await setUp({ claimName: "Voucher 14", candidates: [taxi], inOtherClaimsCount: 0 });

    component.cancel();

    expect(matDialogRef.close).toHaveBeenCalledWith(undefined);
  });
});
