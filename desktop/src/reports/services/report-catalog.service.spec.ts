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

  // Every custom field cuts, so a currency one is a dimension AND a measure; only
  // measuring is type-restricted.
  it("adds every custom field as a dimension and currency ones as measures too", () => {
    customFieldService.getPagedCustomFields.mockReturnValue(
      of({
        data: [
          { id: 7, name: "HST", type: CustomFieldType.Currency },
          { id: 8, name: "Vendor", type: CustomFieldType.Text },
          { id: 9, name: "Reimbursed", type: CustomFieldType.Boolean },
        ],
        totalCount: 3,
      })
    );

    service.load();

    expect(service.measures().find((field) => field.key === "custom_7")?.label).toBe("HST");
    expect(service.dimensions().find((field) => field.key === "custom_7")?.label).toBe("HST");
    expect(service.dimensions().find((field) => field.key === "custom_8")?.label).toBe("Vendor");
    expect(service.dimensions().find((field) => field.key === "custom_9")?.label).toBe("Reimbursed");
    expect(service.measures().some((field) => field.key === "custom_8")).toBe(false);
  });

  it("marks custom fields so the builder can badge them", () => {
    customFieldService.getPagedCustomFields.mockReturnValue(
      of({ data: [{ id: 7, name: "HST", type: CustomFieldType.Currency }], totalCount: 1 })
    );

    service.load();

    expect(service.dimensions().find((field) => field.key === "custom_7")?.isCustom).toBe(true);
    expect(service.measures().find((field) => field.key === "custom_7")?.isCustom).toBe(true);
    expect(service.dimensions().find((field) => field.key === "category")?.isCustom).toBeUndefined();
  });

  // Grouping by a raw date buckets on the exact instant (one bucket per receipt), so
  // a date custom field also offers the calendar-period fields the backend derives.
  it("offers day/month/year dimensions for a date custom field", () => {
    customFieldService.getPagedCustomFields.mockReturnValue(
      of({ data: [{ id: 4, name: "Due Date", type: CustomFieldType.Date }], totalCount: 1 })
    );

    service.load();

    const keys = service.dimensions().map((field) => field.key);
    expect(keys).toEqual(
      expect.arrayContaining(["custom_4", "custom_4_day", "custom_4_month", "custom_4_year"])
    );
    expect(service.dimensions().find((field) => field.key === "custom_4_month")?.label).toBe(
      "Due Date (Month)"
    );
    // Period fields are dimensions only — there is nothing to sum in a calendar month.
    expect(service.measures().some((field) => field.key.startsWith("custom_4"))).toBe(false);
  });

  it("derives period fields only for date custom fields", () => {
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

    expect(service.dimensions().some((field) => field.key.endsWith("_month"))).toBe(true); // date_month, a built-in
    expect(service.dimensions().some((field) => field.key === "custom_7_month")).toBe(false);
    expect(service.dimensions().some((field) => field.key === "custom_8_month")).toBe(false);
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
