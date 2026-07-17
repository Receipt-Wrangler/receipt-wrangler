import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

import '../functions/report_actions.dart';
import '../screens/report_preview_screen.dart';

/// A single saved report template row. The per-row action buttons are gated
/// **only** on the server-computed `allowedActions` array (read/generate/delete),
/// mirroring the desktop list — they must NOT be additionally AND-ed with a
/// client permission check, since `allowedActions` already bakes in the base /
/// `*All` report permissions, the per-group ceiling, and the per-template matrix.
class ReportListItem extends StatelessWidget {
  const ReportListItem({
    super.key,
    required this.template,
    required this.onChanged,
  });

  final api.ReportTemplate template;

  /// Called after a mutating action (delete) so the list can refresh.
  final VoidCallback onChanged;

  bool _can(String action) => template.allowedActions?.contains(action) ?? false;

  String get _subtitle {
    final formats =
        template.configuration.formats.map((f) => f.name.toUpperCase()).join(', ');
    return formats.isEmpty ? 'No formats' : formats;
  }

  @override
  Widget build(BuildContext context) {
    return ListTile(
      key: ValueKey('report-${template.id}'),
      title: Text(template.name),
      subtitle: Text(_subtitle),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (_can('read'))
            IconButton(
              key: ValueKey('report-preview-${template.id}'),
              icon: const Icon(Icons.visibility),
              tooltip: 'Preview',
              onPressed: () => Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (_) => ReportPreviewScreen(template: template),
                ),
              ),
            ),
          if (_can('generate'))
            IconButton(
              key: ValueKey('report-generate-${template.id}'),
              icon: const Icon(Icons.play_arrow),
              tooltip: 'Generate report',
              onPressed: () => generateAndSaveReport(context, template),
            ),
          if (_can('delete'))
            IconButton(
              key: ValueKey('report-delete-${template.id}'),
              icon: const Icon(Icons.delete),
              tooltip: 'Delete',
              onPressed: () =>
                  confirmAndDeleteReport(context, template, onChanged),
            ),
        ],
      ),
    );
  }
}
