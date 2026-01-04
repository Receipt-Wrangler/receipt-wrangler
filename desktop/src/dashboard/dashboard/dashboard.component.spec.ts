import { CommonModule } from "@angular/common";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatCardModule } from "@angular/material/card";
import { MatDialogModule } from "@angular/material/dialog";
import { MatListModule } from "@angular/material/list";
import { ActivatedRoute, Params } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { BehaviorSubject } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { DashboardState } from "src/store/dashboard.state";
import { ApiModule, Dashboard } from "../../open-api";
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
    expect(component.dashboards).toEqual(dashboards);
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

    expect(component.selectedDashboard).toEqual(dashboards[1]);
  });
});
