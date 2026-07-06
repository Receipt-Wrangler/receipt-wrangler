import 'package:openapi/openapi.dart';

/// Resolved show/require flags for the four quick-scan fields (paid-by, status,
/// categories, tags), derived from a group's [GroupReceiptSettings]. Mirrors the
/// backend's `resolveQuickScanFields` defaults so the per-image form (visibility +
/// validators) and the submit path (which fields to send / require) stay aligned.
///
/// Null [settings] (no group selected yet) falls back to the backend defaults:
/// paid-by/status shown, categories/tags hidden.
class QuickScanFieldConfig {
  final bool showPaidBy;
  final bool requirePaidBy;
  final bool showStatus;
  final bool requireStatus;
  final bool showCategories;
  final bool requireCategories;
  final bool showTags;
  final bool requireTags;

  const QuickScanFieldConfig({
    required this.showPaidBy,
    required this.requirePaidBy,
    required this.showStatus,
    required this.requireStatus,
    required this.showCategories,
    required this.requireCategories,
    required this.showTags,
    required this.requireTags,
  });
}

QuickScanFieldConfig resolveQuickScanFieldConfig(GroupReceiptSettings? settings) {
  final showPaidBy = settings?.quickScanPaidByEnabled ?? true;
  final showStatus = settings?.quickScanStatusEnabled ?? true;
  final showCategories = settings?.quickScanCategoriesEnabled ?? false;
  final showTags = settings?.quickScanTagsEnabled ?? false;

  return QuickScanFieldConfig(
    showPaidBy: showPaidBy,
    requirePaidBy: showPaidBy && (settings?.quickScanPaidByRequired ?? true),
    showStatus: showStatus,
    requireStatus: showStatus && (settings?.quickScanStatusRequired ?? true),
    showCategories: showCategories,
    requireCategories:
        showCategories && (settings?.quickScanCategoriesRequired ?? false),
    showTags: showTags,
    requireTags: showTags && (settings?.quickScanTagsRequired ?? false),
  );
}
