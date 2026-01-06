import { FormOption } from "../../interfaces/form-option.interface";
import { ChartGrouping } from "../../open-api/index";

export const chartGroupingOptions: FormOption[] = [
  {
    value: ChartGrouping.Categories,
    displayValue: "Categories",
  },
  {
    value: ChartGrouping.Tags,
    displayValue: "Tags",
  },
  {
    value: ChartGrouping.Paidby,
    displayValue: "Paid By",
  }
];
