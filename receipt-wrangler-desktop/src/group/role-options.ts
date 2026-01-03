import { FormOption } from "src/interfaces/form-option.interface";
import { formatStatus } from "src/utils";
import { GroupRole, UserRole } from "../open-api";

export const GROUP_ROLE_OPTIONS: FormOption[] = Object.keys(GroupRole).map(
  (key) => {
    const value = (GroupRole as any)[key];
    return {
      value: value,
      displayValue: formatStatus(value),
    };
  }
);

export const USER_ROLE_OPTIONS: FormOption[] = Object.keys(UserRole).map(
  (key) => {
    const value = (UserRole as any)[key];
    return {
      value: value,
      displayValue: formatStatus(value),
    };
  }
);
