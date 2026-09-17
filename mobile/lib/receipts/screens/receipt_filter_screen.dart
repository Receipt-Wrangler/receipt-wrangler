import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/add_receipt_filter_sheet.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_filter_condition_card.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_filter_condition_editor.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter_options.dart';

/// Opens the filter screen for the receipts list of [groupId].
///
/// A pushed route rather than a bottom sheet, for two reasons: it keeps the
/// modal stack at the depth the app already ships (editor -> picker, the same as
/// Quick Scan -> category picker), and a real Scaffold puts "Apply Filter" in
/// `bottomNavigationBar`, which reserves its space, instead of fighting
/// `Scaffold.bottomSheet`, which floats over the last card.
Future<void> showReceiptFilterScreen(BuildContext context, String groupId) {
  return Navigator.of(context).push(
    MaterialPageRoute(builder: (_) => ReceiptFilterScreen(groupId: groupId)),
  );
}

class ReceiptFilterScreen extends StatefulWidget {
  const ReceiptFilterScreen({super.key, required this.groupId});

  final String groupId;

  @override
  State<ReceiptFilterScreen> createState() => _ReceiptFilterScreen();
}

class _ReceiptFilterScreen extends State<ReceiptFilterScreen> {
  late final ReceiptListModel _receiptListModel =
      Provider.of<ReceiptListModel>(context, listen: false);

  late final GroupModel _groupModel =
      Provider.of<GroupModel>(context, listen: false);

  /// The conditions being edited.
  ///
  /// A copy of the applied filter, committed only by "Apply Filter" -- so
  /// backing out (the X, the system back gesture) discards rather than
  /// half-applying, and the list never refetches mid-edit.
  late Map<String, ReceiptFilterCondition> _draft =
      Map.of(_receiptListModel.filter);

  /// The fields with a condition, in [receiptFilterFields] order rather than
  /// the order they happen to have been added in, so the screen looks the same
  /// however the user got there.
  List<ReceiptFilterField> get _activeFields => receiptFilterFields
      .where((field) => _draft.containsKey(field.key))
      .toList();

  bool get _showGroupField => isAllGroupId(_groupModel, widget.groupId);

  Future<void> _addCondition() async {
    final key = await showAddReceiptFilterSheet(
      context,
      usedKeys: _draft.keys.toSet(),
      showGroupField: _showGroupField,
    );

    if (key == null || !mounted) {
      return;
    }

    final field = receiptFilterFieldByKey(key);
    if (field != null) {
      await _editCondition(field);
    }
  }

  Future<void> _editCondition(ReceiptFilterField field) async {
    final result = await showReceiptFilterConditionEditor(
      context,
      field: field,
      groupId: widget.groupId,
      existing: _draft[field.key],
    );

    // A dismissed editor returns null, which means "leave the draft alone" --
    // distinct from Remove, which drops the condition.
    if (result == null || !mounted) {
      return;
    }

    setState(() {
      if (result.removed) {
        _draft.remove(field.key);
      } else if (result.condition != null) {
        _draft[field.key] = result.condition!;
      }
    });
  }

  void _removeCondition(ReceiptFilterField field) {
    setState(() => _draft.remove(field.key));
  }

  void _reset() {
    setState(() => _draft = {});
  }

  void _apply() {
    _receiptListModel.setFilter(_draft, true);
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    final activeFields = _activeFields;

    return ScreenWrapper(
      // A slate canvas is what makes the white condition cards read as raised
      // rather than as outlined boxes on the same sheet of paper.
      backgroundColor: slate50,
      appBarWidget: AppBar(
        automaticallyImplyLeading: false,
        leading: IconButton(
          key: const ValueKey("receipt-filter-close"),
          icon: const Icon(Icons.close),
          tooltip: "Close",
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: const Text("Filter"),
        actions: [
          IconButton(
            key: const ValueKey("receipt-filter-reset"),
            icon: const Icon(Icons.restart_alt),
            tooltip: "Reset filter",
            onPressed: _draft.isEmpty ? null : _reset,
          ),
        ],
      ),
      bottomNavigationBarWidget: BottomSubmitButton(
        key: const ValueKey("receipt-filter-apply"),
        buttonText: "Apply Filter",
        onPressed: _apply,
      ),
      bodyPadding: const EdgeInsets.fromLTRB(12, 12, 12, 0),
      // ScreenWrapper hands its child straight to the body with no scrolling of
      // its own, and ten conditions overflow a phone.
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            if (activeFields.isEmpty)
              const Padding(
                key: ValueKey("receipt-filter-empty"),
                padding: EdgeInsets.symmetric(vertical: 36, horizontal: 16),
                child: Text("No conditions yet — everything is showing.",
                    textAlign: TextAlign.center,
                    style: TextStyle(color: slate500)),
              )
            else ...[
              Padding(
                key: const ValueKey("receipt-filter-count"),
                padding: const EdgeInsets.fromLTRB(4, 2, 4, 10),
                child: Text(
                  activeFields.length == 1
                      ? "1 CONDITION"
                      : "${activeFields.length} CONDITIONS",
                  style: const TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    letterSpacing: 0.72,
                    color: slate500,
                  ),
                ),
              ),
              ...activeFields.map((field) => ReceiptFilterConditionCard(
                    key: ValueKey("receipt-filter-card-${field.key}"),
                    field: field,
                    condition: _draft[field.key]!,
                    onTap: () => _editCondition(field),
                    onRemove: () => _removeCondition(field),
                  )),
            ],
            const SizedBox(height: 2),
            _AddFilterButton(
              key: const ValueKey("receipt-filter-add"),
              onPressed: _addCondition,
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }
}

/// The design's dashed "+ Add filter" placeholder.
///
/// Flutter has no dashed border, so the outline is painted: a rounded rect
/// walked with a [PathMetric] and stroked in alternating on/off runs.
class _AddFilterButton extends StatelessWidget {
  const _AddFilterButton({super.key, required this.onPressed});

  final VoidCallback onPressed;

  @override
  Widget build(BuildContext context) {
    // The app's normal blue, not the darker accent: `receipt_form.dart` colours
    // its add icons `primaryColor` and M3 gives every text / outlined button
    // label `colorScheme.primary`, so this is the convention already in place.
    final accent = Theme.of(context).colorScheme.primary;

    return SizedBox(
      height: 52,
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onPressed,
          borderRadius: BorderRadius.circular(14),
          child: CustomPaint(
            painter: const _DashedBorderPainter(
              color: borderSlate,
              radius: 14,
              strokeWidth: 1.5,
              dash: 6,
              gap: 4,
            ),
            child: Center(
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.add, size: 20, color: accent),
                  const SizedBox(width: 8),
                  Text(
                    "Add filter",
                    style: TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: accent,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _DashedBorderPainter extends CustomPainter {
  const _DashedBorderPainter({
    required this.color,
    required this.radius,
    required this.strokeWidth,
    required this.dash,
    required this.gap,
  });

  final Color color;
  final double radius;
  final double strokeWidth;
  final double dash;
  final double gap;

  @override
  void paint(Canvas canvas, Size size) {
    final inset = strokeWidth / 2;
    final path = Path()
      ..addRRect(RRect.fromRectAndRadius(
        Rect.fromLTWH(inset, inset, size.width - strokeWidth,
            size.height - strokeWidth),
        Radius.circular(radius),
      ));

    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth;

    for (final metric in path.computeMetrics()) {
      var distance = 0.0;
      while (distance < metric.length) {
        final end = (distance + dash).clamp(0.0, metric.length);
        canvas.drawPath(metric.extractPath(distance, end), paint);
        distance = end + gap;
      }
    }
  }

  @override
  bool shouldRepaint(_DashedBorderPainter oldDelegate) =>
      oldDelegate.color != color ||
      oldDelegate.radius != radius ||
      oldDelegate.strokeWidth != strokeWidth ||
      oldDelegate.dash != dash ||
      oldDelegate.gap != gap;
}
