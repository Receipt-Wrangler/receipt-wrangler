import 'package:flutter/material.dart';

class ContextModel extends ChangeNotifier {
  BuildContext? _shellContext;

  BuildContext? get shellContext => _shellContext;

  void setShellContext(BuildContext context) {
    _shellContext = context;
  }

  /// Resolves a live [BuildContext] for opening a modal sheet. Prefers the
  /// cached shell context when it is set and still mounted, otherwise falls
  /// back to [fallback] (the caller's own context).
  ///
  /// The shell context is only ever populated by the receipt-form screen, and
  /// even there it can refer to a deactivated element after navigation. Flows
  /// that never mount that screen -- notably **Quick Scan**, which opens its
  /// sheet straight from the bottom-nav Add menu -- leave it null, so passing
  /// it directly to `Navigator.of(...)` throws. Always route sheet-opening
  /// through this helper.
  BuildContext resolveSheetContext(BuildContext fallback) {
    final shell = _shellContext;
    return (shell != null && shell.mounted) ? shell : fallback;
  }
}
