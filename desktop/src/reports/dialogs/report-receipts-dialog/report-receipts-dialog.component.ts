import {
  ChangeDetectionStrategy,
  Component,
  computed,
  DestroyRef,
  Inject,
  inject,
  signal,
} from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { Receipt, ReportRequestCommand, ReportService } from "../../../open-api";
import { ReportBuilderValue } from "../../models/report-command.mapper";
import {
  formatPeriodRange,
  reportPeriodDateFieldLabel,
  resolvePeriodRange,
  toReportPeriodDateField,
} from "../../models/report-period.util";

export interface ReportReceiptsDialogData {
  // The command the preview ran, so the list covers exactly what it counted.
  command: ReportRequestCommand;
  // The builder's period, for the subtitle only.
  period: ReportBuilderValue["period"];
  // The report's true covered count (from the preview), shown in the subtitle;
  // falls back to the list's own total when absent.
  receiptCount?: number;
}

/**
 * Lists the receipts a report covers. The server resolves the report's filter and
 * period (on its own clock and time zone) and runs the report's own query, so the
 * list matches the preview's count; resolving the period here would use the
 * browser's time zone instead. Read-only — it exists so a user can sanity-check
 * what's flowing into the report.
 */
@Component({
  selector: "app-report-receipts-dialog",
  templateUrl: "./report-receipts-dialog.component.html",
  styleUrls: ["./report-receipts-dialog.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportReceiptsDialogComponent {
  private readonly reportService = inject(ReportService);
  private readonly dialogRef = inject(MatDialogRef<ReportReceiptsDialogComponent>);
  private readonly destroyRef = inject(DestroyRef);

  public readonly loading = signal<boolean>(true);
  public readonly receipts = signal<Receipt[]>([]);
  // Set when the list fails to load.
  public readonly hasError = signal<boolean>(false);
  // The receipt being inspected; null shows the list, non-null the breakdown.
  public readonly selected = signal<Receipt | null>(null);

  public readonly periodLabel: string;
  private readonly providedCount?: number;
  private readonly totalCount = signal<number>(0);
  // Subtitle count: the report's true total when known, else the list's own total.
  public readonly count = computed(() => this.providedCount ?? this.totalCount());

  constructor(@Inject(MAT_DIALOG_DATA) data: ReportReceiptsDialogData) {
    const range = formatPeriodRange(
      resolvePeriodRange(data.period.preset, data.period.startDate, data.period.endDate)
    );
    const dateField = reportPeriodDateFieldLabel(toReportPeriodDateField(data.period.dateField));
    this.periodLabel = `${range} on ${dateField}`;
    this.providedCount = data.receiptCount;
    this.load(data.command);
  }

  public viewReceipt(receipt: Receipt): void {
    this.selected.set(receipt);
  }

  public backToList(): void {
    this.selected.set(null);
  }

  /** Opens the receipt's full page in a new tab (read-only drill-in stays open). */
  public openFullReceipt(receipt: Receipt): void {
    window.open(`/receipts/${receipt.id}/view`, "_blank");
  }

  public close(): void {
    this.dialogRef.close();
  }

  private load(command: ReportRequestCommand): void {
    this.reportService
      .getReportReceipts(command)
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe({
        next: (paged) => {
          this.receipts.set((paged.data ?? []) as unknown as Receipt[]);
          this.totalCount.set(paged.totalCount ?? 0);
          this.loading.set(false);
        },
        error: () => {
          this.hasError.set(true);
          this.loading.set(false);
        },
      });
  }
}
