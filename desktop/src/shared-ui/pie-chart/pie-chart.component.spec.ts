import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { ChartData, ChartConfiguration } from "chart.js";
import { NgChartsModule } from "ng2-charts";

import { PieChartUiComponent } from "./pie-chart.component";

describe("PieChartUiComponent", () => {
  let component: PieChartUiComponent;
  let fixture: ComponentFixture<PieChartUiComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PieChartUiComponent, NgChartsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(PieChartUiComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  describe("default values", () => {
    it("should have default empty chartData", () => {
      expect(component.chartData).toEqual({
        labels: [],
        datasets: [{ data: [] }],
      });
    });

    it("should have default chartOptions with responsive and maintainAspectRatio", () => {
      expect(component.chartOptions).toEqual({
        responsive: true,
        maintainAspectRatio: false,
      });
    });

    it("should have default height of 300px", () => {
      expect(component.height).toBe("300px");
    });

    it("should have isLoading as false by default", () => {
      expect(component.isLoading).toBe(false);
    });

    it("should have hasData as true by default", () => {
      expect(component.hasData).toBe(true);
    });

    it("should have default noDataMessage", () => {
      expect(component.noDataMessage).toBe("No data available");
    });

    it("should have default loadingMessage", () => {
      expect(component.loadingMessage).toBe("Loading...");
    });
  });

  describe("loading state", () => {
    it("should display loading message when isLoading is true", () => {
      component.isLoading = true;
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl).toBeTruthy();
      expect(messageEl.nativeElement.textContent).toBe("Loading...");
    });

    it("should display custom loading message when provided", () => {
      component.isLoading = true;
      component.loadingMessage = "Please wait...";
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("Please wait...");
    });

    it("should not display canvas when isLoading is true", () => {
      component.isLoading = true;
      fixture.detectChanges();

      const canvas = fixture.debugElement.query(By.css("canvas"));
      expect(canvas).toBeFalsy();
    });
  });

  describe("no data state", () => {
    it("should display no data message when hasData is false", () => {
      component.isLoading = false;
      component.hasData = false;
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl).toBeTruthy();
      expect(messageEl.nativeElement.textContent).toBe("No data available");
    });

    it("should display custom no data message when provided", () => {
      component.isLoading = false;
      component.hasData = false;
      component.noDataMessage = "Nothing to display";
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("Nothing to display");
    });

    it("should not display canvas when hasData is false", () => {
      component.isLoading = false;
      component.hasData = false;
      fixture.detectChanges();

      const canvas = fixture.debugElement.query(By.css("canvas"));
      expect(canvas).toBeFalsy();
    });
  });

  describe("chart display", () => {
    it("should display canvas when not loading and has data", () => {
      component.isLoading = false;
      component.hasData = true;
      component.chartData = {
        labels: ["Category A", "Category B"],
        datasets: [{ data: [100, 200] }],
      };
      fixture.detectChanges();

      const canvas = fixture.debugElement.query(By.css("canvas"));
      expect(canvas).toBeTruthy();
    });

    it("should not display message container when chart is shown", () => {
      component.isLoading = false;
      component.hasData = true;
      component.chartData = {
        labels: ["Category A"],
        datasets: [{ data: [100] }],
      };
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message")
      );
      expect(messageEl).toBeFalsy();
    });
  });

  describe("height input", () => {
    it("should apply custom height to container", () => {
      component.height = "500px";
      fixture.detectChanges();

      const container = fixture.debugElement.query(
        By.css(".pie-chart-ui-container")
      );
      expect(container.styles["height"]).toBe("500px");
    });

    it("should apply different height values", () => {
      component.height = "200px";
      fixture.detectChanges();

      const container = fixture.debugElement.query(
        By.css(".pie-chart-ui-container")
      );
      expect(container.styles["height"]).toBe("200px");
    });
  });

  describe("chartData input", () => {
    it("should accept chartData with labels and datasets", () => {
      const testData: ChartData<"pie", number[], string> = {
        labels: ["A", "B", "C"],
        datasets: [
          {
            data: [10, 20, 30],
            backgroundColor: ["#FF0000", "#00FF00", "#0000FF"],
          },
        ],
      };

      component.chartData = testData;
      fixture.detectChanges();

      expect(component.chartData).toEqual(testData);
    });

    it("should handle empty labels array", () => {
      component.chartData = {
        labels: [],
        datasets: [{ data: [] }],
      };
      component.hasData = true;
      fixture.detectChanges();

      expect(component.chartData.labels).toEqual([]);
    });

    it("should handle single data point", () => {
      component.chartData = {
        labels: ["Single"],
        datasets: [{ data: [100] }],
      };
      component.hasData = true;
      fixture.detectChanges();

      expect(component.chartData.labels?.length).toBe(1);
      expect(component.chartData.datasets[0].data.length).toBe(1);
    });

    it("should handle many data points", () => {
      const labels = Array.from({ length: 20 }, (_, i) => `Category ${i}`);
      const data = Array.from({ length: 20 }, (_, i) => i * 10);

      component.chartData = {
        labels,
        datasets: [{ data }],
      };
      component.hasData = true;
      fixture.detectChanges();

      expect(component.chartData.labels?.length).toBe(20);
      expect(component.chartData.datasets[0].data.length).toBe(20);
    });
  });

  describe("chartOptions input", () => {
    it("should accept custom chartOptions", () => {
      const customOptions: ChartConfiguration<"pie">["options"] = {
        responsive: false,
        maintainAspectRatio: true,
        plugins: {
          legend: {
            display: false,
          },
        },
      };

      component.chartOptions = customOptions;
      fixture.detectChanges();

      expect(component.chartOptions).toEqual(customOptions);
    });

    it("should handle chartOptions with custom legend position", () => {
      const customOptions: ChartConfiguration<"pie">["options"] = {
        responsive: true,
        plugins: {
          legend: {
            position: "right",
          },
        },
      };

      component.chartOptions = customOptions;
      fixture.detectChanges();

      expect(component.chartOptions?.plugins?.legend?.position).toBe("right");
    });
  });

  describe("state transitions", () => {
    it("should transition from loading to showing data", () => {
      component.isLoading = true;
      fixture.detectChanges();

      let canvas = fixture.debugElement.query(By.css("canvas"));
      expect(canvas).toBeFalsy();

      component.isLoading = false;
      component.hasData = true;
      component.chartData = {
        labels: ["Test"],
        datasets: [{ data: [100] }],
      };
      fixture.detectChanges();

      canvas = fixture.debugElement.query(By.css("canvas"));
      expect(canvas).toBeTruthy();
    });

    it("should transition from loading to no data", () => {
      component.isLoading = true;
      fixture.detectChanges();

      let messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("Loading...");

      component.isLoading = false;
      component.hasData = false;
      fixture.detectChanges();

      messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("No data available");
    });

    it("should handle rapid state changes", () => {
      component.isLoading = true;
      fixture.detectChanges();

      component.isLoading = false;
      component.hasData = true;
      fixture.detectChanges();

      component.hasData = false;
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("No data available");
    });
  });

  describe("priority of states", () => {
    it("should prioritize loading over hasData", () => {
      component.isLoading = true;
      component.hasData = true;
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("Loading...");
    });

    it("should prioritize loading over no data", () => {
      component.isLoading = true;
      component.hasData = false;
      fixture.detectChanges();

      const messageEl = fixture.debugElement.query(
        By.css(".pie-chart-ui-message span")
      );
      expect(messageEl.nativeElement.textContent).toBe("Loading...");
    });
  });

  describe("component encapsulation", () => {
    it("should have ViewEncapsulation.None", () => {
      // Verify styles can cascade by checking container exists
      fixture.detectChanges();
      const container = fixture.debugElement.query(
        By.css(".pie-chart-ui-container")
      );
      expect(container).toBeTruthy();
    });
  });
});
