import { Pipe, PipeTransform } from "@angular/core";
import {
  dateOperationOptions,
  listOperationOptions,
  numberOperationOptions,
  textOperationOptions,
  usersOperationOptions
} from "src/constants";
import { FilterOperation } from "src/open-api";

@Pipe({
    name: "operations",
    standalone: false
})
export class OperationsPipe implements PipeTransform {
  private displayValues: { [key: string]: string } = {
    [FilterOperation.Contains]: "Contains",
    [FilterOperation.Equals]: "Equals",
    [FilterOperation.GreaterThan]: "Greater than",
    [FilterOperation.LessThan]: "Less than",
    [FilterOperation.Between]: "Between",
    [FilterOperation.WithinCurrentMonth]: "Within current month"
  };

  public transform(type: string, display: boolean): string[] {
    let operationOptions: FilterOperation[] = [];

    switch (type) {
      case "date":
        return dateOperationOptions.map((option) => this.getDisplayValue(option, display));

      case "text":
        return textOperationOptions.map((option) => this.getDisplayValue(option, display));

      case "number":
        return numberOperationOptions.map((option) => this.getDisplayValue(option, display));

      case "list":
        return listOperationOptions.map((option) => this.getDisplayValue(option, display));

      case "users":
        return usersOperationOptions.map((option) => this.getDisplayValue(option, display));

      default:
        return [];
    }
  }

  private getDisplayValue(option: FilterOperation | null, display: boolean): string {
    if (display) {
      return this.displayValues?.[option ?? ""] ?? "";
    } else {
      return option ?? "";
    }
  }
}
