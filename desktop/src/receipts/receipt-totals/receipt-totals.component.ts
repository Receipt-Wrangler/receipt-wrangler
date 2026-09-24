import { Component, computed, input, output } from "@angular/core";
import { MatChipsModule } from "@angular/material/chips";
import { Group, ReceiptSummary, ReceiptSummaryPosition, ReceiptSummaryRow } from "../../open-api";
import { PipesModule } from "../../pipes";
import { formatStatus } from "../../utils";

/**
 * The block of totals under the receipts table: a receipt count and amount total
 * over the whole current filter, then the same figures for each status the group
 * has configured.
 *
 * Deliberately presentational — it fetches nothing and knows nothing about the
 * filter. `receipts-table` owns the request (on its own refresh stream, which
 * skips paging and sorting because neither changes the figures) and hands the
 * result down.
 *
 * Named `app-receipt-totals`, not anything with "summary" in it: `app-summary-card`
 * already renders on this same page as "Selected Receipt summary" — the
 * who-owes-whom settlement card — and two components called summary on one screen
 * is a trap for whoever reads this next.
 */
@Component({
  selector: "app-receipt-totals",
  standalone: true,
  imports: [MatChipsModule, PipesModule],
  templateUrl: "./receipt-totals.component.html",
  styleUrl: "./receipt-totals.component.scss",
})
export class ReceiptTotalsComponent {
  public readonly summary = input<ReceiptSummary | undefined>(undefined);

  /**
   * Groups whose configuration the viewer may pick between. Only populated on the
   * synthetic "All" group, which spans several groups and so has no configuration
   * of its own; empty everywhere else, where the group being viewed decides.
   */
  public readonly configGroups = input<Group[]>([]);

  public readonly selectedConfigGroupId = input<number | undefined>(undefined);

  public readonly configGroupSelected = output<number>();

  /**
   * Where the block sits relative to the table. The parent decides WHICH anchor renders
   * it; this only tells the component which of its own edges faces the table, so it keeps
   * owning its vertical rhythm instead of the receipts page reaching in with a margin.
   */
  public readonly position = input<ReceiptSummaryPosition>(ReceiptSummaryPosition.Bottom);

  public readonly isTop = computed(() => this.position() === ReceiptSummaryPosition.Top);

  public readonly show = computed(() => !!this.summary()?.enabled);

  /**
   * The overall row first, then one per configured status. Flattened here rather
   * than rendered through two loops (or a template outlet) so the row markup
   * exists once.
   */
  public readonly rows = computed(() => {
    const summary = this.summary();
    if (!summary?.enabled) {
      return [];
    }

    return [summary.overall, ...(summary.statuses ?? [])];
  });

  /**
   * The chip row is a choice, so it renders only when there is one to make. With
   * exactly one configured group the caption below says whose configuration
   * produced the numbers instead — same information, no false affordance.
   */
  public readonly showConfigGroupChips = computed(() => this.configGroups().length > 1);

  public readonly soleConfigGroupName = computed(() => {
    const groups = this.configGroups();
    return groups.length === 1 ? groups[0].name : "";
  });

  public rowLabel(row: ReceiptSummaryRow): string {
    // The overall row carries no status, and its label names the filter rather than a status.
    return row.status ? `${formatStatus(row.status)} Receipts` : "All Receipts";
  }

  public receiptCountLabel(row: ReceiptSummaryRow): string {
    return row.receiptCount === 1 ? "1 receipt" : `${row.receiptCount} receipts`;
  }
}
