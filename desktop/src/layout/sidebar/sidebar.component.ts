import { Component, computed, ViewEncapsulation } from "@angular/core";
import { toSignal } from "@angular/core/rxjs-interop";
import { MatDialog } from "@angular/material/dialog";
import { Data, NavigationEnd, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { filter, map, startWith, switchMap, take } from "rxjs";
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

  // True when the active route opts into a full-height, no-padding content frame
  // (route data `fullHeight`), so the routed page owns its own internal scrolling
  // instead of the shell's block-flow + p-4 padding. Used by the Report Builder.
  public readonly isContentFullHeight = toSignal(
    this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.deepestRouteData()["fullHeight"] === true),
    ),
    { initialValue: false },
  );

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

  // Reports are reachable with read OR the readAll bypass (the *hasAppPermission
  // directive is single-key AND-only, so the OR is resolved through the selector).
  protected readonly canViewReports = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppReportsRead,
      Permission.AppReportsReadAll,
    ])
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

  // Walks from the router state root to the deepest activated child so the
  // fullHeight flag is picked up wherever it is declared in the route tree.
  private deepestRouteData(): Data {
    let route = this.router.routerState.snapshot.root;
    while (route.firstChild) {
      route = route.firstChild;
    }
    return route.data;
  }

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
