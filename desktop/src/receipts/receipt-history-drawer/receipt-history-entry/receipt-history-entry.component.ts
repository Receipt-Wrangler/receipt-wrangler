import { CurrencyPipe } from "@angular/common";
import { Component, LOCALE_ID, computed, inject, input, signal } from "@angular/core";
import { rxResource } from "@angular/core/rxjs-interop";
import { MatIconModule } from "@angular/material/icon";
import { Store } from "@ngxs/store";
import { ReceiptHistoryEntry, ReceiptService } from "../../../open-api";
import { formatAppDate, AppDatePipe } from "../../../pipes/app-date.pipe";
import { CustomCurrencyPipe } from "../../../pipes/custom-currency.pipe";
import { UserState } from "../../../store/user.state";
import { SystemSettingsState } from "../../../store/system-settings.state";
import {
  HISTORY_COLLAPSE_AFTER,
  HistoryValueFormatters,
  UserLookup,
  buildHistoryChangeRow,
  historyActionVerb,
  historyActorDisplay,
  historyBeforeStateNotice,
  historySnapshotsJson,
} from "../receipt-history.utils";
import { AvatarComponent } from "../../../avatar/avatar/avatar.component";

/**
 * One save in the receipt history timeline: who and when, the field-level
 * changes (collapsed after {@link HISTORY_COLLAPSE_AFTER}), a failure box or an
 * older-version notice, and the raw JSON on demand.
 *
 * The JSON is fetched lazily — the history list carries no snapshots, so the
 * request is only made the first time "View JSON" is opened, and only once:
 * hiding it keeps the loaded value.
 */
@Component({
  selector: "app-receipt-history-entry",
  standalone: true,
  imports: [AvatarComponent, MatIconModule, AppDatePipe],
  // CustomCurrencyPipe needs CurrencyPipe;
  // neither is providedIn: "root" (the receipts-table precedent).
  providers: [CurrencyPipe, CustomCurrencyPipe],
  templateUrl: "./receipt-history-entry.component.html",
  styleUrl: "./receipt-history-entry.component.scss",
})
export class ReceiptHistoryEntryComponent {
  public readonly receiptId = input.required<number>();

  public readonly entry = input.required<ReceiptHistoryEntry>();

  private readonly store = inject(Store);

  private readonly receiptService = inject(ReceiptService);

  private readonly customCurrencyPipe = inject(CustomCurrencyPipe);

  private readonly locale = inject(LOCALE_ID);

  private readonly users = this.store.selectSignal(UserState.users);

  private readonly timeZone = this.store.selectSignal(SystemSettingsState.timeZone);

  private readonly findUser = computed<UserLookup>(() => {
    const users = this.users();
    return (userId) => (userId ? users.find((user) => user.id === userId) : undefined);
  });

  public readonly actor = computed(() => historyActorDisplay(this.entry().actor, this.findUser()));

  public readonly verb = computed(() => historyActionVerb(this.entry().action));

  public readonly failed = computed(() => this.entry().action === "FAILED" || !!this.entry().errorMessage);

  public readonly notice = computed(() =>
    historyBeforeStateNotice(this.entry(), this.timeZone(), this.locale)
  );

  public readonly rows = computed(() => {
    const zone = this.timeZone();
    const formatters: HistoryValueFormatters = {
      currency: (value) => this.customCurrencyPipe.transform(value),
      calendarDate: (value) => formatAppDate(value, "mediumDate", "calendar", zone, this.locale),
      instant: (value) => formatAppDate(value, "medium", "instant", zone, this.locale),
      findUser: this.findUser(),
    };

    return (this.entry().changes ?? []).map((change) => buildHistoryChangeRow(change, formatters));
  });

  public readonly expanded = signal(false);

  public readonly hiddenCount = computed(() => Math.max(0, this.rows().length - HISTORY_COLLAPSE_AFTER));

  public readonly visibleRows = computed(() =>
    this.expanded() ? this.rows() : this.rows().slice(0, HISTORY_COLLAPSE_AFTER)
  );

  public readonly changesId = computed(() => `rw-history-changes-${this.entry().systemTaskId}`);

  public readonly jsonId = computed(() => `rw-history-json-${this.entry().systemTaskId}`);

  public readonly jsonOpen = signal(false);

  // Flipped once and never back, so the resource's params settle after the first
  // open and "Hide JSON" / "View JSON" never refetch.
  private readonly jsonRequested = signal(false);

  public readonly snapshots = rxResource({
    params: () =>
      this.jsonRequested()
        ? { receiptId: this.receiptId(), systemTaskId: this.entry().systemTaskId }
        : undefined,
    stream: ({ params }) => this.receiptService.getReceiptHistorySnapshots(params.receiptId, params.systemTaskId),
  });

  public readonly json = computed(() =>
    this.snapshots.hasValue() ? historySnapshotsJson(this.snapshots.value()) : null
  );

  public toggleExpanded(): void {
    this.expanded.update((expanded) => !expanded);
  }

  public toggleJson(): void {
    const open = !this.jsonOpen();
    this.jsonOpen.set(open);
    if (!open) {
      return;
    }

    if (!this.jsonRequested()) {
      this.jsonRequested.set(true);
    } else if (this.snapshots.status() === "error") {
      this.snapshots.reload();
    }
  }
}
