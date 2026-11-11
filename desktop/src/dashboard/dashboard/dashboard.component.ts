import { Component, OnInit, signal } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { tap } from "rxjs";
import { DEFAULT_HOST_CLASS } from "src/constants";
import { DashboardState } from "src/store/dashboard.state";
import { Dashboard, WidgetType } from "../../open-api";
import { GroupState } from "../../store";
import { SummaryCardComponent } from "../../shared-ui/summary-card/summary-card.component";
import { FilteredReceiptsComponent } from "../filtered-receipts/filtered-receipts.component";
import { ActivityComponent } from "../activity/activity.component";
import { PieChartComponent } from "../pie-chart/pie-chart.component";
import { ReportWidgetComponent } from "../report-widget/report-widget.component";

@UntilDestroy()
@Component({
    selector: "app-dashboard",
    templateUrl: "./dashboard.component.html",
    styleUrls: ["./dashboard.component.scss"],
    host: DEFAULT_HOST_CLASS,
    imports: [SummaryCardComponent, FilteredReceiptsComponent, ActivityComponent, PieChartComponent, ReportWidgetComponent]
})
export class DashboardComponent implements OnInit {
  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  public dashboards = signal<Dashboard[]>([]);

  public selectedDashboard = signal<Dashboard | undefined>(undefined);

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
          this.dashboards.set(this.store.selectSnapshot(
            DashboardState.getDashboardsByGroupId(groupId)
          ));
          this.setSelectedDashboard();
        })
      )
      .subscribe();
  }

  private setSelectedDashboard(): void {
    const selectedDashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );
    this.selectedDashboard.set(this.dashboards().find(
      (dashboard) => dashboard?.id?.toString() === selectedDashboardId
    ));
  }
}
