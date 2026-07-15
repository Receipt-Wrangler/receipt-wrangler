import { ReportRequestCommand } from "../../open-api";
import { reportFilename } from "./report-runner.service";

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
