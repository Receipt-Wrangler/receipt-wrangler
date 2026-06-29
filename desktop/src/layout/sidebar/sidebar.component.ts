import { Component, computed, ViewEncapsulation } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { switchMap, take } from "rxjs";
import { LayoutState } from "src/store/layout.state";
import { SetPage } from "src/store/receipt-table.actions";
import { AboutComponent } from "../../about/about/about.component";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { ImportFormComponent } from "../../import/import-form/import-form.component";
import { AuthService, Group, GroupStatus, Permission } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState, GroupState, Logout, SetSelectedGroupId } from "../../store";
import { hasAll } from "../../utils/permission.utils";

@Component({
    selector: "app-sidebar",
    templateUrl: "./sidebar.component.html",
    styleUrls: ["./sidebar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class SidebarComponent {
  constructor(
    private authService: AuthService,
    private matDialog: MatDialog,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {}

  public loggedInUser = this.store.selectSignal(AuthState.loggedInUser);

  public isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);

  public isSidebarOpen = this.store.selectSignal(LayoutState.isSidebarOpen);

  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  protected readonly selectedGroupIdNumber = computed(() => {
    const id = Number.parseInt(this.selectedGroupId() ?? "");
    return Number.isNaN(id) ? 0 : id;
  });

  private allGroups = this.store.selectSignal(GroupState.groups);

  public groups = computed(() =>
    this.allGroups().filter((g: Group) => g.status !== GroupStatus.Archived)
  );

  protected readonly Permission = Permission;

  private readonly appPermissions = this.store.selectSignal(
    AuthState.appPermissions
  );

  private readonly groupPermissions = this.store.selectSignal(
    AuthState.groupPermissions
  );

  // Mirrors the three speed-dial sub-button gates (Add Receipt, Quick Scan, Add
  // Group), so the plus FAB is shown iff at least one sub-button would render.
  protected readonly canAddAnything = computed(() => {
    const groupPerms =
      this.groupPermissions()?.[this.selectedGroupIdNumber()] ?? [];
    return (
      hasAll(this.appPermissions(), Permission.AppGroupsCreate) ||
      hasAll(groupPerms, Permission.GroupReceiptsCreate) ||
      hasAll(groupPerms, Permission.GroupReceiptsQuickScan)
    );
  });

  protected readonly canViewSystemSettings = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppSystemSettingsRead,
      Permission.AppPromptsRead,
      Permission.AppReceiptProcessingSettingsRead,
      Permission.AppSystemEmailsRead,
      Permission.AppSystemTasksRead,
    ])
  );

  protected readonly canViewUserSettings = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppAccountRead,
      Permission.AppUserPreferencesRead,
      Permission.AppApiKeysRead,
    ])
  );

  public addButtonExpanded: boolean | null = null;

  public groupClicked(groupId: number): void {
    this.store.dispatch(new SetSelectedGroupId(groupId.toString()));
    this.store.dispatch(new SetPage(1));
    const dashboardLink = this.store.selectSnapshot(GroupState.dashboardLink);
    this.router.navigate([dashboardLink]);
  }

  public toggleAddButtonExpanded(): void {
    this.addButtonExpanded = !this.addButtonExpanded;
  }

  public logout(): void {
    this.authService
      .logout()
      .pipe(
        take(1),
        switchMap(() => this.store.dispatch(new Logout())),
        switchMap(() => this.router.navigate(["/"])),
      )
      .subscribe();
  }

  public openImportDialog(): void {
    this.matDialog.open(ImportFormComponent, DEFAULT_DIALOG_CONFIG);
  }

  public openAboutDialog(): void {
    this.matDialog.open(AboutComponent, DEFAULT_DIALOG_CONFIG);
  }
}
