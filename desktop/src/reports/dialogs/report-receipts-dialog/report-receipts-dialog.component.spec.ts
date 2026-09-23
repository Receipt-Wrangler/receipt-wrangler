import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { of } from "rxjs";
import { RECEIPT_DATE_FILTER_FIELDS } from "src/constants";
import { FilterOperation, Receipt, ReceiptService, ReportPeriod } from "../../../open-api";
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
  receiptService: { getReceiptsForGroup: jest.Mock };
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
          period: { preset: ReportPeriod.PresetEnum.ThisMonth, startDate: null, endDate: null, dateField: "date" },
          receiptCount: 5,
          ...data,
        } as ReportReceiptsDialogData,
      },
    ],
    schemas: [NO_ERRORS_SCHEMA],
  });
  const fixture = TestBed.createComponent(ReportReceiptsDialogComponent);
  return { fixture, component: fixture.componentInstance, receiptService };
}

const MAY_2026 = {
  preset: ReportPeriod.PresetEnum.Custom,
  startDate: new Date(2026, 4, 1),
  endDate: new Date(2026, 4, 31),
};

const DATE_KEYS = RECEIPT_DATE_FILTER_FIELDS.map((field) => field.key);

// A condition on every date field, so a test can see which one the period replaced.
function sentinelDateFilter(): Record<string, { operation: FilterOperation; value: string }> {
  return Object.fromEntries(
    DATE_KEYS.map((key) => [key, { operation: FilterOperation.Equals, value: `sentinel-${key}` }])
  );
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

  it.each(DATE_KEYS)("narrows the listed receipts on %s, leaving the other date fields alone", (dateField) => {
    const { receiptService } = configure({
      filter: sentinelDateFilter(),
      period: { ...MAY_2026, dateField },
    });

    const filter = receiptService.getReceiptsForGroup.mock.calls[0][1].filter;
    for (const key of DATE_KEYS) {
      if (key === dateField) {
        expect(filter[key]).toEqual({
          operation: FilterOperation.Between,
          value: [new Date(2026, 4, 1), new Date(2026, 4, 31, 23, 59, 59, 999)],
        });
      } else {
        expect(filter[key]).toEqual(sentinelDateFilter()[key]);
      }
    }
  });

  it("narrows on the receipt date when the period names no date field", () => {
    const { receiptService } = configure({
      period: { ...MAY_2026 } as ReportReceiptsDialogData["period"],
    });

    const filter = receiptService.getReceiptsForGroup.mock.calls[0][1].filter;
    expect(filter.date.operation).toBe(FilterOperation.Between);
    expect(filter.resolvedDate).toBeUndefined();
    expect(filter.createdAt).toBeUndefined();
  });

  // A custom range ends at local midnight; on the timestamped Added At column
  // that would drop everything added during the range's last day.
  it("extends the range to the end of its last day", () => {
    const { receiptService } = configure({ period: { ...MAY_2026, dateField: "createdAt" } });

    const [start, end] = receiptService.getReceiptsForGroup.mock.calls[0][1].filter.createdAt.value;
    expect(start).toEqual(new Date(2026, 4, 1, 0, 0, 0, 0));
    expect(end).toEqual(new Date(2026, 4, 31, 23, 59, 59, 999));
  });

  it("names the date field in the subtitle", () => {
    const { component } = configure({ period: { ...MAY_2026, dateField: "createdAt" } });
    expect(component.periodLabel).toBe("2026-05-01 to 2026-05-31 on Added At");
  });

  it("falls back to the loaded count when no true count is provided", () => {
    const { component } = configure({ receiptCount: undefined });
    expect(component.count()).toBe(2);
  });
});
