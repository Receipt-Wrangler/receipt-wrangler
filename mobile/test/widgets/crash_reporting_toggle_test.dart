import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/profile/widgets/crash_reporting_toggle.dart';

void main() {
  Widget host(Widget child) => MaterialApp(home: Scaffold(body: child));

  SwitchListTile theSwitch(WidgetTester tester) =>
      tester.widget<SwitchListTile>(find.byType(SwitchListTile));

  testWidgets('reflects the injected initial enabled state', (tester) async {
    await tester.pumpWidget(host(CrashReportingToggle(
      readEnabled: () => false,
      setEnabled: (_) async {},
    )));

    expect(theSwitch(tester).value, isFalse);
  });

  testWidgets('persists the new value via setEnabled', (tester) async {
    bool? saved;
    await tester.pumpWidget(host(CrashReportingToggle(
      readEnabled: () => true,
      setEnabled: (value) async => saved = value,
    )));

    await tester.tap(find.byType(SwitchListTile));
    await tester.pumpAndSettle();

    expect(saved, isFalse);
    expect(theSwitch(tester).value, isFalse);
  });

  testWidgets('reverts the switch and shows a snackbar when setEnabled throws',
      (tester) async {
    await tester.pumpWidget(host(CrashReportingToggle(
      readEnabled: () => true,
      setEnabled: (_) async => throw Exception('boom'),
    )));

    await tester.tap(find.byType(SwitchListTile));
    await tester.pump(); // let _onChanged run through its catch
    await tester.pump(const Duration(milliseconds: 750)); // snackbar entrance

    expect(theSwitch(tester).value, isTrue); // reverted to the prior value
    expect(find.text("Couldn't update crash reporting setting"), findsOneWidget);
  });

  testWidgets('disables the switch while a toggle is in flight', (tester) async {
    final gate = Completer<void>();
    await tester.pumpWidget(host(CrashReportingToggle(
      readEnabled: () => true,
      setEnabled: (_) => gate.future,
    )));

    await tester.tap(find.byType(SwitchListTile));
    await tester.pump(); // _isLoading = true

    expect(theSwitch(tester).onChanged, isNull); // disabled mid-flight

    gate.complete();
    await tester.pumpAndSettle();

    expect(theSwitch(tester).onChanged, isNotNull); // re-enabled
  });
}
