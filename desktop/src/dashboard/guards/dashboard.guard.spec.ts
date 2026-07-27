import { TestBed } from "@angular/core/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { DashboardState } from "src/store/dashboard.state";
import { Dashboard } from "../../open-api";
import { GroupState, SetSelectedDashboardId } from "../../store";
import { dashboardGuard } from "./dashboard.guard";

describe("dashboardGuard", () => {
  let store: Store;
  let navigateSpy: jest.SpyInstance;

  const executeGuard = (dashboardId: string) =>
    TestBed.runInInjectionContext(() =>
      dashboardGuard({ params: { dashboardId } } as any, {} as any)
    );

  const dashboard = (id: number): Dashboard => ({
    id,
    name: `dashboard-${id}`,
    groupId: 1,
    userId: 1,
  });

  const seedDashboards = (dashboards: Dashboard[]) => {
    store.reset({
      ...store.snapshot(),
      groups: { ...store.snapshot().groups, selectedGroupId: "1" },
      dashboards: {
        dashboards: dashboards.length ? { "1": dashboards } : {},
      },
    });
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule,
        NgxsModule.forRoot([GroupState, DashboardState]),
      ],
    });
    store = TestBed.inject(Store);
    navigateSpy = jest
      .spyOn(TestBed.inject(Router), "navigate")
      .mockResolvedValue(true);
  });

  it("allows activation and selects the dashboard when it exists in the group", () => {
    seedDashboards([dashboard(5)]);
    const dispatchSpy = jest.spyOn(store, "dispatch");

    expect(executeGuard("5")).toBe(true);
    expect(dispatchSpy).toHaveBeenCalledWith(new SetSelectedDashboardId("5"));
    expect(navigateSpy).not.toHaveBeenCalled();
  });

  it("redirects to the group base when the dashboard id is not in the group", () => {
    seedDashboards([dashboard(5)]);

    expect(executeGuard("999")).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith(["/dashboard/group/1"]);
  });

  it("redirects when the group's dashboards have not been loaded yet", () => {
    // Pre-fix failure mode: the resolver must populate DashboardState before the
    // child route activates (SetDashboardsForGroup now awaits the HTTP), otherwise
    // a valid id is bounced. With an empty store even a real id redirects.
    seedDashboards([]);

    expect(executeGuard("5")).toBe(false);
    expect(navigateSpy).toHaveBeenCalledWith(["/dashboard/group/1"]);
  });
});
