import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { CanActivateFn, Router, UrlTree } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Permission } from "../open-api";
import { AuthState } from "../store/auth.state";
import { SetPermissions } from "../store/auth.state.actions";
import { GroupState } from "../store/group.state";
import { settingsLandingGuard } from "./settings-landing.guard";

describe("settingsLandingGuard", () => {
  const executeGuard: CanActivateFn = (...params) =>
    TestBed.runInInjectionContext(() => settingsLandingGuard(...params));
  let store: Store;
  let router: Router;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState, GroupState])],
      providers: [provideZonelessChangeDetection()],
    });

    store = TestBed.inject(Store);
    router = TestBed.inject(Router);
  });

  const run = (): UrlTree => executeGuard({} as any, {} as any) as UrlTree;

  it("redirects to user-profile when account is readable", () => {
    store.dispatch(new SetPermissions([Permission.AppAccountRead], {}));

    expect(run().toString()).toBe("/settings/user-profile/view");
  });

  it("falls through to user-preferences when account is not readable", () => {
    store.dispatch(new SetPermissions([Permission.AppUserPreferencesRead], {}));

    expect(run().toString()).toBe("/settings/user-preferences/view");
  });

  it("falls through to api-keys when only it is readable", () => {
    store.dispatch(new SetPermissions([Permission.AppApiKeysRead], {}));

    expect(run().toString()).toBe("/settings/api-keys/view");
  });

  it("redirects to the dashboard when no tab is readable", () => {
    store.dispatch(new SetPermissions([], {}));

    const dashboard = store.selectSnapshot(GroupState.dashboardLink);
    expect(run().toString()).toBe(router.parseUrl(dashboard).toString());
  });
});
