import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { ActivatedRouteSnapshot, CanActivateFn, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { AuthState } from "../store/auth.state";
import { SetPermissions } from "../store/auth.state.actions";
import { GroupState } from "../store/group.state";
import { appPermissionGuard } from "./app-permission.guard";

describe("appPermissionGuard", () => {
  const executeGuard: CanActivateFn = (...params) =>
    TestBed.runInInjectionContext(() => appPermissionGuard(...params));
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

  const route = (appPermissions: string[]) =>
    ({ data: { appPermissions } }) as unknown as ActivatedRouteSnapshot;

  it("allows when the user holds any required app permission", () => {
    store.dispatch(new SetPermissions(["app.users.read"], {}));

    const result = executeGuard(route(["app.users.read"]), {} as any);

    expect(result).toBe(true);
    expect(navigateSpy).not.toHaveBeenCalled();
  });

  it("denies and redirects to the dashboard when the user lacks the permission", () => {
    store.dispatch(new SetPermissions(["app.roles.read"], {}));

    const result = executeGuard(route(["app.users.read"]), {} as any);

    expect(result).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith([
      store.selectSnapshot(GroupState.dashboardLink),
    ]);
  });

  it("denies when no permissions are required (deny by default)", () => {
    store.dispatch(new SetPermissions(["app.users.read"], {}));

    const result = executeGuard(route([]), {} as any);

    expect(result).toBe(false);
    expect(navigateSpy).toHaveBeenCalled();
  });
});
