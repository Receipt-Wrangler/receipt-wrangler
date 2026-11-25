import { CommonModule } from "@angular/common";
import { Component, ViewEncapsulation, input } from "@angular/core";
import { ChartConfiguration, ChartData, Plugin } from "chart.js";
import { BaseChartDirective } from "ng2-charts";

@Component({
  selector: "app-pie-chart-ui",
  templateUrl: "./pie-chart.component.html",
  styleUrls: ["./pie-chart.component.scss"],
  standalone: true,
  imports: [CommonModule, BaseChartDirective],
  encapsulation: ViewEncapsulation.None,
})
export class PieChartUiComponent {
  public readonly chartData = input<ChartData<"pie", number[], string>>({
    labels: [],
    datasets: [{ data: [] }],
});

  public readonly chartOptions = input<ChartConfiguration<"pie">["options"]>({
    responsive: true,
    maintainAspectRatio: false,
});

  /** Plugins for this chart only, rather than registered with Chart.js for every chart. */
  public readonly plugins = input<Plugin<"pie">[]>([]);

  public readonly height = input<string>("300px");

  public readonly isLoading = input<boolean>(false);

  public readonly hasData = input<boolean>(true);

  public readonly noDataMessage = input<string>("No data available");

  public readonly loadingMessage = input<string>("Loading...");
}
