import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

/// [conditions] encoded and put through the real generated serializer.
///
/// Asserting against the serialized map rather than the builder is what proves
/// the `JsonObject` wrapping is right: the Go handler type-asserts each value
/// (`.(string)`, `.(float64)`, `.([]interface{})`) with no comma-ok, so a wrong
/// shape is a 500 rather than an ignored condition. A Dart-side check would not
/// catch that.
Map<String, dynamic> serializeFilter(
    Map<String, ReceiptFilterCondition> conditions) {
  final filter = buildReceiptPagedRequestFilter(conditions).build();

  return Map<String, dynamic>.from(api.standardSerializers.serializeWith(
      api.ReceiptPagedRequestFilter.serializer, filter) as Map);
}

/// The `{operation, value}` pair serialized for [key], or null when the encoder
/// left that field unset.
Map<String, dynamic>? serializedField(
    Map<String, ReceiptFilterCondition> conditions, String key) {
  final raw = serializeFilter(conditions)[key];
  return raw == null ? null : Map<String, dynamic>.from(raw as Map);
}

/// `buildUserView` / `buildGroup` / `buildGroupMember` live in
/// `receipt_form_test_helpers.dart`; only the category and tag builders are new
/// here.
api.Category buildCategory(int id, String name) =>
    (api.CategoryBuilder()
          ..id = id
          ..name = name)
        .build();

api.Tag buildTag(int id, String name) => (api.TagBuilder()
      ..id = id
      ..name = name)
    .build();
