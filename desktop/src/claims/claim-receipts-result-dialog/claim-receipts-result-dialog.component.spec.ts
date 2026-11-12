import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from "@angular/material/dialog";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { SnackbarService } from "../../services/snackbar.service";
import {
  ClaimReceiptsResultDialogComponent,
  ClaimReceiptsResultDialogData,
} from "./claim-receipts-result-dialog.component";
import { reportClaimReceiptsResult } from "./report-claim-receipts-result";

describe("ClaimReceiptsResultDialogComponent", () => {
  let component: ClaimReceiptsResultDialogComponent;
  let fixture: ComponentFixture<ClaimReceiptsResultDialogComponent>;

  const data: ClaimReceiptsResultDialogData = {
    headerText: "Some receipts could not be added",
    summary: "Added 1 of 2 receipts to Voucher 14.",
    failures: [{ name: "Hotel Lima", message: "This receipt is already in claim 'Voucher 13'" }],
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ClaimReceiptsResultDialogComponent],
      providers: [
        provideZonelessChangeDetection(),
        { provide: MatDialogRef, useValue: { close: jest.fn() } },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    })
      .overrideComponent(ClaimReceiptsResultDialogComponent, {
        set: { imports: [], schemas: [CUSTOM_ELEMENTS_SCHEMA] },
      })
      .compileComponents();

    fixture = TestBed.createComponent(ClaimReceiptsResultDialogComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it("shows the summary and each refused receipt by name with its reason", () => {
    const element: HTMLElement = fixture.nativeElement;

    expect(element.querySelector("[data-testid='claim-receipts-result-summary']")?.textContent?.trim()).toBe(
      "Added 1 of 2 receipts to Voucher 14."
    );
    expect(
      Array.from(element.querySelectorAll("[data-testid='claim-receipts-result-failures'] li")).map((item) =>
        item.textContent?.replace(/\s+/g, " ").trim()
      )
    ).toEqual(["Hotel Lima: This receipt is already in claim 'Voucher 13'"]);
  });

  it("closes through the shared footer's dismiss-only mode", () => {
    const footer = fixture.nativeElement.querySelector("app-dialog-footer");

    footer.dispatchEvent(new CustomEvent("cancelClicked"));

    expect(component.matDialogRef.close).toHaveBeenCalled();
  });
});

describe("reportClaimReceiptsResult", () => {
  const matDialog = { open: jest.fn() } as unknown as MatDialog;
  const snackbarService = { success: jest.fn() } as unknown as SnackbarService;

  beforeEach(() => jest.clearAllMocks());

  it("is a snackbar when every receipt went through", () => {
    reportClaimReceiptsResult(matDialog, snackbarService, {
      action: "add",
      claimName: "Voucher 14",
      requested: [{ id: 1, name: "Taxi" }, { id: 2, name: "Hotel" }],
      result: { receiptIds: [1, 2], failures: [] },
    });

    expect(snackbarService.success).toHaveBeenCalledWith("Added 2 receipts to Voucher 14");
    expect(matDialog.open).not.toHaveBeenCalled();
  });

  it("names each refusal, falling back to the id for a receipt it cannot name", () => {
    reportClaimReceiptsResult(matDialog, snackbarService, {
      action: "remove",
      claimName: "Voucher 14",
      requested: [{ id: 1, name: "Taxi" }, { id: 2 }],
      result: {
        receiptIds: [],
        failures: [
          { receiptId: 1, message: "Not on this claim" },
          { receiptId: 2, message: "Not on this claim" },
        ],
      },
    });

    expect(snackbarService.success).not.toHaveBeenCalled();
    expect(matDialog.open).toHaveBeenCalledWith(ClaimReceiptsResultDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      data: {
        headerText: "Some receipts could not be removed",
        summary: "Removed 0 of 2 receipts from Voucher 14.",
        failures: [
          { name: "Taxi", message: "Not on this claim" },
          { name: "Receipt #2", message: "Not on this claim" },
        ],
      },
    });
  });
});
