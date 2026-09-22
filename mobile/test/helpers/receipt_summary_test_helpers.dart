import 'package:built_collection/built_collection.dart';
import 'package:openapi/openapi.dart' as api;

/// Builders for summary responses, shared by the widget tests and the demo harness so
/// the two cannot drift apart.

api.ReceiptSummaryCustomFieldTotal buildSummaryCustomFieldTotal({
  required int customFieldId,
  required String name,
  required String total,
}) =>
    (api.ReceiptSummaryCustomFieldTotalBuilder()
          ..customFieldId = customFieldId
          ..name = name
          ..total = total)
        .build();

api.ReceiptSummaryRow buildSummaryRow({
  /// The overall row carries ReceiptStatus.empty -- the generated field is non-nullable,
  /// so there is no "no status" to express other than the empty member.
  api.ReceiptStatus status = api.ReceiptStatus.empty,
  int receiptCount = 0,
  String total = '0.00',
  List<api.ReceiptSummaryCustomFieldTotal> customFieldTotals = const [],
}) =>
    (api.ReceiptSummaryRowBuilder()
          ..status = status
          ..receiptCount = receiptCount
          ..total = total
          ..customFieldTotals =
              ListBuilder<api.ReceiptSummaryCustomFieldTotal>(customFieldTotals))
        .build();

api.ReceiptSummary buildReceiptSummary({
  bool enabled = true,
  int configurationGroupId = 1,
  api.ReceiptSummaryPosition position = api.ReceiptSummaryPosition.BOTTOM,
  api.ReceiptSummaryRow? overall,
  List<api.ReceiptSummaryRow> statuses = const [],
}) =>
    (api.ReceiptSummaryBuilder()
          ..enabled = enabled
          ..configurationGroupId = configurationGroupId
          ..position = position
          ..overall = (overall ?? buildSummaryRow(receiptCount: 0)).toBuilder()
          ..statuses = ListBuilder<api.ReceiptSummaryRow>(statuses))
        .build();
