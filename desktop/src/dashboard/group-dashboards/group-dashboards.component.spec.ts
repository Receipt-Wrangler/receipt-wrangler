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
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";
import { DirectivesModule } from "../../directives/directives.module";
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
        NgxsModule.forRoot([GroupState, DashboardState, AuthState]),
        PipesModule,
        ButtonModule,
        DirectivesModule,
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
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
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
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
        selectedGroupId: "1",
      },
      dashboards: {
        dashboards: {
          "1": dashboards,
        },
      },
    });

    component.ngOnInit();

    expect(component.dashboards()).toEqual(dashboards);
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
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
        selectedGroupId: "1",
      },
      dashboards: {
        dashboards: {
          "1": dashboards,
        },
      },
    });
    component.ngOnInit();

    expect(component.dashboards()).toEqual(dashboards);

    store.reset({
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
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

    expect(component.dashboards()).toEqual(newDashboards);
  });

  it("should not navigate to selected dashboard", () => {
    const routerSpy = jest.spyOn(TestBed.inject(Router), "navigateByUrl");
    store.reset({
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
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

  describe("dashboard CRUD permission gating", () => {
    const addButton = () =>
      fixture.nativeElement.querySelector("app-add-button");
    const editButton = () =>
      fixture.nativeElement.querySelector("app-edit-button");
    const deleteButton = () =>
      fixture.nativeElement.querySelector("app-delete-button");

    // The edit/delete controls only render once a dashboard is selected; seed
    // the selection so the permission gate is the only variable under test.
    const selectDashboard = () => {
      store.reset({
        ...store.snapshot(),
        groups: {
          ...store.snapshot().groups,
          selectedGroupId: "1",
          selectedDashboardId: "5",
        },
      });
    };

    const render = async () => {
      fixture.detectChanges();
      await fixture.whenStable();
    };

    it("shows the Add Dashboard button only with group.dashboards.create", async () => {
      store.dispatch(new SetPermissions([], { 1: ["group.dashboards.read"] }));
      await render();
      expect(addButton()).toBeFalsy();

      store.dispatch(new SetPermissions([], { 1: ["group.dashboards.create"] }));
      await render();
      expect(addButton()).toBeTruthy();
    });

    it("shows the Edit Dashboard button only with group.dashboards.update", async () => {
      selectDashboard();
      store.dispatch(new SetPermissions([], { 1: ["group.dashboards.read"] }));
      await render();
      expect(editButton()).toBeFalsy();

      store.dispatch(new SetPermissions([], { 1: ["group.dashboards.update"] }));
      await render();
      expect(editButton()).toBeTruthy();
    });

    it("shows the Delete Dashboard button only with group.dashboards.delete", async () => {
      selectDashboard();
      store.dispatch(new SetPermissions([], { 1: ["group.dashboards.update"] }));
      await render();
      expect(deleteButton()).toBeFalsy();

      store.dispatch(new SetPermissions([], { 1: ["group.dashboards.delete"] }));
      await render();
      expect(deleteButton()).toBeTruthy();
    });
  });
});
