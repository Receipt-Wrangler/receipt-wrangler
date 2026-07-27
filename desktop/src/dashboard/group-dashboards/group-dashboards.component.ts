import { Component, computed, effect } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "src/constants";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { DashboardState } from "src/store/dashboard.state";
import { AddDashboardToGroup, DeleteDashboardFromGroup, UpdateDashBoardForGroup, } from "src/store/dashboard.state.actions";
import { Dashboard, DashboardService, Permission } from "../../open-api";
import { SnackbarService } from "../../services";
import { GroupState, SetSelectedDashboardId } from "../../store";
import { DashboardFormComponent } from "../dashboard-form/dashboard-form.component";

// Stable empty reference so the dashboards computed doesn't emit a fresh [] on
// every recompute (which would needlessly re-run the auto-select effect).
const EMPTY_DASHBOARDS: Dashboard[] = [];

@UntilDestroy()
@Component({
    selector: "app-group-dashboards",
    templateUrl: "./group-dashboards.component.html",
    styleUrls: ["./group-dashboards.component.scss"],
    standalone: false
})
export class GroupDashboardsComponent {
  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  // Dashboard CRUD buttons gate on the group-scoped permissions via
  // *hasGroupPermission, which needs a numeric group id (selectedGroupId is a
  // string). Falls back to 0 (no permissions for a non-existent group → hidden).
  protected readonly Permission = Permission;

  protected readonly selectedGroupIdNum = computed(
    () => +(this.selectedGroupId() ?? 0)
  );

  public selectedDashboardId = this.store.selectSignal(
    GroupState.selectedDashboardId
  );

  private dashboardsByGroup = this.store.selectSignal(DashboardState.dashboards);

  // Derived reactively from the store so the chip list re-renders whenever the
  // current group's dashboards land — including when the resolver's fetch
  // resolves after this (reused) component has already reacted to the group
  // switch. A snapshot read here would be untracked and reintroduce the
  // "dashboards don't load on switch" bug.
  public dashboards = computed<Dashboard[]>(
    () => this.dashboardsByGroup()[this.selectedGroupId()] ?? EMPTY_DASHBOARDS
  );

  constructor(
    private dashboardService: DashboardService,
    private matDialog: MatDialog,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store
  ) {
    // Once the current group's dashboards are available, ensure a dashboard is
    // selected and its outlet is showing. Depends only on dashboards();
    // selectedDashboardId is read as an untracked snapshot and the effect's only
    // write (SetSelectedDashboardId) changes neither dashboards() nor
    // selectedGroupId, so it cannot re-trigger itself.
    effect(() => {
      const dashboards = this.dashboards();
      const selectedDashboardId = this.store.selectSnapshot(
        GroupState.selectedDashboardId
      );

      if (selectedDashboardId) {
        this.navigateToDashboard(+selectedDashboardId);
      } else if (dashboards.length > 0) {
        this.setSelectedDashboardId(dashboards[0].id);
        this.navigateToDashboard(dashboards[0].id);
      }
    });
  }

  public navigateToDashboard(dashboardId: number): void {
    const selectedGroupId = this.store.selectSnapshot(
      GroupState.selectedGroupId
    );

    setTimeout(() => {
      this.router.navigateByUrl(
        `/dashboard/group/${selectedGroupId}/${dashboardId}`
      );
    }, 0);
  }

  public openDashboardDialog(isCreate?: boolean): void {
    const dialogRef = this.matDialog.open(
      DashboardFormComponent,
      { ...DEFAULT_DIALOG_CONFIG, width: "75%" }
    );
    const selectedDashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );

    if (!isCreate) {
      const dashboard = this.dashboards().find(
        (d) => d.id === +selectedDashboardId
      );

      dialogRef.componentInstance.dashboard = dashboard;
      dialogRef.componentInstance.headerText = `Edit Dashboard ${dashboard?.name}`;
    } else {
      dialogRef.componentInstance.headerText = "Add a Dashboard";
    }

    dialogRef
      .afterClosed()
      .pipe(
        untilDestroyed(this),
        tap((dashboard) => {
          const index = this.dashboards().findIndex(
            (d) => d.id === dashboard?.id
          );
          const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
          if (dashboard && index < 0) {
            this.store.dispatch(new AddDashboardToGroup(groupId, dashboard));
          } else if (dashboard && index > -1) {
            this.store.dispatch(
              new UpdateDashBoardForGroup(groupId, dashboard.id, dashboard)
            );
          }
        })
      )
      .subscribe();
  }

  public setSelectedDashboardId(dashboardId: number): void {
    this.store.dispatch(new SetSelectedDashboardId(dashboardId?.toString()));
  }

  public openDeleteConfirmationDialog(): void {
    const dialogRef = this.matDialog.open(ConfirmationDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      panelClass: "overflow-scroll",
    });
    const dashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );
    const selectedDashboard = this.dashboards().find(
      (d) => d.id.toString() === dashboardId
    );

    dialogRef.componentInstance.headerText = "Delete Dashboard";
    dialogRef.componentInstance.dialogContent = `Are you sure you want to delete dashboard "${selectedDashboard?.name}"? This action is irreversable.`;

    dialogRef
      .afterClosed()
      .pipe(
        untilDestroyed(this),
        tap((confirmed) => {
          if (confirmed) {
            this.dashboardService
              .deleteDashboard(+dashboardId)
              .pipe(
                take(1),
                tap(() => {
                  this.snackbarService.success(
                    "Successfully deleted dashboard"
                  );
                  const dashboardLink = this.store.selectSnapshot(
                    GroupState.dashboardLink
                  );
                  this.store.dispatch(
                    new DeleteDashboardFromGroup(
                      this.store.selectSnapshot(GroupState.selectedGroupId),
                      +dashboardId
                    )
                  );
                  this.store.dispatch(new SetSelectedDashboardId(undefined));
                  this.router.navigateByUrl(dashboardLink);
                })
              )
              .subscribe();
          }
        })
      )
      .subscribe();
  }
}
