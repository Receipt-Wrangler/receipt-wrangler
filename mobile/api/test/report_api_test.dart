import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ReportApi
void main() {
  final instance = Openapi().getReportApi();

  group(ReportApi, () {
    // Save a report template
    //
    // Saves the current report configuration as a reusable template. App-scoped (requires app.reports.create) — it persists a configuration and touches no group's receipts.
    //
    //Future<ReportTemplate> createReportTemplate(ReportRequestCommand reportRequestCommand) async
    test('test createReportTemplate', () async {
      // TODO
    });

    // Delete a report template
    //
    // Deletes a saved report template by id. Requires app.reports.delete.
    //
    //Future deleteReportTemplate(int id) async
    test('test deleteReportTemplate', () async {
      // TODO
    });

    // Duplicate a report template
    //
    // Duplicates a saved report template by id, returning the new copy. App-scoped (requires app.reports.duplicate).
    //
    //Future<ReportTemplate> duplicateReportTemplate(int id) async
    test('test duplicateReportTemplate', () async {
      // TODO
    });

    // Generate a report
    //
    // Generates a report over the receipts of one or more groups and streams it back as a file, or a zip when several formats are requested. Runs under the caller's group permissions and requires group.reports.read in every covered group.
    //
    //Future<Uint8List> generateReport(ReportRequestCommand reportRequestCommand) async
    test('test generateReport', () async {
      // TODO
    });

    // Generate a report from a saved template
    //
    // Generates a saved report template by id and streams the file. Enforces the per-template generate grant (generateAll, or the per-group ceiling plus the matrix), loading the stored configuration server-side.
    //
    //Future<Uint8List> generateReportFromTemplate(int id) async
    test('test generateReportFromTemplate', () async {
      // TODO
    });

    // Get a report template
    //
    // Returns a saved report template by id. Access is resolved per template (readAll, or the per-group ceiling plus the per-template matrix).
    //
    //Future<ReportTemplate> getReportTemplate(int id) async
    test('test getReportTemplate', () async {
      // TODO
    });

    // Get report template options
    //
    // Returns every report template as a lightweight {id, name, groupIds} option for the role-form access matrix. Gated on app.roles.read (the admin role editor may not personally hold report permissions), not a report permission.
    //
    //Future<BuiltList<ReportTemplateOption>> getReportTemplateOptions() async
    test('test getReportTemplateOptions', () async {
      // TODO
    });

    // Get paged report templates
    //
    // Returns a paged, sorted list of saved report templates. App-scoped (requires app.reports.read).
    //
    //Future<PagedData> getReportTemplates(PagedRequestCommand pagedRequestCommand) async
    test('test getReportTemplates', () async {
      // TODO
    });

    // Preview a report
    //
    // Renders the current report configuration as HTML for the builder's live preview, along with the number of receipts it covers. Requires the app-level app.reports.read (or app.reports.readAll) and group.reports.read in every covered group. The preview is a row-capped sample of the engine's own output.
    //
    //Future<ReportPreviewResponse> previewReport(ReportRequestCommand reportRequestCommand) async
    test('test previewReport', () async {
      // TODO
    });

    // Update a report template
    //
    // Updates a saved report template in place, replacing its name and stored configuration. App-scoped (requires app.reports.update).
    //
    //Future<ReportTemplate> updateReportTemplate(int id, ReportRequestCommand reportRequestCommand) async
    test('test updateReportTemplate', () async {
      // TODO
    });

  });
}
