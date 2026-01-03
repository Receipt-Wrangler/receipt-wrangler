import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { SortDirection } from "../open-api";
import { SystemTaskTableState } from "./system-task-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "./system-task-table.state.actions";

describe("SystemTaskTableState", () => {
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([SystemTaskTableState])]
    });

    store = TestBed.inject(Store);
  });

  it("should set page", () => {
    store.dispatch(new SetPage(2));
    const state = store.selectSnapshot(SystemTaskTableState.state);
    expect(state.page).toBe(2);
  });

  it("should set page size", () => {
    store.dispatch(new SetPageSize(100));
    const state = store.selectSnapshot(SystemTaskTableState.state);
    expect(state.pageSize).toBe(100);
  });

  it("should set order by", () => {
    store.dispatch(new SetOrderBy("type"));
    const state = store.selectSnapshot(SystemTaskTableState.state);
    expect(state.orderBy).toBe("type");
  });

  it("should set sort direction", () => {
    store.dispatch(new SetSortDirection(SortDirection.Asc));
    const state = store.selectSnapshot(SystemTaskTableState.state);
    expect(state.sortDirection).toBe(SortDirection.Asc);
  });
});
