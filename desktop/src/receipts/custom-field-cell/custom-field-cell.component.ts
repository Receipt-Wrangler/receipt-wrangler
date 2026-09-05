import { Component, computed, input } from "@angular/core";
import { CustomField, CustomFieldType, CustomFieldValue, Receipt } from "../../open-api";

/**
 * Renders one custom field's value for a receipt row in the receipts table.
 *
 * The read-only counterpart of `app-custom-field`, which switches on the same
 * `CustomFieldType` to build the *editable* control on the receipt form. A cell
 * cannot reuse that component because there is no FormGroup here — the table
 * renders plain API records.
 */
@Component({
  selector: "app-custom-field-cell",
  standalone: false,
  templateUrl: "./custom-field-cell.component.html",
})
export class CustomFieldCellComponent {
  public readonly receipt = input<Receipt | undefined>(undefined);

  public readonly customField = input<CustomField | undefined>(undefined);

  /**
   * The value the row shows for this field.
   *
   * Nothing stops a receipt carrying several values for one field, so a winner
   * has to be picked. The lowest id among the values that actually hold
   * something wins — the same rule the backend sorts by and the reporting engine
   * reads by, so the cell, the sort order and a report never disagree.
   */
  public readonly value = computed<CustomFieldValue | undefined>(() => {
    const customField = this.customField();
    if (!customField) {
      return undefined;
    }

    return ((this.receipt()?.customFields ?? []) as CustomFieldValue[])
      .filter(
        (value) =>
          value.customFieldId === customField.id &&
          this.hasValue(value, customField.type)
      )
      .sort((a, b) => a.id - b.id)[0];
  });

  /**
   * The raw currency value, as the string the API sends it as. Narrowed here
   * rather than in the template: the cell only renders once a value is set, but
   * the template cannot know that, and the currency pipe takes no undefined.
   */
  public readonly currencyValue = computed<string>(
    () => this.value()?.currencyValue ?? ""
  );

  public readonly selectedOptionValue = computed<string>(() => {
    const selectValue = this.value()?.selectValue;
    const option = this.customField()?.options?.find(
      (candidate) => candidate.id === selectValue
    );

    return option?.value ?? "";
  });

  protected readonly CustomFieldType = CustomFieldType;

  private hasValue(value: CustomFieldValue, type: CustomFieldType): boolean {
    switch (type) {
      case CustomFieldType.Text:
        return value.stringValue !== null && value.stringValue !== undefined;
      case CustomFieldType.Date:
        return value.dateValue !== null && value.dateValue !== undefined;
      case CustomFieldType.Select:
        return value.selectValue !== null && value.selectValue !== undefined;
      case CustomFieldType.Currency:
        return value.currencyValue !== null && value.currencyValue !== undefined;
      case CustomFieldType.Boolean:
        return value.booleanValue !== null && value.booleanValue !== undefined;
      default:
        return false;
    }
  }
}
