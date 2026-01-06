import { CommonModule } from "@angular/common";
import { Component, Input, OnInit, OnChanges, SimpleChanges } from "@angular/core";
import { ChartConfiguration, ChartData } from "chart.js";
import { NgChartsModule } from "ng2-charts";
import { take, tap } from "rxjs";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { ChartGrouping, PieChartData, PieChartDataCommand, Widget, WidgetService } from "../../open-api";

@Component({
  selector: "app-pie-chart",
  templateUrl: "./pie-chart.component.html",
  styleUrls: ["./pie-chart.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule, NgChartsModule]
})
export class PieChartComponent implements OnInit, OnChanges {
  @Input() public widget!: Widget;
  @Input() public groupId?: number;

  public pieChartData: ChartData<"pie", number[], string> = {
    labels: [],
    datasets: [
      {
        data: [],
        backgroundColor: [
          "#FF6384",
          "#36A2EB",
          "#FFCE56",
          "#4BC0C0",
          "#9966FF",
          "#FF9F40",
          "#E7E9ED",
          "#7C4DFF",
          "#FF5252",
          "#64FFDA",
          "#FFD740",
          "#448AFF",
        ],
      },
    ],
  };

  public pieChartOptions: ChartConfiguration<"pie">["options"] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: "bottom",
      },
      tooltip: {
        callbacks: {
          label: (context) => {
            const label = context.label || "";
            const value = context.parsed || 0;
            return `${label}: $${value.toFixed(2)}`;
          },
        },
      },
    },
  };

  public isLoading: boolean = true;
  public hasData: boolean = false;

  constructor(private widgetService: WidgetService) {}

  public ngOnInit(): void {
    this.loadData();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["groupId"] && !changes["groupId"].firstChange) {
      this.loadData();
    }
  }

  private loadData(): void {
    if (!this.groupId || !this.widget?.configuration) {
      this.isLoading = false;
      return;
    }

    const config = this.widget.configuration as { chartGrouping?: ChartGrouping; filter?: any };
    if (!config.chartGrouping) {
      this.isLoading = false;
      return;
    }

    const command: PieChartDataCommand = {
      chartGrouping: config.chartGrouping,
      filter: config.filter,
    };

    this.isLoading = true;
    this.widgetService
      .getPieChartData(this.groupId, command)
      .pipe(
        take(1),
        tap((response: PieChartData) => {
          this.updateChartData(response);
          this.isLoading = false;
        })
      )
      .subscribe();
  }

  private updateChartData(data: PieChartData): void {
    if (!data.data || data.data.length === 0) {
      this.hasData = false;
      return;
    }

    this.hasData = true;
    const labels = data.data.map((point) => point.label || "Unknown");
    const values = data.data.map((point) => point.value || 0);

    this.pieChartData = {
      labels: labels,
      datasets: [
        {
          data: values,
          backgroundColor: this.pieChartData.datasets[0].backgroundColor,
        },
      ],
    };
  }

  public getChartGroupingLabel(): string {
    const config = this.widget?.configuration as { chartGrouping?: ChartGrouping };
    switch (config?.chartGrouping) {
      case ChartGrouping.Categories:
        return "Categories";
      case ChartGrouping.Tags:
        return "Tags";
      case ChartGrouping.Paidby:
        return "Paid By";
      default:
        return "Unknown";
    }
  }
}
