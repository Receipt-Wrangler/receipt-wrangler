import { CommonModule } from "@angular/common";
import { Component, Input, ViewEncapsulation } from "@angular/core";
import { ChartConfiguration, ChartData } from "chart.js";
import { NgChartsModule } from "ng2-charts";

@Component({
  selector: "app-pie-chart-ui",
  templateUrl: "./pie-chart.component.html",
  styleUrls: ["./pie-chart.component.scss"],
  standalone: true,
  imports: [CommonModule, NgChartsModule],
  encapsulation: ViewEncapsulation.None,
})
export class PieChartUiComponent {
  @Input() public chartData: ChartData<"pie", number[], string> = {
    labels: [],
    datasets: [{ data: [] }],
  };

  @Input() public chartOptions: ChartConfiguration<"pie">["options"] = {
    responsive: true,
    maintainAspectRatio: false,
  };

  @Input() public height: string = "300px";

  @Input() public isLoading: boolean = false;

  @Input() public hasData: boolean = true;

  @Input() public noDataMessage: string = "No data available";

  @Input() public loadingMessage: string = "Loading...";
}
