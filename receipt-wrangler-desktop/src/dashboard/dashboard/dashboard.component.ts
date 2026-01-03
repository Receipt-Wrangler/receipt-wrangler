import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Select, Store } from "@ngxs/store";
import { Observable, tap } from "rxjs";
import { DEFAULT_HOST_CLASS } from "src/constants";
import { DashboardState } from "src/store/dashboard.state";
import { Dashboard, WidgetType } from "../../open-api";
import { GroupState } from "../../store";

@UntilDestroy()
@Component({
    selector: "app-dashboard",
    templateUrl: "./dashboard.component.html",
    styleUrls: ["./dashboard.component.scss"],
    host: DEFAULT_HOST_CLASS,
    standalone: false
})
export class DashboardComponent implements OnInit {
  @Select(GroupState.selectedGroupId)
  public selectedGroupId!: Observable<string>;

  public dashboards: Dashboard[] = [];

  public selectedDashboard?: Dashboard;

  public widgetType = WidgetType;

  constructor(
    private activatedRoute: ActivatedRoute,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.listenForDashboardChanges();
    this.listenForParamChanges();
  }

  private listenForParamChanges(): void {
    this.activatedRoute.params
      .pipe(
        tap(() => {
          this.setSelectedDashboard();
        })
      )
      .subscribe();
  }

  private listenForDashboardChanges(): void {
    this.store
      .select(DashboardState.dashboards)
      .pipe(
        untilDestroyed(this),
        tap(() => {
          const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
          this.dashboards = this.store.selectSnapshot(
            DashboardState.getDashboardsByGroupId(groupId)
          );
          this.setSelectedDashboard();
        })
      )
      .subscribe();
  }

  private setSelectedDashboard(): void {
    const selectedDashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );
    this.selectedDashboard = this.dashboards.find(
      (dashboard) => dashboard?.id?.toString() === selectedDashboardId
    );
  }
}
