import { CommonModule, CurrencyPipe } from "@angular/common";
import { Component, OnInit, OnChanges, SimpleChanges, input, signal } from "@angular/core";
import { Chart, ChartConfiguration, ChartData } from "chart.js";
import ChartDataLabels from "chartjs-plugin-datalabels";
import { take, tap } from "rxjs";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { PipesModule } from "../../pipes/pipes.module";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { ChartGrouping, PieChartData, PieChartDataCommand, Widget, WidgetService } from "../../open-api";
import { RECEIPT_DATE_FILTER_KEYS } from "../../constants/receipt-filter-fields.constant";
import { toDateWireFilter } from "../../utils/date-wire-filter";

// Register the datalabels plugin
Chart.register(ChartDataLabels);

@Component({
  selector: "app-pie-chart",
  templateUrl: "./pie-chart.component.html",
  styleUrls: ["./pie-chart.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule, PipesModule],
  providers: [CurrencyPipe, CustomCurrencyPipe],
})
export class PieChartComponent implements OnInit, OnChanges {
  public readonly widget = input.required<Widget>();
  public readonly groupId = input<number>();

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

  public pieChartOptions: ChartConfiguration<"pie">["options"];

  public isLoading = signal(true);
  public hasData = signal(false);

  constructor(
    private widgetService: WidgetService,
    private customCurrencyPipe: CustomCurrencyPipe,
  ) {
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
              const total = context.dataset.data.reduce((a: number, b: number) => a + Math.abs(b), 0);
              const percentage = total > 0 ? ((Math.abs(value) / total) * 100).toFixed(1) : "0";
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
            const total = context.dataset.data.reduce((a: number, b: number) => a + Math.abs(b), 0);
            const percentage = total > 0 ? ((Math.abs(value) / total) * 100).toFixed(1) : "0";
            // Only show label if percentage is > 5% to avoid cluttering small slices
            return parseFloat(percentage) > 5 ? `${percentage}%` : "";
          },
        },
      },
    };
  }

  public ngOnInit(): void {
    this.loadData();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["groupId"] && !changes["groupId"].firstChange) {
      this.loadData();
    }
  }

  private loadData(): void {
    const groupId = this.groupId();
    const widget = this.widget();
    if (!groupId || !widget?.configuration) {
      this.isLoading.set(false);
      return;
    }

    const config = widget.configuration as { chartGrouping?: ChartGrouping; filter?: any };
    if (!config.chartGrouping) {
      this.isLoading.set(false);
      return;
    }

    const command: PieChartDataCommand = {
      chartGrouping: config.chartGrouping,
      // The picked calendar day, not a browser-local instant: see toDateWireFilter.
      filter: toDateWireFilter(config.filter, RECEIPT_DATE_FILTER_KEYS),
    };

    this.isLoading.set(true);
    this.widgetService
      .getPieChartData(groupId, command)
      .pipe(
        take(1),
        tap((response: PieChartData) => {
          this.updateChartData(response);
          this.isLoading.set(false);
        })
      )
      .subscribe();
  }

  private updateChartData(data: PieChartData): void {
    if (!data.data || data.data.length === 0) {
      this.hasData.set(false);
      return;
    }

    this.hasData.set(true);
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
    const config = this.widget()?.configuration as { chartGrouping?: ChartGrouping };
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
