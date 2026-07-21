import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { of } from "rxjs";
import { Receipt, ReceiptService, ReportPeriod } from "../../../open-api";
import { PipesModule } from "../../../pipes";
import {
  ReportReceiptsDialogComponent,
  ReportReceiptsDialogData,
} from "./report-receipts-dialog.component";

const receipts: Receipt[] = [
  { id: 7, name: "Coffee", date: "2026-07-03", amount: "12.50", status: "RESOLVED", groupId: 1, paidByUserId: 1, categories: [], tags: [] } as Receipt,
  { id: 8, name: "Lunch", date: "2026-07-05", amount: "30.00", status: "OPEN", groupId: 1, paidByUserId: 2, categories: [], tags: [] } as Receipt,
];

function configure(data: Partial<ReportReceiptsDialogData> = {}): {
  fixture: ComponentFixture<ReportReceiptsDialogComponent>;
  component: ReportReceiptsDialogComponent;
} {
  const receiptService = {
    getReceiptsForGroup: jest.fn(() => of({ data: receipts, totalCount: receipts.length })),
  };
  TestBed.configureTestingModule({
    declarations: [ReportReceiptsDialogComponent],
    imports: [PipesModule],
    providers: [
      provideZonelessChangeDetection(),
      { provide: ReceiptService, useValue: receiptService },
      { provide: MatDialogRef, useValue: { close: jest.fn() } },
      {
        provide: MAT_DIALOG_DATA,
        useValue: {
          groupIds: ["1"],
          filter: {},
          period: { preset: ReportPeriod.PresetEnum.ThisMonth, startDate: null, endDate: null },
          receiptCount: 5,
          ...data,
        } as ReportReceiptsDialogData,
      },
    ],
    schemas: [NO_ERRORS_SCHEMA],
  });
  const fixture = TestBed.createComponent(ReportReceiptsDialogComponent);
  return { fixture, component: fixture.componentInstance };
}

describe("ReportReceiptsDialogComponent", () => {
  afterEach(() => TestBed.resetTestingModule());

  it("loads the covered receipts and starts on the list view", () => {
    const { component } = configure();
    expect(component.receipts().length).toBe(2);
    expect(component.selected()).toBeNull();
    // The subtitle count is the report's true total, not the loaded sample.
    expect(component.count()).toBe(5);
  });

  it("opens a receipt's breakdown and returns to the list", () => {
    const { component } = configure();
    component.viewReceipt(receipts[0]);
    expect(component.selected()).toBe(receipts[0]);
    component.backToList();
    expect(component.selected()).toBeNull();
  });

  it("opens the full receipt page in a new tab", () => {
    const { component } = configure();
    const openSpy = jest.spyOn(window, "open").mockReturnValue(null);
    component.openFullReceipt(receipts[0]);
    expect(openSpy).toHaveBeenCalledWith("/receipts/7/view", "_blank");
    openSpy.mockRestore();
  });

  it("falls back to the loaded count when no true count is provided", () => {
    const { component } = configure({ receiptCount: undefined });
    expect(component.count()).toBe(2);
  });
});
