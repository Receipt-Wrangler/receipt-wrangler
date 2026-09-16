import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/multi-select-field.dart';

/// Stand-in for the api models the real fields carry (`api.Category` /
/// `api.Tag`), so the generic parameter and `itemDisplayName` are genuinely
/// exercised rather than defaulted to a String.
class _Thing {
  const _Thing(this.name);

  final String name;
}

void main() {
  const fieldKey = ValueKey('multi-select-field');
  const fieldName = 'things';

  const alpha = _Thing('Alpha');
  const beta = _Thing('Beta');

  /// Rebuilds the field's owner. The real consumers answer a value change with
  /// `setState(() => ...setValue(list))`: FormBuilder's `setValue` does not
  /// call `setState` itself (only `didChange` does), so it is the owner's
  /// rebuild that re-runs the field's builder against the new value.
  late StateSetter rebuildOwner;

  /// Mounts a [MultiSelectField] inside a FormBuilder, mirroring how the
  /// receipt / quick-scan forms host it. The surrounding SizedBox pins the
  /// field's width so offset-based taps are computed against a stable rect.
  Future<GlobalKey<FormBuilderState>> pumpField(
    WidgetTester tester, {
    List<_Thing>? initialValue,
    VoidCallback? onTap,
    void Function(List<_Thing>)? onRemove,
    bool? required,
  }) async {
    final formKey = GlobalKey<FormBuilderState>();

    await tester.pumpWidget(
      MaterialApp(
        // Matches the app's global decoration theme (lib/main.dart), which is
        // what gives the field its full-width bordered box.
        theme: ThemeData(
          inputDecorationTheme: const InputDecorationTheme(
            border: OutlineInputBorder(),
          ),
        ),
        home: Scaffold(
          body: Center(
            child: SizedBox(
              width: 400,
              child: FormBuilder(
                key: formKey,
                child: StatefulBuilder(
                  builder: (context, setState) {
                    rebuildOwner = setState;
                    return MultiSelectField<_Thing>(
                      key: fieldKey,
                      name: fieldName,
                      label: 'Categories',
                      itemName: 'Categories',
                      itemDisplayName: (thing) => thing.name,
                      initialValue: initialValue,
                      onTap: onTap,
                      onRemove: onRemove,
                      required: required,
                    );
                  },
                ),
              ),
            ),
          ),
        ),
      ),
    );
    await tester.pump();

    return formKey;
  }

  Finder ownGestureDetector() => find.descendant(
        of: find.byKey(fieldKey),
        matching: find.byType(GestureDetector),
      );

  Rect fieldRect(WidgetTester tester) => tester.getRect(
        find.descendant(
          of: find.byKey(fieldKey),
          matching: find.byType(InputDecorator),
        ),
      );

  /// The X on a named chip. Icons.cancel is the chip's default delete icon,
  /// and the glyph desktop's matChipRemove button uses.
  Finder removeButton(String label) => find.descendant(
        of: find.widgetWithText(InputChip, label),
        matching: find.byIcon(Icons.cancel),
      );

  group('tap target', () {
    testWidgets('the whole field is one opaque tap target', (tester) async {
      await pumpField(tester, onTap: () {});

      final detector = tester.widget<GestureDetector>(ownGestureDetector());
      expect(detector.behavior, HitTestBehavior.opaque);

      // The detector must cover the decorated box, not just its inner Wrap.
      expect(
        tester.getRect(ownGestureDetector()),
        fieldRect(tester),
      );
    });

    testWidgets('a tap near the right edge opens the picker', (tester) async {
      var taps = 0;
      await pumpField(tester, onTap: () => taps++);

      final rect = fieldRect(tester);
      await tester.tapAt(Offset(rect.right - 4, rect.center.dy));
      await tester.pump();

      expect(taps, 1);
    });

    testWidgets('a tap on the label row opens the picker', (tester) async {
      var taps = 0;
      await pumpField(tester, onTap: () => taps++);

      final rect = fieldRect(tester);
      await tester.tapAt(Offset(rect.left + 12, rect.top + 4));
      await tester.pump();

      expect(taps, 1);
    });

    testWidgets('a tap in the bottom-left padding gutter opens the picker',
        (tester) async {
      var taps = 0;
      await pumpField(tester, onTap: () => taps++);

      final rect = fieldRect(tester);
      await tester.tapAt(Offset(rect.left + 4, rect.bottom - 4));
      await tester.pump();

      expect(taps, 1);
    });

    testWidgets('the edges stay tappable once chips are rendered',
        (tester) async {
      var taps = 0;
      await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () => taps++,
      );

      final rect = fieldRect(tester);
      await tester.tapAt(Offset(rect.right - 4, rect.center.dy));
      await tester.pump();
      await tester.tapAt(Offset(rect.left + 4, rect.bottom - 4));
      await tester.pump();

      expect(taps, 2);
    });

    testWidgets('a tap in the gap between two chips opens the picker',
        (tester) async {
      var taps = 0;
      await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () => taps++,
      );

      final first = tester.getRect(find.widgetWithText(InputChip, 'Alpha'));
      final second = tester.getRect(find.widgetWithText(InputChip, 'Beta'));
      // The chips are separated by a childless SizedBox spacer, which the old
      // deferToChild detector could not hit-test.
      expect(second.left, greaterThan(first.right));

      await tester.tapAt(
        Offset((first.right + second.left) / 2, first.center.dy),
      );
      await tester.pump();

      expect(taps, 1);
    });

    testWidgets('a tap on the placeholder text still opens the picker',
        (tester) async {
      // The e2e specs locate the field by this text and tap it directly.
      var taps = 0;
      await pumpField(tester, onTap: () => taps++);

      await tester.tap(find.text('No Categories selected'));
      await tester.pump();

      expect(taps, 1);
    });

    testWidgets('a tap on a chip still opens the picker', (tester) async {
      var taps = 0;
      await pumpField(
        tester,
        initialValue: const [alpha],
        onTap: () => taps++,
      );

      await tester.tap(find.widgetWithText(InputChip, 'Alpha'));
      await tester.pumpAndSettle();

      expect(taps, 1);
    });
  });

  group('remove button', () {
    testWidgets('every chip carries one when onRemove is given',
        (tester) async {
      await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () {},
        onRemove: (_) {},
      );

      expect(find.byIcon(Icons.cancel), findsNWidgets(2));
    });

    testWidgets('there is none without onRemove', (tester) async {
      await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () {},
      );

      expect(find.byIcon(Icons.cancel), findsNothing);
    });

    testWidgets('tapping one reports the remaining values in order',
        (tester) async {
      List<_Thing>? remaining;
      await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () {},
        onRemove: (value) => remaining = value,
      );

      await tester.tap(removeButton('Alpha'));
      await tester.pump();

      expect(remaining, const [beta]);
    });

    testWidgets('tapping one does not also open the picker', (tester) async {
      // The remove button is an InkWell nested inside the field's opaque
      // GestureDetector; the inner recognizer has to win the gesture arena, or
      // removing a chip would open the sheet on top of it.
      var taps = 0;
      await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () => taps++,
        onRemove: (_) {},
      );

      await tester.tap(removeButton('Beta'));
      await tester.pump();

      expect(taps, 0);
    });

    testWidgets('removal is by index, so duplicates drop one at a time',
        (tester) async {
      // Identical const instances: removing by value would drop both.
      List<_Thing>? remaining;
      await pumpField(
        tester,
        initialValue: const [alpha, alpha],
        onTap: () {},
        onRemove: (value) => remaining = value,
      );

      await tester.tap(removeButton('Alpha').first);
      await tester.pump();

      expect(remaining, const [alpha]);
    });

    testWidgets('removing the last chip reports an empty list, not null',
        (tester) async {
      List<_Thing>? remaining;
      await pumpField(
        tester,
        initialValue: const [alpha],
        onTap: () {},
        onRemove: (value) => remaining = value,
      );

      await tester.tap(removeButton('Alpha'));
      await tester.pump();

      expect(remaining, isEmpty);
    });

    testWidgets('the chip goes away once the owner writes the value back',
        (tester) async {
      // Mirrors the real consumers (receipt form / quick scan / split sheet),
      // which all answer the callback with setValue. The field renders from
      // its own value, so a consumer that dropped the callback would leave the
      // chip on screen.
      late GlobalKey<FormBuilderState> formKey;
      formKey = await pumpField(
        tester,
        initialValue: const [alpha, beta],
        onTap: () {},
        onRemove: (remaining) => rebuildOwner(
          () => formKey.currentState!.fields[fieldName]!.setValue(remaining),
        ),
      );

      await tester.tap(removeButton('Alpha'));
      await tester.pump();

      expect(find.byType(InputChip), findsOneWidget);
      expect(find.widgetWithText(InputChip, 'Beta'), findsOneWidget);
      expect(formKey.currentState!.fields[fieldName]!.value, const [beta]);
    });
  });

  group('rendering', () {
    testWidgets('renders the label', (tester) async {
      await pumpField(tester, onTap: () {});

      expect(find.text('Categories'), findsOneWidget);
    });

    testWidgets('renders the placeholder and no chips when empty',
        (tester) async {
      await pumpField(tester, onTap: () {});

      expect(find.text('No Categories selected'), findsOneWidget);
      expect(find.byType(InputChip), findsNothing);
    });

    testWidgets('renders one chip per value, labelled and in order',
        (tester) async {
      await pumpField(tester, initialValue: const [alpha, beta], onTap: () {});

      expect(find.byType(InputChip), findsNWidgets(2));
      expect(find.text('No Categories selected'), findsNothing);

      final labels = tester
          .widgetList<InputChip>(find.byType(InputChip))
          .map((chip) => (chip.label as Text).data)
          .toList();
      expect(labels, ['Alpha', 'Beta']);
    });

    testWidgets('renders an empty initial list as the placeholder',
        (tester) async {
      await pumpField(tester, initialValue: const [], onTap: () {});

      expect(find.text('No Categories selected'), findsOneWidget);
      expect(find.byType(InputChip), findsNothing);
    });
  });

  group('form integration', () {
    testWidgets('registers its initial value under the field name',
        (tester) async {
      final formKey =
          await pumpField(tester, initialValue: const [alpha], onTap: () {});

      expect(formKey.currentState!.fields[fieldName]!.value, const [alpha]);
    });

    testWidgets('swaps the placeholder for chips when the value changes',
        (tester) async {
      final formKey = await pumpField(tester, onTap: () {});

      formKey.currentState!.fields[fieldName]!
          .didChange(const <_Thing>[alpha, beta]);
      await tester.pump();

      expect(find.byType(InputChip), findsNWidgets(2));
      expect(
        formKey.currentState!.fields[fieldName]!.value,
        const [alpha, beta],
      );
    });

    testWidgets('falls back to the placeholder when cleared', (tester) async {
      final formKey =
          await pumpField(tester, initialValue: const [alpha], onTap: () {});

      formKey.currentState!.fields[fieldName]!.didChange(const <_Thing>[]);
      await tester.pump();
      expect(find.text('No Categories selected'), findsOneWidget);

      formKey.currentState!.fields[fieldName]!.didChange(null);
      await tester.pump();
      expect(find.text('No Categories selected'), findsOneWidget);
      expect(find.byType(InputChip), findsNothing);
    });
  });

  group('required validator', () {
    testWidgets('rejects an empty value when required', (tester) async {
      final formKey = await pumpField(tester, required: true, onTap: () {});

      expect(formKey.currentState!.validate(), isFalse);

      formKey.currentState!.fields[fieldName]!.didChange(const <_Thing>[alpha]);
      await tester.pump();

      expect(formKey.currentState!.validate(), isTrue);
    });

    testWidgets('accepts an empty value when required is false',
        (tester) async {
      final formKey = await pumpField(tester, required: false, onTap: () {});

      expect(formKey.currentState!.validate(), isTrue);
    });

    testWidgets('accepts an empty value when required is unset',
        (tester) async {
      final formKey = await pumpField(tester, onTap: () {});

      expect(formKey.currentState!.validate(), isTrue);
    });
  });

  group('view mode (no onTap)', () {
    testWidgets('installs no tap surface of its own', (tester) async {
      await pumpField(tester);

      expect(ownGestureDetector(), findsNothing);
    });

    testWidgets('taps anywhere on the field are inert', (tester) async {
      await pumpField(tester);

      final rect = fieldRect(tester);
      await tester.tapAt(Offset(rect.right - 4, rect.center.dy));
      await tester.tapAt(rect.center);
      await tester.tap(find.text('No Categories selected'));
      await tester.pump();

      expect(tester.takeException(), isNull);
    });

    testWidgets('still renders the label and its chips', (tester) async {
      await pumpField(tester, initialValue: const [alpha, beta]);

      expect(find.text('Categories'), findsOneWidget);
      expect(find.byType(InputChip), findsNWidgets(2));
    });

    testWidgets('renders no remove buttons', (tester) async {
      // The wrappers pass a null onRemove alongside the null onTap, matching
      // desktop hiding matChipRemove when readonly.
      await pumpField(tester, initialValue: const [alpha, beta]);

      expect(find.byIcon(Icons.cancel), findsNothing);
    });
  });
}
