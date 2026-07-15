import { Injectable, inject } from "@angular/core";
import { Observable, tap } from "rxjs";
import {
  ReportPreviewResponse,
  ReportRequestCommand,
  ReportService,
} from "../../open-api";
import { downloadFile } from "../../utils/file";

/**
 * Runs a report configuration against the backend for the builder: a one-shot
 * generate-and-download, and a live preview. Both feed the same
 * ReportRequestCommand to the same synchronous endpoints; neither needs NGXS
 * state (mirrors ReceiptExportService).
 */
@Injectable({ providedIn: "root" })
export class ReportRunnerService {
  private readonly reportService = inject(ReportService);

  /**
   * Generates the report and triggers a browser download of the returned file (or
   * zip). Emits once when the download has been kicked off, so the caller can drop
   * its in-flight state.
   */
  public generateAndDownload(command: ReportRequestCommand): Observable<Blob> {
    return this.reportService
      .generateReport(command)
      .pipe(tap((blob) => downloadFile(blob, reportFilename(command))));
  }

  /** Renders the current configuration to preview HTML + the covered receipt count. */
  public preview(command: ReportRequestCommand): Observable<ReportPreviewResponse> {
    return this.reportService.previewReport(command);
  }
}

/**
 * Computes the download filename from the report name and selected formats — a
 * single format gets that extension, several get a `.zip` (matching how the
 * backend bundles them). Mirrors the backend's own filename sanitization.
 */
export function reportFilename(command: ReportRequestCommand): string {
  const base =
    (command.name ?? "").trim().replace(/[^\w]+/g, "_").replace(/^_+|_+$/g, "") ||
    "report";
  const formats = command.formats ?? [];
  if (formats.length > 1) {
    return base + ".zip";
  }
  return base + "." + (formats[0] ?? "csv");
}
