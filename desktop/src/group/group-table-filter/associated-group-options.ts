import { FormOption } from "../../interfaces/form-option.interface";
import { AssociatedGroup } from "../../open-api";

export const associatedGroupOptions: FormOption[] = [
  {
    value: AssociatedGroup.Mine,
    displayValue: "My Groups",
  },
  {
    value: AssociatedGroup.All,
    displayValue: "All Groups",
  }
];
