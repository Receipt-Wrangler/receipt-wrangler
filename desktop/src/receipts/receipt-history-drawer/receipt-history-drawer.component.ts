import { DOCUMENT } from "@angular/common";
import {
  Component,
  DestroyRef,
  ElementRef,
  LOCALE_ID,
  afterNextRender,
  computed,
  inject,
  input,
  output,
  signal,
  viewChild,
} from "@angular/core";
import { rxResource } from "@angular/core/rxjs-interop";
import { MatDialog } from "@angular/material/dialog";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { Store } from "@ngxs/store";
import { ButtonModule } from "../../button";
import { ReceiptService } from "../../open-api";
import { SystemSettingsState } from "../../store/system-settings.state";
import { ReceiptHistoryEntryComponent } from "./receipt-history-entry/receipt-history-entry.component";
import { groupHistoryByDay, historySummary } from "./receipt-history.utils";

let nextDrawerId = 0;

/**
 * The receipt change history: a right-hand drawer over the receipt view page,
 * listing the receipt's saves newest first, grouped by app-zone day.
 *
 * The parent renders it only while it is open, so every open is a fresh
 * instance and a fresh request — the history is never stale by more than the
 * time the drawer has been open.
 *
 * It is `position: fixed` with its top aligned to the bottom of the sticky app
 * header, measured once when it opens. It stacks inside `mat-drawer-content`'s
 * stacking context (`z-index: 1`), so the CDK overlay container — dialogs,
 * menus, tooltips and snackbars — still paints above it.
 */
@Component({
  selector: "app-receipt-history-drawer",
  standalone: true,
  imports: [ButtonModule, MatProgressSpinnerModule, ReceiptHistoryEntryComponent],
  templateUrl: "./receipt-history-drawer.component.html",
  styleUrl: "./receipt-history-drawer.component.scss",
  host: {
    role: "complementary",
    "data-testid": "receipt-history-drawer",
    "[attr.aria-labelledby]": "titleId",
    "[style.top.px]": "frame()?.top ?? null",
    "[style.max-width]": "maxWidth()",
    "(document:keydown.escape)": "onEscape($event)",
  },
})
export class ReceiptHistoryDrawerComponent {
  public readonly receiptId = input.required<number>();

  public readonly closed = output<void>();

  public readonly titleId = `rw-history-title-${nextDrawerId++}`;

  private readonly document = inject(DOCUMENT);

  private readonly host = inject<ElementRef<HTMLElement>>(ElementRef);

  private readonly matDialog = inject(MatDialog);

  private readonly receiptService = inject(ReceiptService);

  private readonly locale = inject(LOCALE_ID);

  private readonly timeZone = inject(Store).selectSignal(SystemSettingsState.timeZone);

  private readonly closeButton = viewChild("closeButton", { read: ElementRef });

  public readonly history = rxResource({
    params: () => ({ receiptId: this.receiptId() }),
    stream: ({ params }) => this.receiptService.getReceiptHistory(params.receiptId),
  });

  public readonly summary = computed(() =>
    this.history.hasValue() ? historySummary(this.history.value()) : ""
  );

  public readonly days = computed(() =>
    this.history.hasValue()
      ? groupHistoryByDay(this.history.value().entries ?? [], this.timeZone(), this.locale)
      : []
  );

  /**
   * Where the app header ends: its bottom (the drawer's top) and its left edge
   * (the content area's, so the drawer never slides under an open sidebar on a
   * narrow screen). Null until measured, when the stylesheet's fallback applies.
   */
  public readonly frame = signal<{ top: number; left: number } | null>(null);

  public readonly maxWidth = computed(() => {
    const frame = this.frame();
    return frame ? `calc(100vw - ${frame.left}px)` : null;
  });

  constructor() {
    // The toggle that opened the drawer, to hand focus back to on close.
    const opener = this.document.activeElement;

    afterNextRender(() => {
      this.measureFrame();
      this.closeButton()?.nativeElement.querySelector("button")?.focus();
    });

    inject(DestroyRef).onDestroy(() => this.restoreFocus(opener));
  }

  public onEscape(event: Event): void {
    // A dialog, menu or autocomplete above the drawer handles its own Esc
    // (Material marks it handled); that press must not also close the drawer.
    if (event.defaultPrevented || this.matDialog.openDialogs.length > 0) {
      return;
    }

    this.closed.emit();
  }

  public close(): void {
    this.closed.emit();
  }

  public retry(): void {
    this.history.reload();
  }

  // The app header is sticky at the top of the scrolling content area, so its
  // bottom edge does not move as the page scrolls; measuring once is enough.
  private measureFrame(): void {
    const rect = this.document.querySelector("app-header")?.getBoundingClientRect();
    if (rect && rect.height > 0) {
      this.frame.set({ top: Math.round(rect.bottom), left: Math.max(0, Math.round(rect.left)) });
    }
  }

  // Only when focus was in the drawer (or was lost with it): a user who moved on
  // to something else on the page keeps their place.
  private restoreFocus(opener: Element | null): void {
    const active = this.document.activeElement;
    const focusWasInDrawer =
      !active || active === this.document.body || this.host.nativeElement.contains(active);
    if (focusWasInDrawer && opener instanceof HTMLElement && opener.isConnected) {
      opener.focus();
    }
  }
}
