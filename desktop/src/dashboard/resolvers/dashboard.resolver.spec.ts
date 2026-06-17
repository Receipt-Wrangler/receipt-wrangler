import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { SetDashboardsForGroup } from "src/store/dashboard.state.actions";
import { Dashboard, DashboardService, Permission } from "../../open-api";
import { AuthState } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
import { dashboardResolverFn } from "./dashboard.resolver";

describe("dashboardResolver", () => {
  let store: Store;

  const executeResolver: ResolveFn<Observable<Dashboard[]>> = (
    ...resolverParameters
  ) =>
    TestBed.runInInjectionContext(() =>
      dashboardResolverFn(...resolverParameters)
    );

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState])],
      providers: [
        DashboardService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });
    store = TestBed.inject(Store);
  });

  it("should be created", () => {
    expect(executeResolver).toBeTruthy();
  });

  it("should attempt to get dashboards by group id when permitted", () => {
    store.dispatch(
      new SetPermissions([], { 1: [Permission.GroupDashboardsRead] })
    );
    const dispatchSpy = jest.spyOn(store, "dispatch");

    executeResolver({ params: { groupId: "1" } } as any, {} as any);

    expect(dispatchSpy).toHaveBeenCalledWith(new SetDashboardsForGroup("1"));
  });

  it("does not fetch dashboards without group.dashboards.read", (done) => {
    store.dispatch(new SetPermissions([], {}));
    const dispatchSpy = jest.spyOn(store, "dispatch");

    const result = executeResolver(
      { params: { groupId: "1" } } as any,
      {} as any
    ) as Observable<Dashboard[]>;

    expect(dispatchSpy).not.toHaveBeenCalledWith(new SetDashboardsForGroup("1"));
    result.subscribe((dashboards) => {
      expect(dashboards).toEqual([]);
      done();
    });
  });
});
