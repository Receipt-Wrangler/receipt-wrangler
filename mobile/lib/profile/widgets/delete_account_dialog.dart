import 'package:flutter/material.dart';

Future<String?> showDeleteAccountDialog(BuildContext context) {
  return showDialog<String?>(
    context: context,
    builder: (context) => const _DeleteAccountDialog(),
  );
}

class _DeleteAccountDialog extends StatefulWidget {
  const _DeleteAccountDialog();

  @override
  State<_DeleteAccountDialog> createState() => _DeleteAccountDialogState();
}

class _DeleteAccountDialogState extends State<_DeleteAccountDialog> {
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Delete Account'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text.rich(
            TextSpan(
              children: [
                TextSpan(text: 'This action is '),
                TextSpan(
                  text: 'permanent',
                  style: TextStyle(fontWeight: FontWeight.bold),
                ),
                TextSpan(
                  text:
                      ' and cannot be undone. All of your data, including receipts, group memberships, and preferences will be deleted.',
                ),
              ],
            ),
          ),
          const SizedBox(height: 12),
          const Text('Please enter your password to confirm account deletion.'),
          const SizedBox(height: 16),
          TextField(
            // Same key the login form's password field uses -- the two never
            // share a tree, and matching them keeps both password fields
            // addressable the same way, as their eye toggles already are.
            key: const ValueKey('password-field'),
            controller: _passwordController,
            obscureText: _obscurePassword,
            decoration: InputDecoration(
              labelText: 'Password',
              suffixIcon: IconButton(
                key: const ValueKey('password-visibility-toggle'),
                icon: Icon(
                  _obscurePassword ? Icons.visibility : Icons.visibility_off,
                ),
                tooltip: _obscurePassword ? 'Show Password' : 'Hide Password',
                onPressed: () {
                  setState(() {
                    _obscurePassword = !_obscurePassword;
                  });
                },
              ),
            ),
            onChanged: (_) => setState(() {}),
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context, null),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: _passwordController.text.isEmpty
              ? null
              : () => Navigator.pop(context, _passwordController.text),
          style: ElevatedButton.styleFrom(
            backgroundColor: Theme.of(context).colorScheme.error,
            foregroundColor: Theme.of(context).colorScheme.onError,
          ),
          child: const Text('Delete Account'),
        ),
      ],
    );
  }
}
