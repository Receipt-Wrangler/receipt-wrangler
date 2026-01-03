import { FormOption } from "src/interfaces/form-option.interface";
import { formatStatus } from "src/utils";
import { GroupStatus, ItemStatus, ReceiptStatus } from "../open-api";

export const RECEIPT_STATUS_OPTIONS: FormOption[] = Object.keys(
  ReceiptStatus
).filter(k => !!(ReceiptStatus as any)[k]).map((key) => {
  const value = (ReceiptStatus as any)[key];
  return {
    value: value,
    displayValue: formatStatus(value),
  };
});

export const RECEIPT_ITEM_STATUS_OPTIONS: FormOption[] = Object.keys(
  ItemStatus
).map((key) => {
  const value = (ItemStatus as any)[key];
  return {
    value: value,
    displayValue: formatStatus(value),
  };
});

export const GROUP_STATUS_OPTIONS: FormOption[] = Object.keys(GroupStatus).map(
  (key) => {
    const value = (GroupStatus as any)[key];
    return {
      value: value,
      displayValue: formatStatus(value),
    };
  }
);
