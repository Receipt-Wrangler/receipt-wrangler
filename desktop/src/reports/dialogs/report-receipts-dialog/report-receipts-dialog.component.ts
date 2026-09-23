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
import { endOfDay, startOfDay } from "date-fns";
import { catchError, forkJoin, map, of, take } from "rxjs";
import {
  FilterOperation,
  Receipt,
  ReceiptPagedRequestCommand,
  ReceiptPagedRequestFilter,
  ReceiptService,
} from "../../../open-api";
import { ReportBuilderValue } from "../../models/report-command.mapper";
import {
  formatPeriodRange,
  reportPeriodDateFieldLabel,
  resolvePeriodRange,
  toReportPeriodDateField,
} from "../../models/report-period.util";

export interface ReportReceiptsDialogData {
  groupIds: string[];
  filter: ReceiptPagedRequestFilter;
  period: ReportBuilderValue["period"];
  // The report's true covered count (from the preview), shown in the subtitle;
  // falls back to the loaded list length when absent.
  receiptCount?: number;
}

// Bounds the drill-in fetch per group; the count chip still reports the true total.
const DRILL_IN_PAGE_SIZE = 200;

/**
 * Lists the receipts a report covers: the report's filter narrowed to the resolved
 * period on the date field it covers, fetched across every scope group and merged.
 * Read-only — it exists so a user can sanity-check what's flowing into the report.
 */
@Component({
  selector: "app-report-receipts-dialog",
  templateUrl: "./report-receipts-dialog.component.html",
  styleUrls: ["./report-receipts-dialog.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportReceiptsDialogComponent {
  private readonly receiptService = inject(ReceiptService);
  private readonly dialogRef = inject(MatDialogRef<ReportReceiptsDialogComponent>);
  private readonly destroyRef = inject(DestroyRef);

  public readonly loading = signal<boolean>(true);
  public readonly receipts = signal<Receipt[]>([]);
  // Set when any group's fetch fails, so the list can warn it may be incomplete.
  public readonly hasError = signal<boolean>(false);
  // The receipt being inspected; null shows the list, non-null the breakdown.
  public readonly selected = signal<Receipt | null>(null);

  public readonly periodLabel: string;
  private readonly providedCount?: number;
  // Subtitle count: the report's true total when known, else the loaded count.
  public readonly count = computed(() => this.providedCount ?? this.receipts().length);

  constructor(@Inject(MAT_DIALOG_DATA) data: ReportReceiptsDialogData) {
    const range = formatPeriodRange(
      resolvePeriodRange(data.period.preset, data.period.startDate, data.period.endDate)
    );
    const dateField = reportPeriodDateFieldLabel(toReportPeriodDateField(data.period.dateField));
    this.periodLabel = `${range} on ${dateField}`;
    this.providedCount = data.receiptCount;
    this.load(data);
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

  private load(data: ReportReceiptsDialogData): void {
    const range = resolvePeriodRange(data.period.preset, data.period.startDate, data.period.endDate);
    // Like the report, the period replaces whatever condition its date field held
    // and leaves the other date fields' conditions alone. The bounds span whole
    // days as the server's do: a custom range ends at local midnight, which would
    // drop its last day on the timestamped Added At and Resolved Date columns.
    const filter: ReceiptPagedRequestFilter = { ...data.filter };
    filter[toReportPeriodDateField(data.period.dateField)] = {
      operation: FilterOperation.Between,
      value: [startOfDay(range.start), endOfDay(range.end)],
    };
    const command: ReceiptPagedRequestCommand = {
      page: 1,
      pageSize: DRILL_IN_PAGE_SIZE,
      filter,
    };

    if (data.groupIds.length === 0) {
      this.loading.set(false);
      return;
    }

    const requests = data.groupIds.map((groupId) =>
      this.receiptService.getReceiptsForGroup(Number.parseInt(groupId, 10), command).pipe(
        take(1),
        map((paged) => (paged.data ?? []) as unknown as Receipt[]),
        catchError(() => {
          this.hasError.set(true);
          return of<Receipt[]>([]);
        })
      )
    );

    forkJoin(requests)
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((results) => {
        this.receipts.set(results.flat());
        this.loading.set(false);
      });
  }
}
