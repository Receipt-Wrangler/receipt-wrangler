import 'package:flutter/material.dart';

import '../../constants/receipt_entry.dart';
import '../../shared/functions/receipt_entry_availability.dart';

/// Explains why tapping Scan landed on the manual form instead of the camera.
///
/// Shown only when the manual form was reached *that* way — a deliberate
/// "Add Manual Receipt" passes no reason, so it never appears there. Dismissible
/// because the reason is usually install-wide and the user can do nothing about
/// it themselves.
class QuickScanUnavailableBanner extends StatefulWidget {
  const QuickScanUnavailableBanner({
    super.key,
    required this.reason,
    this.groupName,
  });

  final QuickScanBlockedReason reason;

  /// The group the missing permission applies to, when the form was opened from
  /// inside a single group.
  final String? groupName;

  /// Reads the reason a receipt-entry navigation attached to the route, or null
  /// when the form was opened some other way.
  static QuickScanUnavailableBanner? fromRouteExtra(Object? extra) {
    if (extra is! Map) {
      return null;
    }
    final reason = extra[quickScanBlockedReasonExtraKey];
    if (reason is! QuickScanBlockedReason) {
      return null;
    }
    final groupName = extra[quickScanBlockedGroupExtraKey];
    return QuickScanUnavailableBanner(
      reason: reason,
      groupName: groupName is String ? groupName : null,
    );
  }

  @override
  State<QuickScanUnavailableBanner> createState() =>
      _QuickScanUnavailableBanner();
}

class _QuickScanUnavailableBanner extends State<QuickScanUnavailableBanner> {
  bool _dismissed = false;

  String get _message {
    if (widget.reason == QuickScanBlockedReason.aiDisabled) {
      return quickScanAiDisabledMessage;
    }
    final groupName = widget.groupName;
    return groupName == null
        ? quickScanNoPermissionMessage
        : quickScanNoPermissionMessageForGroup(groupName);
  }

  @override
  Widget build(BuildContext context) {
    if (_dismissed) {
      return const SizedBox.shrink();
    }

    final colors = Theme.of(context).colorScheme;

    return Container(
      key: const ValueKey("quick-scan-unavailable-banner"),
      margin: const EdgeInsets.only(top: 12, bottom: 4),
      padding: const EdgeInsets.fromLTRB(14, 12, 4, 12),
      decoration: BoxDecoration(
        color: colors.secondaryContainer,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(Icons.info_outline, size: 20, color: colors.onSecondaryContainer),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  quickScanUnavailableTitle,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        color: colors.onSecondaryContainer,
                      ),
                ),
                const SizedBox(height: 3),
                Text(
                  _message,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: colors.onSecondaryContainer,
                      ),
                ),
              ],
            ),
          ),
          IconButton(
            key: const ValueKey("quick-scan-unavailable-banner-dismiss"),
            icon: const Icon(Icons.close, size: 18),
            color: colors.onSecondaryContainer,
            tooltip: "Dismiss",
            onPressed: () => setState(() => _dismissed = true),
          ),
        ],
      ),
    );
  }
}
