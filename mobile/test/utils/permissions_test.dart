import 'dart:async';

import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/utils/permissions.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();
  final messenger =
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;

  const permissions =
      MethodChannel('flutter.baseflow.com/permissions/methods');
  const gal = MethodChannel('gal');

  late int permissionCalls;
  late int galCalls;
  late Completer<Map<int, int>> permissionGate;
  late Completer<bool> galGate;

  setUp(() {
    permissionCalls = 0;
    galCalls = 0;
    permissionGate = Completer<Map<int, int>>();
    galGate = Completer<bool>();

    messenger.setMockMethodCallHandler(permissions, (call) async {
      if (call.method == 'requestPermissions') {
        permissionCalls += 1;
        return permissionGate.future;
      }
      if (call.method == 'checkPermissionStatus' ||
          call.method == 'checkServiceStatus') {
        return 1;
      }
      return null;
    });
    messenger.setMockMethodCallHandler(gal, (call) async {
      if (call.method == 'requestAccess') {
        galCalls += 1;
        return galGate.future;
      }
      if (call.method == 'hasAccess') return true;
      return null;
    });
  });

  tearDown(() {
    messenger.setMockMethodCallHandler(permissions, null);
    messenger.setMockMethodCallHandler(gal, null);
  });

  test(
      'concurrent requestPermissions calls share a single in-flight platform request',
      () async {
    // Fire two callers before the platform side resolves.
    final first = requestPermissions();
    final second = requestPermissions();

    // Pump one microtask cycle so the first call has reached the channel.
    await Future<void>.delayed(Duration.zero);

    expect(permissionCalls, 1,
        reason: 'second caller should not re-fire the platform channel');
    expect(galCalls, 0,
        reason: 'gal request only fires after permission resolves');

    // Resolve the platform side; both callers complete with the same value.
    permissionGate.complete(<int, int>{});
    galGate.complete(true);

    await Future.wait([first, second]);

    expect(permissionCalls, 1);
    expect(galCalls, 1);
  });

  test('after the in-flight Future resolves, a subsequent call re-fires',
      () async {
    // Replace the gated handlers with immediate-returning ones for this
    // test so we can sequence two complete calls.
    messenger.setMockMethodCallHandler(permissions, (call) async {
      if (call.method == 'requestPermissions') {
        permissionCalls += 1;
        return <int, int>{};
      }
      return null;
    });
    messenger.setMockMethodCallHandler(gal, (call) async {
      if (call.method == 'requestAccess') {
        galCalls += 1;
        return true;
      }
      return null;
    });

    await requestPermissions();
    expect(permissionCalls, 1);

    await requestPermissions();
    expect(permissionCalls, 2,
        reason: 'a fresh call after settle should hit the platform again');
  });
}
