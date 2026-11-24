import { InteractivityChecker } from "@angular/cdk/a11y";
import { OverlayContainer } from "@angular/cdk/overlay";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { ToggleIsSidebarOpen } from "src/store/layout.state.actions";
import { ApiModule, NotificationsService } from "../../open-api";
import { SetAuthState, SetPermissions } from "../../store/auth.state.actions";
import { SetSelectedGroupId } from "../../store/group.state.actions";
import { provideTestStore } from "../../store/test-store";
import { HeaderComponent } from "./header.component";
import { provideRouter } from "@angular/router";

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
    // The test store persists auth (incl. permissions) to localStorage; clear it so
    // granted permissions don't leak across tests in this file.
    localStorage.clear();

    await TestBed.configureTestingModule({
    imports: [ApiModule,
        MatDialogModule,
        MatSnackBarModule,
        HeaderComponent,],
    providers: [provideTestStore(), provideRouter([]), provideZonelessChangeDetection(), provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()],
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
    const store = vi.spyOn(TestBed.inject(Store), "dispatch");
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
    const countSpy = vi
      .spyOn(notificationsService, "getNotificationCount")
      .mockReturnValue(of(0) as any);

    logIn();
    await fixture.whenStable();

    expect(countSpy).not.toHaveBeenCalled();
  });

  it("should fetch the notification count once permission is present", async () => {
    const countSpy = vi
      .spyOn(notificationsService, "getNotificationCount")
      .mockReturnValue(of(3) as any);

    logIn();
    grantNotificationsRead();
    TestBed.flushEffects();
    await fixture.whenStable();

    expect(countSpy).toHaveBeenCalledTimes(1);
    expect(component.notificationCount()).toBe(3);
  });

  describe("notifications panel", () => {
    let overlayContainer: OverlayContainer;

    const bellButton = (): HTMLButtonElement =>
      fixture.nativeElement.querySelector('[data-testid="header-notifications-button"] button');

    const panel = (): HTMLElement | null =>
      overlayContainer
        .getContainerElement()
        .querySelector('[data-testid="header-notifications-panel"]');

    const clickBell = async () => {
      bellButton().click();
      await fixture.whenStable();
    };

    // The overlay's keyboard dispatcher reads the legacy keyCode, which jsdom
    // leaves at 0 unless it is defined explicitly.
    const pressEscape = async (target: HTMLElement) => {
      const event = new KeyboardEvent("keydown", { key: "Escape", bubbles: true });
      Object.defineProperty(event, "keyCode", { get: () => 27 });
      target.dispatchEvent(event);
      await fixture.whenStable();
    };

    beforeEach(async () => {
      overlayContainer = TestBed.inject(OverlayContainer);
      vi.spyOn(notificationsService, "getNotificationCount").mockReturnValue(of(0) as any);
      vi.spyOn(notificationsService, "getNotificationsForuser").mockReturnValue(of([]) as any);

      logIn();
      grantNotificationsRead();
      await fixture.whenStable();
    });

    it("opens from the bell, and a second click closes it", async () => {
      expect(panel()).toBeNull();

      await clickBell();

      expect(panel()).not.toBeNull();
      expect(panel()!.textContent).toContain("All caught up!");
      expect(bellButton().getAttribute("aria-expanded")).toBe("true");

      await clickBell();

      expect(panel()).toBeNull();
      expect(bellButton().getAttribute("aria-expanded")).toBe("false");
    });

    it("moves focus into the panel when it opens", async () => {
      bellButton().focus();

      await clickBell();

      expect(document.activeElement).toBe(panel());
    });

    it("closes on Escape and returns focus to the bell", async () => {
      bellButton().focus();
      await clickBell();

      await pressEscape(panel()!);

      expect(panel()).toBeNull();
      expect(component.notificationsOpen()).toBe(false);
      expect(document.activeElement).toBe(bellButton());
    });

    describe("Tab", () => {
      const pressTab = async (target: HTMLElement): Promise<KeyboardEvent> => {
        const event = new KeyboardEvent("keydown", { key: "Tab", bubbles: true, cancelable: true });
        target.dispatchEvent(event);
        await fixture.whenStable();
        return event;
      };

      // jsdom lays nothing out, so the checker finds nothing visible and hence
      // nothing tabbable; a spec that needs tabbable content says so.
      const treatEverythingAsTabbable = () => {
        const checker = TestBed.inject(InteractivityChecker);
        vi.spyOn(checker, "isFocusable").mockReturnValue(true);
        vi.spyOn(checker, "isTabbable").mockReturnValue(true);
      };

      it("stays on the panel when nothing in it is tabbable", async () => {
        await clickBell();

        const event = await pressTab(panel()!);

        expect(event.defaultPrevented).toBe(true);
        expect(document.activeElement).toBe(panel());
      });

      it("brings focus back to the panel once it has fallen to the body", async () => {
        treatEverythingAsTabbable();
        await clickBell();
        panel()!.blur();
        expect(document.activeElement).toBe(document.body);

        const event = await pressTab(document.body);

        expect(event.defaultPrevented).toBe(true);
        expect(document.activeElement).toBe(panel());
      });

      it("is left to the focus trap while focus is inside and there is somewhere to go", async () => {
        treatEverythingAsTabbable();
        await clickBell();

        const event = await pressTab(panel()!);

        expect(event.defaultPrevented).toBe(false);
      });
    });

    // Whether the click also reaches what is under it needs real hit-testing,
    // which jsdom does not do: e2e/notifications-popover.spec.ts covers that.
    it("closes when something else in the header is clicked", async () => {
      const dispatch = vi.spyOn(store, "dispatch");
      await clickBell();

      fixture.nativeElement.querySelector('[data-testid="sidebar-toggle"]').click();
      await fixture.whenStable();

      expect(panel()).toBeNull();
      expect(dispatch).toHaveBeenCalledWith(new ToggleIsSidebarOpen());
    });
  });
});
