import 'package:openapi/openapi.dart' as api;
import 'package:uuid/uuid.dart';

class FormItem {
  final String formId;
  final String name;
  final String amount;
  final int chargedToUserId;
  final int receiptId;
  final api.ItemStatus status;
  final List<api.Category> categories;
  final List<api.Tag> tags;

  FormItem({
    required this.formId,
    required String name,
    required String amount,
    required int chargedToUserId,
    required int receiptId,
    required api.ItemStatus status,
    required this.categories,
    required this.tags,
  })  : name = name,
        amount = amount,
        chargedToUserId = chargedToUserId,
        receiptId = receiptId,
        status = status;

  static FormItem fromItem(api.Item item) {
    return FormItem(
      formId: Uuid().v4(),
      name: item.name,
      amount: item.amount,
      chargedToUserId: item.chargedToUserId,
      receiptId: item.receiptId,
      status: item.status,
      categories: item.categories?.toList() ?? [],
      tags: item.tags?.toList() ?? [],
    );
  }

  static List<FormItem> fromItems(List<api.Item> items) {
    return items.map((item) => FormItem.fromItem(item)).toList();
  }

  static String buildItemNameName(FormItem item) {
    return "items.${item.formId}.name";
  }

  static String buildItemAmountName(FormItem item) {
    return "items.${item.formId}.amount";
  }

  static String buildItemStatusName(FormItem item) {
    return "items.${item.formId}.status";
  }

  static String buildItemCategoryName(FormItem item) {
    return "items.${item.formId}.category";
  }

  static String buildItemTagName(FormItem item) {
    return "items.${item.formId}.tag";
  }
}
