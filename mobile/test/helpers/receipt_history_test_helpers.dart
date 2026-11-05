import 'package:openapi/openapi.dart' as api;

/// Receipt history fixtures, built from JSON through the real generated
/// deserializer so a test exercises the same contract the screen receives.

api.ReceiptHistory historyFromJson(Map<String, Object?> json) =>
    api.standardSerializers.deserializeWith(api.ReceiptHistory.serializer, json)!;

api.ReceiptHistoryEntry entryFromJson(Map<String, Object?> json) => api
    .standardSerializers
    .deserializeWith(api.ReceiptHistoryEntry.serializer, json)!;

api.ReceiptHistoryChange changeFromJson(Map<String, Object?> json) => api
    .standardSerializers
    .deserializeWith(api.ReceiptHistoryChange.serializer, json)!;

Map<String, Object?> historyJson(
  List<Map<String, Object?>> entries, {
  int changeCount = 0,
  int saveCount = 0,
  int contributorCount = 0,
}) =>
    {
      'entries': entries,
      'changeCount': changeCount,
      'saveCount': saveCount,
      'contributorCount': contributorCount,
    };

Map<String, Object?> entryJson({
  int systemTaskId = 1,
  String action = 'UPDATED',
  String recordedAt = '2026-10-13T15:00:00Z',
  Map<String, Object?> actor = const {'type': 'USER', 'userId': 2},
  bool hasSnapshots = true,
  String? errorMessage,
  String? beforeState,
  String? beforeSourceType,
  String? beforeSourceRecordedAt,
  List<Map<String, Object?>> changes = const [],
}) =>
    {
      'systemTaskId': systemTaskId,
      'action': action,
      'recordedAt': recordedAt,
      'actor': actor,
      'hasSnapshots': hasSnapshots,
      if (errorMessage != null) 'errorMessage': errorMessage,
      if (beforeState != null) 'beforeState': beforeState,
      if (beforeSourceType != null) 'beforeSourceType': beforeSourceType,
      if (beforeSourceRecordedAt != null)
        'beforeSourceRecordedAt': beforeSourceRecordedAt,
      'changes': changes,
    };

Map<String, Object?> textChangeJson(String field, {String? from, String? to}) => {
      'field': field,
      'kind': from != null && to != null
          ? 'CHANGED'
          : (to != null ? 'ADDED' : 'REMOVED'),
      if (from != null) 'from': {'type': 'TEXT', 'value': from},
      if (to != null) 'to': {'type': 'TEXT', 'value': to},
    };
