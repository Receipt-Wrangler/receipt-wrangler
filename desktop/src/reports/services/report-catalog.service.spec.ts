import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { of, throwError } from "rxjs";
import { CustomFieldService, CustomFieldType } from "../../open-api";
import { ReportCatalogService } from "./report-catalog.service";

describe("ReportCatalogService", () => {
  let service: ReportCatalogService;
  let customFieldService: { getPagedCustomFields: jest.Mock };

  beforeEach(() => {
    customFieldService = { getPagedCustomFields: jest.fn() };
    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        ReportCatalogService,
        { provide: CustomFieldService, useValue: customFieldService },
      ],
    });
    service = TestBed.inject(ReportCatalogService);
  });

  it("starts with the built-in dimensions and the amount measure", () => {
    expect(service.dimensions().some((field) => field.key === "category")).toBe(true);
    expect(service.measures().map((field) => field.key)).toEqual(["amount"]);
  });

  it("adds currency custom fields as measures and other kinds as dimensions", () => {
    customFieldService.getPagedCustomFields.mockReturnValue(
      of({
        data: [
          { id: 7, name: "HST", type: CustomFieldType.Currency },
          { id: 8, name: "Vendor", type: CustomFieldType.Text },
        ],
        totalCount: 2,
      })
    );

    service.load();

    expect(service.measures().find((field) => field.key === "custom_7")?.label).toBe("HST");
    expect(service.dimensions().find((field) => field.key === "custom_8")?.label).toBe("Vendor");
    expect(service.measures().some((field) => field.key === "custom_8")).toBe(false);
  });

  it("keeps only the built-ins when the custom-field lookup fails", () => {
    customFieldService.getPagedCustomFields.mockReturnValue(throwError(() => new Error("403")));

    service.load();

    expect(service.dimensions().some((field) => field.key.startsWith("custom_"))).toBe(false);
    expect(service.measures().length).toBe(1);
  });

  it("loads the custom-field pool only once", () => {
    customFieldService.getPagedCustomFields.mockReturnValue(of({ data: [], totalCount: 0 }));
    service.load();
    service.load();
    expect(customFieldService.getPagedCustomFields).toHaveBeenCalledTimes(1);
  });
});
