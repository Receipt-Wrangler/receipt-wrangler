import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/service/crash_reporting.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/profile/widgets/delete_account_dialog.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

class UserProfileScreen extends StatelessWidget {
  const UserProfileScreen({super.key});

  Future<void> _handleDeleteAccount(BuildContext context) async {
    final password = await showDeleteAccountDialog(context);
    if (password == null) return;

    try {
      await OpenApiClient.client.getUserApi().deleteAccount(
            deleteAccountCommand: (api.DeleteAccountCommandBuilder()
                  ..password = password)
                .build(),
          );
      final authModel = Provider.of<AuthModel>(context, listen: false);
      await authModel.purgeTokens();

      if (context.mounted) {
        showSuccessSnackbar(context, 'Account deleted successfully.');
        context.go('/login');
      }
    } on DioException catch (e) {
      if (context.mounted) {
        showApiErrorSnackbar(context, e);
      }
    } catch (e) {
      if (context.mounted) {
        showErrorSnackbar(context, 'An error occurred while deleting account.');
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return ScreenWrapper(
      appBarWidget: const TopAppBar(
        titleText: 'User Profile',
        leadingArrowPop: true,
        hideAvatar: true,
      ),
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 24),
            const Text(
              'Privacy',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const _CrashReportingToggle(),
            const SizedBox(height: 24),
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                border: Border.all(
                  color: Theme.of(context).colorScheme.error,
                ),
                borderRadius: BorderRadius.circular(8),
                color: Theme.of(context).colorScheme.error.withValues(alpha: 0.05),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Danger Zone',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                      color: Theme.of(context).colorScheme.error,
                    ),
                  ),
                  const SizedBox(height: 8),
                  const Text(
                    'Account deletion is irreversible. This will permanently remove all associated data including receipts, group memberships, and preferences.',
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () => _handleDeleteAccount(context),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Theme.of(context).colorScheme.error,
                        foregroundColor: Theme.of(context).colorScheme.onError,
                      ),
                      child: const Text('Delete Account'),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Opt-out toggle for crash & error reporting (on by default). Flips Sentry on
/// the fly via [setCrashReportingEnabled] — no restart.
class _CrashReportingToggle extends StatefulWidget {
  const _CrashReportingToggle();

  @override
  State<_CrashReportingToggle> createState() => _CrashReportingToggleState();
}

class _CrashReportingToggleState extends State<_CrashReportingToggle> {
  late bool _enabled = isCrashReportingEnabled();

  @override
  Widget build(BuildContext context) {
    return SwitchListTile(
      contentPadding: EdgeInsets.zero,
      value: _enabled,
      title: const Text('Crash & error reporting'),
      subtitle: const Text(
        'Sends anonymous crash and error reports so bugs can be fixed. No '
        'personal data, no usage tracking, no advertising. You can turn this '
        'off anytime.',
      ),
      onChanged: (value) async {
        setState(() => _enabled = value);
        await setCrashReportingEnabled(value);
      },
    );
  }
}
