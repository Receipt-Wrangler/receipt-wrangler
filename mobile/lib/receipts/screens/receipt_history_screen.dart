import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';

import '../../client/client.dart';
import '../../constants/colors.dart';
import '../../models/user_model.dart';
import '../../shared/widgets/screen_wrapper.dart';
import '../../utils/currency.dart';
import '../../utils/receipt_history.dart';
import '../widgets/receipt_history_app_bar.dart';
import '../widgets/receipt_history_entry_card.dart';

/// A receipt's change history, newest first, grouped by day.
///
/// Pushed with `Navigator.push` from the receipt screen's ⋮ menu, never routed:
/// the receipt form derives view/edit from the live URL, so a route of its own
/// would flip an open edit form underneath to view. See
/// `ReceiptAppBarActionBuilder`.
///
/// The server builds the history, computes the diff and redacts it to the
/// viewer; this screen only formats it (`utils/receipt_history.dart`).
class ReceiptHistoryScreen extends StatefulWidget {
  const ReceiptHistoryScreen({super.key, required this.receiptId});

  final int receiptId;

  static const loadingKey = ValueKey('receipt-history-loading');

  static const errorKey = ValueKey('receipt-history-error');

  static const retryKey = ValueKey('receipt-history-retry');

  static const emptyKey = ValueKey('receipt-history-empty');

  static const dayKey = ValueKey('receipt-history-day');

  @override
  State<ReceiptHistoryScreen> createState() => _ReceiptHistoryScreenState();
}

class _ReceiptHistoryScreenState extends State<ReceiptHistoryScreen> {
  late Future<api.ReceiptHistory> _history;

  @override
  void initState() {
    super.initState();
    _history = _load();
  }

  Future<api.ReceiptHistory> _load() {
    return OpenApiClient.client
        .getReceiptApi()
        .getReceiptHistory(receiptId: widget.receiptId)
        .then((response) => response.data!);
  }

  // Block bodies on purpose: an arrow would return the assigned Future from
  // the setState callback, which Flutter rejects.
  void _reload() {
    final next = _load();
    setState(() {
      _history = next;
    });
  }

  // RefreshIndicator spins until this completes; a failure shows in the
  // builder's error state rather than escaping as an unhandled error here.
  Future<void> _refresh() async {
    final next = _load();
    setState(() {
      _history = next;
    });
    try {
      await next;
    } catch (_) {}
  }

  @override
  Widget build(BuildContext context) {
    final userModel = Provider.of<UserModel>(context);
    final formatters = HistoryValueFormatters(
      currency: (value) => formatCurrency(context, value) ?? value,
      findUser: (userId) =>
          userId == null ? null : userModel.getUserById(userId.toString()),
    );

    return FutureBuilder<api.ReceiptHistory>(
      future: _history,
      builder: (context, snapshot) {
        final history = snapshot.data;
        return ScreenWrapper(
          backgroundColor: slate50,
          appBarWidget: ReceiptHistoryAppBar(
            onBack: () => Navigator.of(context).pop(),
            summary: history == null ? null : historySummary(history),
          ),
          child: _buildBody(snapshot, formatters),
        );
      },
    );
  }

  Widget _buildBody(
    AsyncSnapshot<api.ReceiptHistory> snapshot,
    HistoryValueFormatters formatters,
  ) {
    if (snapshot.hasError) {
      return Center(
        key: ReceiptHistoryScreen.errorKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text("Couldn't load this receipt's history."),
            TextButton(
              key: ReceiptHistoryScreen.retryKey,
              onPressed: _reload,
              child: const Text('Try again'),
            ),
          ],
        ),
      );
    }

    final history = snapshot.data;
    if (snapshot.connectionState != ConnectionState.done || history == null) {
      return const Center(
        key: ReceiptHistoryScreen.loadingKey,
        child: CircularProgressIndicator(),
      );
    }

    if (history.entries.isEmpty) {
      return const Center(
        key: ReceiptHistoryScreen.emptyKey,
        child: Text('No changes have been recorded for this receipt yet.'),
      );
    }

    final days = groupHistoryByDay(history.entries);
    return RefreshIndicator(
      onRefresh: _refresh,
      child: ListView(
        padding: const EdgeInsets.only(bottom: 24),
        children: [
          for (final day in days) ...[
            Padding(
              padding: const EdgeInsets.fromLTRB(4, 18, 4, 10),
              child: Text(
                day.label.toUpperCase(),
                key: ReceiptHistoryScreen.dayKey,
                style: const TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  letterSpacing: 0.78,
                  color: slate500,
                ),
              ),
            ),
            for (final entry in day.entries)
              Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: ReceiptHistoryEntryCard(
                  key: ValueKey('receipt-history-entry-${entry.systemTaskId}'),
                  entry: entry,
                  formatters: formatters,
                ),
              ),
          ],
        ],
      ),
    );
  }
}
