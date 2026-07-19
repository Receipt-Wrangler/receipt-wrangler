import 'package:built_collection/built_collection.dart';
import 'package:built_value/json_object.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/reports/functions/report_actions.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';
import 'package:webview_flutter/webview_flutter.dart';

import '../../../client/client.dart';
import '../constants/text_styles.dart';

/// Reads the pinned report template id from a dashboard widget's untyped
/// `configuration` blob (`{ reportTemplateId: <int> }`, authored on desktop).
/// Returns null when the key is absent or not a whole number — a malformed
/// fractional value must not silently truncate to (and load) the wrong
/// template; the widget then surfaces an error state.
int? reportTemplateIdFromConfig(BuiltMap<String, JsonObject?>? config) {
  final value = config?['reportTemplateId'];
  if (value == null) return null;
  try {
    final id = value.asNum;
    if (id % 1 != 0) return null; // reject fractional ids
    return id.toInt();
  } catch (_) {
    return null;
  }
}

/// Whether the widget should show its Download button, gated **only** on the
/// server-computed [allowedActions] (never AND-ed with a client permission
/// check) — the same contract `report_list_item.dart` uses. Download reuses the
/// enforcing generate path, so `generate` is the gating action.
bool reportWidgetCanDownload(BuiltList<String>? allowedActions) =>
    allowedActions?.contains('generate') ?? false;

/// View-only dashboard widget that pins a saved report template and renders it,
/// mirroring the desktop report widget (`desktop/src/dashboard/report-widget/`).
/// It asks the server to render the whole template
/// (`POST /report/template/{id}/render`) and drops the returned self-contained
/// HTML into a WebView — the same rendering path as the report preview screen.
/// A Download button is shown **only** when the server's `allowedActions`
/// includes `generate` (never AND-ed with a client permission check), matching
/// desktop and `report_list_item.dart`.
class ReportWidget extends StatefulWidget {
  const ReportWidget({super.key, required this.dashboardWidget});

  final api.Widget dashboardWidget;

  @override
  State<ReportWidget> createState() => _ReportWidgetState();
}

class _ReportWidgetState extends State<ReportWidget> {
  static const _errorMessage = "Couldn't load this report.";

  late final WebViewController _controller;
  int? _templateId;
  bool _loading = true;
  bool _downloading = false;
  String? _error;
  BuiltList<String>? _allowedActions;

  bool get _canDownload => reportWidgetCanDownload(_allowedActions);

  @override
  void initState() {
    super.initState();
    _templateId = reportTemplateIdFromConfig(widget.dashboardWidget.configuration);
    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.disabled)
      ..setBackgroundColor(Colors.white)
      ..setNavigationDelegate(NavigationDelegate(
        // Clear the spinner only once the page has finished painting, not when
        // loadHtmlString returns (which resolves before the WebView renders) —
        // otherwise the report flashes blank while it lays out.
        onPageFinished: (_) {
          if (mounted) setState(() => _loading = false);
        },
      ));
    _load();
  }

  Future<void> _load() async {
    final id = _templateId;
    if (id == null) {
      setState(() {
        _loading = false;
        _error = _errorMessage;
      });
      return;
    }
    try {
      final response =
          await OpenApiClient.client.getReportApi().renderReportTemplate(id: id);
      if (!mounted) return;
      final data = response.data;
      setState(() => _allowedActions = data?.allowedActions);
      // A revoked/deleted template comes back as restricted-notice HTML at 200
      // with empty allowedActions, so we always have HTML to render and never
      // need a client-side "restricted" branch. _loading is cleared by the
      // controller's onPageFinished once the HTML has rendered.
      await _controller.loadHtmlString(data?.html ?? '');
    } on DioException catch (e) {
      if (!mounted) return;
      setState(() {
        _loading = false;
        _error = _errorMessage;
      });
      showApiErrorSnackbar(context, e);
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _loading = false;
        _error = _errorMessage;
      });
    }
  }

  /// Fetches the full template (for its filename/config) then reuses the shared
  /// generate-and-share helper — the same path desktop's `downloadTemplateById`
  /// takes (get template → generate). The gate above is server-authoritative;
  /// the generate call re-enforces per-template access, so a stale button at
  /// worst 403s.
  Future<void> _download(int id) async {
    if (_downloading) return;
    setState(() => _downloading = true);
    try {
      final template =
          (await OpenApiClient.client.getReportApi().getReportTemplate(id: id))
              .data;
      if (template != null && mounted) {
        await generateAndSaveReport(context, template);
      }
    } on DioException catch (e) {
      if (mounted) showApiErrorSnackbar(context, e);
    } finally {
      if (mounted) setState(() => _downloading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 10),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Expanded(
              child: Text(
                widget.dashboardWidget.name ?? 'Report',
                style: dashboardWidgetNameStyle,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            if (_canDownload)
              IconButton(
                key: const ValueKey('report-widget-download'),
                icon: const Icon(Icons.download),
                tooltip: 'Download',
                onPressed:
                    _downloading ? null : () => _download(_templateId!),
              ),
          ],
        ),
        const SizedBox(height: 10),
        Expanded(child: _buildContent()),
      ],
    );
  }

  Widget _buildContent() {
    if (_error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 32),
            const SizedBox(height: 8),
            Text(_error!),
          ],
        ),
      );
    }
    // The WebView stays mounted so it can load/paint while the spinner overlays
    // it; the spinner is cleared by onPageFinished once the content is visible.
    return Stack(
      children: [
        WebViewWidget(controller: _controller),
        if (_loading) const Center(child: CircularProgressIndicator()),
      ],
    );
  }
}
