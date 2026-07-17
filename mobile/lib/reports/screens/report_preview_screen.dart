import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:webview_flutter/webview_flutter.dart';

import '../../utils/snackbar.dart';
import '../functions/report_actions.dart';

/// Renders the HTML preview of a report template (mirroring the desktop
/// builder's live preview) in a WebView. Pushed via `Navigator.push`; the HTML
/// is a self-contained, row-capped sample loaded from `POST /report/preview`, so
/// no network access is needed to render it.
class ReportPreviewScreen extends StatefulWidget {
  const ReportPreviewScreen({super.key, required this.template});

  final api.ReportTemplate template;

  @override
  State<ReportPreviewScreen> createState() => _ReportPreviewScreenState();
}

class _ReportPreviewScreenState extends State<ReportPreviewScreen> {
  late final WebViewController _controller;
  bool _loading = true;
  String? _error;
  int? _receiptCount;

  @override
  void initState() {
    super.initState();
    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.disabled)
      ..setBackgroundColor(Colors.white)
      ..setNavigationDelegate(NavigationDelegate(
        // Clear the spinner only once the page has actually finished painting,
        // not when loadHtmlString returns (which resolves before the WebView has
        // rendered) — otherwise the preview flashes blank while it lays out.
        onPageFinished: (_) {
          if (mounted) setState(() => _loading = false);
        },
      ));
    _load();
  }

  Future<void> _load() async {
    try {
      final preview = await fetchReportPreview(widget.template);
      if (!mounted) return;
      if (preview == null) {
        setState(() {
          _loading = false;
          _error = 'No preview available';
        });
        return;
      }
      setState(() => _receiptCount = preview.receiptCount);
      // _loading is cleared by the controller's onPageFinished callback once the
      // HTML has rendered, so the spinner covers both the fetch and the paint.
      await _controller.loadHtmlString(preview.html);
    } on DioException catch (e) {
      if (!mounted) return;
      setState(() {
        _loading = false;
        _error = 'Unable to load preview';
      });
      showApiErrorSnackbar(context, e);
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _loading = false;
        _error = 'Unable to load preview';
      });
    }
  }

  String get _title {
    if (_receiptCount != null) {
      final noun = _receiptCount == 1 ? 'receipt' : 'receipts';
      return 'Preview · $_receiptCount $noun';
    }
    return 'Preview';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(_title)),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_error != null) {
      return Center(child: Text(_error!));
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
