import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { ActivatedRouteSnapshot, CanActivateFn, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Permission } from "../../open-api";
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";
import { GroupState } from "../../store/group.state";
import { groupDashboardReadGuard } from "./group-dashboard-read.guard";

describe("groupDashboardReadGuard", () => {
  const executeGuard: CanActivateFn = (...params) =>
    TestBed.runInInjectionContext(() => groupDashboardReadGuard(...params));
  let store: Store;
  let navigateSpy: jest.SpyInstance;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState, GroupState])],
      providers: [provideZonelessChangeDetection()],
    });

    store = TestBed.inject(Store);
    navigateSpy = jest
      .spyOn(TestBed.inject(Router), "navigate")
      .mockResolvedValue(true);
  });

  it("allows when the user holds group.dashboards.read for the route group", () => {
    store.dispatch(
      new SetPermissions([], { 1: [Permission.GroupDashboardsRead] })
    );

    const route = {
      params: { groupId: "1" },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(true);
    expect(navigateSpy).not.toHaveBeenCalled();
  });

  it("denies and redirects to the group's receipt list when lacking the permission", () => {
    store.dispatch(new SetPermissions([], { 1: ["group.receipts.read"] }));

    const route = {
      params: { groupId: "1" },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith(["/receipts/group", "1"]);
  });
});
