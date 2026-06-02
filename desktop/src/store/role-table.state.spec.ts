import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { RoleTableState } from "./role-table.state";
import {
  SetOrderBy,
  SetPage,
  SetPageSize,
  SetScope,
  SetSortDirection,
} from "./role-table.state.actions";

describe("RoleTableState", () => {
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([RoleTableState])],
    });

    store = TestBed.inject(Store);
  });

  it("defaults to the first page sorted by name across all scopes", () => {
    const state = store.selectSnapshot(RoleTableState.state);
    expect(state.page).toBe(1);
    expect(state.pageSize).toBe(50);
    expect(state.orderBy).toBe("name");
    expect(state.sortDirection).toBe("asc");
    expect(store.selectSnapshot(RoleTableState.scope)).toBe("all");
  });

  it("should set page", () => {
    store.dispatch(new SetPage(2));
    expect(store.selectSnapshot(RoleTableState.state).page).toBe(2);
  });

  it("should set page size", () => {
    store.dispatch(new SetPageSize(100));
    expect(store.selectSnapshot(RoleTableState.state).pageSize).toBe(100);
  });

  it("should set order by", () => {
    store.dispatch(new SetOrderBy("name"));
    expect(store.selectSnapshot(RoleTableState.state).orderBy).toBe("name");
  });

  it("should set sort direction", () => {
    store.dispatch(new SetSortDirection("desc"));
    expect(store.selectSnapshot(RoleTableState.state).sortDirection).toBe("desc");
  });

  it("resets to the first page when the scope changes", () => {
    store.dispatch(new SetPage(4));
    store.dispatch(new SetScope("group"));

    expect(store.selectSnapshot(RoleTableState.scope)).toBe("group");
    expect(store.selectSnapshot(RoleTableState.state).page).toBe(1);
  });
});
