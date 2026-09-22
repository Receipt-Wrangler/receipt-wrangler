import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';

/// The bar is shared by the receipt form, Quick Scan, the multi-select picker
/// and the receipt filter, so the restyle has to leave its behaviour alone.
void main() {
  const buttonKey = ValueKey("submit");

  Future<LoadingModel> pumpBar(
    WidgetTester tester, {
    VoidCallback? onPressed,
    bool? disabled,
    String? buttonText,
  }) async {
    final loadingModel = LoadingModel();

    await tester.pumpWidget(
      ChangeNotifierProvider<LoadingModel>.value(
        value: loadingModel,
        child: MaterialApp(
          theme: buildAppTheme(),
          home: Scaffold(
            bottomNavigationBar: BottomSubmitButton(
              key: buttonKey,
              onPressed: onPressed ?? () {},
              disabled: disabled,
              buttonText: buttonText,
            ),
          ),
        ),
      ),
    );
    await tester.pump();

    return loadingModel;
  }

  FilledButton buttonOf(WidgetTester tester) => tester.widget<FilledButton>(
        find.descendant(
            of: find.byKey(buttonKey), matching: find.byType(FilledButton)),
      );

  testWidgets("labels itself, defaulting to Submit", (tester) async {
    await pumpBar(tester);
    expect(find.text("Submit"), findsOneWidget);

    await pumpBar(tester, buttonText: "Apply Filter");
    expect(find.text("Apply Filter"), findsOneWidget);
    expect(find.text("Submit"), findsNothing);
  });

  testWidgets("fires onPressed", (tester) async {
    var taps = 0;
    await pumpBar(tester, onPressed: () => taps++);

    await tester.tap(find.byKey(buttonKey));
    await tester.pump();

    expect(taps, 1);
  });

  testWidgets("disabled: true takes the callback away", (tester) async {
    var taps = 0;
    await pumpBar(tester, onPressed: () => taps++, disabled: true);

    expect(buttonOf(tester).onPressed, isNull);

    await tester.tap(find.byKey(buttonKey), warnIfMissed: false);
    await tester.pump();
    expect(taps, 0);
  });

  testWidgets("a null disabled is not a disabled button", (tester) async {
    await pumpBar(tester);
    expect(buttonOf(tester).onPressed, isNotNull);
  });

  testWidgets("a loading request swaps the label for a spinner and disables it",
      (tester) async {
    final loadingModel = await pumpBar(tester, buttonText: "Apply Filter");

    expect(find.byType(CircularProgressIndicator), findsNothing);

    loadingModel.setIsLoading(true);
    await tester.pump();

    expect(find.byType(CircularProgressIndicator), findsOneWidget);
    expect(find.text("Apply Filter"), findsNothing);
    expect(buttonOf(tester).onPressed, isNull);

    loadingModel.setIsLoading(false);
    await tester.pump();

    expect(find.byType(CircularProgressIndicator), findsNothing);
    expect(find.text("Apply Filter"), findsOneWidget);
    expect(buttonOf(tester).onPressed, isNotNull);
  });

  testWidgets("reserves exactly the height submitButtonSpacing clears",
      (tester) async {
    // Where the bar floats as a Scaffold.bottomSheet, a scrollable above it
    // clears it with submitButtonSpacing. The two are derived from one constant
    // precisely so they cannot drift -- this pins the constant to the render.
    await pumpBar(tester);

    expect(tester.getSize(find.byKey(buttonKey)).height,
        moreOrLessEquals(bottomSubmitBarHeight, epsilon: 0.5));
  });
}
