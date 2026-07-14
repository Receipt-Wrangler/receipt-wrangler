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
import { catchError, forkJoin, map, of, take } from "rxjs";
import {
  FilterOperation,
  Receipt,
  ReceiptPagedRequestCommand,
  ReceiptPagedRequestFilter,
  ReceiptService,
  ReportPeriod,
} from "../../../open-api";
import { formatPeriodRange, resolvePeriodRange } from "../../models/report-period.util";

export interface ReportReceiptsDialogData {
  groupIds: string[];
  filter: ReceiptPagedRequestFilter;
  period: { preset: ReportPeriod.PresetEnum; startDate: Date | null; endDate: Date | null };
  // The report's true covered count (from the preview), shown in the subtitle;
  // falls back to the loaded list length when absent.
  receiptCount?: number;
}

// Bounds the drill-in fetch per group; the count chip still reports the true total.
const DRILL_IN_PAGE_SIZE = 200;

/**
 * Lists the receipts a report covers: the report's filter narrowed to the resolved
 * period, fetched across every scope group and merged. Read-only — it exists so a
 * user can sanity-check what's flowing into the report.
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
  // The receipt being inspected; null shows the list, non-null the breakdown.
  public readonly selected = signal<Receipt | null>(null);

  public readonly periodLabel: string;
  private readonly providedCount?: number;
  // Subtitle count: the report's true total when known, else the loaded count.
  public readonly count = computed(() => this.providedCount ?? this.receipts().length);

  constructor(@Inject(MAT_DIALOG_DATA) data: ReportReceiptsDialogData) {
    this.periodLabel = formatPeriodRange(
      resolvePeriodRange(data.period.preset, data.period.startDate, data.period.endDate)
    );
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
    const command: ReceiptPagedRequestCommand = {
      page: 1,
      pageSize: DRILL_IN_PAGE_SIZE,
      filter: {
        ...data.filter,
        date: { operation: FilterOperation.Between, value: [range.start, range.end] },
      },
    };

    if (data.groupIds.length === 0) {
      this.loading.set(false);
      return;
    }

    const requests = data.groupIds.map((groupId) =>
      this.receiptService.getReceiptsForGroup(Number.parseInt(groupId, 10), command).pipe(
        take(1),
        map((paged) => (paged.data ?? []) as unknown as Receipt[]),
        catchError(() => of<Receipt[]>([]))
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
