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
      ..setBackgroundColor(Colors.white);
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
      await _controller.loadHtmlString(preview.html);
      if (!mounted) return;
      setState(() {
        _loading = false;
        _receiptCount = preview.receiptCount;
      });
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
    if (_loading) {
      return const Center(child: CircularProgressIndicator());
    }
    if (_error != null) {
      return Center(child: Text(_error!));
    }
    return WebViewWidget(controller: _controller);
  }
}
