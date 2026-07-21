import 'package:built_collection/built_collection.dart';
import 'package:openapi/openapi.dart' as api;

/// Builds a minimal-but-valid [api.ReportRequestCommand] for tests. Only [name]
/// and [formats] usually matter to the code under test; the other required
/// fields (groupIds/period/detail) are set so `build()` succeeds.
api.ReportRequestCommand buildReportRequestCommand({
  String? name = 'Test Report',
  List<api.ReportRequestCommandFormatsEnum> formats = const [
    api.ReportRequestCommandFormatsEnum.csv,
  ],
}) {
  return api.ReportRequestCommand((b) => b
    ..name = name
    ..groupIds = ListBuilder<String>(<String>['1'])
    ..period.preset = api.ReportPeriodPresetEnum.thisMonth
    ..detail.mode = api.ReportDetailModeEnum.records
    ..formats =
        ListBuilder<api.ReportRequestCommandFormatsEnum>(formats));
}

/// Builds an [api.ReportTemplate] with the given [allowedActions] (the
/// server-computed per-row action set) and format list, for widget tests.
api.ReportTemplate buildReportTemplate({
  int id = 1,
  String name = 'Test Report',
  List<String>? allowedActions,
  List<api.ReportRequestCommandFormatsEnum> formats = const [
    api.ReportRequestCommandFormatsEnum.csv,
  ],
}) {
  return api.ReportTemplate((b) {
    b
      ..id = id
      ..createdAt = '2024-01-01T00:00:00Z'
      ..configurationVersion = 1
      ..name = name
      ..configuration
          .replace(buildReportRequestCommand(name: name, formats: formats));
    if (allowedActions != null) {
      b.allowedActions = ListBuilder<String>(allowedActions);
    }
  });
}
