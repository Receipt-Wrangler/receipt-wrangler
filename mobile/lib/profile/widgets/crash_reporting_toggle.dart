import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/service/crash_reporting.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

/// Opt-out toggle for crash & error reporting (on by default). Flips Sentry on
/// the fly via [setCrashReportingEnabled] — no restart.
///
/// [readEnabled] / [setEnabled] default to the real functions; they exist as
/// seams so the widget can be tested without touching the native Sentry SDK.
class CrashReportingToggle extends StatefulWidget {
  const CrashReportingToggle({
    super.key,
    this.readEnabled = isCrashReportingEnabled,
    this.setEnabled = setCrashReportingEnabled,
  });

  final bool Function() readEnabled;
  final Future<void> Function(bool) setEnabled;

  @override
  State<CrashReportingToggle> createState() => _CrashReportingToggleState();
}

class _CrashReportingToggleState extends State<CrashReportingToggle> {
  late bool _enabled = widget.readEnabled();
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    return SwitchListTile(
      contentPadding: EdgeInsets.zero,
      value: _enabled,
      // Disabled while a toggle is in flight so a rapid double-tap can't start a
      // second concurrent SDK init/close (a null onChanged disables the switch).
      onChanged: _isLoading ? null : _onChanged,
      title: const Text('Crash & error reporting'),
      subtitle: const Text(
        'Sends anonymous crash and error reports so bugs can be fixed. No '
        'personal data, no usage tracking, no advertising. You can turn this '
        'off anytime.',
      ),
    );
  }

  Future<void> _onChanged(bool value) async {
    final previous = _enabled;
    setState(() {
      _enabled = value;
      _isLoading = true;
    });
    try {
      await widget.setEnabled(value);
    } catch (_) {
      // Toggling the SDK failed — revert the switch so it keeps matching the
      // persisted preference / actual SDK state, and surface the error.
      if (mounted) {
        setState(() => _enabled = previous);
        showErrorSnackbar(context, "Couldn't update crash reporting setting");
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }
}
