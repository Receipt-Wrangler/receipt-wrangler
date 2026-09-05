import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { DEFAULT_RECEIPT_TABLE_COLUMNS } from "src/interfaces";
import { SetColumnConfig, SetReceiptFilterData } from "src/store/receipt-table.actions";
import { ReceiptTableState } from "src/store/receipt-table.state";
import { ApiModule, CustomField, CustomFieldType, Permission, Receipt } from "../../open-api";
import { AuthState } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
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
        NgxsModule.forRoot([ReceiptTableState, AuthState]),
        ReactiveFormsModule,
        MatSnackBarModule,
        MatTooltipModule,
        MatDialogModule,
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
      component.customFields = [customField(7, "Vendor")];
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
      component.customFields = [];
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
      component.customFields = [];
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

    // The sort is persisted per browser and shared across accounts, so the table
    // can start up asking the API to order by a column that no longer exists -
    // which the API rejects outright, failing the very first load.
    it("resets a sort on a custom field that no longer exists", () => {
      component.customFields = [];
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
      component.customFields = [customField_(7)];
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
});
