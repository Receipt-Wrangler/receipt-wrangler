import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { of } from "rxjs";
import { PagedRequestCommand, ReportRequestCommand, ReportService } from "../../open-api";
import { downloadFile } from "../../utils/file";
import { ReportRunnerService, reportFilename } from "./report-runner.service";

jest.mock("../../utils/file", () => ({ downloadFile: jest.fn() }));

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
    createReportTemplate: jest.Mock;
    getReportTemplates: jest.Mock;
    getReportTemplate: jest.Mock;
    duplicateReportTemplate: jest.Mock;
    deleteReportTemplate: jest.Mock;
    generateReport: jest.Mock;
  };
  let service: ReportRunnerService;

  beforeEach(() => {
    (downloadFile as jest.Mock).mockClear();
    reportService = {
      createReportTemplate: jest.fn(() => of({ id: 1 })),
      getReportTemplates: jest.fn(() => of({ data: [], totalCount: 0 })),
      getReportTemplate: jest.fn(() => of({ id: 3 })),
      duplicateReportTemplate: jest.fn(() => of({ id: 4 })),
      deleteReportTemplate: jest.fn(() => of(undefined)),
      generateReport: jest.fn(() => of(new Blob())),
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

  it("generateFromTemplate generates the stored config and triggers a download", () => {
    const cmd = command("Run", ["pdf"]);
    service.generateFromTemplate(cmd).subscribe();
    expect(reportService.generateReport).toHaveBeenCalledWith(cmd);
    expect(downloadFile).toHaveBeenCalled();
  });
});
