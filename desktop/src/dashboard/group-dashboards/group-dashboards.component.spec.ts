import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed, } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute, Params, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { BehaviorSubject } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { DashboardState } from "src/store/dashboard.state";
import { ButtonModule } from "../../button";
import { Dashboard, DashboardService } from "../../open-api";
import { GroupState, SetSelectedDashboardId } from "../../store";
import { GroupDashboardsComponent } from "./group-dashboards.component";

describe("GroupDashboardsComponent", () => {
  let component: GroupDashboardsComponent;
  let fixture: ComponentFixture<GroupDashboardsComponent>;
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [GroupDashboardsComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [PipesModule,
        MatDialogModule,
        NgxsModule.forRoot([GroupState, DashboardState]),
        PipesModule,
        ButtonModule,
        MatSnackBarModule],
      providers: [
        DashboardService,
        {
          provide: ActivatedRoute,
          useValue: {
            params: new BehaviorSubject<Params>({}),
            snapshot: {
              data: {
                dashboards: [],
              },
            },
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    store = TestBed.inject(Store);
    store.reset({
      groups: {
        selectedGroupId: "1",
        groups: [],
      },
    });
    fixture = TestBed.createComponent(GroupDashboardsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should set dashboards with dashboards", () => {
    const dashboards: Dashboard[] = [
      {
        id: 1,
        name: "test",
        groupId: 1,
        userId: 1,
      },
    ];
    store.reset({
      groups: {
        selectedGroupId: "1",
      },
      dashboards: {
        dashboards: {
          "1": dashboards,
        },
      },
    });

    component.ngOnInit();

    expect(component.dashboards).toEqual(dashboards);
  });

  it("should set dashboards with dashboards on seleced group id change", () => {
    const dashboards: Dashboard[] = [
      {
        id: 1,
        name: "test",
        groupId: 1,
        userId: 1,
      },
    ];
    const newDashboards: Dashboard[] = [
      {
        id: 2,
        name: "test",
        groupId: 1,
        userId: 1,
      },
    ];
    const activatedRoute = TestBed.inject(ActivatedRoute);
    store.reset({
      groups: {
        selectedGroupId: "1",
      },
      dashboards: {
        dashboards: {
          "1": dashboards,
        },
      },
    });
    component.ngOnInit();

    expect(component.dashboards).toEqual(dashboards);

    store.reset({
      groups: {
        selectedGroupId: "2",
      },
      dashboards: {
        dashboards: {
          "2": newDashboards,
        },
      },
    });

    (activatedRoute.params as any).next({
      dashboardId: 2,
    });

    expect(component.dashboards).toEqual(newDashboards);
  });

  it("should not navigate to selected dashboard", () => {
    const routerSpy = jest.spyOn(TestBed.inject(Router), "navigateByUrl");
    store.reset({
      groups: {
        selectedDashboardId: undefined,
      },
    });

    component.ngOnInit();

    expect(routerSpy).toHaveBeenCalledTimes(0);
  });

  it("should set selected dashboard id", () => {
    const store = TestBed.inject(Store);
    const storeSpy = jest.spyOn(store, "dispatch");

    component.setSelectedDashboardId(1);

    expect(storeSpy).toHaveBeenCalledWith(new SetSelectedDashboardId("1"));
  });
});
