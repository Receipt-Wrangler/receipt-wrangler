import { CommonModule, CurrencyPipe } from "@angular/common";
import { Component, Input, OnInit, OnChanges, SimpleChanges } from "@angular/core";
import { Chart, ChartConfiguration, ChartData } from "chart.js";
import ChartDataLabels from "chartjs-plugin-datalabels";
import { take, tap } from "rxjs";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { ChartGrouping, PieChartData, PieChartDataCommand, Widget, WidgetService } from "../../open-api";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";

// Register the datalabels plugin
Chart.register(ChartDataLabels);

@Component({
  selector: "app-pie-chart",
  templateUrl: "./pie-chart.component.html",
  styleUrls: ["./pie-chart.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule],
  providers: [CustomCurrencyPipe, CurrencyPipe]
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

  public pieChartOptions!: ChartConfiguration<"pie">["options"];

  public isLoading: boolean = true;
  public hasData: boolean = false;

  constructor(
    private widgetService: WidgetService,
    private customCurrencyPipe: CustomCurrencyPipe
  ) {}

  private initializeChartOptions(): void {
    this.pieChartOptions = {
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
              const total = context.dataset.data.reduce((a: number, b: number) => a + b, 0);
              const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : "0";
              const formattedValue = this.customCurrencyPipe.transform(value);
              return `${label}: ${formattedValue} (${percentage}%)`;
            },
          },
        },
        datalabels: {
          color: "#fff",
          font: {
            weight: "bold",
            size: 12,
          },
          formatter: (value: number, context: any) => {
            const total = context.dataset.data.reduce((a: number, b: number) => a + b, 0);
            const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : "0";
            // Only show label if percentage is > 5% to avoid cluttering small slices
            return parseFloat(percentage) > 5 ? `${percentage}%` : "";
          },
        },
      },
    };
  }

  public ngOnInit(): void {
    this.initializeChartOptions();
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

    // Sort data alphabetically by label for consistent colors and orientation
    const sortedData = [...data.data].sort((a, b) => {
      const labelA = (a.label || "Unknown").toLowerCase();
      const labelB = (b.label || "Unknown").toLowerCase();
      return labelA.localeCompare(labelB);
    });

    const labels = sortedData.map((point) => point.label || "Unknown");
    const values = sortedData.map((point) => point.value || 0);

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
