import {
  ChangeDetectionStrategy,
  Component,
  Inject,
  computed,
  inject,
  signal,
} from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormBuilder, FormGroup } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { startWith } from "rxjs";
import { ReportColumn } from "../../../open-api";
import {
  aggregateNeedsMeasure,
  ReportField,
  ReportFieldOption,
  toFieldOptions,
} from "../../models/report-catalog.constants";
import { ReportColumnValue } from "../../models/report-command.mapper";
import { deriveColumnName, validateFormulaExpr } from "../../models/report-column.util";
import { nextReportRowId } from "../../models/report-form.factory";

export interface ColumnPickerDialogData {
  dimensions: ReportField[];
  measures: ReportField[];
  existingColumns: ReportColumnValue[];
  column?: ReportColumnValue;
  /**
   * Opens the dialog with its field fixed: the kind step is unreachable (no Back)
   * and the Field picker is read-only, so only the label can be changed. Used to
   * rename the column a grouping level renders as — the level's field is chosen
   * in the Grouping section, not here.
   */
  lockField?: boolean;
}

type PickerStep = "kind" | "dim" | "agg" | "formula";

const STEP_TITLES: Record<PickerStep, string> = {
  kind: "Add a column",
  dim: "Dimension column",
  agg: "Aggregate column",
  formula: "Formula column",
};

const STEP_SUBTITLES: Record<PickerStep, string> = {
  kind: "What kind of column is this?",
  dim: "A label you can group / read by",
  agg: "A number summed across rows",
  formula: "Computed from other columns",
};

// A locked-field dialog is always on the dimension step, but it is not adding a
// dimension column — it renames the column a grouping level already produces.
const LOCKED_TITLE = "Grouping column";
const LOCKED_SUBTITLE = "The column this grouping level adds to the report";

/**
 * The column builder: pick a kind (dimension / aggregate / formula), then configure
 * it. Returns a ReportColumnValue on save (a new machine name is derived from the
 * label for a new column; an edited column keeps its name so formulas that
 * reference it stay valid). Formula validation is lightweight, inline feedback;
 * the backend is the authoritative validator.
 */
@Component({
  selector: "app-column-picker-dialog",
  templateUrl: "./column-picker-dialog.component.html",
  styleUrls: ["./column-picker-dialog.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ColumnPickerDialogComponent {
  private readonly formBuilder = inject(FormBuilder);
  private readonly dialogRef =
    inject<MatDialogRef<ColumnPickerDialogComponent, ReportColumnValue>>(MatDialogRef);

  public readonly step = signal<PickerStep>("kind");
  public readonly ReportColumnKind = ReportColumn.KindEnum;
  public readonly aggregateFunctions = [
    ReportColumn.AggFuncEnum.Sum,
    ReportColumn.AggFuncEnum.Count,
    ReportColumn.AggFuncEnum.Avg,
    ReportColumn.AggFuncEnum.Min,
    ReportColumn.AggFuncEnum.Max,
  ];
  public readonly operators = ["+", "-", "*", "/", "(", ")"];

  public readonly dimensions: ReportField[];
  public readonly measures: ReportField[];
  public readonly dimensionOptions: ReportFieldOption[];
  public readonly measureOptions: ReportFieldOption[];
  public readonly referenceableColumns: ReportColumnValue[];

  private readonly availableNames: string[];
  private readonly editingId: string | null;

  public readonly pickerForm: FormGroup = this.formBuilder.group({
    label: "",
    field: "",
    measure: "",
    aggFunc: ReportColumn.AggFuncEnum.Sum as ReportColumn.AggFuncEnum,
    expr: "",
  });

  public readonly exprValue = toSignal(
    this.pickerForm.get("expr")!.valueChanges.pipe(startWith("")),
    { initialValue: "" }
  );
  private readonly formTick = toSignal(this.pickerForm.valueChanges, { initialValue: null });

  public readonly formulaStatus = computed(() =>
    validateFormulaExpr(this.exprValue() ?? "", this.availableNames)
  );

  /** True when the dialog only edits the label (see ColumnPickerDialogData.lockField). */
  public readonly lockField: boolean;

  public readonly title = computed(() =>
    this.lockField ? LOCKED_TITLE : STEP_TITLES[this.step()]
  );
  public readonly subtitle = computed(() =>
    this.lockField ? LOCKED_SUBTITLE : STEP_SUBTITLES[this.step()]
  );

  public readonly canSave = computed<boolean>(() => {
    this.formTick();
    const value = this.pickerForm.getRawValue();
    switch (this.step()) {
      case "dim":
        return !!value.field && !!value.label.trim();
      case "agg":
        return (
          !!value.label.trim() &&
          (value.aggFunc === ReportColumn.AggFuncEnum.Count || !!value.measure)
        );
      case "formula":
        return !!value.label.trim() && this.formulaStatus().ok;
      default:
        return false;
    }
  });

  constructor(@Inject(MAT_DIALOG_DATA) data: ColumnPickerDialogData) {
    this.dimensions = data.dimensions;
    this.measures = data.measures;
    this.dimensionOptions = toFieldOptions(data.dimensions);
    this.measureOptions = toFieldOptions(data.measures);
    this.referenceableColumns = data.existingColumns.filter(
      (column) => column.kind !== ReportColumn.KindEnum.Dimension
    );
    this.availableNames = data.existingColumns.map((column) => column.name);
    this.editingId = data.column?.id ?? null;
    this.lockField = data.lockField === true;

    if (data.column) {
      this.seedFromColumn(data.column);
    }
    if (this.lockField) {
      // Disabled rather than merely presented read-only, so the field cannot be
      // changed through the control either. save() reads getRawValue(), so the
      // fixed field still rides the result.
      this.pickerForm.get("field")!.disable({ emitEvent: false });
    }

    // Auto-suggest the label on user-driven field/measure/function changes (seeds
    // above are emitEvent:false, so they never clobber an edited column's label).
    this.pickerForm.get("field")!.valueChanges.pipe(takeUntilDestroyed()).subscribe((key: string) =>
      this.pickerForm.get("label")!.setValue(this.labelForDimension(key))
    );
    this.pickerForm.get("measure")!.valueChanges.pipe(takeUntilDestroyed()).subscribe((key: string) =>
      this.pickerForm.get("label")!.setValue(this.labelForMeasure(key))
    );
    this.pickerForm.get("aggFunc")!.valueChanges.pipe(takeUntilDestroyed()).subscribe((fn: ReportColumn.AggFuncEnum) =>
      this.pickerForm.get("label")!.setValue(this.suggestAggregateLabel(fn))
    );
  }

  public get currentAggFunc(): ReportColumn.AggFuncEnum {
    return this.pickerForm.get("aggFunc")!.value;
  }

  public get aggregateNeedsMeasure(): boolean {
    return aggregateNeedsMeasure(this.currentAggFunc);
  }

  public pickDimension(): void {
    const first = this.dimensions[0];
    this.pickerForm.patchValue({ field: first?.key ?? "" }, { emitEvent: false });
    this.pickerForm.get("label")!.setValue(first?.label ?? "", { emitEvent: false });
    this.step.set("dim");
  }

  public pickAggregate(): void {
    const first = this.measures[0];
    this.pickerForm.patchValue(
      { aggFunc: ReportColumn.AggFuncEnum.Sum, measure: first?.key ?? "" },
      { emitEvent: false }
    );
    this.pickerForm.get("label")!.setValue(first?.label ?? "", { emitEvent: false });
    this.step.set("agg");
  }

  public pickFormula(): void {
    this.pickerForm.patchValue({ expr: "" }, { emitEvent: false });
    this.pickerForm.get("label")!.setValue("", { emitEvent: false });
    this.step.set("formula");
  }

  public selectFunction(fn: ReportColumn.AggFuncEnum): void {
    this.pickerForm.get("aggFunc")!.setValue(fn);
  }

  // The expression is built only by clicking column/operator chips (no typing), so
  // it is edited as a list of space-separated tokens: appending a chip, removing
  // the last token, or clearing. This keeps the stored value cleanly normalized.
  public insertColumn(name: string): void {
    this.setExprTokens([...this.exprTokens(), name]);
  }

  public insertOperator(operator: string): void {
    this.setExprTokens([...this.exprTokens(), operator]);
  }

  public removeLastToken(): void {
    const tokens = this.exprTokens();
    tokens.pop();
    this.setExprTokens(tokens);
  }

  public clearExpression(): void {
    this.pickerForm.get("expr")!.setValue("");
  }

  private exprTokens(): string[] {
    return ((this.pickerForm.get("expr")!.value as string) ?? "").trim().split(/\s+/).filter(Boolean);
  }

  private setExprTokens(tokens: string[]): void {
    this.pickerForm.get("expr")!.setValue(tokens.join(" "));
  }

  public back(): void {
    this.step.set("kind");
  }

  public save(): void {
    if (!this.canSave()) {
      return;
    }
    const value = this.pickerForm.getRawValue();
    const kind = this.stepKind();
    const label = value.label.trim();
    // Editing keeps the column's existing name so any formula referencing it stays
    // valid; a new column derives a fresh identifier from its label.
    const name = this.editingId
      ? this.editingColumnName
      : deriveColumnName(label, this.availableNames);

    const column: ReportColumnValue = {
      id: this.editingId ?? nextReportRowId(),
      kind,
      name,
      label,
    };
    if (kind === ReportColumn.KindEnum.Dimension) {
      column.field = value.field;
    } else if (kind === ReportColumn.KindEnum.Aggregate) {
      column.aggFunc = value.aggFunc;
      if (value.aggFunc !== ReportColumn.AggFuncEnum.Count) {
        column.measure = value.measure;
      }
    } else {
      column.expr = (value.expr as string).trim();
    }
    this.dialogRef.close(column);
  }

  public cancel(): void {
    this.dialogRef.close();
  }

  private seedFromColumn(column: ReportColumnValue): void {
    this.editingColumnName = column.name;
    this.pickerForm.patchValue(
      {
        label: column.label,
        field: column.field ?? "",
        measure: column.measure ?? "",
        aggFunc: column.aggFunc ?? ReportColumn.AggFuncEnum.Sum,
        expr: column.expr ?? "",
      },
      { emitEvent: false }
    );
    this.step.set(
      column.kind === ReportColumn.KindEnum.Dimension
        ? "dim"
        : column.kind === ReportColumn.KindEnum.Aggregate
          ? "agg"
          : "formula"
    );
  }

  private editingColumnName = "";

  private stepKind(): ReportColumn.KindEnum {
    switch (this.step()) {
      case "dim":
        return ReportColumn.KindEnum.Dimension;
      case "agg":
        return ReportColumn.KindEnum.Aggregate;
      default:
        return ReportColumn.KindEnum.Formula;
    }
  }

  private labelForDimension(key: string): string {
    return this.dimensions.find((field) => field.key === key)?.label ?? key;
  }

  private labelForMeasure(key: string): string {
    return this.measures.find((field) => field.key === key)?.label ?? key;
  }

  private suggestAggregateLabel(fn: ReportColumn.AggFuncEnum): string {
    if (fn === ReportColumn.AggFuncEnum.Count) {
      return "Count";
    }
    return this.labelForMeasure(this.pickerForm.get("measure")!.value);
  }
}
