import { ReceiptStatus } from "../open-api";

export interface QuickScanCommand {
  file: File;
  groupId: number;
  status: ReceiptStatus;
  paidByUserId: number;
}
