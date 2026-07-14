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

describe("ColumnPickerDialogComponent", () => {
  afterEach(() => TestBed.resetTestingModule());

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
