import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { DEFAULT_RECEIPT_TABLE_COLUMNS } from "src/interfaces";
import { FilterOperation, Group } from "../open-api";
import { GroupState } from "./group.state";
import { RemoveGroup, SetSelectedGroupId } from "./group.state.actions";
import { defaultReceiptFilter, ReceiptTableState } from "./receipt-table.state";

describe("GroupState receipt-filter reset on group change", () => {
  let store: Store;

  const groups = [
    { id: 1, name: "Group One" },
    { id: 2, name: "Group Two" },
  ] as Group[];

  const nonDefaultFilter = {
    ...defaultReceiptFilter,
    name: { operation: FilterOperation.Contains, value: "coffee" },
  } as any;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([GroupState, ReceiptTableState])],
    }).compileComponents();

    store = TestBed.inject(Store);
  });

  function seed(selectedGroupId: string): void {
    store.reset({
      groups: {
        groups,
        selectedGroupId,
        selectedDashboardId: "",
      },
      receiptTable: {
        page: 3,
        pageSize: 50,
        orderBy: "created_at",
        sortDirection: "desc",
        filter: nonDefaultFilter,
        columnConfig: DEFAULT_RECEIPT_TABLE_COLUMNS,
      },
    });
  }

  it("resets the receipt filter and page when switching to a different group", () => {
    seed("1");

    store.dispatch(new SetSelectedGroupId("2"));

    expect(store.selectSnapshot(ReceiptTableState.filterData).filter).toEqual(
      defaultReceiptFilter
    );
    expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(1);
    expect(store.selectSnapshot(GroupState.selectedGroupId)).toEqual("2");
  });

  it("preserves the receipt filter and page when re-selecting the same group", () => {
    seed("1");

    store.dispatch(new SetSelectedGroupId("1"));

    expect(store.selectSnapshot(ReceiptTableState.filterData).filter).toEqual(
      nonDefaultFilter
    );
    expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(3);
  });

  it("resets when an empty payload resolves to a different group than the current", () => {
    // Current group is 2; an empty payload resolves to groups[0] (id 1), a change.
    seed("2");

    store.dispatch(new SetSelectedGroupId(""));

    expect(store.selectSnapshot(GroupState.selectedGroupId)).toEqual("1");
    expect(store.selectSnapshot(ReceiptTableState.filterData).filter).toEqual(
      defaultReceiptFilter
    );
    expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(1);
  });

  it("resets the receipt filter and page when the active group is deleted", () => {
    // Selected group is 2; deleting it falls back to groups[0] (id 1), a change.
    seed("2");

    store.dispatch(new RemoveGroup("2"));

    expect(store.selectSnapshot(GroupState.selectedGroupId)).toEqual("1");
    expect(store.selectSnapshot(ReceiptTableState.filterData).filter).toEqual(
      defaultReceiptFilter
    );
    expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(1);
  });

  it("preserves the receipt filter and page when a non-active group is deleted", () => {
    // Selected group is 1; deleting group 2 leaves the selection unchanged.
    seed("1");

    store.dispatch(new RemoveGroup("2"));

    expect(store.selectSnapshot(GroupState.selectedGroupId)).toEqual("1");
    expect(store.selectSnapshot(ReceiptTableState.filterData).filter).toEqual(
      nonDefaultFilter
    );
    expect(store.selectSnapshot(ReceiptTableState.page)).toEqual(3);
  });
});
