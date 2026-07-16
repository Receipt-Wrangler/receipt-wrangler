import { ChangeDetectionStrategy, Component, OnInit, computed, inject, input, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormArray, FormControl, FormGroup } from "@angular/forms";
import { Store } from "@ngxs/store";
import { RECEIPT_STATUS_OPTIONS } from "src/constants";
import { AuthState, GroupState } from "src/store";
import { OperationsPipe } from "../../shared-ui/receipt-filter/operations.pipe";
import { Category, FilterOperation, Tag } from "../../open-api";

/** A single addable filter field and how to render its value editor. */
interface FilterFieldDef {
  field: string;
  label: string;
  type: "text" | "number" | "date" | "users" | "list";
  isCurrency?: boolean;
  optionsKey?: "categories" | "tags" | "status";
  optionValueKey?: string;
  optionFilterKey?: string;
  optionDisplayKey?: string;
}

const FILTER_FIELDS: FilterFieldDef[] = [
  { field: "name", label: "Name", type: "text" },
  { field: "amount", label: "Amount", type: "number", isCurrency: true },
  { field: "date", label: "Date", type: "date" },
  { field: "paidBy", label: "Paid by", type: "users" },
  { field: "categories", label: "Categories", type: "list", optionsKey: "categories", optionValueKey: "id", optionFilterKey: "name", optionDisplayKey: "name" },
  { field: "tags", label: "Tags", type: "list", optionsKey: "tags", optionValueKey: "id", optionFilterKey: "name", optionDisplayKey: "name" },
  { field: "status", label: "Status", type: "list", optionsKey: "status", optionValueKey: "value", optionFilterKey: "value", optionDisplayKey: "displayValue" },
];

/**
 * The builder's inline filter: the design's add-a-filter chip layout over the
 * exact ReceiptPagedRequestFilter contract. It reuses the shared filter form
 * (built by buildReceiptFilterForm in the root form, including BETWEEN handling)
 * and the SharedUiModule OperationsPipe, so it can't drift from the receipts
 * filter — only the presentation differs. Category/tag options are the union of the
 * user's group catalogs, so the picker is available regardless of the chosen scope.
 */
@Component({
  selector: "app-report-filters",
  templateUrl: "./report-filters.component.html",
  styleUrls: ["./report-filters.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportFiltersComponent implements OnInit {
  public readonly form = input.required<FormGroup>();

  private readonly store = inject(Store);
  private readonly operationsPipe = new OperationsPipe();
  private readonly activeFieldKeys = signal<string[]>([]);

  public readonly FilterOperation = FilterOperation;
  public readonly receiptStatusOptions = RECEIPT_STATUS_OPTIONS;

  // The "Add filter…" picker is a stateless dispatcher: app-select is
  // FormControl-driven (no change output), so a scratch control feeds each pick
  // into addFilter and is reset to the blank option afterward.
  public readonly addFilterControl = new FormControl<string | null>(null);

  constructor() {
    this.addFilterControl.valueChanges.pipe(takeUntilDestroyed()).subscribe((field) => {
      if (!field) {
        return;
      }
      this.addFilter(field);
      this.addFilterControl.setValue(null, { emitEvent: false });
    });
  }

  public ngOnInit(): void {
    // On open-in-builder the form arrives pre-seeded from the saved template, but the
    // visible-row set lives in a local signal. Show a row for every field the stored
    // filter actually set (non-empty operation); untouched fields stay collapsed behind
    // "Add filter…". This mirrors the addFilter/removeFilter invariant exactly.
    const active = FILTER_FIELDS.filter((def) => this.isFilterActive(def.field)).map((def) => def.field);
    this.activeFieldKeys.set(active);
  }

  private isFilterActive(field: string): boolean {
    const operation = this.filterGroup.get(field + ".operation")?.value;
    return operation != null && operation !== FilterOperation.Empty;
  }

  private readonly groups = this.store.selectSignal(GroupState.groupsWithoutAll);
  public readonly categories = computed<Category[]>(() =>
    unionCatalog(this.groups(), (id) => this.store.selectSnapshot(AuthState.groupCategories(id)))
  );
  public readonly tags = computed<Tag[]>(() =>
    unionCatalog(this.groups(), (id) => this.store.selectSnapshot(AuthState.groupTags(id)))
  );

  public readonly activeFields = computed<FilterFieldDef[]>(() => {
    const active = this.activeFieldKeys();
    return active
      .map((key) => FILTER_FIELDS.find((def) => def.field === key))
      .filter((def): def is FilterFieldDef => !!def);
  });

  public readonly addableFields = computed<FilterFieldDef[]>(() => {
    const active = new Set(this.activeFieldKeys());
    return FILTER_FIELDS.filter((def) => !active.has(def.field));
  });

  public get filterGroup(): FormGroup {
    return this.form().get("filter") as FormGroup;
  }

  public operationOptions(type: string): string[] {
    return this.operationsPipe.transform(type, false);
  }

  public operationDisplay(type: string): string[] {
    return this.operationsPipe.transform(type, true);
  }

  public operationValueOf(field: string): FilterOperation {
    return this.filterGroup.get(field + ".operation")?.value;
  }

  public optionsFor(def: FilterFieldDef): unknown[] {
    switch (def.optionsKey) {
      case "categories":
        return this.categories();
      case "tags":
        return this.tags();
      case "status":
        return this.receiptStatusOptions;
      default:
        return [];
    }
  }

  public addFilter(field: string): void {
    const def = FILTER_FIELDS.find((candidate) => candidate.field === field);
    if (!def) {
      return;
    }
    const [firstOperation] = this.operationsPipe.transform(def.type, false);
    this.filterGroup.get(field + ".operation")?.setValue(firstOperation ?? FilterOperation.Empty);
    this.activeFieldKeys.update((keys) => [...keys, field]);
  }

  public removeFilter(field: string): void {
    this.filterGroup.get(field + ".operation")?.setValue(FilterOperation.Empty);
    const valueControl = this.filterGroup.get(field + ".value");
    if (valueControl instanceof FormArray) {
      valueControl.clear();
    } else {
      valueControl?.setValue(null);
    }
    this.activeFieldKeys.update((keys) => keys.filter((key) => key !== field));
  }
}

function unionCatalog<T extends { id?: number }>(
  groups: { id?: number }[],
  read: (groupId: number) => T[]
): T[] {
  const merged = new Map<number, T>();
  for (const group of groups) {
    if (group.id == null) {
      continue;
    }
    for (const item of read(group.id) ?? []) {
      if (item.id != null) {
        merged.set(item.id, item);
      }
    }
  }
  return [...merged.values()];
}
