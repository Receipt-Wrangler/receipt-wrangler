import 'package:flutter/material.dart';

/// Confirmation dialog for deleting a saved report template. Returns `true` when
/// the user confirms, `null`/`false` otherwise. Mirrors the desktop confirmation
/// copy ("Delete Report Template" / irreversible).
Future<bool?> showReportDeleteDialog(BuildContext context, String name) {
  return showDialog<bool>(
    context: context,
    builder: (context) => AlertDialog(
      title: const Text('Delete Report Template'),
      content: Text(
        'Are you sure you want to delete "$name"? This action is irreversible.',
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context, false),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () => Navigator.pop(context, true),
          style: ElevatedButton.styleFrom(
            backgroundColor: Theme.of(context).colorScheme.error,
            foregroundColor: Theme.of(context).colorScheme.onError,
          ),
          child: const Text('Delete'),
        ),
      ],
    ),
  );
}
