import { Component, computed, effect, signal, untracked } from "@angular/core";
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
import { NgbPopover } from "@ng-bootstrap/ng-bootstrap";
import { MatProgressBar } from "@angular/material/progress-bar";
import { NotificationsListComponent } from "../../notifications/notifications-list/notifications-list.component";

@Component({
    selector: "app-header",
    templateUrl: "./header.component.html",
    styleUrls: ["./header.component.scss"],
    imports: [MatIconButton, MatTooltip, MatIcon, HasGroupPermissionDirective, RouterLinkActive, RouterLink, HasAppPermissionDirective, SearchbarComponent, ButtonComponent, NgbPopover, MatProgressBar, NotificationsListComponent]
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

  protected readonly Permission = Permission;

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
}
