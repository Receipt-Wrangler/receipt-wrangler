import { ChangeDetectionStrategy, Component, computed, inject, input, signal } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { DomSanitizer, SafeHtml } from "@angular/platform-browser";
import { DEFAULT_DIALOG_CONFIG } from "src/constants/dialog.constant";
import {
  ReportReceiptsDialogComponent,
  ReportReceiptsDialogData,
} from "../dialogs/report-receipts-dialog/report-receipts-dialog.component";

/**
 * The live preview pane: renders the engine's HTML for the current configuration
 * inside a sandboxed iframe (scripts disabled; the HTML is trusted engine output
 * but isolated as defense-in-depth) and shows the covered receipt count, which
 * opens the receipts drill-in.
 */
@Component({
  selector: "app-report-preview-panel",
  templateUrl: "./report-preview-panel.component.html",
  styleUrls: ["./report-preview-panel.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportPreviewPanelComponent {
  public readonly form = input.required<FormGroup>();
  public readonly html = input<string>("");
  public readonly receiptCount = input<number>(0);
  public readonly loading = input<boolean>(false);
  public readonly error = input<boolean>(false);
  public readonly ready = input<boolean>(false);

  private readonly sanitizer = inject(DomSanitizer);
  private readonly dialog = inject(MatDialog);

  public readonly frameHeight = signal<number>(560);
  public readonly safeHtml = computed<SafeHtml>(() =>
    this.sanitizer.bypassSecurityTrustHtml(this.html())
  );

  /** Sizes the iframe to its content (same-origin srcdoc, no scripts allowed). */
  public onFrameLoad(iframe: HTMLIFrameElement): void {
    const doc = iframe.contentDocument;
    if (doc?.documentElement) {
      this.frameHeight.set(doc.documentElement.scrollHeight + 24);
    }
  }

  public openReceipts(): void {
    if (this.receiptCount() === 0) {
      return;
    }
    const value = this.form().getRawValue();
    const data: ReportReceiptsDialogData = {
      groupIds: value.scope,
      filter: value.filter,
      period: value.period,
    };
    this.dialog.open(ReportReceiptsDialogComponent, { ...DEFAULT_DIALOG_CONFIG, width: "72%", data });
  }
}
