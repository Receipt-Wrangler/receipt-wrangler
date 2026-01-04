class ReceiptSortOption {
  final String columnName;
  final String displayLabel;

  const ReceiptSortOption({
    required this.columnName,
    required this.displayLabel,
  });
}

const List<ReceiptSortOption> receiptSortOptions = [
  ReceiptSortOption(columnName: "created_at", displayLabel: "Added At"),
  ReceiptSortOption(columnName: "date", displayLabel: "Receipt Date"),
  ReceiptSortOption(columnName: "name", displayLabel: "Name"),
  ReceiptSortOption(columnName: "amount", displayLabel: "Amount"),
  ReceiptSortOption(columnName: "paid_by_user_id", displayLabel: "Paid By"),
  ReceiptSortOption(columnName: "status", displayLabel: "Status"),
  ReceiptSortOption(columnName: "resolved_date", displayLabel: "Resolved Date"),
];
