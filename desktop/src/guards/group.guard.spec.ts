import { TestBed } from "@angular/core/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { GroupState, SetSelectedDashboardId, SetSelectedGroupId } from "../store";
import { GroupGuard } from "./group.guard";

describe("GroupGuard", () => {
  let guard: GroupGuard;
  let store: Store;
  let navigateSpy: jest.SpyInstance;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([GroupState]), RouterTestingModule],
    });

    navigateSpy = jest.spyOn(TestBed.inject(Router), "navigate");
    navigateSpy.mockReturnValue(Promise.resolve(true));
    store = TestBed.inject(Store);
    guard = TestBed.inject(GroupGuard);
  });

  it("should be created", () => {
    expect(guard).toBeTruthy();
  });

  it("should return true", () => {
    store.reset({ groups: { groups: [{ id: "1" }] } });
    const result = guard.canActivate(
      { params: { groupId: "1" } } as any,
      {} as any
    );

    expect(result).toBe(true);
  });

  it("should return false", () => {
    store.reset({ groups: { groups: [{ id: "1" }] } });
    let storeSpy = jest.spyOn(store, "dispatch");

    const result = guard.canActivate(
      {
        params: { groupId: "2" },
        data: {
          groupGuardBasePath: "dashboard/group",
        },
      } as any,
      {} as any
    );

    expect(result).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith(["dashboard/group/1"]);
    expect(storeSpy).toHaveBeenCalledWith(new SetSelectedGroupId("1"));
  });

  it("should reset selected dashboard id", () => {
    store.reset({ groups: { groups: [{ id: "1" }], selectedGroupId: "3" } });
    let storeSpy = jest.spyOn(store, "dispatch");

    const result = guard.canActivate(
      {
        params: { groupId: "1" },
        data: {
          groupGuardBasePath: "dashboard/group",
        },
      } as any,
      {} as any
    );

    expect(result).toBe(true);
    expect(storeSpy).toHaveBeenCalledWith(
      new SetSelectedDashboardId(undefined)
    );
  });

  it("should reset selected dashboard id when group not found", () => {
    store.reset({ groups: { groups: [{ id: "1" }], selectedGroupId: "3" } });
    let storeSpy = jest.spyOn(store, "dispatch");

    const result = guard.canActivate(
      {
        params: { groupId: "70" },
        data: {
          groupGuardBasePath: "dashboard/group",
        },
      } as any,
      {} as any
    );

    expect(result).toBe(false);
    expect(storeSpy).toHaveBeenCalledWith(
      new SetSelectedDashboardId(undefined)
    );
  });

  it("should not reset selected dashboard id when group not found", () => {
    store.reset({ groups: { groups: [{ id: "1" }], selectedGroupId: "70" } });
    let storeSpy = jest.spyOn(store, "dispatch");

    const result = guard.canActivate(
      {
        params: { groupId: "70" },
        data: {
          groupGuardBasePath: "dashboard/group",
        },
      } as any,
      {} as any
    );

    expect(result).toBe(false);
    expect(storeSpy).toHaveBeenCalledWith(
      new SetSelectedDashboardId(undefined)
    );
  });
});
