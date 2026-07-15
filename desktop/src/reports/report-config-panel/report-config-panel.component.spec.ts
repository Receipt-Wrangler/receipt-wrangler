import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { CustomFieldService, ReportColumn, ReportDetail, ReportPeriod } from "../../open-api";
import { buildColumnGroup } from "../models/report-form.factory";
import { ReportConfigPanelComponent } from "./report-config-panel.component";

const GROUPS = [
  { id: 1, name: "Household" },
  { id: 2, name: "Roommates" },
];

function buildForm(formBuilder: FormBuilder): FormGroup {
  return formBuilder.group({
    name: formBuilder.control("My Report"),
    scope: formBuilder.array([]),
    groupBy: formBuilder.array([]),
    detail: formBuilder.group({
      mode: formBuilder.control<string>(ReportDetail.ModeEnum.Aggregate),
      by: formBuilder.control("category"),
    }),
    columns: formBuilder.array([]),
    filter: formBuilder.group({}),
    period: formBuilder.group({
      preset: formBuilder.control<string>(ReportPeriod.PresetEnum.ThisMonth),
      startDate: formBuilder.control<Date | null>(null),
      endDate: formBuilder.control<Date | null>(null),
    }),
    document: formBuilder.group({ intro: formBuilder.control("") }),
  });
}

describe("ReportConfigPanelComponent", () => {
  let fixture: ComponentFixture<ReportConfigPanelComponent>;
  let component: ReportConfigPanelComponent;
  let form: FormGroup;
  let formBuilder: FormBuilder;
  let dialog: { open: jest.Mock };

  beforeEach(async () => {
    dialog = { open: jest.fn() };

    await TestBed.configureTestingModule({
      declarations: [ReportConfigPanelComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: Store, useValue: { selectSignal: () => signal(GROUPS) } },
        { provide: MatDialog, useValue: dialog },
        { provide: CustomFieldService, useValue: { getPagedCustomFields: jest.fn(() => of({ data: [], totalCount: 0 })) } },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    })
      // This suite exercises the component's public API, not its template; an empty
      // template keeps the harness free of the shared pipes/child components.
      .overrideComponent(ReportConfigPanelComponent, { set: { template: "" } })
      .compileComponents();

    formBuilder = TestBed.inject(FormBuilder);
    form = buildForm(formBuilder);
    fixture = TestBed.createComponent(ReportConfigPanelComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("form", form);
    component.ngOnInit();
  });

  function columnsArray(): FormArray {
    return form.get("columns") as FormArray;
  }

  // ---- add-grouping control wiring --------------------------------------

  it("adds a grouping level when the add-select control is set, then resets the control", () => {
    const groupBy = form.get("groupBy") as FormArray;
    expect(groupBy.length).toBe(0);

    component.addGroupControl.setValue("paid_by");

    expect(groupBy.length).toBe(1);
    expect(groupBy.at(0).value).toBe("paid_by");
    expect(component.addGroupControl.value).toBeNull();
  });

  it("ignores the blank (reset) option", () => {
    component.addGroupControl.setValue(null);
    expect((form.get("groupBy") as FormArray).length).toBe(0);
  });

  // ---- scope -------------------------------------------------------------

  it("openAddGroups sets the scope and scopeChips resolves the group names", () => {
    dialog.open.mockReturnValue({ afterClosed: () => of(["1", "2"]) });

    component.openAddGroups();

    const chips = component.scopeChips();
    expect(chips.map((chip) => chip.name)).toEqual(["Household", "Roommates"]);
    expect(chips[0].initials).toBe("HO");
  });

  it("openAddGroups is a no-op when the dialog is cancelled", () => {
    dialog.open.mockReturnValue({ afterClosed: () => of(undefined) });
    component.openAddGroups();
    expect(component.scopeChips().length).toBe(0);
  });

  it("removeScope drops the group at the index", () => {
    dialog.open.mockReturnValue({ afterClosed: () => of(["1", "2"]) });
    component.openAddGroups();

    component.removeScope(0);

    expect(component.scopeChips().map((chip) => chip.name)).toEqual(["Roommates"]);
  });

  // ---- grouping ----------------------------------------------------------

  it("adds, reorders, and removes grouping levels and narrows addableDimensions", () => {
    component.addGroupBy("paid_by");
    component.addGroupBy("tag");
    expect(component.groupByLevels().map((level) => level.label)).toEqual(["Paid By", "Tag"]);

    const addable = component.addableDimensions().map((field) => field.key);
    expect(addable).toContain("category");
    expect(addable).not.toContain("paid_by");
    expect(addable).not.toContain("tag");

    component.moveGroupBy(0, 1);
    expect(component.groupByLevels().map((level) => level.label)).toEqual(["Tag", "Paid By"]);

    component.removeGroupBy(0);
    expect(component.groupByLevels().map((level) => level.label)).toEqual(["Paid By"]);
  });

  it("addGroupBy ignores an empty key", () => {
    component.addGroupBy("");
    expect(component.groupByLevels().length).toBe(0);
  });

  // ---- detail mode -------------------------------------------------------

  it("setDetailMode updates the form and the detailMode getter", () => {
    component.setDetailMode(ReportDetail.ModeEnum.Records);
    expect(form.get("detail.mode")!.value).toBe(ReportDetail.ModeEnum.Records);
    expect(component.detailMode).toBe(ReportDetail.ModeEnum.Records);
  });

  // ---- columns -----------------------------------------------------------

  it("openColumnPicker appends a column from the dialog result", () => {
    dialog.open.mockReturnValue({
      afterClosed: () =>
        of({ id: "x", kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count }),
    });

    component.openColumnPicker();

    expect(columnsArray().length).toBe(1);
    expect(columnsArray().at(0).get("label")!.value).toBe("Count");
  });

  it("moveColumn reorders and removeColumn drops a column", () => {
    columnsArray().push(buildColumnGroup(formBuilder, { kind: ReportColumn.KindEnum.Dimension, name: "A", label: "A", field: "category" }));
    columnsArray().push(buildColumnGroup(formBuilder, { kind: ReportColumn.KindEnum.Aggregate, name: "B", label: "B", aggFunc: ReportColumn.AggFuncEnum.Count }));

    component.moveColumn(0, 1);
    expect(columnsArray().at(0).get("label")!.value).toBe("B");

    component.removeColumn(0);
    expect(columnsArray().length).toBe(1);
    expect(columnsArray().at(0).get("label")!.value).toBe("A");
  });

  it("columnRows derives the disabled state and reason for an invalid aggregate dimension", () => {
    columnsArray().push(buildColumnGroup(formBuilder, { kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" }));
    columnsArray().push(buildColumnGroup(formBuilder, { kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" }));

    // Aggregate by tag with no grouping — the Category dimension column is invalid.
    form.get("detail.by")!.setValue("tag");

    const rows = component.columnRows();
    const category = rows.find((row) => row.label === "Category")!;
    expect(category.disabled).toBe(true);
    expect(category.disabledReason).toContain("hidden");

    const total = rows.find((row) => row.label === "Total")!;
    expect(total.disabled).toBe(false);
  });

  // ---- parameters + document --------------------------------------------

  it("periodLabel resolves the preset, and reports a custom range with no dates", () => {
    expect(component.periodLabel()).toBeTruthy();

    form.get("period.preset")!.setValue(ReportPeriod.PresetEnum.Custom);
    expect(component.periodLabel()).toBe("a custom range");
  });

  it("insertVariable appends tokens to the document intro", () => {
    component.insertVariable("{{period}}");
    expect(form.get("document.intro")!.value).toBe("{{period}}");

    component.insertVariable("{{generatedAt}}");
    expect(form.get("document.intro")!.value).toBe("{{period}} {{generatedAt}}");
  });
});
