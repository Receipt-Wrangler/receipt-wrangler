import { Pipe, PipeTransform } from "@angular/core";
import { WidgetType } from "../open-api/index";

@Pipe({
    name: "widgetType",
    standalone: false
})
export class WidgetTypePipe implements PipeTransform {
  transform(value: WidgetType): string {
    switch (value) {
      case WidgetType.FilteredReceipts:
        return "Filtered Receipts";
      case WidgetType.GroupSummary:
        return "Group Summary";
    }

    return "";
  }
}
