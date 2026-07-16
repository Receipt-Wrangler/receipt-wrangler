import { CommonModule } from "@angular/common";
import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { of, Subject, throwError } from "rxjs";
import {
  ReportColumn,
  ReportDetail,
  ReportPeriod,
  ReportRequestCommand,
  ReportTemplate,
} from "../../open-api";
import { SnackbarService } from "../../services";
import { ReportRunnerService } from "../services/report-runner.service";
import { ReportTemplateListComponent } from "./report-template-list.component";

function makeTemplate(id: number, name: string): ReportTemplate {
  return {
    id,
    name,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-02-01T00:00:00Z",
    configurationVersion: 1,
    configuration: {
      name,
      groupIds: ["1"],
      period: { preset: ReportPeriod.PresetEnum.ThisMonth },
      filter: {},
      groupBy: ["group"],
      detail: { mode: ReportDetail.ModeEnum.Records },
      columns: [{ kind: ReportColumn.KindEnum.Dimension, name: "Name", label: "Name", field: "name" }],
      subtotals: true,
      grandTotals: true,
      formats: [ReportRequestCommand.FormatsEnum.Csv, ReportRequestCommand.FormatsEnum.Pdf],
    },
  };
}

describe("ReportTemplateListComponent", () => {
  let fixture: ComponentFixture<ReportTemplateListComponent>;
  let component: ReportTemplateListComponent;
  let runner: {
    listTemplates: jest.Mock;
    duplicateTemplate: jest.Mock;
    deleteTemplate: jest.Mock;
    generateFromTemplate: jest.Mock;
    getTemplate: jest.Mock;
  };
  let router: { navigate: jest.Mock };
  let snackbar: { success: jest.Mock; error: jest.Mock };
  let dialogResult: boolean;
  let dialogInstance: { headerText: string; dialogContent: string };

  const templates = [makeTemplate(1, "Alpha"), makeTemplate(2, "Beta")];

  beforeEach(async () => {
    dialogResult = true;
    dialogInstance = { headerText: "", dialogContent: "" };
    runner = {
      listTemplates: jest.fn(() => of({ data: templates, totalCount: templates.length })),
      duplicateTemplate: jest.fn(() => of(makeTemplate(3, "Alpha duplicate"))),
      deleteTemplate: jest.fn(() => of(undefined)),
      generateFromTemplate: jest.fn(() => of(new Blob())),
      getTemplate: jest.fn(),
    };
    router = { navigate: jest.fn() };
    snackbar = { success: jest.fn(), error: jest.fn() };

    const store = {
      selectSignal: jest.fn(() => signal([])), // GroupState.groupsWithoutAll
      selectSnapshot: jest.fn(() => ({ page: 1, pageSize: 50, orderBy: "updated_at", sortDirection: "desc" })),
      select: jest.fn(() => of(1)),
      dispatch: jest.fn(),
    };
    const matDialog = {
      open: jest.fn(() => ({
        componentInstance: dialogInstance,
        afterClosed: () => of(dialogResult),
      })),
    };

    await TestBed.configureTestingModule({
      declarations: [ReportTemplateListComponent],
      imports: [CommonModule],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ReportRunnerService, useValue: runner },
        { provide: Store, useValue: store },
        { provide: Router, useValue: router },
        { provide: SnackbarService, useValue: snackbar },
        { provide: MatDialog, useValue: matDialog },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(ReportTemplateListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    await fixture.whenStable();
  });

  it("loads templates into the table on init", () => {
    expect(runner.listTemplates).toHaveBeenCalled();
    expect(component.totalCount()).toBe(2);
    expect(component.dataSource().data.length).toBe(2);
    expect(component.loaded()).toBe(true);
  });

  it("New Report navigates to the builder", () => {
    component.newReport();
    expect(router.navigate).toHaveBeenCalledWith(["/reports/new"]);
  });

  it("generate runs the stored configuration and clears the in-flight id", () => {
    component.generate(templates[0]);
    expect(runner.generateFromTemplate).toHaveBeenCalledWith(templates[0].configuration);
    // of() completes synchronously, so finalize resets the flag.
    expect(component.generatingId()).toBeNull();
  });

  it("generate ignores a second call while one is in flight", () => {
    // A non-completing observable keeps the first generate in flight.
    const pending = new Subject<Blob>();
    runner.generateFromTemplate.mockReturnValue(pending.asObservable());

    component.generate(templates[0]);
    expect(component.generatingId()).toBe(1);

    // A second click while the first runs must not fire another request.
    component.generate(templates[1]);
    expect(runner.generateFromTemplate).toHaveBeenCalledTimes(1);

    pending.complete();
    expect(component.generatingId()).toBeNull();
  });

  it("duplicate copies the template, toasts, and reloads", () => {
    runner.listTemplates.mockClear();
    component.duplicate(templates[0]);
    expect(runner.duplicateTemplate).toHaveBeenCalledWith(1);
    expect(snackbar.success).toHaveBeenCalledWith("Template duplicated");
    expect(runner.listTemplates).toHaveBeenCalled();
  });

  it("duplicate ignores a second call while one is in flight", () => {
    // A non-completing observable keeps the first duplicate in flight.
    const pending = new Subject<ReportTemplate>();
    runner.duplicateTemplate.mockReturnValue(pending.asObservable());

    component.duplicate(templates[0]);
    expect(component.duplicatingId()).toBe(1);

    // A second click while the first runs must not fire another request.
    component.duplicate(templates[1]);
    expect(runner.duplicateTemplate).toHaveBeenCalledTimes(1);

    pending.complete();
    expect(component.duplicatingId()).toBeNull();
  });

  it("delete confirms then deletes and reloads", () => {
    dialogResult = true;
    runner.listTemplates.mockClear();
    component.delete(templates[0]);
    expect(dialogInstance.headerText).toBe("Delete Report Template");
    expect(runner.deleteTemplate).toHaveBeenCalledWith(1);
    expect(snackbar.success).toHaveBeenCalledWith("Template deleted");
    expect(runner.listTemplates).toHaveBeenCalled();
  });

  it("delete does nothing when the dialog is dismissed", () => {
    dialogResult = false;
    component.delete(templates[0]);
    expect(runner.deleteTemplate).not.toHaveBeenCalled();
  });

  it("derives display summaries from the stored configuration", () => {
    expect(component.columnCountFor(templates[0])).toBe(1);
    expect(component.formatChipsFor(templates[0])).toEqual(["CSV", "PDF"]);
    expect(component.detailSummaryFor(templates[0])).toBe("Record-level");
    expect(component.groupingSummaryFor(templates[0])).toBe("Group");
  });

  it("summarizes scope by group count (unresolved groups fall back to the raw id)", () => {
    const scoped = (ids: string[]): ReportTemplate => ({
      ...templates[0],
      configuration: { ...templates[0].configuration, groupIds: ids },
    });
    expect(component.scopeSummary(scoped([]))).toBe("—");
    expect(component.scopeSummary(scoped(["1", "2"]))).toBe("1, 2");
    expect(component.scopeSummary(scoped(["1", "2", "3"]))).toBe("3 groups");
  });

  it("getTableData still flips loaded on a fetch error so the empty state doesn't hang", () => {
    component.loaded.set(false);
    runner.listTemplates.mockReturnValue(throwError(() => new Error("boom")));

    component.getTableData();

    expect(component.loaded()).toBe(true);
  });
});
