import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { PipesModule } from "../../../pipes";
import { ReportColumn } from "../../../open-api";
import { ColumnPickerDialogComponent, ColumnPickerDialogData } from "./column-picker-dialog.component";

function configure(data: ColumnPickerDialogData): {
  fixture: ComponentFixture<ColumnPickerDialogComponent>;
  component: ColumnPickerDialogComponent;
  close: jest.Mock;
} {
  const close = jest.fn();
  TestBed.configureTestingModule({
    declarations: [ColumnPickerDialogComponent],
    imports: [PipesModule],
    providers: [
      provideZonelessChangeDetection(),
      FormBuilder,
      { provide: MatDialogRef, useValue: { close } },
      { provide: MAT_DIALOG_DATA, useValue: data },
    ],
    schemas: [NO_ERRORS_SCHEMA],
  });
  const fixture = TestBed.createComponent(ColumnPickerDialogComponent);
  return { fixture, component: fixture.componentInstance, close };
}

const baseData: ColumnPickerDialogData = {
  dimensions: [{ key: "category", label: "Category" }],
  measures: [{ key: "amount", label: "Amount" }],
  existingColumns: [],
};

// A picker seeded with two aggregate columns a formula can reference by name.
const formulaData: ColumnPickerDialogData = {
  ...baseData,
  existingColumns: [
    { id: "t", kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
    { id: "c", kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
  ],
};

// The same picker over a catalog that includes custom fields. A currency custom
// field is both a dimension and a measure, so it appears in both lists.
const customFieldData: ColumnPickerDialogData = {
  dimensions: [
    { key: "category", label: "Category" },
    { key: "custom_7", label: "HST", isCustom: true },
    { key: "custom_8", label: "Vendor", isCustom: true },
  ],
  measures: [
    { key: "amount", label: "Amount" },
    { key: "custom_7", label: "HST", isCustom: true },
  ],
  existingColumns: [],
};

describe("ColumnPickerDialogComponent", () => {
  afterEach(() => TestBed.resetTestingModule());

  it("badges custom fields in the field and measure pickers", async () => {
    const { fixture, component } = configure(customFieldData);
    await fixture.whenStable();

    const dimensions = new Map(component.dimensionOptions.map((o) => [o.value, o.badge]));
    expect(dimensions.get("custom_7")).toBe("Custom");
    expect(dimensions.get("custom_8")).toBe("Custom");
    expect(dimensions.get("category")).toBe("");

    const measures = new Map(component.measureOptions.map((o) => [o.value, o.badge]));
    expect(measures.get("custom_7")).toBe("Custom");
    expect(measures.get("amount")).toBe("");
  });

  it("saves an aggregate over a custom currency measure", async () => {
    const { fixture, component, close } = configure(customFieldData);
    await fixture.whenStable();

    component.pickAggregate();
    await fixture.whenStable();
    component.selectFunction(ReportColumn.AggFuncEnum.Sum);
    component.pickerForm.get("measure")!.setValue("custom_7");
    await fixture.whenStable();

    // The label is suggested from the custom field's name.
    expect(component.pickerForm.get("label")!.value).toContain("HST");
    component.save();

    expect(close).toHaveBeenCalledWith(
      expect.objectContaining({
        kind: ReportColumn.KindEnum.Aggregate,
        aggFunc: ReportColumn.AggFuncEnum.Sum,
        measure: "custom_7",
      })
    );
  });

  it("saves a dimension column over a custom field", async () => {
    const { fixture, component, close } = configure(customFieldData);
    await fixture.whenStable();

    component.pickDimension();
    await fixture.whenStable();
    component.pickerForm.get("field")!.setValue("custom_8");
    await fixture.whenStable();

    expect(component.pickerForm.get("label")!.value).toBe("Vendor");
    component.save();

    expect(close).toHaveBeenCalledWith(
      expect.objectContaining({ kind: ReportColumn.KindEnum.Dimension, field: "custom_8" })
    );
  });

  it("opens on the kind step and cannot save yet", async () => {
    const { fixture, component } = configure(baseData);
    await fixture.whenStable();
    expect(component.step()).toBe("kind");
    expect(component.canSave()).toBe(false);
  });

  it("configures and saves an aggregate column", async () => {
    const { fixture, component, close } = configure(baseData);
    await fixture.whenStable();

    component.pickAggregate();
    await fixture.whenStable();
    expect(component.step()).toBe("agg");
    expect(component.currentAggFunc).toBe(ReportColumn.AggFuncEnum.Sum);
    expect(component.canSave()).toBe(true);

    component.save();
    expect(close).toHaveBeenCalledWith(
      expect.objectContaining({
        kind: ReportColumn.KindEnum.Aggregate,
        name: "Amount",
        label: "Amount",
        aggFunc: ReportColumn.AggFuncEnum.Sum,
        measure: "amount",
      })
    );
  });

  it("drops the measure for a COUNT aggregate", async () => {
    const { fixture, component, close } = configure(baseData);
    await fixture.whenStable();
    component.pickAggregate();
    await fixture.whenStable();
    component.selectFunction(ReportColumn.AggFuncEnum.Count);
    await fixture.whenStable();
    component.save();
    const saved = close.mock.calls[0][0];
    expect(saved.aggFunc).toBe(ReportColumn.AggFuncEnum.Count);
    expect(saved.measure).toBeUndefined();
  });

  it("configures and saves a dimension column", async () => {
    const { fixture, component, close } = configure(baseData);
    await fixture.whenStable();

    component.pickDimension();
    await fixture.whenStable();
    expect(component.step()).toBe("dim");
    // The field and label preset to the first dimension.
    expect(component.pickerForm.get("field")!.value).toBe("category");
    expect(component.canSave()).toBe(true);

    component.save();
    expect(close).toHaveBeenCalledWith(
      expect.objectContaining({
        kind: ReportColumn.KindEnum.Dimension,
        field: "category",
        label: "Category",
      })
    );
  });

  it("builds a formula from clicked column and operator chips", async () => {
    const { fixture, component } = configure(formulaData);
    await fixture.whenStable();
    component.pickFormula();
    component.insertColumn("Total");
    component.insertOperator("/");
    component.insertColumn("Count");
    await fixture.whenStable();
    expect(component.pickerForm.get("expr")!.value).toBe("Total / Count");
    expect(component.formulaStatus().ok).toBe(true);
  });

  it("removes the last token and clears the built expression", async () => {
    const { fixture, component } = configure(formulaData);
    await fixture.whenStable();
    component.pickFormula();
    component.insertColumn("Total");
    component.insertOperator("/");
    component.insertColumn("Count");
    component.removeLastToken();
    expect(component.pickerForm.get("expr")!.value).toBe("Total /");
    component.clearExpression();
    expect(component.pickerForm.get("expr")!.value).toBe("");
  });

  it("saves the formula built from chips", async () => {
    const { fixture, component, close } = configure(formulaData);
    await fixture.whenStable();
    component.pickFormula();
    component.pickerForm.get("label")!.setValue("Average");
    component.insertColumn("Total");
    component.insertOperator("/");
    component.insertColumn("Count");
    await fixture.whenStable();
    expect(component.canSave()).toBe(true);
    component.save();
    expect(close).toHaveBeenCalledWith(
      expect.objectContaining({ kind: ReportColumn.KindEnum.Formula, label: "Average", expr: "Total / Count" })
    );
  });

  it("blocks a formula that references an unknown column", async () => {
    const { fixture, component } = configure(baseData);
    await fixture.whenStable();
    component.pickFormula();
    component.pickerForm.get("label")!.setValue("Total");
    component.pickerForm.get("expr")!.setValue("Subtotal + Hst");
    await fixture.whenStable();
    expect(component.formulaStatus().ok).toBe(false);
    expect(component.canSave()).toBe(false);
  });

  it("keeps an edited column's existing name", async () => {
    const { fixture, component, close } = configure({
      ...baseData,
      column: {
        id: "x1",
        kind: ReportColumn.KindEnum.Aggregate,
        name: "Total",
        label: "Total",
        aggFunc: ReportColumn.AggFuncEnum.Sum,
        measure: "amount",
      },
    });
    await fixture.whenStable();
    expect(component.step()).toBe("agg");
    component.pickerForm.get("label")!.setValue("Grand Total");
    await fixture.whenStable();
    component.save();
    expect(close).toHaveBeenCalledWith(
      expect.objectContaining({ id: "x1", name: "Total", label: "Grand Total" })
    );
  });
});
