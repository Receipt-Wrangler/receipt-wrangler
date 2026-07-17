import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:path_provider/path_provider.dart';
import 'package:provider/provider.dart';
import 'package:share_plus/share_plus.dart';

import '../../client/client.dart';
import '../../models/loading_model.dart';
import '../../utils/snackbar.dart';
import '../widgets/report_delete_dialog.dart';
import 'report_filename.dart';

/// Fetches the HTML preview for [template] (`POST /report/preview`) using its
/// stored configuration. Mirrors desktop's read-scoped preview — the response is
/// a self-contained, row-capped HTML sample plus a receipt count.
Future<api.ReportPreviewResponse?> fetchReportPreview(
    api.ReportTemplate template) async {
  final response = await OpenApiClient.client
      .getReportApi()
      .previewReport(reportRequestCommand: template.configuration);
  return response.data;
}

/// Generates [template] in its saved formats (`POST /report/template/{id}/generate`,
/// which loads the config server-side) and opens the OS share / "Save to Files"
/// sheet with the downloaded file. Multiple formats come back zipped; the filename
/// is derived to match desktop.
Future<void> generateAndSaveReport(
    BuildContext context, api.ReportTemplate template) async {
  final loadingModel = Provider.of<LoadingModel>(context, listen: false);
  loadingModel.setIsLoading(true);
  try {
    final response = await OpenApiClient.client
        .getReportApi()
        .generateReportFromTemplate(id: template.id);
    final bytes = response.data;
    if (bytes == null) {
      loadingModel.setIsLoading(false);
      if (context.mounted) {
        showErrorSnackbar(context, 'Report generation returned no data');
      }
      return;
    }

    final filename = reportFilename(template.configuration);
    final dir = await getTemporaryDirectory();
    final file = File('${dir.path}/$filename');
    await file.writeAsBytes(bytes);

    loadingModel.setIsLoading(false);
    await SharePlus.instance.share(ShareParams(files: [XFile(file.path)]));
  } on DioException catch (e) {
    loadingModel.setIsLoading(false);
    if (context.mounted) showApiErrorSnackbar(context, e);
  } catch (e) {
    loadingModel.setIsLoading(false);
    if (context.mounted) showErrorSnackbar(context, 'Failed to generate report');
  }
}

/// Prompts for confirmation, then deletes [template]
/// (`DELETE /report/template/{id}`). On success shows a snackbar and calls
/// [onDeleted] so the caller can refresh the list.
Future<void> confirmAndDeleteReport(
  BuildContext context,
  api.ReportTemplate template,
  VoidCallback onDeleted,
) async {
  final confirmed = await showReportDeleteDialog(context, template.name);
  if (confirmed != true) return;

  try {
    await OpenApiClient.client
        .getReportApi()
        .deleteReportTemplate(id: template.id);
    if (context.mounted) {
      showSuccessSnackbar(context, 'Template deleted');
    }
    onDeleted();
  } on DioException catch (e) {
    if (context.mounted) showApiErrorSnackbar(context, e);
  }
}
