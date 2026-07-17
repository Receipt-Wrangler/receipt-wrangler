import 'package:openapi/openapi.dart' as api;

/// Derives the download filename for a generated report from its stored
/// configuration. Faithful port of the desktop client's `reportFilename`
/// (`desktop/src/reports/services/report-runner.service.ts`): the sanitized
/// report name, then `.zip` when the template outputs multiple formats (the
/// backend zips them) or `.<format>` for a single format.
String reportFilename(api.ReportRequestCommand configuration) {
  final sanitized = (configuration.name ?? '')
      .trim()
      .replaceAll(RegExp(r'[^\w]+'), '_')
      .replaceAll(RegExp(r'^_+|_+$'), '');
  final base = sanitized.isEmpty ? 'report' : sanitized;

  final formats = configuration.formats;
  if (formats.length > 1) {
    return '$base.zip';
  }
  final ext = formats.isNotEmpty ? formats.first.name : 'csv';
  return '$base.$ext';
}
