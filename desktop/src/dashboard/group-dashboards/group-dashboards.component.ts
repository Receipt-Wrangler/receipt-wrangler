import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Select, Store } from "@ngxs/store";
import { Observable, take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "src/constants";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { DashboardState } from "src/store/dashboard.state";
import { AddDashboardToGroup, DeleteDashboardFromGroup, UpdateDashBoardForGroup, } from "src/store/dashboard.state.actions";
import { Dashboard, DashboardService } from "../../open-api";
import { SnackbarService } from "../../services";
import { GroupState, SetSelectedDashboardId } from "../../store";
import { DashboardFormComponent } from "../dashboard-form/dashboard-form.component";

@UntilDestroy()
@Component({
    selector: "app-group-dashboards",
    templateUrl: "./group-dashboards.component.html",
    styleUrls: ["./group-dashboards.component.scss"],
    standalone: false
})
export class GroupDashboardsComponent implements OnInit {
  constructor(
    private dashboardService: DashboardService,
    private matDialog: MatDialog,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store
  ) {}

  @Select(GroupState.selectedGroupId)
  public selectedGroupId!: Observable<string>;

  @Select(GroupState.selectedDashboardId)
  public selectedDashboardId!: Observable<string>;

  public dashboards: Dashboard[] = [];

  public ngOnInit(): void {
    this.setDashboards();
  }

  private checkForSelectedDashboard(): void {
    const selectedDashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );

    if (selectedDashboardId) {
      this.navigateToDashboard(+selectedDashboardId);
      return;
    } else if (this.dashboards.length > 0) {
      this.setSelectedDashboardId(this.dashboards[0].id);
      this.navigateToDashboard(this.dashboards[0].id);
    }
  }

  private setDashboards(): void {
    this.store
      .select(GroupState.selectedGroupId)
      .pipe(
        untilDestroyed(this),
        tap((groupId) => {
          this.refreshDashboards(groupId);
          this.checkForSelectedDashboard();
        })
      )
      .subscribe();
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
      const dashboard = this.dashboards.find(
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
          const index = this.dashboards.findIndex(
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
          this.refreshDashboards(groupId);
        })
      )
      .subscribe();
  }

  public setSelectedDashboardId(dashboardId: number): void {
    this.store.dispatch(new SetSelectedDashboardId(dashboardId?.toString()));
  }

  private refreshDashboards(groupId: string): void {
    this.store
      .select(DashboardState.getDashboardsByGroupId(groupId))
      .pipe(
        take(1),
        tap((dashboards) => {
          this.dashboards = dashboards;
        })
      )
      .subscribe();
  }

  public openDeleteConfirmationDialog(): void {
    const dialogRef = this.matDialog.open(ConfirmationDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      panelClass: "overflow-scroll",
    });
    const dashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );
    const selectedDashboard = this.dashboards.find(
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
                  this.refreshDashboards(
                    this.store.selectSnapshot(GroupState.selectedGroupId)
                  );
                })
              )
              .subscribe();
          }
        })
      )
      .subscribe();
  }
}
