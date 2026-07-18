import { CommonModule } from "@angular/common";
import { ChangeDetectionStrategy, Component, DestroyRef, computed, effect, inject, input, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { MatIconModule } from "@angular/material/icon";
import { DomSanitizer, SafeHtml } from "@angular/platform-browser";
import { EMPTY, Subscription, catchError, finalize, tap } from "rxjs";
import { ButtonModule } from "../../button/index";
import { Widget } from "../../open-api";
import { ReportRunnerService } from "../../reports/services/report-runner.service";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

/**
 * A view-only dashboard widget that renders a saved report template's HTML, exactly
 * as the report preview does — a sandboxed iframe (scripts disabled) fed the engine
 * output. The pinned template is stored as `{ reportTemplateId }` in the widget
 * configuration; the server renders the FULL dataset and re-resolves access on every
 * load, so when the user's access is revoked (or the template is deleted) the
 * endpoint returns restricted-notice HTML at 200 — the widget always drops whatever
 * HTML it gets into the iframe, no special "restricted" branch. A download button
 * shows only when the server-computed allowedActions include "generate".
 */
@Component({
  selector: "app-report-widget",
  templateUrl: "./report-widget.component.html",
  styleUrls: ["./report-widget.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule, ButtonModule, MatIconModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportWidgetComponent {
  public readonly widget = input.required<Widget>();

  private readonly runner = inject(ReportRunnerService);
  private readonly sanitizer = inject(DomSanitizer);
  private readonly destroyRef = inject(DestroyRef);

  public readonly html = signal<string>("");
  public readonly allowedActions = signal<string[]>([]);
  public readonly isLoading = signal<boolean>(true);
  public readonly hasError = signal<boolean>(false);
  public readonly isDownloading = signal<boolean>(false);
  public readonly frameHeight = signal<number>(320);

  public readonly reportTemplateId = computed<number | undefined>(() => {
    const config = this.widget()?.configuration as { reportTemplateId?: number } | undefined;
    return config?.reportTemplateId ?? undefined;
  });

  // Gate the download button purely on the server-computed allowedActions (which
  // already bake in the base/"*All" report permissions, the per-group ceiling, and
  // the per-template matrix) — never AND-ed with a client permission check.
  public readonly canDownload = computed<boolean>(
    () => !!this.reportTemplateId() && this.allowedActions().includes("generate")
  );

  public readonly safeHtml = computed<SafeHtml>(() =>
    this.sanitizer.bypassSecurityTrustHtml(this.html())
  );

  constructor() {
    // Load whenever the pinned template changes (the effect tracks the
    // widget-config-derived id). Replaces ngOnChanges for this input-watching load.
    // onCleanup cancels the prior render before the next one starts (and on destroy),
    // so a slower in-flight request for the old template can never overwrite the new
    // one — takeUntilDestroyed alone only cancels on destroy, not on re-run.
    effect((onCleanup) => {
      const subscription = this.load(this.reportTemplateId());
      onCleanup(() => subscription?.unsubscribe());
    });
  }

  /** Sizes the iframe to its content; the stage caps the height and scrolls. */
  public onFrameLoad(iframe: HTMLIFrameElement): void {
    const doc = iframe.contentDocument;
    if (doc?.documentElement) {
      this.frameHeight.set(doc.documentElement.scrollHeight + 8);
    }
  }

  public download(): void {
    const id = this.reportTemplateId();
    if (!id || this.isDownloading()) {
      return;
    }
    this.isDownloading.set(true);
    this.runner
      .downloadTemplateById(id)
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        catchError(() => EMPTY), // the HTTP interceptor surfaces the failure toast
        finalize(() => this.isDownloading.set(false))
      )
      .subscribe();
  }

  private load(id: number | undefined): Subscription | undefined {
    // Clear the previous template's output up front so a swap can't briefly show the
    // old report or (worse) enable download using the old template's authorization
    // result while the new one loads.
    this.html.set("");
    this.allowedActions.set([]);

    if (!id) {
      this.isLoading.set(false);
      this.hasError.set(true);
      return;
    }
    this.isLoading.set(true);
    this.hasError.set(false);
    return this.runner
      .renderTemplate(id)
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        tap((response) => {
          this.html.set(response.html ?? "");
          this.allowedActions.set(response.allowedActions ?? []);
        }),
        catchError(() => {
          this.hasError.set(true);
          return EMPTY;
        }),
        finalize(() => this.isLoading.set(false))
      )
      .subscribe();
  }
}
