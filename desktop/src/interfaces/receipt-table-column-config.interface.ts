export interface ReceiptTableColumnConfig {
  matColumnDef: string;
  visible: boolean;
  order: number;
}

export interface ReceiptTableColumns {
  columns: ReceiptTableColumnConfig[];
}

export const DEFAULT_RECEIPT_TABLE_COLUMNS: ReceiptTableColumnConfig[] = [
  { matColumnDef: 'created_at', visible: true, order: 0 },
  { matColumnDef: 'date', visible: true, order: 1 },
  { matColumnDef: 'name', visible: true, order: 2 },
  { matColumnDef: 'paid_by_user_id', visible: true, order: 3 },
  { matColumnDef: 'amount', visible: true, order: 4 },
  { matColumnDef: 'categories', visible: true, order: 5 },
  { matColumnDef: 'tags', visible: true, order: 6 },
  { matColumnDef: 'status', visible: true, order: 7 },
  { matColumnDef: 'resolved_date', visible: true, order: 8 },
];