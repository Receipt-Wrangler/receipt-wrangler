import { CommonModule } from "@angular/common";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatCardModule } from "@angular/material/card";
import { MatDialogModule } from "@angular/material/dialog";
import { MatListModule } from "@angular/material/list";
import { ActivatedRoute, Params } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { BehaviorSubject, of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { DashboardState } from "src/store/dashboard.state";
import { ApiModule, Dashboard, UserService, Widget, WidgetType } from "../../open-api";
import { SummaryCardComponent } from "../../shared-ui/summary-card/summary-card.component";
import { GroupState } from "../../store";
import { DashboardRoutingModule } from "../dashboard-routing.module";
import { DashboardComponent } from "./dashboard.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("DashboardComponent", () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;
  let dashboards: Dashboard[];
  let store: Store;

  beforeEach(async () => {
    dashboards = [
      {
        id: 1,
        userId: 1,
        name: "Test Dashboard",
        widgets: [],
      } as Dashboard,
      {
        id: 2,
        userId: 1,
        name: "Test Dashboard 2",
        widgets: [],
      } as Dashboard,
    ];
    await TestBed.configureTestingModule({
    declarations: [DashboardComponent, SummaryCardComponent],
    imports: [ApiModule,
        CommonModule,
        DashboardRoutingModule,
        MatCardModule,
        MatDialogModule,
        MatListModule,
        NgxsModule.forRoot([GroupState, DashboardState]),
        PipesModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {
                params: new BehaviorSubject<Params>({ id: "1" }),
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    store = TestBed.inject(Store);

    store.reset({
      ...store.snapshot(),
      dashboards: {
        dashboards: {
          "1": dashboards,
        },
      },
    });
    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should set dashboards", () => {
    store.reset({
      ...store.snapshot(),
      groups: {
        selectedGroupId: "1",
      },
    });
    component.ngOnInit();
    expect(component.dashboards()).toEqual(dashboards);
  });

  it("should set selected dashboard", () => {
    const store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      groups: {
        selectedGroupId: "1",
        selectedDashboardId: "2",
      },
    });
    component.ngOnInit();

    expect(component.selectedDashboard()).toEqual(dashboards[1]);
  });

  // Regression: every widget used to double-fetch on load. The `dashboards` slice
  // emits twice (the cached replay, then the resolver's refetch with brand-new
  // widget object references but the SAME widget id), and the widget loop had no
  // track key, so emission 2 destroyed and recreated every widget → each re-ran its
  // init fetch. `@for ... track widget.id` keeps the instance, so it does not
  // re-fetch. Proven here with the (non-report) GROUP_SUMMARY widget, whose effect
  // calls UserService.getAmountOwedForUser.
  it("fetches a widget's data only once when the dashboards slice re-emits with the same widget id", async () => {
    const owedSpy = jest
      .spyOn(TestBed.inject(UserService), "getAmountOwedForUser")
      .mockReturnValue(of({}) as any);

    // A fresh dashboard + fresh widget object each call (new references), same id.
    const summaryDashboards = (): { [groupId: string]: Dashboard[] } => ({
      "1": [
        {
          id: 1,
          userId: 1,
          name: "Summary",
          groupId: 1,
          widgets: [
            {
              id: 10,
              dashboardId: 1,
              widgetType: WidgetType.GroupSummary,
              name: "S",
            } as Widget,
          ],
        } as Dashboard,
      ],
    });

    // Emission 1 — the cached dashboard mounts the summary widget → one fetch.
    // The summary card fetches from a constructor effect, so flush effects and let
    // the fixture stabilize before asserting the spy count (zoneless CD).
    store.reset({
      ...store.snapshot(),
      groups: { selectedGroupId: "1", selectedDashboardId: "1" },
      dashboards: { dashboards: summaryDashboards() },
    });
    TestBed.flushEffects();
    await fixture.whenStable();

    expect(owedSpy).toHaveBeenCalledTimes(1);

    // Emission 2 — the resolver's refetch replaces the slice with new object
    // references (same widget id 10). The widget must NOT re-fetch.
    store.reset({
      ...store.snapshot(),
      groups: { selectedGroupId: "1", selectedDashboardId: "1" },
      dashboards: { dashboards: summaryDashboards() },
    });
    TestBed.flushEffects();
    await fixture.whenStable();

    expect(owedSpy).toHaveBeenCalledTimes(1);
  });
});
