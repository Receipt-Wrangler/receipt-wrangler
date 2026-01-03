import { FormOption } from "../../interfaces/form-option.interface";
import { AssociatedApiKeys } from "../../open-api";

export const associatedApiKeyOptions: FormOption[] = [
  {
    value: AssociatedApiKeys.Mine,
    displayValue: "My API Keys",
  },
  {
    value: AssociatedApiKeys.All,
    displayValue: "All API Keys",
  }
];