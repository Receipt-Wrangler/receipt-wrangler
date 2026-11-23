import { A11yModule, InteractivityChecker } from "@angular/cdk/a11y";
import { hasModifierKey } from "@angular/cdk/keycodes";
import { ConnectedPosition, OverlayModule } from "@angular/cdk/overlay";
import { DOCUMENT } from "@angular/common";
import {
  afterNextRender,
  Component,
  computed,
  effect,
  ElementRef,
  inject,
  Injector,
  signal,
  untracked,
  viewChild,
} from "@angular/core";
import { Router, RouterLinkActive, RouterLink } from "@angular/router";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { LayoutState } from "src/store/layout.state";
import { ToggleIsSidebarOpen } from "src/store/layout.state.actions";
import { AuthService, NotificationsService, Permission } from "../../open-api";
import { AuthState, GroupState } from "../../store";
import { MatIconButton } from "@angular/material/button";
import { MatTooltip } from "@angular/material/tooltip";
import { MatIcon } from "@angular/material/icon";
import { HasGroupPermissionDirective } from "../../directives/has-group-permission.directive";
import { HasAppPermissionDirective } from "../../directives/has-app-permission.directive";
import { SearchbarComponent } from "../../searchbar/searchbar/searchbar.component";
import { ButtonComponent } from "../../button/button/button.component";
import { MatProgressBar } from "@angular/material/progress-bar";
import { NotificationsListComponent } from "../../notifications/notifications-list/notifications-list.component";

@Component({
    selector: "app-header",
    templateUrl: "./header.component.html",
    styleUrls: ["./header.component.scss"],
    imports: [A11yModule, MatIconButton, MatTooltip, MatIcon, HasGroupPermissionDirective, RouterLinkActive, RouterLink, HasAppPermissionDirective, SearchbarComponent, ButtonComponent, MatProgressBar, NotificationsListComponent, OverlayModule]
})
export class HeaderComponent {
  public isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);

  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  public selectedGroupIdNumber = computed(() => Number(this.selectedGroupId()));

  public loggedInUser = this.store.selectSignal(AuthState.loggedInUser);

  public showProgressBar = this.store.selectSignal(LayoutState.showProgressBar);

  public userPreferences = this.store.selectSignal(AuthState.userPreferences);

  public receiptHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.receiptListLink)];
  });

  public dashboardHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.dashboardLink)];
  });

  public settingsBaseHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.settingsLinkBase) + "/view"];
  });

  public groupName = computed(() => {
    const groupId = this.selectedGroupId();
    const group = this.store.selectSnapshot(GroupState.getGroupById(groupId));
    return group?.name as string ?? "";
  });

  public notificationCount = signal<number | undefined>(undefined);

  public readonly notificationsOpen = signal(false);

  /**
   * Under the bell, end-aligned to it, 8px down: where the ng-bootstrap popover
   * this replaced sat. No flipped fallback — the bell lives in the sticky top
   * header, and the panel's arrow only points up.
   */
  public readonly notificationsPositions: ConnectedPosition[] = [
    { originX: "end", originY: "bottom", overlayX: "end", overlayY: "top", offsetY: 8 },
  ];

  protected readonly Permission = Permission;

  private readonly notificationsPanel = viewChild<ElementRef<HTMLElement>>("notificationsPanel");

  /** What had focus when the panel opened (the bell), for Escape to return to. */
  private focusBeforeNotifications: HTMLElement | null = null;

  private readonly document = inject(DOCUMENT);

  private readonly injector = inject(Injector);

  private readonly interactivityChecker = inject(InteractivityChecker);

  private readonly canReadNotifications = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppNotificationsRead)
  );

  constructor(
    private authService: AuthService,
    private notificationsService: NotificationsService,
    private router: Router,
    private store: Store
  ) {
    this.listenForLoggedInUser();
  }

  private listenForLoggedInUser(): void {
    let fetched = false;
    effect(() => {
      const loggedIn = this.isLoggedIn();
      const canRead = this.canReadNotifications();
      if (!loggedIn) {
        fetched = false;
        return;
      }
      if (canRead && !fetched) {
        fetched = true;
        untracked(() => {
          this.notificationsService.getNotificationCount().pipe(
            take(1),
            tap((n) => {
              this.notificationCount.set(n > 0 ? n : undefined);
            })
          ).subscribe();
        });
      }
    });
  }

  public toggleSidebar(): void {
    this.store.dispatch(new ToggleIsSidebarOpen());
  }

  /**
   * The overlay renders at the end of <body>, not beside the bell, so Tab from
   * the bell would never reach it: opening moves focus into the panel, and
   * `cdkTrapFocus` keeps Tab inside. Deliberately not `cdkTrapFocusAutoCapture`,
   * which restores focus on every close — after a click elsewhere in the header
   * that would pull focus back off the control the user just clicked.
   */
  public toggleNotifications(): void {
    if (this.notificationsOpen()) {
      this.closeNotifications();
      return;
    }

    this.focusBeforeNotifications = this.document.activeElement as HTMLElement | null;
    this.notificationsOpen.set(true);
    afterNextRender(() => this.notificationsPanel()?.nativeElement.focus(), {
      injector: this.injector,
    });
  }

  public closeNotifications(): void {
    this.notificationsOpen.set(false);
  }

  /**
   * Escape closes through the overlay's own handler; this only returns focus.
   *
   * Tab is normally `cdkTrapFocus`'s job, but its anchors can only hand focus to
   * something tabbable inside the panel, and the CDK we pin has no fallback when
   * there is nothing (an empty list: "Mark all as read" is disabled), so Tab
   * walked out of the open panel. Deleting a notification also drops focus to
   * <body>. In both cases keep focus on the panel itself. The overlay receives
   * every keydown while it is open, wherever focus is.
   */
  public onNotificationsKeydown(event: KeyboardEvent): void {
    if (event.key === "Escape" && !hasModifierKey(event)) {
      this.focusBeforeNotifications?.focus();
      return;
    }

    const panel = this.notificationsPanel()?.nativeElement;
    if (event.key !== "Tab" || !panel) {
      return;
    }

    const focusIsInside = panel.contains(this.document.activeElement);
    if (!focusIsInside || !this.hasTabbableContent(panel)) {
      event.preventDefault();
      panel.focus();
    }
  }

  /** The same test `cdkTrapFocus` uses to find somewhere to send Tab. */
  private hasTabbableContent(panel: HTMLElement): boolean {
    return Array.from(panel.querySelectorAll<HTMLElement>("*")).some(
      (element) =>
        this.interactivityChecker.isFocusable(element) &&
        this.interactivityChecker.isTabbable(element)
    );
  }
}
