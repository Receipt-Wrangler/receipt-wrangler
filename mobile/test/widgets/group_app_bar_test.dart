import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/nav/group/group_app_bar.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

import '../helpers/auth_test_helpers.dart';

api.Claims _claims(String displayName) =>
    api.Claims((b) => b..displayName = displayName);

api.Group _group({required int id, required String name}) {
  // NOTE: api.Group requires several non-nullable fields. We mock instead of
  // building a real one so the test isn't coupled to the full builder graph.
  final mock = MockGroup();
  when(() => mock.id).thenReturn(id);
  when(() => mock.name).thenReturn(name);
  return mock;
}

Future<void> _pumpAt(
  WidgetTester tester, {
  required String initialLocation,
  required api.Group? Function(String groupId) lookup,
}) async {
  final groupModel = MockGroupModel();
  when(() => groupModel.getGroupById(any())).thenAnswer(
    (invocation) => lookup(invocation.positionalArguments.first as String),
  );

  final authModel = MockAuthModel();
  when(() => authModel.claims).thenReturn(_claims('Admin'));

  final router = GoRouter(
    initialLocation: initialLocation,
    routes: [
      GoRoute(
        path: '/groups/:groupId/x',
        builder: (_, __) => const Scaffold(appBar: GroupAppBar()),
      ),
      GoRoute(
        path: '/no-group',
        builder: (_, __) => const Scaffold(appBar: GroupAppBar()),
      ),
    ],
  );

  await tester.pumpWidget(MultiProvider(
    providers: [
      // Plain Provider.value avoids ChangeNotifierProvider trying to
      // addListener on mocktail Mocks (which don't stub addListener).
      // GroupAppBar reads with listen: false; UserAvatar's listen: true
      // on AuthModel still works because Provider.of(listen: true) only
      // sets up inherited-widget dependency, it doesn't subscribe.
      Provider<AuthModel>.value(value: authModel),
      ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
      Provider<GroupModel>.value(value: groupModel),
      ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
    ],
    child: MaterialApp.router(routerConfig: router),
  ));
  await tester.pump();
}

Finder _titleInside(GroupAppBar _, String text) => find.descendant(
      of: find.byType(GroupAppBar),
      matching: find.text(text),
    );

void main() {
  setUpAll(() {
    registerFallbackValue('');
    // Allow Provider<T>.value with mocktail Mocks of ChangeNotifier subclasses.
    // We use plain Provider (not ChangeNotifierProvider) on purpose because
    // the test never relies on listener-based rebuilds.
    Provider.debugCheckInvalidValueType = null;
  });

  testWidgets('group resolved with name containing "receipt" → name verbatim',
      (tester) async {
    final group = _group(id: 1, name: 'My Receipts');
    await _pumpAt(
      tester,
      initialLocation: '/groups/1/x',
      lookup: (id) => id == '1' ? group : null,
    );
    expect(_titleInside(const GroupAppBar(), 'My Receipts'), findsOneWidget);
  });

  testWidgets('group resolved with arbitrary name → "Receipts" suffix appended',
      (tester) async {
    final group = _group(id: 3, name: 'Trips');
    await _pumpAt(
      tester,
      initialLocation: '/groups/3/x',
      lookup: (id) => id == '3' ? group : null,
    );
    expect(_titleInside(const GroupAppBar(), 'Trips Receipts'), findsOneWidget);
  });

  testWidgets('groupId in route but unknown to model → fallback "Receipts"',
      (tester) async {
    await _pumpAt(
      tester,
      initialLocation: '/groups/2/x',
      lookup: (_) => null,
    );
    expect(_titleInside(const GroupAppBar(), 'Receipts'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });

  testWidgets('no groupId param → getGroupId falls back to "0", null group, '
      'fallback "Receipts" rendered without crash', (tester) async {
    await _pumpAt(
      tester,
      initialLocation: '/no-group',
      lookup: (_) => null,
    );
    expect(_titleInside(const GroupAppBar(), 'Receipts'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });
}
