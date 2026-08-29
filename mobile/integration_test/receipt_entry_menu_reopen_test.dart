// The scan slot has to keep working after its menu closes.
//
// Regression: the long-press action used to run on `onLongPress`, which fires at
// the 500ms mark with the finger still down. Opening the menu route there left
// the gesture recognizer holding a pointer it never saw released, and that slice
// then ignored every later tap and hold — a slot that worked exactly once per
// app launch. Running on `onLongPressUp` instead lets the gesture finish first.
//
// This lives in the e2e suite on purpose: a minimal widget harness delivers the
// pointer-up cleanly enough that the recognizer recovers, so the equivalent
// widget test passes even with the bug present. Only the real app reproduces it.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';

import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('the scan slot survives its menu being opened and dismissed',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);
    await pumpUntilFound(tester, scanNavSlot());

    // Three cycles: one open would pass even with the bug, since it was the
    // *second* interaction that died.
    for (var i = 1; i <= 3; i++) {
      await tester.longPress(scanNavSlot());
      await pumpUntilFound(tester, find.text(addManualReceiptLabel),
          timeout: const Duration(seconds: 5));

      await tester.tapAt(const Offset(10, 10));
      await pumpUntilGone(tester, find.text(addManualReceiptLabel));
    }

    // A plain tap must still reach the destination's action afterwards — the
    // dead slice swallowed taps as well as holds.
    await tester.tap(scanNavSlot());
    await pumpUntilFound(tester, find.text('Name'));
  });
}
