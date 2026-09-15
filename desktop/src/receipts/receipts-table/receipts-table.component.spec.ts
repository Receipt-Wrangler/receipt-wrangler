import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatChipsModule } from "@angular/material/chips";
import { MatDialogModule } from "@angular/material/dialog";
import { MatMenuModule } from "@angular/material/menu";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { ActivatedRoute, provideRouter } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of, Subject, throwError } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { DEFAULT_RECEIPT_TABLE_COLUMNS } from "src/interfaces";
import { SetColumnConfig, SetReceiptFilterData } from "src/store/receipt-table.actions";
import { ReceiptTableState } from "src/store/receipt-table.state";
import { MonthStepperComponent } from "../../shared-ui/month-stepper/month-stepper.component";
import { ApiModule, CustomField, CustomFieldType, FilterOperation, Group, Permission, Receipt, ReceiptStatus, ReceiptSummary } from "../../open-api";
import { ReceiptFilterService } from "../../services/receipt-filter.service";
import { AuthState, GroupState, UserState } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
import { SetGroups } from "../../store/group.state.actions";
import { SetQuickDateField, SetReceiptFilter, SetSummaryConfigGroupId } from "../../store/receipt-table.actions";
import { ReceiptsTableComponent } from "./receipts-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

const customField_ = (id: number): CustomField =>
  ({ id, name: "Vendor", type: CustomFieldType.Text } as CustomField);

describe("ReceiptsTableComponent", () => {
  let component: ReceiptsTableComponent;
  let fixture: ComponentFixture<ReceiptsTableComponent>;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [ReceiptsTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        NgxsModule.forRoot([ReceiptTableState, AuthState, GroupState, UserState]),
        ReactiveFormsModule,
        MatSnackBarModule,
        MatTooltipModule,
        MatDialogModule,
        MatMenuModule,
        MatChipsModule,
        MonthStepperComponent,
        PipesModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        categories: [],
                        tags: [],
                    },
                },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
        provideZonelessChangeDetection(),
        provideRouter([]),
    ]
}).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(ReceiptsTableComponent);
    component = fixture.componentInstance;
    Object.defineProperty(component, 'table', {
      value: () => ({
        selection: {},
        changed: of(undefined),
      }),
    });
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("gates each header action on its own group permission", () => {
    store.dispatch(
      new SetPermissions([], {
        5: [
          Permission.GroupReceiptsCreate,
          Permission.GroupReceiptsQuickScan,
          Permission.GroupEmailPoll,
        ],
      })
    );
    component.groupId = "5";

    (component as any).setCanEdit();

    expect(component.canCreate).toEqual(true);
    expect(component.canQuickScan).toEqual(true);
    expect(component.canPollEmail).toEqual(true);
    // None of those imply update, which gates Edit / Bulk Status Update.
    expect(component.canEdit).toEqual(false);
  });

  it("keeps canEdit tied to group.receipts.update only", () => {
    store.dispatch(
      new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] })
    );
    component.groupId = "5";

    (component as any).setCanEdit();

    expect(component.canEdit).toEqual(true);
    expect(component.canCreate).toEqual(false);
    expect(component.canQuickScan).toEqual(false);
    expect(component.canPollEmail).toEqual(false);
  });

  describe("custom field columns", () => {
    const customField = (id: number, name: string): CustomField =>
      ({ id, name, type: CustomFieldType.Text } as CustomField);

    // setColumns reads the cell templates through viewChild.required, which never
    // resolves for a component this spec never renders.
    const stubCellTemplates = (): void => {
      for (const cell of [
        "createdAtCell", "dateCell", "nameCell", "paidByCell", "amountCell",
        "categoryCell", "tagCell", "statusCell", "resolvedDateCell",
        "customFieldCell", "actionsCell",
      ]) {
        Object.defineProperty(component, cell, { value: () => ({}) });
      }
    };

    beforeEach(() => {
      stubCellTemplates();
    });

    it("builds a sortable column per custom field", () => {
      component.customFields.set([customField(7, "Vendor")]);
      store.dispatch(
        new SetColumnConfig([
          ...DEFAULT_RECEIPT_TABLE_COLUMNS,
          { matColumnDef: "custom_7", visible: true, order: 9 },
        ])
      );

      (component as any).setColumns();

      const column = component
        .columns()
        .find((col) => col.matColumnDef === "custom_7");
      expect(column?.columnHeader).toEqual("Vendor");
      expect(column?.sortable).toEqual(true);
    });

    // mat-table throws on a displayed id it has no definition for, and a config
    // naming a since-deleted custom field is ordinary rather than exotic.
    it("never displays a column it could not resolve", () => {
      component.customFields.set([]);
      store.dispatch(
        new SetColumnConfig([
          ...DEFAULT_RECEIPT_TABLE_COLUMNS,
          { matColumnDef: "custom_99", visible: true, order: 9 },
        ])
      );

      (component as any).setColumns();

      expect(component.displayedColumns()).not.toContain("custom_99");
      for (const displayed of component.displayedColumns()) {
        if (displayed === "select") {
          continue;
        }
        expect(
          component.columns().some((col) => col.matColumnDef === displayed)
        ).toEqual(true);
      }
    });
  });

  describe("reconcileColumnConfig", () => {
    it("drops a persisted column whose custom field no longer exists", () => {
      component.customFields.set([]);
      store.dispatch(
        new SetColumnConfig([
          ...DEFAULT_RECEIPT_TABLE_COLUMNS,
          { matColumnDef: "custom_99", visible: true, order: 9 },
        ])
      );

      (component as any).reconcileColumnConfig();

      expect(
        store
          .selectSnapshot(ReceiptTableState.columnConfig)
          .map((col) => col.matColumnDef)
      ).not.toContain("custom_99");
    });

    // A user without app.custom-fields.read resolves an EMPTY catalog, which is not
    // the same as "no custom fields exist". The configuration is persisted per
    // browser and shared across accounts, so dropping against it would destroy an
    // administrator's saved layout the moment a colleague opens the page here.
    it("keeps a persisted custom field column when the catalog is unavailable", () => {
      component.customFields.set([]);
      component.customFieldsAvailable.set(false);
      store.dispatch(
        new SetColumnConfig([
          ...DEFAULT_RECEIPT_TABLE_COLUMNS,
          { matColumnDef: "custom_99", visible: true, order: 9 },
        ])
      );

      (component as any).reconcileColumnConfig();

      expect(
        store
          .selectSnapshot(ReceiptTableState.columnConfig)
          .map((col) => col.matColumnDef)
      ).toContain("custom_99");
    });

    // Follows from the column surviving: the sort is reset off the reconciled
    // columns, so preserving one without the other would still discard the sort.
    it("keeps a sort on a custom field when the catalog is unavailable", () => {
      component.customFields.set([]);
      component.customFieldsAvailable.set(false);
      store.dispatch(
        new SetColumnConfig([
          ...DEFAULT_RECEIPT_TABLE_COLUMNS,
          { matColumnDef: "custom_99", visible: true, order: 9 },
        ])
      );
      const filterData = store.selectSnapshot(ReceiptTableState.filterData);
      store.dispatch(
        new SetReceiptFilterData({ ...filterData, orderBy: "custom_99", sortDirection: "asc" })
      );

      (component as any).reconcileColumnConfig();

      expect(store.selectSnapshot(ReceiptTableState.filterData).orderBy).toEqual(
        "custom_99"
      );
    });

    // The sort is persisted per browser and shared across accounts, so the table
    // can start up asking the API to order by a column that no longer exists -
    // which the API rejects outright, failing the very first load.
    it("resets a sort on a custom field that no longer exists", () => {
      component.customFields.set([]);
      const filterData = store.selectSnapshot(ReceiptTableState.filterData);
      store.dispatch(
        new SetReceiptFilterData({ ...filterData, orderBy: "custom_99", sortDirection: "asc" })
      );

      (component as any).reconcileColumnConfig();

      expect(store.selectSnapshot(ReceiptTableState.filterData).orderBy).toEqual(
        "created_at"
      );
    });

    it("keeps a sort on a custom field that still exists", () => {
      component.customFields.set([customField_(7)]);
      const filterData = store.selectSnapshot(ReceiptTableState.filterData);
      store.dispatch(
        new SetReceiptFilterData({ ...filterData, orderBy: "custom_7", sortDirection: "asc" })
      );

      (component as any).reconcileColumnConfig();

      expect(store.selectSnapshot(ReceiptTableState.filterData).orderBy).toEqual(
        "custom_7"
      );
    });
  });

  it("should map selected ids from selecton", () => {
    const selectedReceipts: Receipt[] = [
      {
        id: 1,
      } as Receipt,
      {
        id: 2,
      } as Receipt,
    ];
    Object.defineProperty(component, 'table', {
      value: () => ({
        selection: {
          changed: of({
            source: {
              selected: selectedReceipts,
            },
          }),
        },
      }),
    });
    component.ngAfterViewInit();

    expect(component.selectedReceiptIds()).toEqual([1, 2]);
  });
  describe("quick date filtering and filter chips", () => {
    let refetch: jest.SpyInstance;

    const setFilter = (filter: Record<string, unknown>) =>
      store.dispatch(new SetReceiptFilter(filter as any));

    beforeEach(() => {
      refetch = jest
        .spyOn(TestBed.inject(ReceiptFilterService), "getPagedReceiptsForGroups")
        .mockReturnValue(of({ data: [], totalCount: 0 } as any));
    });

    it("reads All time with no date filter", () => {
      expect(component.stepperMonth()).toBeNull();
      expect(component.stepperLabel()).toEqual("All time");
    });

    it("writes the picked month as a BETWEEN over the whole month", () => {
      component.monthSelected({ year: 2026, month: 8 });

      const date = (store.selectSnapshot(ReceiptTableState.filterData).filter as any).date;
      const [start, end] = date.value as Date[];
      expect(date.operation).toEqual(FilterOperation.Between);
      expect(start.getMonth()).toEqual(8);
      expect(start.getDate()).toEqual(1);
      expect(end.getDate()).toEqual(30);

      expect(component.stepperLabel()).toEqual("September 2026");
      expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(1);
      expect(refetch).toHaveBeenCalled();
    });

    // The quick control IS the date filter, so it replaces whatever was there.
    it("overrides an existing date filter rather than sitting beside it", () => {
      setFilter({ date: { operation: FilterOperation.GreaterThan, value: new Date(2020, 0, 1) } });

      component.monthSelected({ year: 2026, month: 8 });

      const date = (store.selectSnapshot(ReceiptTableState.filterData).filter as any).date;
      expect(date.operation).toEqual(FilterOperation.Between);
    });

    it("clears the date filter for all time", () => {
      component.monthSelected({ year: 2026, month: 8 });
      component.allTimeSelected();

      expect(component.stepperLabel()).toEqual("All time");
      expect(component.filterChips()).toEqual([]);
    });

    // A date filter the stepper cannot express stays visible as a chip so it can
    // still be seen and cleared.
    it("reads Custom and shows a date chip for a filter it cannot express", () => {
      setFilter({ date: { operation: FilterOperation.WithinCurrentMonth, value: null } });

      expect(component.stepperMonth()).toBeNull();
      expect(component.stepperLabel()).toEqual("Custom");
      expect(component.filterChips()).toEqual([
        { key: "date", label: "Receipt Date within current month" },
      ]);
    });

    // The chip row used to omit the field the stepper was naming. It is the only
    // place that says WHICH date column is filtered now that the target is
    // selectable, so it no longer makes exceptions.
    it("chips the month the stepper is showing", () => {
      component.monthSelected({ year: 2026, month: 8 });

      expect(component.filterChips().map((chip) => chip.key)).toEqual(["date"]);
    });

    describe("choosing which date field to filter on", () => {
      it("defaults to the receipt date", () => {
        expect(component.quickDateField()).toEqual("date");
        expect(component.quickDateFieldLabel()).toEqual("Receipt Date");
      });

      it("offers exactly the date fields, labelled as the dialog labels them", () => {
        expect(component.dateFilterFields.map((field) => [field.key, field.label])).toEqual([
          ["date", "Receipt Date"],
          ["resolvedDate", "Resolved Date"],
          ["createdAt", "Added At"],
        ]);
      });

      it("writes the picked month to the selected field, not to date", () => {
        component.quickDateFieldSelected("resolvedDate");
        component.monthSelected({ year: 2026, month: 8 });

        const filter = store.selectSnapshot(ReceiptTableState.filterData).filter as any;
        expect(filter.resolvedDate.operation).toEqual(FilterOperation.Between);
        expect(filter.date).toEqual({ operation: null, value: null });

        expect(component.quickDateFieldLabel()).toEqual("Resolved Date");
        expect(component.stepperLabel()).toEqual("September 2026");
      });

      it("reads the stepper label and Custom fallback off the selected field", () => {
        setFilter({
          date: { operation: FilterOperation.Between, value: [new Date(2026, 8, 1), new Date(2026, 8, 30)] },
          createdAt: { operation: FilterOperation.WithinCurrentMonth, value: null },
        });

        expect(component.stepperLabel()).toEqual("September 2026");

        component.quickDateFieldSelected("createdAt");

        expect(component.stepperMonth()).toBeNull();
        expect(component.stepperLabel()).toEqual("Custom");
      });

      // Switching the target changes no condition, only which one the stepper
      // describes — so nothing the user set elsewhere is destroyed, and the
      // abandoned condition stays visible and clearable as its own chip.
      it("leaves the previous field's condition applied, with its chip", () => {
        component.monthSelected({ year: 2026, month: 8 });
        refetch.mockClear();

        component.quickDateFieldSelected("resolvedDate");

        const filter = store.selectSnapshot(ReceiptTableState.filterData).filter as any;
        expect(filter.date.operation).toEqual(FilterOperation.Between);
        expect(component.stepperLabel()).toEqual("All time");
        expect(component.filterChips().map((chip) => chip.key)).toEqual(["date"]);
        expect(refetch).not.toHaveBeenCalled();
      });

      it("follows a field chosen outside the component", () => {
        store.dispatch(new SetQuickDateField("createdAt"));

        expect(component.quickDateFieldLabel()).toEqual("Added At");
      });
    });

    it("builds a chip per active field, resolving ids to names", () => {
      component.categories.set([{ id: 3, name: "Groceries" }] as any);
      setFilter({
        name: { operation: FilterOperation.Contains, value: "whole" },
        categories: { operation: FilterOperation.Contains, value: [3] },
        status: { operation: FilterOperation.Contains, value: [ReceiptStatus.Open] },
      });

      expect(component.filterChips()).toEqual([
        { key: "name", label: "Name contains whole" },
        { key: "categories", label: "Categories contains Groceries" },
        { key: "status", label: "Status contains Open" },
      ]);
    });

    // Each refresh used to be its own subscription, so the last RESPONSE won
    // rather than the last REQUEST — and the quick date arrows are one click
    // apart, which is what makes this reachable.
    it("lets a newer refresh supersede an in-flight one", () => {
      const superseded = new Subject<any>();
      const latest = new Subject<any>();
      refetch.mockReturnValueOnce(superseded).mockReturnValueOnce(latest);

      component.monthSelected({ year: 2026, month: 8 });
      component.monthSelected({ year: 2026, month: 9 });

      // The first request resolving late must not repaint the table.
      superseded.next({ data: [{ id: 1 }], totalCount: 1 });
      expect(component.totalCount()).toEqual(0);

      latest.next({ data: [{ id: 2 }], totalCount: 2 });
      expect(component.totalCount()).toEqual(2);
    });

    // switchMap completes the outer stream on an error unless the inner one
    // swallows it, which would silently kill every refresh after the first
    // failure.
    it("keeps refreshing after a failed request", () => {
      refetch.mockReturnValueOnce(throwError(() => new Error("boom")));

      component.monthSelected({ year: 2026, month: 8 });

      refetch.mockReturnValue(of({ data: [{ id: 3 }], totalCount: 3 }));
      component.monthSelected({ year: 2026, month: 9 });

      expect(component.totalCount()).toEqual(3);
    });

    it("clears exactly the field whose chip was dismissed", () => {
      setFilter({
        name: { operation: FilterOperation.Contains, value: "whole" },
        status: { operation: FilterOperation.Contains, value: [ReceiptStatus.Open] },
      });

      component.filterChipCleared("status");

      const filter = store.selectSnapshot(ReceiptTableState.filterData).filter as any;
      expect(filter.status).toEqual({ operation: null, value: [] });
      expect(filter.name).toEqual({ operation: FilterOperation.Contains, value: "whole" });
      expect(component.filterChips().map((chip) => chip.key)).toEqual(["name"]);
      expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(1);
      expect(refetch).toHaveBeenCalled();
    });
  });
});

/**
 * The summary rides its OWN refresh stream, and the point of that stream is what it does not do:
 * paging and sorting change neither the filter nor the figures, so re-requesting an unpaged
 * aggregate for them would be the most expensive no-op in the app.
 */
describe("ReceiptsTableComponent receipt summary", () => {
  let component: ReceiptsTableComponent;
  let fixture: ComponentFixture<ReceiptsTableComponent>;
  let store: Store;
  let receiptFilterService: ReceiptFilterService;
  let summaryCalls: { groupId: string; configurationGroupId: number }[];

  function summaryGroup(id: number, name: string, enabled: boolean, isAllGroup = false): Group {
    return {
      id,
      name,
      isAllGroup,
      groupReceiptSettings: { receiptSummaryEnabled: enabled },
    } as Group;
  }

  const summaryResponse: ReceiptSummary = {
    enabled: true,
    configurationGroupId: 1,
    overall: {
      status: "" as ReceiptStatus,
      receiptCount: 2,
      total: "20.00",
      customFieldTotals: [],
    },
    statuses: [],
  } as ReceiptSummary;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReceiptsTableComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        ApiModule,
        NgxsModule.forRoot([ReceiptTableState, AuthState, GroupState, UserState]),
        ReactiveFormsModule,
        MatSnackBarModule,
        MatTooltipModule,
        MatDialogModule,
        MatMenuModule,
        MatChipsModule,
        MonthStepperComponent,
        PipesModule,
      ],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { data: { categories: [], tags: [] } } },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
        provideZonelessChangeDetection(),
        provideRouter([]),
      ],
    }).compileComponents();

    store = TestBed.inject(Store);
    receiptFilterService = TestBed.inject(ReceiptFilterService);

    summaryCalls = [];
    jest
      .spyOn(receiptFilterService, "getReceiptSummaryForGroup")
      .mockImplementation((groupId: string, configurationGroupId: number) => {
        summaryCalls.push({ groupId, configurationGroupId });
        return of(summaryResponse);
      });
    jest
      .spyOn(receiptFilterService, "getPagedReceiptsForGroups")
      .mockReturnValue(of({ data: [], totalCount: 0 } as any));

    fixture = TestBed.createComponent(ReceiptsTableComponent);
    component = fixture.componentInstance;
    Object.defineProperty(component, "table", {
      value: () => ({ selection: {}, changed: of(undefined) }),
    });
  });

  it("requests the summary on a filter change but not on a page or sort change", async () => {
    store.dispatch(new SetGroups([summaryGroup(1, "Alpha", true)]));
    component.groupId = "1";

    component.getFilteredReceipts();
    expect(summaryCalls.length).toEqual(1);

    (component as any).getFilteredReceiptsPage();
    expect(summaryCalls.length).toEqual(1);

    component.updatePageData({ pageIndex: 2, pageSize: 50, length: 100 } as any);
    expect(summaryCalls.length).toEqual(1);

    // And a genuine filter change asks again.
    component.getFilteredReceipts();
    expect(summaryCalls.length).toEqual(2);
  });

  it("stores the summary it receives", () => {
    store.dispatch(new SetGroups([summaryGroup(1, "Alpha", true)]));
    component.groupId = "1";

    component.getFilteredReceipts();

    expect(component.summary()).toEqual(summaryResponse);
  });

  // The guard lives inside the switchMap, so an install that has not opted in pays no request.
  it("issues no request when no configured group resolves", () => {
    store.dispatch(new SetGroups([summaryGroup(1, "Alpha", false, true)]));
    component.groupId = "1";

    component.getFilteredReceipts();

    expect(summaryCalls.length).toEqual(0);
    expect(component.summary()).toBeUndefined();
  });

  // Same contract as the table's stream: an error must not complete the outer subscription, or
  // every later refresh dies silently.
  it("survives a failed summary request and refreshes again after", () => {
    store.dispatch(new SetGroups([summaryGroup(1, "Alpha", true)]));
    component.groupId = "1";

    jest
      .spyOn(receiptFilterService, "getReceiptSummaryForGroup")
      .mockReturnValueOnce(throwError(() => new Error("boom")))
      .mockReturnValue(of(summaryResponse));

    component.getFilteredReceipts();
    expect(component.summary()).toBeUndefined();

    component.getFilteredReceipts();
    expect(component.summary()).toEqual(summaryResponse);
  });

  describe("on the All group", () => {
    beforeEach(() => {
      store.dispatch(
        new SetGroups([
          summaryGroup(99, "All", false, true),
          summaryGroup(2, "Bravo", true),
          summaryGroup(1, "Alpha", true),
          summaryGroup(3, "Charlie", false),
        ])
      );
      component.groupId = "99";
    });

    it("offers every enabled group alphabetically and defaults to the first", () => {
      expect(component.summaryConfigGroupOptions().map((g) => g.name)).toEqual(["Alpha", "Bravo"]);
      expect(component.summaryConfigGroupId()).toEqual(1);
    });

    it("honours a persisted pick and sends it as the configuration group", () => {
      store.dispatch(new SetSummaryConfigGroupId(2));

      component.getFilteredReceipts();

      expect(summaryCalls).toEqual([{ groupId: "99", configurationGroupId: 2 }]);
    });

    // A persisted group that has since disabled its summary, or that the user has left.
    it("falls back when the persisted pick is stale", () => {
      store.dispatch(new SetSummaryConfigGroupId(3));

      expect(component.summaryConfigGroupId()).toEqual(1);
    });

    it("re-requests only the summary when the chip changes", () => {
      component.summaryConfigGroupSelected(2);

      expect(summaryCalls).toEqual([{ groupId: "99", configurationGroupId: 2 }]);
    });
  });

  it("offers no configuration chips on a real group", () => {
    store.dispatch(
      new SetGroups([summaryGroup(1, "Alpha", true), summaryGroup(2, "Bravo", true)])
    );
    component.groupId = "1";

    expect(component.summaryConfigGroupOptions()).toEqual([]);
    expect(component.summaryConfigGroupId()).toEqual(1);
  });
});
