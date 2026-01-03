import { Pipe, PipeTransform } from "@angular/core";
import { CustomFieldType } from "../open-api/index";


@Pipe({
  name: "customFieldType",
  standalone: false
})
export class CustomFieldTypePipe implements PipeTransform {

  public transform(type: CustomFieldType): string {
    switch (type) {
      case CustomFieldType.Date:
        return "Date";
      case CustomFieldType.Select:
        return "Select";
      case CustomFieldType.Text:
        return "Text";
      case CustomFieldType.Currency:
        return "Currency";
      case CustomFieldType.Boolean:
        return "Boolean";
      default:
        return "";
    }
  }
}
