import { UserPreferences } from '../open-api';
import { ReceiptTableColumnConfig } from './receipt-table-column-config.interface';

export interface ExtendedUserPreferences extends UserPreferences {
  receiptTableColumns?: ReceiptTableColumnConfig[];
}