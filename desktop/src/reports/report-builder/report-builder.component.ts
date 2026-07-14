import {
  ChangeDetectionStrategy,
  Component,
  DestroyRef,
  computed,
  inject,
  signal,
} from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormArray, FormBuilder, FormGroup } from "@angular/forms";
import { UntilDestroy } from "@ngneat/until-destroy";
import { EMPTY, catchError, debounceTime, finalize, startWith, switchMap, take, tap } from "rxjs";
import { ReportPeriod, ReportRequestCommand } from "../../open-api";
import { ReportBuilderValue, toReportRequestCommand } from "../models/report-command.mapper";
import { buildReportForm } from "../models/report-form.factory";
import { ReportCatalogService } from "../services/report-catalog.service";
import { ReportRunnerService } from "../services/report-runner.service";

/**
 * The Report Builder screen: owns the report form, drives the debounced live
 * preview, and runs the synchronous generate-and-download. The config panel,
 * preview panel, and generate bar are presentational children bound to this form.
 * It is @UntilDestroy()-decorated because the shared receipt-filter form builder
 * (reused inside buildReportForm) registers untilDestroyed subscriptions.
 */
@UntilDestroy()
@Component({
  selector: "app-report-builder",
  templateUrl: "./report-builder.component.html",
  styleUrls: ["./report-builder.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportBuilderComponent {
  private readonly formBuilder = inject(FormBuilder);
  private readonly runner = inject(ReportRunnerService);
  private readonly catalogService = inject(ReportCatalogService);
  private readonly destroyRef = inject(DestroyRef);

  public readonly form: FormGroup = buildReportForm(this.formBuilder, this);

  public readonly previewHtml = signal<string>("");
  public readonly receiptCount = signal<number>(0);
  public readonly previewLoading = signal<boolean>(false);
  public readonly previewError = signal<boolean>(false);
  public readonly generating = signal<boolean>(false);

  // A tick that recomputes the run-readiness guards whenever the form changes.
  private readonly formTick = toSignal(this.form.valueChanges.pipe(startWith(null)), {
    initialValue: null,
  });

  /** Preview needs a group, at least one column, and (for custom) both dates. */
  public readonly canPreview = computed<boolean>(() => {
    this.formTick();
    return this.isRunnable();
  });

  /** Generate additionally needs at least one output format selected. */
  public readonly canGenerate = computed<boolean>(() => {
    this.formTick();
    return this.isRunnable() && this.selectedFormatCount() > 0;
  });

  constructor() {
    this.catalogService.load();

    this.form.valueChanges
      .pipe(
        startWith(null),
        tap(() => {
          if (this.isRunnable()) {
            this.previewLoading.set(true);
          }
        }),
        debounceTime(450),
        switchMap(() => {
          if (!this.isRunnable()) {
            this.previewHtml.set("");
            this.receiptCount.set(0);
            this.previewLoading.set(false);
            return EMPTY;
          }
          return this.runner.preview(this.currentCommand()).pipe(
            catchError(() => {
              this.previewError.set(true);
              this.previewLoading.set(false);
              return EMPTY;
            })
          );
        }),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe((preview) => {
        this.previewHtml.set(preview.html);
        this.receiptCount.set(preview.receiptCount);
        this.previewError.set(false);
        this.previewLoading.set(false);
      });
  }

  /** Generates the report and downloads it (a single file, or a zip of formats). */
  public generate(): void {
    if (!this.canGenerate() || this.generating()) {
      return;
    }
    this.generating.set(true);
    this.runner
      .generateAndDownload(this.currentCommand())
      .pipe(
        take(1),
        catchError(() => EMPTY), // the HTTP interceptor surfaces the failure toast
        finalize(() => this.generating.set(false)),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe();
  }

  private currentCommand(): ReportRequestCommand {
    return toReportRequestCommand(this.form.getRawValue() as ReportBuilderValue);
  }

  private isRunnable(): boolean {
    const scope = this.form.get("scope") as FormArray;
    const columns = this.form.get("columns") as FormArray;
    if (scope.length === 0 || columns.length === 0) {
      return false;
    }
    const period = this.form.get("period")!.value;
    if (period.preset === ReportPeriod.PresetEnum.Custom) {
      return !!period.startDate && !!period.endDate;
    }
    return true;
  }

  private selectedFormatCount(): number {
    const formats = this.form.get("formats")!.value as Record<string, boolean>;
    return Object.values(formats).filter(Boolean).length;
  }
}
