import { Injectable, inject } from "@angular/core";
import { Observable, tap } from "rxjs";
import {
  PagedData,
  PagedRequestCommand,
  ReportPreviewResponse,
  ReportRequestCommand,
  ReportService,
  ReportTemplate,
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

  /** Saves the current configuration as a new reusable report template. */
  public saveTemplate(command: ReportRequestCommand): Observable<ReportTemplate> {
    return this.reportService.createReportTemplate(command);
  }

  /** Overwrites an existing template in place with the current configuration. */
  public updateTemplate(id: number, command: ReportRequestCommand): Observable<ReportTemplate> {
    return this.reportService.updateReportTemplate(id, command);
  }

  /** A page of saved report templates for the list. */
  public listTemplates(command: PagedRequestCommand): Observable<PagedData> {
    return this.reportService.getReportTemplates(command);
  }

  /** One saved template by id — used to hydrate the builder for editing. */
  public getTemplate(id: number): Observable<ReportTemplate> {
    return this.reportService.getReportTemplate(id);
  }

  /** Copies a saved template into a new one (name suffixed by the backend). */
  public duplicateTemplate(id: number): Observable<ReportTemplate> {
    return this.reportService.duplicateReportTemplate(id);
  }

  /** Deletes a saved template. */
  public deleteTemplate(id: number): Observable<unknown> {
    return this.reportService.deleteReportTemplate(id);
  }

  /**
   * Runs a saved template by id: the server loads its stored configuration and
   * enforces the per-template generate grant (which the ad-hoc /generate cannot).
   * The download filename is derived from the template's own configuration.
   */
  public generateFromTemplateById(template: ReportTemplate): Observable<Blob> {
    return this.reportService
      .generateReportFromTemplate(template.id)
      .pipe(tap((blob) => downloadFile(blob, reportFilename(template.configuration))));
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
