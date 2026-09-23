import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { NgxsModule } from "@ngxs/store";
import { FilterOperation, SystemTask, SystemTaskPagedRequestFilter, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { SystemEmailTaskTableService } from "../../services/system-email-task-table.service";
import { SystemEmailTaskTableState } from "../../store/system-email-task-table.state";

import { ReceiptUpdateDiffDialogComponent } from "../receipt-update-diff-dialog/receipt-update-diff-dialog.component";
import { TaskTableComponent } from "./task-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("TaskTableComponent", () => {
  let component: TaskTableComponent;
  let fixture: ComponentFixture<TaskTableComponent>;
  let httpTesting: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [TaskTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [NgxsModule.forRoot([SystemEmailTaskTableState])],
    providers: [
        {
            provide: TABLE_SERVICE_INJECTION_TOKEN,
            useClass: SystemEmailTaskTableService
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
})
      .compileComponents();

    httpTesting = TestBed.inject(HttpTestingController);
    fixture = TestBed.createComponent(TaskTableComponent);
    component = fixture.componentInstance;
  });

  afterEach(() => {
    httpTesting.verify();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("sends the filter the provider reports with the paged request", () => {
    const filter = {
      type: { operation: FilterOperation.Contains, value: ["QUICK_SCAN"] },
    } as unknown as SystemTaskPagedRequestFilter;
    fixture.componentRef.setInput("filterProvider", () => filter);

    component.getTableData();

    const request = httpTesting.expectOne("/api/systemTask/getPagedSystemTasks");
    expect(request.request.body.filter).toEqual(filter);
    request.flush({ data: [], totalCount: 0 });
  });

  // A local-midnight Date would serialize as an instant, which the server
  // resolves to a day in its own zone -- so the request must carry the calendar
  // day the user picked. See toSystemTaskWireFilter.
  it("sends the date fields as local calendar days", () => {
    fixture.componentRef.setInput("filterProvider", () => ({
      startedAt: { operation: FilterOperation.Equals, value: new Date(2026, 8, 22, 0, 0, 0, 0) },
      endedAt: {
        operation: FilterOperation.Between,
        value: [new Date(2026, 8, 10, 0, 0, 0, 0), new Date(2026, 8, 12, 0, 0, 0, 0)],
      },
    }) as unknown as SystemTaskPagedRequestFilter);

    component.getTableData();

    const request = httpTesting.expectOne("/api/systemTask/getPagedSystemTasks");
    expect((request.request.body.filter.startedAt as any).value).toBe("2026-09-22");
    expect((request.request.body.filter.endedAt as any).value).toEqual(["2026-09-10", "2026-09-12"]);
    request.flush({ data: [], totalCount: 0 });
  });

  // The page dispatches a filter change and refreshes in the same synchronous
  // turn, before change detection can push a new input value in. Reading the
  // filter at request time is what makes a cleared chip actually clear.
  it("reads the filter at request time, not when the provider was bound", () => {
    let filter: any = { type: { operation: FilterOperation.Contains, value: ["QUICK_SCAN"] } };
    fixture.componentRef.setInput("filterProvider", () => filter);

    component.getTableData();
    httpTesting.expectOne("/api/systemTask/getPagedSystemTasks").flush({ data: [], totalCount: 0 });

    filter = { type: { operation: null, value: [] } };
    component.getTableData();

    const request = httpTesting.expectOne("/api/systemTask/getPagedSystemTasks");
    expect(request.request.body.filter).toEqual(filter);
    request.flush({ data: [], totalCount: 0 });
  });

  // The two embedded task tables never bind it, so the key must be absent
  // rather than sent as null -- the API then sees a zero-value filter and adds
  // no predicates.
  it("omits the filter key when no provider is bound", () => {
    component.getTableData();

    const request = httpTesting.expectOne("/api/systemTask/getPagedSystemTasks");
    expect(request.request.body.filter).toBeUndefined();
    request.flush({ data: [], totalCount: 0 });
  });

  // Two chip clears in quick succession put two fetches in flight; the last
  // request must win, not the last response.
  it("supersedes an in-flight refresh rather than letting it repaint the table", () => {
    component.getTableData();
    component.getTableData();

    const requests = httpTesting.match("/api/systemTask/getPagedSystemTasks");
    expect(requests.length).toBe(2);

    // The first was cancelled by switchMap, so only the second can land.
    expect(requests[0].cancelled).toBe(true);
    requests[1].flush({ data: [], totalCount: 3 });

    expect(component.totalCount()).toBe(3);
  });
  describe("receipt update rows", () => {
    const before = { id: 1, updatedAt: "2026-09-22T10:00:00Z", name: "Costco", amount: "42.1" };
    const after = { id: 1, updatedAt: "2026-09-23T10:00:00Z", name: "Costco Wholesale", amount: "45.6" };

    const task = (id: number, type: SystemTaskType, resultDescription: string): SystemTask =>
      ({ id, type, status: SystemTaskStatus.Succeeded, resultDescription } as SystemTask);

    const loadRows = (rows: SystemTask[]) => {
      component.getTableData();
      httpTesting.expectOne("/api/systemTask/getPagedSystemTasks").flush({ data: rows, totalCount: rows.length });
    };

    // What UpdateReceipt stores: each side is itself a JSON string.
    const stored = JSON.stringify({ before: JSON.stringify(before), after: JSON.stringify(after) });

    it("parses a receipt update and lists the keys it changed, ignoring updatedAt", () => {
      loadRows([task(1, SystemTaskType.ReceiptUpdated, stored)]);

      const row = component.receiptUpdates().get(1);
      expect(row?.snapshots).toEqual({ before, after });
      expect(row?.changedKeys).toEqual(["name", "amount"]);
    });

    it("leaves a failed update and other task types on the plain description path", () => {
      loadRows([
        task(1, SystemTaskType.ReceiptUpdated, "record not found"),
        task(2, SystemTaskType.QuickScan, stored),
      ]);

      expect(component.receiptUpdates().size).toBe(0);
    });

    it("opens the diff dialog with the parsed snapshots", () => {
      const dialog = TestBed.inject(MatDialog);
      const open = jest.spyOn(dialog, "open").mockReturnValue(undefined as any);

      component.openReceiptUpdateDialog({ before, after });

      expect(open).toHaveBeenCalledWith(
        ReceiptUpdateDiffDialogComponent,
        expect.objectContaining({ data: { snapshots: { before, after } }, width: "90vw" }),
      );
    });
  });
});
