import type { Mock } from "vitest";
import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { of } from "rxjs";
import { PagedRequestCommand, ReportRequestCommand, ReportService } from "../../open-api";
import { captureDownloads } from "../../testing/download.testing";
import { ReportRunnerService, reportFilename } from "./report-runner.service";

function command(name: string, formats: ReportRequestCommand.FormatsEnum[]): ReportRequestCommand {
  return {
    name,
    groupIds: ["1"],
    period: { preset: "this_month" } as ReportRequestCommand["period"],
    detail: { mode: "records" } as ReportRequestCommand["detail"],
    columns: [],
    formats,
  };
}

describe("reportFilename", () => {
  it("uses the single format's extension", () => {
    expect(reportFilename(command("My Report", ["csv"]))).toBe("My_Report.csv");
  });

  it("uses .zip for multiple formats", () => {
    expect(reportFilename(command("My Report", ["csv", "pdf"]))).toBe("My_Report.zip");
  });

  it("sanitizes non-word characters and trims underscores", () => {
    expect(reportFilename(command("Q2/2026 Expenses!", ["xlsx"]))).toBe("Q2_2026_Expenses.xlsx");
  });

  it("falls back to 'report' when the name is empty", () => {
    expect(reportFilename(command("   ", ["pdf"]))).toBe("report.pdf");
  });
});

describe("ReportRunnerService", () => {
  let reportService: {
    createReportTemplate: Mock;
    updateReportTemplate: Mock;
    getReportTemplates: Mock;
    getReportTemplate: Mock;
    duplicateReportTemplate: Mock;
    deleteReportTemplate: Mock;
    generateReport: Mock;
    generateReportFromTemplate: Mock;
    renderReportTemplate: Mock;
  };
  let service: ReportRunnerService;

  beforeEach(() => {
    reportService = {
      createReportTemplate: vi.fn(() => of({ id: 1 })),
      updateReportTemplate: vi.fn(() => of({ id: 2 })),
      getReportTemplates: vi.fn(() => of({ data: [], totalCount: 0 })),
      getReportTemplate: vi.fn(() => of({ id: 3 })),
      duplicateReportTemplate: vi.fn(() => of({ id: 4 })),
      deleteReportTemplate: vi.fn(() => of(undefined)),
      generateReport: vi.fn(() => of(new Blob())),
      generateReportFromTemplate: vi.fn(() => of(new Blob())),
      renderReportTemplate: vi.fn(() => of({ html: "", receiptCount: 0 })),
    };
    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        ReportRunnerService,
        { provide: ReportService, useValue: reportService },
      ],
    });
    service = TestBed.inject(ReportRunnerService);
  });

  it("saveTemplate delegates to createReportTemplate", () => {
    const cmd = command("Template", ["csv"]);
    service.saveTemplate(cmd).subscribe();
    expect(reportService.createReportTemplate).toHaveBeenCalledWith(cmd);
  });

  it("updateTemplate delegates to updateReportTemplate with the id", () => {
    const cmd = command("Template", ["csv"]);
    service.updateTemplate(2, cmd).subscribe();
    expect(reportService.updateReportTemplate).toHaveBeenCalledWith(2, cmd);
  });

  it("listTemplates delegates to getReportTemplates", () => {
    const paged: PagedRequestCommand = { page: 1, pageSize: 10 };
    service.listTemplates(paged).subscribe();
    expect(reportService.getReportTemplates).toHaveBeenCalledWith(paged);
  });

  it("getTemplate delegates to getReportTemplate", () => {
    service.getTemplate(3).subscribe();
    expect(reportService.getReportTemplate).toHaveBeenCalledWith(3);
  });

  it("duplicateTemplate delegates to duplicateReportTemplate", () => {
    service.duplicateTemplate(4).subscribe();
    expect(reportService.duplicateReportTemplate).toHaveBeenCalledWith(4);
  });

  it("deleteTemplate delegates to deleteReportTemplate", () => {
    service.deleteTemplate(5).subscribe();
    expect(reportService.deleteReportTemplate).toHaveBeenCalledWith(5);
  });

  it("generateFromTemplateById generates by id and triggers a download named from the config", () => {
    const template = {
      id: 7,
      name: "Run",
      createdAt: "2026-01-01T00:00:00Z",
      configuration: command("Run", ["pdf"]),
      configurationVersion: 1,
    };
    const downloads = captureDownloads();
    service.generateFromTemplateById(template).subscribe();
    expect(reportService.generateReportFromTemplate).toHaveBeenCalledWith(7);
    expect(downloads).toEqual([{ blob: expect.any(Blob), filename: "Run.pdf" }]);
  });

  it("renderTemplate delegates to renderReportTemplate", () => {
    service.renderTemplate(9).subscribe();
    expect(reportService.renderReportTemplate).toHaveBeenCalledWith(9);
  });

  it("downloadTemplateById resolves the template then generates + downloads by id", () => {
    reportService.getReportTemplate.mockReturnValue(
      of({
        id: 9,
        name: "Widget",
        createdAt: "2026-01-01T00:00:00Z",
        configuration: command("Widget", ["csv"]),
        configurationVersion: 1,
      })
    );

    const downloads = captureDownloads();

    service.downloadTemplateById(9).subscribe();

    expect(reportService.getReportTemplate).toHaveBeenCalledWith(9);
    expect(reportService.generateReportFromTemplate).toHaveBeenCalledWith(9);
    expect(downloads).toEqual([{ blob: expect.any(Blob), filename: "Widget.csv" }]);
  });
});
