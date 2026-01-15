import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, SimpleChange } from "@angular/core";
import { ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { of, throwError } from "rxjs";
import { ChartGrouping, PieChartData, Widget, WidgetService, WidgetType } from "../../open-api";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

import { PieChartComponent } from "./pie-chart.component";

describe("PieChartComponent", () => {
  let component: PieChartComponent;
  let fixture: ComponentFixture<PieChartComponent>;
  let widgetService: jest.Mocked<WidgetService>;

  const mockWidget: Widget = {
    id: 1,
    name: "Test Pie Chart",
    widgetType: WidgetType.PieChart,
    configuration: {
      chartGrouping: ChartGrouping.Categories,
    },
  };

  const mockPieChartData: PieChartData = {
    data: [
      { label: "Category A", value: 100 },
      { label: "Category B", value: 200 },
      { label: "Category C", value: 150 },
    ],
  };

  beforeEach(async () => {
    const widgetServiceMock = {
      getPieChartData: jest.fn(),
    };

    await TestBed.configureTestingModule({
      imports: [PieChartComponent, SharedUiModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        { provide: WidgetService, useValue: widgetServiceMock },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    widgetService = TestBed.inject(WidgetService) as jest.Mocked<WidgetService>;
    widgetService.getPieChartData.mockReturnValue(of(mockPieChartData));

    fixture = TestBed.createComponent(PieChartComponent);
    component = fixture.componentInstance;
    component.widget = mockWidget;
    component.groupId = 1;
  });

  it("should create", () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  describe("initialization", () => {
    it("should have default isLoading as true", () => {
      expect(component.isLoading).toBe(true);
    });

    it("should have default hasData as false", () => {
      expect(component.hasData).toBe(false);
    });

    it("should have default empty pieChartData", () => {
      expect(component.pieChartData.labels).toEqual([]);
      expect(component.pieChartData.datasets[0].data).toEqual([]);
    });

    it("should have backgroundColor array in pieChartData", () => {
      expect(component.pieChartData.datasets[0].backgroundColor).toBeDefined();
      expect(
        (component.pieChartData.datasets[0].backgroundColor as string[]).length
      ).toBeGreaterThan(0);
    });
  });

  describe("ngOnInit", () => {
    it("should call loadData on init", fakeAsync(() => {
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).toHaveBeenCalledWith(1, {
        chartGrouping: ChartGrouping.Categories,
        filter: undefined,
      });
    }));

    it("should update chart data after loading", fakeAsync(() => {
      fixture.detectChanges();
      tick();

      expect(component.isLoading).toBe(false);
      expect(component.hasData).toBe(true);
      expect(component.pieChartData.labels).toEqual([
        "Category A",
        "Category B",
        "Category C",
      ]);
      expect(component.pieChartData.datasets[0].data).toEqual([100, 200, 150]);
    }));

    it("should not call service if groupId is not set", fakeAsync(() => {
      component.groupId = undefined;
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).not.toHaveBeenCalled();
      expect(component.isLoading).toBe(false);
    }));

    it("should not call service if widget configuration is not set", fakeAsync(() => {
      component.widget = { ...mockWidget, configuration: undefined };
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).not.toHaveBeenCalled();
      expect(component.isLoading).toBe(false);
    }));

    it("should not call service if chartGrouping is not set", fakeAsync(() => {
      component.widget = { ...mockWidget, configuration: {} };
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).not.toHaveBeenCalled();
      expect(component.isLoading).toBe(false);
    }));
  });

  describe("ngOnChanges", () => {
    it("should reload data when groupId changes", fakeAsync(() => {
      fixture.detectChanges();
      tick();
      widgetService.getPieChartData.mockClear();

      component.ngOnChanges({
        groupId: new SimpleChange(1, 2, false),
      });
      tick();

      expect(widgetService.getPieChartData).toHaveBeenCalled();
    }));

    it("should not reload data on first change", fakeAsync(() => {
      fixture.detectChanges();
      tick();
      widgetService.getPieChartData.mockClear();

      component.ngOnChanges({
        groupId: new SimpleChange(undefined, 1, true),
      });
      tick();

      expect(widgetService.getPieChartData).not.toHaveBeenCalled();
    }));

    it("should not reload data when other inputs change", fakeAsync(() => {
      fixture.detectChanges();
      tick();
      widgetService.getPieChartData.mockClear();

      component.ngOnChanges({
        widget: new SimpleChange(null, mockWidget, false),
      });
      tick();

      expect(widgetService.getPieChartData).not.toHaveBeenCalled();
    }));
  });

  describe("loadData", () => {
    it("should pass filter from widget configuration", fakeAsync(() => {
      const filterConfig = {
        chartGrouping: ChartGrouping.Tags,
        filter: { status: { value: ["OPEN"], operation: "equals" } },
      };
      component.widget = { ...mockWidget, configuration: filterConfig };
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).toHaveBeenCalledWith(1, {
        chartGrouping: ChartGrouping.Tags,
        filter: filterConfig.filter,
      });
    }));

    it("should handle empty response data", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(of({ data: [] }));
      fixture.detectChanges();
      tick();

      expect(component.hasData).toBe(false);
      expect(component.isLoading).toBe(false);
    }));

    it("should handle null response data", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(of({ data: undefined } as any));
      fixture.detectChanges();
      tick();

      expect(component.hasData).toBe(false);
    }));

    it("should handle data points with missing labels", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: undefined, value: 100 },
            { label: "Category B", value: 200 },
          ],
        } as any)
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.labels).toEqual(["Unknown", "Category B"]);
    }));

    it("should handle data points with missing values", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "Category A", value: undefined },
            { label: "Category B", value: 200 },
          ],
        } as any)
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.datasets[0].data).toEqual([0, 200]);
    }));
  });

  describe("updateChartData", () => {
    it("should preserve backgroundColor when updating data", fakeAsync(() => {
      const originalColors = component.pieChartData.datasets[0].backgroundColor;
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.datasets[0].backgroundColor).toEqual(
        originalColors
      );
    }));

    it("should handle single data point", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [{ label: "Only One", value: 500 }],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.labels).toEqual(["Only One"]);
      expect(component.pieChartData.datasets[0].data).toEqual([500]);
      expect(component.hasData).toBe(true);
    }));

    it("should handle many data points", fakeAsync(() => {
      const manyDataPoints = Array.from({ length: 20 }, (_, i) => ({
        label: `Category ${i}`,
        value: i * 10,
      }));
      widgetService.getPieChartData.mockReturnValue(
        of({ data: manyDataPoints })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.labels?.length).toBe(20);
      expect(component.pieChartData.datasets[0].data.length).toBe(20);
    }));
  });

  describe("getChartGroupingLabel", () => {
    it("should return 'Categories' for CATEGORIES grouping", () => {
      component.widget = {
        ...mockWidget,
        configuration: { chartGrouping: ChartGrouping.Categories },
      };
      expect(component.getChartGroupingLabel()).toBe("Categories");
    });

    it("should return 'Tags' for TAGS grouping", () => {
      component.widget = {
        ...mockWidget,
        configuration: { chartGrouping: ChartGrouping.Tags },
      };
      expect(component.getChartGroupingLabel()).toBe("Tags");
    });

    it("should return 'Paid By' for PAIDBY grouping", () => {
      component.widget = {
        ...mockWidget,
        configuration: { chartGrouping: ChartGrouping.Paidby },
      };
      expect(component.getChartGroupingLabel()).toBe("Paid By");
    });

    it("should return 'Unknown' for undefined grouping", () => {
      component.widget = {
        ...mockWidget,
        configuration: {},
      };
      expect(component.getChartGroupingLabel()).toBe("Unknown");
    });

    it("should return 'Unknown' for null configuration", () => {
      component.widget = {
        ...mockWidget,
        configuration: undefined,
      };
      expect(component.getChartGroupingLabel()).toBe("Unknown");
    });
  });

  describe("pieChartOptions", () => {
    it("should have responsive set to true", () => {
      expect(component.pieChartOptions?.responsive).toBe(true);
    });

    it("should have maintainAspectRatio set to false", () => {
      expect(component.pieChartOptions?.maintainAspectRatio).toBe(false);
    });

    it("should have legend display set to true", () => {
      expect(component.pieChartOptions?.plugins?.legend?.display).toBe(true);
    });

    it("should have legend position set to bottom", () => {
      expect(component.pieChartOptions?.plugins?.legend?.position).toBe(
        "bottom"
      );
    });

    it("should have tooltip callback defined", () => {
      expect(
        component.pieChartOptions?.plugins?.tooltip?.callbacks?.label
      ).toBeDefined();
    });

    it("should have datalabels configuration", () => {
      expect(component.pieChartOptions?.plugins?.datalabels).toBeDefined();
    });
  });

  describe("tooltip callback", () => {
    it("should format tooltip label correctly", () => {
      const labelCallback =
        component.pieChartOptions?.plugins?.tooltip?.callbacks?.label;
      expect(labelCallback).toBeDefined();

      if (labelCallback) {
        const mockContext = {
          label: "Category A",
          parsed: 100,
          dataset: { data: [100, 200, 200] },
        };
        const result = labelCallback(mockContext as any);
        expect(result).toBe("Category A: $100.00 (20.0%)");
      }
    });

    it("should handle zero total in tooltip", () => {
      const labelCallback =
        component.pieChartOptions?.plugins?.tooltip?.callbacks?.label;

      if (labelCallback) {
        const mockContext = {
          label: "Category A",
          parsed: 0,
          dataset: { data: [0, 0, 0] },
        };
        const result = labelCallback(mockContext as any);
        expect(result).toBe("Category A: $0.00 (0%)");
      }
    });

    it("should handle missing label in tooltip", () => {
      const labelCallback =
        component.pieChartOptions?.plugins?.tooltip?.callbacks?.label;

      if (labelCallback) {
        const mockContext = {
          label: "",
          parsed: 100,
          dataset: { data: [100, 200] },
        };
        const result = labelCallback(mockContext as any);
        expect(result).toBe(": $100.00 (33.3%)");
      }
    });
  });

  describe("datalabels formatter", () => {
    it("should format percentage correctly for large slices", () => {
      const formatter =
        component.pieChartOptions?.plugins?.datalabels?.formatter;
      expect(formatter).toBeDefined();

      if (formatter) {
        const mockContext = {
          dataset: { data: [100, 100] },
        };
        const result = (formatter as Function)(50, mockContext);
        expect(result).toBe("25.0%");
      }
    });

    it("should return empty string for small slices (< 5%)", () => {
      const formatter =
        component.pieChartOptions?.plugins?.datalabels?.formatter;

      if (formatter) {
        const mockContext = {
          dataset: { data: [100, 900] },
        };
        // 100 out of 1000 = 10%
        const result = (formatter as Function)(4, mockContext);
        // 4 out of 1000 = 0.4% which is < 5%
        const smallResult = (formatter as Function)(4, {
          dataset: { data: [4, 996] },
        });
        expect(smallResult).toBe("");
      }
    });

    it("should handle zero total in formatter", () => {
      const formatter =
        component.pieChartOptions?.plugins?.datalabels?.formatter;

      if (formatter) {
        const mockContext = {
          dataset: { data: [0, 0] },
        };
        const result = (formatter as Function)(0, mockContext);
        expect(result).toBe("");
      }
    });
  });

  describe("template rendering", () => {
    it("should display widget name in header", fakeAsync(() => {
      fixture.detectChanges();
      tick();

      const header = fixture.debugElement.query(By.css("h3"));
      expect(header.nativeElement.textContent.trim()).toBe("Test Pie Chart");
    }));

    it("should display default name when widget name is empty", fakeAsync(() => {
      component.widget = { ...mockWidget, name: "" };
      fixture.detectChanges();
      tick();

      const header = fixture.debugElement.query(By.css("h3"));
      expect(header.nativeElement.textContent.trim()).toBe("Pie Chart");
    }));

    it("should display chart grouping badge", fakeAsync(() => {
      fixture.detectChanges();
      tick();

      const badge = fixture.debugElement.query(By.css(".badge"));
      expect(badge.nativeElement.textContent.trim()).toBe("Categories");
    }));

    it("should pass correct props to app-pie-chart-ui", fakeAsync(() => {
      fixture.detectChanges();
      tick();

      const pieChartUi = fixture.debugElement.query(
        By.css("app-pie-chart-ui")
      );
      expect(pieChartUi).toBeTruthy();
    }));
  });

  describe("different chart groupings", () => {
    it("should load data with TAGS grouping", fakeAsync(() => {
      component.widget = {
        ...mockWidget,
        configuration: { chartGrouping: ChartGrouping.Tags },
      };
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).toHaveBeenCalledWith(1, {
        chartGrouping: ChartGrouping.Tags,
        filter: undefined,
      });
    }));

    it("should load data with PAIDBY grouping", fakeAsync(() => {
      component.widget = {
        ...mockWidget,
        configuration: { chartGrouping: ChartGrouping.Paidby },
      };
      fixture.detectChanges();
      tick();

      expect(widgetService.getPieChartData).toHaveBeenCalledWith(1, {
        chartGrouping: ChartGrouping.Paidby,
        filter: undefined,
      });
    }));
  });

  describe("edge cases", () => {
    it("should handle very large values", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "Big", value: 999999999 },
            { label: "Small", value: 1 },
          ],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.datasets[0].data).toEqual([999999999, 1]);
    }));

    it("should handle decimal values", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "A", value: 10.55 },
            { label: "B", value: 20.45 },
          ],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.datasets[0].data).toEqual([10.55, 20.45]);
    }));

    it("should handle zero values", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "Zero", value: 0 },
            { label: "Some", value: 100 },
          ],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.datasets[0].data).toEqual([0, 100]);
      expect(component.hasData).toBe(true);
    }));

    it("should handle negative values gracefully", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "Negative", value: -50 },
            { label: "Positive", value: 100 },
          ],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.datasets[0].data).toEqual([-50, 100]);
    }));

    it("should handle special characters in labels", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "Category & Special <chars>", value: 100 },
            { label: 'With "quotes"', value: 200 },
          ],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.labels).toEqual([
        "Category & Special <chars>",
        'With "quotes"',
      ]);
    }));

    it("should handle unicode in labels", fakeAsync(() => {
      widgetService.getPieChartData.mockReturnValue(
        of({
          data: [
            { label: "日本語", value: 100 },
            { label: "Émojis 🎉", value: 200 },
          ],
        })
      );
      fixture.detectChanges();
      tick();

      expect(component.pieChartData.labels).toEqual(["日本語", "Émojis 🎉"]);
    }));
  });
});
