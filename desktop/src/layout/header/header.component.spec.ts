import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { Store } from "@ngxs/store";
import { NgbPopoverModule } from "@ng-bootstrap/ng-bootstrap";
import { of } from "rxjs";
import { ToggleIsSidebarOpen } from "src/store/layout.state.actions";
import { DirectivesModule } from "../../directives/directives.module";
import { ApiModule, NotificationsService } from "../../open-api";
import { SetAuthState, SetPermissions } from "../../store/auth.state.actions";
import { SetSelectedGroupId } from "../../store/group.state.actions";
import { StoreModule } from "../../store/store.module";
import { HeaderComponent } from "./header.component";

describe("HeaderComponent", () => {
  let component: HeaderComponent;
  let fixture: ComponentFixture<HeaderComponent>;
  let store: Store;
  let notificationsService: NotificationsService;

  const logIn = () =>
    store.dispatch(
      new SetAuthState({ exp: Math.floor(Date.now() / 1000) + 3600 } as any)
    );

  const grantNotificationsRead = () =>
    store.dispatch(new SetPermissions(["app.notifications.read"], {}));

  const bellRendered = (): boolean =>
    !!fixture.nativeElement.querySelector('app-button[icon="notifications"]');

  const searchbarRendered = (): boolean =>
    !!fixture.nativeElement.querySelector("app-searchbar");

  const grantReceiptsSearch = () =>
    store.dispatch(new SetPermissions(["app.receipts.search"], {}));

  const dashboardButtonRendered = (): boolean =>
    !!fixture.nativeElement.querySelector('button[matTooltip="Dashboard"]');

  const selectGroupOne = () => store.dispatch(new SetSelectedGroupId("1"));

  const grantDashboardsRead = () =>
    store.dispatch(new SetPermissions([], { 1: ["group.dashboards.read"] }));

  beforeEach(async () => {
    // StoreModule persists auth (incl. permissions) to localStorage; clear it so
    // granted permissions don't leak across tests in this file.
    localStorage.clear();

    await TestBed.configureTestingModule({
    declarations: [HeaderComponent],
    imports: [ApiModule,
        DirectivesModule,
        MatDialogModule,
        MatSnackBarModule,
        NgbPopoverModule,
        StoreModule,
    ],
    providers: [provideZonelessChangeDetection(), provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()],
    schemas: [NO_ERRORS_SCHEMA],
}).compileComponents();

    store = TestBed.inject(Store);
    notificationsService = TestBed.inject(NotificationsService);
    fixture = TestBed.createComponent(HeaderComponent);
    component = fixture.componentInstance;
    fixture.autoDetectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should toggle sidebar", () => {
    const store = jest.spyOn(TestBed.inject(Store), "dispatch");
    component.toggleSidebar();

    expect(store).toHaveBeenCalledWith(new ToggleIsSidebarOpen());
  });

  it("should hide the notifications bell without app.notifications.read", async () => {
    logIn();
    await fixture.whenStable();

    expect(bellRendered()).toBe(false);
  });

  it("should show the notifications bell with app.notifications.read", async () => {
    logIn();
    grantNotificationsRead();
    await fixture.whenStable();

    expect(bellRendered()).toBe(true);
  });

  it("should hide the search bar without app.receipts.search", async () => {
    logIn();
    await fixture.whenStable();

    expect(searchbarRendered()).toBe(false);
  });

  it("should show the search bar with app.receipts.search", async () => {
    logIn();
    grantReceiptsSearch();
    await fixture.whenStable();

    expect(searchbarRendered()).toBe(true);
  });

  it("should hide the dashboard button without group.dashboards.read for the selected group", async () => {
    logIn();
    selectGroupOne();
    await fixture.whenStable();

    expect(dashboardButtonRendered()).toBe(false);
  });

  it("should show the dashboard button with group.dashboards.read for the selected group", async () => {
    logIn();
    selectGroupOne();
    grantDashboardsRead();
    await fixture.whenStable();

    expect(dashboardButtonRendered()).toBe(true);
  });

  it("should not fetch the notification count without app.notifications.read", async () => {
    const countSpy = jest
      .spyOn(notificationsService, "getNotificationCount")
      .mockReturnValue(of(0));

    logIn();
    await fixture.whenStable();

    expect(countSpy).not.toHaveBeenCalled();
  });

  it("should fetch the notification count once permission is present", async () => {
    const countSpy = jest
      .spyOn(notificationsService, "getNotificationCount")
      .mockReturnValue(of(3));

    logIn();
    grantNotificationsRead();
    TestBed.flushEffects();
    await fixture.whenStable();

    expect(countSpy).toHaveBeenCalledTimes(1);
    expect(component.notificationCount()).toBe(3);
  });
});
