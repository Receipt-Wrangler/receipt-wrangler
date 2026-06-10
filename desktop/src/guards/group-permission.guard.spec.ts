import { TestBed } from "@angular/core/testing";
import { ActivatedRouteSnapshot, CanActivateFn, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { AuthState } from "../store/auth.state";
import { SetPermissions } from "../store/auth.state.actions";
import { GroupState } from "../store/group.state";
import { SetSelectedGroupId } from "../store/group.state.actions";
import { groupPermissionGuard } from "./group-permission.guard";

describe("groupPermissionGuard", () => {
  const executeGuard: CanActivateFn = (...params) =>
    TestBed.runInInjectionContext(() => groupPermissionGuard(...params));
  let store: Store;
  let navigateSpy: jest.SpyInstance;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState, GroupState])],
    });

    store = TestBed.inject(Store);
    navigateSpy = jest
      .spyOn(TestBed.inject(Router), "navigate")
      .mockResolvedValue(true);
  });

  it("allows when the user holds the group permission for the selected group", () => {
    store.dispatch(new SetSelectedGroupId("1"));
    store.dispatch(new SetPermissions([], { 1: ["group.view"] }));

    const route = {
      data: { groupPermission: "group.view" },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(true);
    expect(navigateSpy).not.toHaveBeenCalled();
  });

  it("denies and redirects when the user lacks the group permission", () => {
    store.dispatch(new SetSelectedGroupId("1"));
    store.dispatch(new SetPermissions([], { 1: ["group.receipts.read"] }));

    const route = {
      data: { groupPermission: "group.view" },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith([
      store.selectSnapshot(GroupState.dashboardLink),
    ]);
  });

  it("resolves the group id from the route param when useRouteGroupId is set", () => {
    store.dispatch(new SetSelectedGroupId("99")); // selected group differs from route
    store.dispatch(new SetPermissions([], { 5: ["group.view"] }));

    const route = {
      data: { groupPermission: "group.view", useRouteGroupId: true },
      params: { id: "5" },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(true);
    expect(navigateSpy).not.toHaveBeenCalled();
  });

  it("resolves the group id from the parent route param when not on the param itself", () => {
    store.dispatch(new SetPermissions([], { 7: ["group.view"] }));

    const route = {
      data: { groupPermission: "group.view", useRouteGroupId: true },
      params: {},
      parent: { params: { id: "7" } },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(true);
  });

  it("allows a non-member who holds an orApp fallback permission", () => {
    store.dispatch(new SetSelectedGroupId("3"));
    // No group permissions for group 3, but the app fallback is held.
    store.dispatch(new SetPermissions(["app.groups.read"], {}));

    const route = {
      data: {
        groupPermission: "group.view",
        orAppPermissions: ["app.groups.read"],
      },
    } as unknown as ActivatedRouteSnapshot;

    expect(executeGuard(route, {} as any)).toBe(true);
    expect(navigateSpy).not.toHaveBeenCalled();
  });
});
