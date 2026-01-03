import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { SetDashboardsForGroup } from "src/store/dashboard.state.actions";
import { Dashboard, DashboardService } from "../../open-api";
import { dashboardResolverFn } from "./dashboard.resolver";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("dashboardResolver", () => {
  const executeResolver: ResolveFn<Observable<Dashboard[]>> = (
    ...resolverParameters
  ) =>
    TestBed.runInInjectionContext(() =>
      dashboardResolverFn(...resolverParameters)
    );

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [NgxsModule.forRoot([])],
    providers: [DashboardService, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
  });

  it("should be created", () => {
    expect(executeResolver).toBeTruthy();
  });

  it("should attempt to get dashboards by group id", () => {
    const dispatchSpy = jest.spyOn(TestBed.inject(Store), "dispatch");

    executeResolver(
      {
        params: {
          groupId: "1",
        },
      } as any,
      {} as any
    );

    expect(dispatchSpy).toHaveBeenCalledWith(new SetDashboardsForGroup("1"));
  });
});
