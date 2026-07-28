import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import {
  HttpTestingController,
  provideHttpClientTesting,
} from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { Dashboard } from "../open-api";
import { DashboardState } from "./dashboard.state";
import { SetDashboardsForGroup } from "./dashboard.state.actions";

describe("DashboardState", () => {
  let store: Store;
  let httpMock: HttpTestingController;

  const dashboard = (id: number, groupId: number): Dashboard => ({
    id,
    name: `dashboard-${id}`,
    groupId,
    userId: 1,
  });

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([DashboardState])],
      providers: [
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });
    store = TestBed.inject(Store);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe("SetDashboardsForGroup", () => {
    // Regression guard for the root cause of "dashboards don't load on group
    // switch": the action must return its HTTP observable so NGXS keeps the
    // dispatch (and the route resolver that relies on it) open until the store is
    // populated. The old fire-and-forget `.subscribe()` form completed the
    // dispatch before the request resolved and would fail this test.
    it("patches the group's dashboards only after the service responds (awaits the HTTP)", () => {
      let completed = false;
      store
        .dispatch(new SetDashboardsForGroup("2"))
        .subscribe(() => (completed = true));

      expect(completed).toBe(false);
      expect(
        store.selectSnapshot(DashboardState.getDashboardsByGroupId("2"))
      ).toEqual([]);

      httpMock
        .expectOne((r) => r.method === "GET" && r.url.endsWith("/dashboard/2"))
        .flush([dashboard(1, 2)]);

      expect(completed).toBe(true);
      expect(
        store.selectSnapshot(DashboardState.getDashboardsByGroupId("2"))
      ).toEqual([dashboard(1, 2)]);
    });

    it("does not clobber other groups' dashboards", () => {
      store.reset({
        ...store.snapshot(),
        dashboards: { dashboards: { "1": [dashboard(9, 1)] } },
      });

      store.dispatch(new SetDashboardsForGroup("2"));
      httpMock
        .expectOne((r) => r.method === "GET" && r.url.endsWith("/dashboard/2"))
        .flush([dashboard(1, 2)]);

      expect(
        store.selectSnapshot(DashboardState.getDashboardsByGroupId("1"))
      ).toEqual([dashboard(9, 1)]);
      expect(
        store.selectSnapshot(DashboardState.getDashboardsByGroupId("2"))
      ).toEqual([dashboard(1, 2)]);
    });
  });

  describe("getDashboardsByGroupId", () => {
    it("returns [] for a group with no dashboards", () => {
      expect(
        store.selectSnapshot(DashboardState.getDashboardsByGroupId("999"))
      ).toEqual([]);
    });
  });
});
