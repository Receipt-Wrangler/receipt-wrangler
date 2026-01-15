import { FormOption } from "../../interfaces/form-option.interface";
import { WidgetType } from "../../open-api/index";

export const widgetTypeOptions: FormOption[] = [
  {
    value: WidgetType.FilteredReceipts,
    displayValue: "Filtered Receipts",
  },
  {
    value: WidgetType.GroupSummary,
    displayValue: "Group Summary",
  },
  {
    value: WidgetType.GroupActivity,
    displayValue: "Activity",
  },
  {
    value: WidgetType.PieChart,
    displayValue: "Pie Chart",
  }
];
