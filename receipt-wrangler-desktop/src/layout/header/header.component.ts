import { Component, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Select, Store } from "@ngxs/store";
import { filter, Observable, switchMap, take, tap } from "rxjs";
import { LayoutState } from "src/store/layout.state";
import { ToggleIsSidebarOpen } from "src/store/layout.state.actions";
import { AuthService, GroupRole, NotificationsService, User, UserPreferences } from "../../open-api";
import { AuthState, GroupState } from "../../store";

@UntilDestroy()
@Component({
    selector: "app-header",
    templateUrl: "./header.component.html",
    styleUrls: ["./header.component.scss"],
    standalone: false
})
export class HeaderComponent implements OnInit {
  @Select(AuthState.isLoggedIn) public isLoggedIn!: Observable<boolean>;

  @Select(GroupState.selectedGroupId)
  public selectedGroupId!: Observable<string>;

  @Select(AuthState.loggedInUser) public loggedInUser!: Observable<User>;

  @Select(LayoutState.showProgressBar)
  public showProgressBar!: Observable<boolean>;

  @Select(AuthState.userPreferences)
  public userPreferences!: Observable<UserPreferences>;

  public receiptHeaderLink: string[] = [""];

  public dashboardHeaderLink: string[] = [""];

  public settingsBaseHeaderLink: string[] = [""];

  public groupName = "";

  public groupRoleEnum = GroupRole;

  public notificationCount: number | undefined = undefined;

  constructor(
    private authService: AuthService,
    private notificationsService: NotificationsService,
    private router: Router,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.setGroupData();
    this.listenForLoggedInUser();
  }

  private setGroupData(): void {
    this.selectedGroupId
      .pipe(
        tap((groupId) => {
          this.receiptHeaderLink = [
            this.store.selectSnapshot(GroupState.receiptListLink),
          ];
          this.dashboardHeaderLink = [
            this.store.selectSnapshot(GroupState.dashboardLink),
          ];
          this.settingsBaseHeaderLink = [
            this.store.selectSnapshot(GroupState.settingsLinkBase) + "/view",
          ];
          const newGroup = this.store.selectSnapshot(
            GroupState.getGroupById(groupId)
          );
          this.groupName = newGroup?.name as string;
        })
      )
      .subscribe();
  }

  private listenForLoggedInUser(): void {
    this.isLoggedIn
      .pipe(
        untilDestroyed(this),
        filter((loggedIn) => !!loggedIn),
        switchMap(() => this.updateNotificationCount())
      )
      .subscribe();
  }

  private updateNotificationCount(): Observable<number> {
    return this.notificationsService.getNotificationCount().pipe(
      take(1),
      tap((n) => {
        if (n > 0) {
          this.notificationCount = n;
        } else {
          this.notificationCount = undefined;
        }
      })
    );
  }

  public toggleSidebar(): void {
    this.store.dispatch(new ToggleIsSidebarOpen());
  }
}
