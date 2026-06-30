import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { CanActivateFn, Router, UrlTree } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Permission } from "../open-api";
import { AuthState } from "../store/auth.state";
import { SetPermissions } from "../store/auth.state.actions";
import { GroupState } from "../store/group.state";
import { systemSettingsLandingGuard } from "./system-settings-landing.guard";

describe("systemSettingsLandingGuard", () => {
  const executeGuard: CanActivateFn = (...params) =>
    TestBed.runInInjectionContext(() => systemSettingsLandingGuard(...params));
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

  const run = (): UrlTree =>
    executeGuard({} as any, {} as any) as UrlTree;

  it("redirects to the first readable tab in render order", () => {
    store.dispatch(
      new SetPermissions(
        [Permission.AppPromptsRead, Permission.AppSystemTasksRead],
        {}
      )
    );

    expect(run().toString()).toBe("/system-settings/prompts");
  });

  it("lands on System Settings for a user who can read every tab", () => {
    store.dispatch(
      new SetPermissions(
        [
          Permission.AppSystemSettingsRead,
          Permission.AppReceiptProcessingSettingsRead,
          Permission.AppPromptsRead,
          Permission.AppSystemEmailsRead,
          Permission.AppSystemTasksRead,
        ],
        {}
      )
    );

    expect(run().toString()).toBe("/system-settings/settings/view");
  });

  it("falls through to a later tab when earlier ones are not readable", () => {
    store.dispatch(new SetPermissions([Permission.AppSystemTasksRead], {}));

    expect(run().toString()).toBe("/system-settings/system-tasks");
  });

  it("prefers an earlier-rendered tab over system-emails", () => {
    store.dispatch(
      new SetPermissions(
        [Permission.AppSystemEmailsRead, Permission.AppPromptsRead],
        {}
      )
    );

    expect(run().toString()).toBe("/system-settings/prompts");
  });

  it("redirects to the dashboard when no tab is readable", () => {
    store.dispatch(new SetPermissions([], {}));

    const dashboard = store.selectSnapshot(GroupState.dashboardLink);
    expect(run().toString()).toBe(router.parseUrl(dashboard).toString());
  });
});
