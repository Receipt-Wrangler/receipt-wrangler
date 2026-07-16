import { Component, NO_ERRORS_SCHEMA, provideZonelessChangeDetection, signal, WritableSignal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder, FormGroup } from "@angular/forms";
import { Store } from "@ngxs/store";
import { UntilDestroy } from "@ngneat/until-destroy";
import { FilterOperation, User } from "../../open-api";
import { UserState } from "../../store";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { CURRENT_VIEWER_PAID_BY_ID, ReportFiltersComponent } from "./report-filters.component";

// buildReceiptFilterForm wires untilDestroyed subscriptions, so seeding a realistic
// filter needs an @UntilDestroy()-decorated context — a throwaway host, exactly as the
// receipt filter's own spec and report-form.factory.spec do.
@UntilDestroy()
@Component({ selector: "app-noop", template: "", standalone: false })
class NoopComponent {}

describe("ReportFiltersComponent", () => {
  let fixture: ComponentFixture<ReportFiltersComponent>;
  let component: ReportFiltersComponent;
  let formBuilder: FormBuilder;
  let host: NoopComponent;
  let usersSignal: WritableSignal<User[]>;

  beforeEach(async () => {
    usersSignal = signal<User[]>([]);
    // Route UserState.users to a signal the test controls (for the paid-by picker);
    // every other selector stays empty as before.
    const store = {
      selectSignal: (selector: unknown) => (selector === UserState.users ? usersSignal : signal([])),
      selectSnapshot: () => [],
    };

    await TestBed.configureTestingModule({
      declarations: [ReportFiltersComponent, NoopComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: Store, useValue: store },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    formBuilder = TestBed.inject(FormBuilder);
    host = TestBed.createComponent(NoopComponent).componentInstance;
    fixture = TestBed.createComponent(ReportFiltersComponent);
    component = fixture.componentInstance;
  });

  /**
   * Mounts the component over a root form whose `filter` group is built by the shared
   * builder from `seed` — i.e. the exact form the report builder produces on
   * open-in-builder — then fires ngOnInit, which seeds the visible-row signal.
   * ngOnInit is invoked directly rather than via detectChanges: activeFields() is the
   * exact signal the template's @for iterates, so a full DOM render (and its formGet
   * pipe / child components) isn't needed to prove which rows show — the e2e suite
   * covers the rendered DOM.
   */
  function mountWithFilter(seed: unknown): void {
    const form = new FormGroup({ filter: buildReceiptFilterForm(seed, host) });
    fixture.componentRef.setInput("form", form);
    component.ngOnInit();
  }

  it("adds a filter when the add-select control is set, then resets the control", () => {
    fixture.componentRef.setInput("form", formBuilder.group({ filter: formBuilder.group({}) }));

    expect(component.activeFields().length).toBe(0);

    component.addFilterControl.setValue("name");

    expect(component.activeFields().map((def) => def.field)).toContain("name");
    expect(component.addFilterControl.value).toBeNull();
  });

  it("ignores the blank (reset) option", () => {
    fixture.componentRef.setInput("form", formBuilder.group({ filter: formBuilder.group({}) }));

    component.addFilterControl.setValue(null);
    expect(component.activeFields().length).toBe(0);
  });

  it("shows the row for a hydrated category filter and keeps its selected ids (the reported bug)", () => {
    mountWithFilter({ categories: { operation: FilterOperation.Contains, value: [2, 3] } });

    expect(component.activeFields().map((def) => def.field)).toEqual(["categories"]);
    expect(component.filterGroup.get("categories.value")!.value).toEqual([2, 3]);
    expect(component.filterGroup.get("categories.operation")!.value).toBe(FilterOperation.Contains);
  });

  it("shows a row for every field a saved filter set, across all editor types", () => {
    mountWithFilter({
      name: { operation: FilterOperation.Contains, value: "coffee" },
      amount: { operation: FilterOperation.Between, value: [5, 50] },
      date: { operation: FilterOperation.WithinCurrentMonth }, // value-less operation
      paidBy: { operation: FilterOperation.Contains, value: [11] },
      categories: { operation: FilterOperation.Contains, value: [2, 3] },
      tags: { operation: FilterOperation.Contains, value: [7] },
      status: { operation: FilterOperation.Contains, value: ["OPEN"] },
    });

    // Order follows FILTER_FIELDS declaration order.
    expect(component.activeFields().map((def) => def.field)).toEqual([
      "name",
      "amount",
      "date",
      "paidBy",
      "categories",
      "tags",
      "status",
    ]);
  });

  it("keeps a value-less operation (WITHIN_CURRENT_MONTH) visible", () => {
    mountWithFilter({ date: { operation: FilterOperation.WithinCurrentMonth } });

    expect(component.activeFields().map((def) => def.field)).toEqual(["date"]);
  });

  it("shows no rows for a blank filter (no false positives)", () => {
    mountWithFilter({});

    expect(component.activeFields().length).toBe(0);
    expect(component.addableFields().length).toBeGreaterThan(0);
  });

  it("treats an empty-operation field as inactive", () => {
    mountWithFilter({ name: { operation: FilterOperation.Empty, value: "" } });

    expect(component.activeFields().length).toBe(0);
  });

  it("prepends a 'current viewer (me)' sentinel ahead of the user pool in the paid-by picker", () => {
    fixture.componentRef.setInput("form", formBuilder.group({ filter: formBuilder.group({}) }));
    usersSignal.set([
      { id: 5, username: "amy", displayName: "Amy" } as User,
      { id: 6, username: "ben", displayName: "" } as User,
    ]);

    const options = component.paidByOptions();
    // The sentinel is first, then every user (display-name fallback to username).
    expect(options[0]).toEqual({ id: CURRENT_VIEWER_PAID_BY_ID, displayName: "Current viewer (me)" });
    expect(options.slice(1)).toEqual([
      { id: 5, displayName: "Amy" },
      { id: 6, displayName: "ben" },
    ]);
  });

  it("hydrates a saved 'current viewer' paid-by filter (the -1 sentinel) into the builder", () => {
    mountWithFilter({ paidBy: { operation: FilterOperation.Contains, value: [CURRENT_VIEWER_PAID_BY_ID] } });

    expect(component.activeFields().map((def) => def.field)).toEqual(["paidBy"]);
    // The sentinel round-trips as a plain id, and the picker always carries its option.
    expect(component.filterGroup.get("paidBy.value")!.value).toEqual([CURRENT_VIEWER_PAID_BY_ID]);
    expect(component.paidByOptions()[0].id).toBe(CURRENT_VIEWER_PAID_BY_ID);
  });

  it("still adds and removes filters after a seeded init", () => {
    mountWithFilter({ categories: { operation: FilterOperation.Contains, value: [2] } });
    expect(component.activeFields().map((def) => def.field)).toEqual(["categories"]);

    component.addFilter("name");
    expect(component.activeFields().map((def) => def.field)).toContain("name");

    component.removeFilter("categories");
    expect(component.activeFields().map((def) => def.field)).not.toContain("categories");
    expect(component.filterGroup.get("categories.operation")!.value).toBe(FilterOperation.Empty);
  });
});
