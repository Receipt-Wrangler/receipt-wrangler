import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan_field_config.dart';

// resolveQuickScanFieldConfig is the single source of the quick-scan field
// show/require defaults, shared by the per-image form and the submit path.
// These tests pin the backend-mirroring defaults so the two callers can't drift.
api.GroupReceiptSettings _settings({
  bool paidByEnabled = true,
  bool paidByRequired = true,
  bool statusEnabled = true,
  bool statusRequired = true,
  bool categoriesEnabled = false,
  bool categoriesRequired = false,
  bool tagsEnabled = false,
  bool tagsRequired = false,
}) {
  return (api.GroupReceiptSettingsBuilder()
        ..id = 1
        ..createdAt = ''
        ..groupId = 1
        ..quickScanPaidByEnabled = paidByEnabled
        ..quickScanPaidByRequired = paidByRequired
        ..quickScanStatusEnabled = statusEnabled
        ..quickScanStatusRequired = statusRequired
        ..quickScanCategoriesEnabled = categoriesEnabled
        ..quickScanCategoriesRequired = categoriesRequired
        ..quickScanTagsEnabled = tagsEnabled
        ..quickScanTagsRequired = tagsRequired)
      .build();
}

void main() {
  group('resolveQuickScanFieldConfig', () {
    test('null settings falls back to backend defaults', () {
      final config = resolveQuickScanFieldConfig(null);

      // Paid-by/status shown+required; categories/tags hidden.
      expect(config.showPaidBy, isTrue);
      expect(config.requirePaidBy, isTrue);
      expect(config.showStatus, isTrue);
      expect(config.requireStatus, isTrue);
      expect(config.showCategories, isFalse);
      expect(config.requireCategories, isFalse);
      expect(config.showTags, isFalse);
      expect(config.requireTags, isFalse);
    });

    test('mirrors each required flag when the fields are shown', () {
      // All fields shown so require reflects the persisted required flag directly.
      final config = resolveQuickScanFieldConfig(_settings(
        paidByEnabled: true,
        paidByRequired: false,
        statusEnabled: true,
        statusRequired: true,
        categoriesEnabled: true,
        categoriesRequired: true,
        tagsEnabled: true,
        tagsRequired: false,
      ));

      expect(config.showPaidBy, isTrue);
      expect(config.requirePaidBy, isFalse);
      expect(config.showStatus, isTrue);
      expect(config.requireStatus, isTrue);
      expect(config.showCategories, isTrue);
      expect(config.requireCategories, isTrue);
      expect(config.showTags, isTrue);
      expect(config.requireTags, isFalse);
    });

    test('require is false when the field is hidden even if required flag set', () {
      // A hidden field can never be "required" — require is gated on show.
      final config = resolveQuickScanFieldConfig(_settings(
        paidByEnabled: false,
        paidByRequired: true,
        statusEnabled: false,
        statusRequired: true,
        categoriesEnabled: false,
        categoriesRequired: true,
        tagsEnabled: false,
        tagsRequired: true,
      ));

      expect(config.showPaidBy, isFalse);
      expect(config.requirePaidBy, isFalse);
      expect(config.showStatus, isFalse);
      expect(config.requireStatus, isFalse);
      expect(config.showCategories, isFalse);
      expect(config.requireCategories, isFalse);
      expect(config.showTags, isFalse);
      expect(config.requireTags, isFalse);
    });
  });
}
