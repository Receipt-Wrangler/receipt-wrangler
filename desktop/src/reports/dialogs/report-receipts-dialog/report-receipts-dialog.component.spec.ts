import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { Observable, of, throwError } from "rxjs";
import {
  PagedData,
  Receipt,
  ReportColumn,
  ReportDetail,
  ReportPeriod,
  ReportRequestCommand,
  ReportService,
} from "../../../open-api";
import { PipesModule } from "../../../pipes";
import {
  ReportReceiptsDialogComponent,
  ReportReceiptsDialogData,
} from "./report-receipts-dialog.component";

const receipts: Receipt[] = [
  { id: 7, name: "Coffee", date: "2026-07-03", amount: "12.50", status: "RESOLVED", groupId: 1, paidByUserId: 1, categories: [], tags: [] } as Receipt,
  { id: 8, name: "Lunch", date: "2026-07-05", amount: "30.00", status: "OPEN", groupId: 1, paidByUserId: 2, categories: [], tags: [] } as Receipt,
];

const command: ReportRequestCommand = {
  name: "Report",
  groupIds: ["1", "2"],
  period: { preset: ReportPeriod.PresetEnum.Custom, startDate: "2026-05-01", endDate: "2026-05-31", dateField: "createdAt" },
  filter: {},
  detail: { mode: ReportDetail.ModeEnum.Records },
  columns: [{ kind: ReportColumn.KindEnum.Dimension, name: "Name", label: "Name", field: "name" }],
  formats: [ReportRequestCommand.FormatsEnum.Pdf],
};

const MAY_2026: ReportReceiptsDialogData["period"] = {
  preset: ReportPeriod.PresetEnum.Custom,
  startDate: new Date(2026, 4, 1),
  endDate: new Date(2026, 4, 31),
  dateField: "createdAt",
};

function configure(
  data: Partial<ReportReceiptsDialogData> = {},
  response: Observable<PagedData> = of({ data: receipts, totalCount: 5 } as unknown as PagedData)
): {
  fixture: ComponentFixture<ReportReceiptsDialogComponent>;
  component: ReportReceiptsDialogComponent;
  reportService: { getReportReceipts: jest.Mock };
} {
  const reportService = { getReportReceipts: jest.fn(() => response) };
  TestBed.configureTestingModule({
    declarations: [ReportReceiptsDialogComponent],
    imports: [PipesModule],
    providers: [
      provideZonelessChangeDetection(),
      { provide: ReportService, useValue: reportService },
      { provide: MatDialogRef, useValue: { close: jest.fn() } },
      {
        provide: MAT_DIALOG_DATA,
        useValue: { command, period: MAY_2026, receiptCount: 5, ...data } as ReportReceiptsDialogData,
      },
    ],
    schemas: [NO_ERRORS_SCHEMA],
  });
  const fixture = TestBed.createComponent(ReportReceiptsDialogComponent);
  return { fixture, component: fixture.componentInstance, reportService };
}

describe("ReportReceiptsDialogComponent", () => {
  afterEach(() => TestBed.resetTestingModule());

  it("loads the covered receipts and starts on the list view", () => {
    const { component } = configure();
    expect(component.receipts().length).toBe(2);
    expect(component.loading()).toBe(false);
    expect(component.hasError()).toBe(false);
    expect(component.selected()).toBeNull();
  });

  // The server resolves the period and filter exactly as the report does, so the
  // dialog sends the preview's own command, once, and builds no bounds itself.
  it("asks the server for the receipts the preview's command covers, in one request", () => {
    const { reportService } = configure();
    expect(reportService.getReportReceipts).toHaveBeenCalledTimes(1);
    expect(reportService.getReportReceipts).toHaveBeenCalledWith(command);
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

  it("names the date field in the subtitle", () => {
    const { component } = configure();
    expect(component.periodLabel).toBe("2026-05-01 to 2026-05-31 on Added At");
  });

  it("names the receipt date when the period names no date field", () => {
    const { component } = configure({
      period: { ...MAY_2026, dateField: undefined } as unknown as ReportReceiptsDialogData["period"],
    });
    expect(component.periodLabel).toBe("2026-05-01 to 2026-05-31 on Receipt Date");
  });

  it("shows the report's own count when one is provided", () => {
    const { component } = configure({ receiptCount: 9 });
    expect(component.count()).toBe(9);
  });

  // The list is capped server-side, so its length is not the count; its total is.
  it("falls back to the list's total when no count is provided", () => {
    const { component } = configure({ receiptCount: undefined });
    expect(component.count()).toBe(5);
  });

  it("flags an error and stops loading when the list cannot be fetched", () => {
    const { component } = configure({}, throwError(() => new Error("403")));
    expect(component.hasError()).toBe(true);
    expect(component.loading()).toBe(false);
    expect(component.receipts()).toEqual([]);
  });
});
