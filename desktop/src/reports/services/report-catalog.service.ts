import { Injectable, computed, inject, signal } from "@angular/core";
import { catchError, of, take } from "rxjs";
import { CustomField, CustomFieldService, CustomFieldType } from "../../open-api";
import {
  REPORT_BUILTIN_DIMENSIONS,
  REPORT_BUILTIN_MEASURES,
  ReportField,
} from "../models/report-catalog.constants";

/**
 * Supplies the dimension and measure catalog the builder's dropdowns bind to: the
 * built-in engine fields plus every custom field, keyed `custom_<id>` to match the
 * backend. The backend validates every field key, so this is the picker's option
 * source, not the authorization gate.
 *
 * Every custom field is a dimension, whatever its type — the engine restricts
 * measuring, not cutting, so a currency field is groupable as well as summable and
 * appears in *both* lists. A date field additionally contributes the three
 * calendar-period fields the backend derives for it (`custom_<id>_month` and
 * friends): grouping by the raw date buckets on the exact instant, which puts
 * every receipt in its own group.
 */
@Injectable({ providedIn: "root" })
export class ReportCatalogService {
  private readonly customFieldService = inject(CustomFieldService);
  private readonly customFields = signal<CustomField[]>([]);
  private loaded = false;

  public readonly dimensions = computed<ReportField[]>(() => [
    ...REPORT_BUILTIN_DIMENSIONS,
    ...this.customFields().flatMap(customFieldDimensions),
  ]);

  public readonly measures = computed<ReportField[]>(() => [
    ...REPORT_BUILTIN_MEASURES,
    ...this.customFields()
      .filter((field) => field.type === CustomFieldType.Currency)
      .map(customFieldToField),
  ]);

  /**
   * Loads the custom-field pool once. A user lacking app.custom-fields.read gets a
   * 403 that is swallowed, leaving just the built-in catalog rather than erroring.
   */
  public load(): void {
    if (this.loaded) {
      return;
    }
    this.loaded = true;
    this.customFieldService
      .getPagedCustomFields({ page: 1, pageSize: -1, orderBy: "name", sortDirection: "desc" })
      .pipe(
        take(1),
        catchError(() => of({ data: [], totalCount: 0 }))
      )
      .subscribe((paged) =>
        this.customFields.set((paged.data ?? []) as unknown as CustomField[])
      );
  }
}

function customFieldToField(field: CustomField): ReportField {
  return { key: "custom_" + field.id, label: field.name, isCustom: true };
}

/**
 * The dimensions one custom field offers: itself, plus — for a date field — the
 * day/month/year fields the backend derives from it. The keys and labels mirror
 * `receiptsource.CustomFieldPeriodKeys` / `dateFieldRefs`; a mismatch here is a
 * 400 from the engine, not a silent fallback.
 */
function customFieldDimensions(field: CustomField): ReportField[] {
  const self = customFieldToField(field);
  if (field.type !== CustomFieldType.Date) {
    return [self];
  }
  return [
    self,
    { key: self.key + "_day", label: field.name + " (Day)", isCustom: true },
    { key: self.key + "_month", label: field.name + " (Month)", isCustom: true },
    { key: self.key + "_year", label: field.name + " (Year)", isCustom: true },
  ];
}
