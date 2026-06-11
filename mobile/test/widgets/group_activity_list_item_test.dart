import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_activity_list_item.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';

/// Permission-gating coverage for the activity rerun swipe action
/// (`group_activity_list_item.dart`). The gate is
/// `slideEnabled = canRerun && (canBeRestarted ?? false) && !hasBeenRerun`,
/// where `canRerun` is `group.activities.rerun` for the activity's own group
/// (falling back to the listing group when the activity carries none).
///
/// This lives as a widget test (not e2e) on purpose: the `canBeRestarted`
/// factor is only true for a background task the backend has archived as
/// failed in Redis, which is not deterministically seedable from an e2e
/// fixture. Pumping the widget directly lets us isolate and assert the
/// permission factor for every combination.
const _groupId = 42;

/// A listing-context group id that never appears in the permission map --
/// stands in for the synthetic "all groups" group.
const _allGroupId = 99;

const _itemKey = ValueKey('activity-under-test');

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockSystemTaskApi extends Mock implements api.SystemTaskApi {}

PermissionsModel _permissions(List<api.Permission> groupPerms) {
  final model = PermissionsModel();
  model.setPermissions(
    BuiltList<api.Permission>(),
    BuiltMap<String, BuiltList<api.Permission>>({
      '$_groupId': BuiltList<api.Permission>(groupPerms),
    }),
  );
  return model;
}

api.Activity _activity({required bool canBeRestarted}) => api.Activity(
      (b) => b
        ..id = 1
        ..type = api.SystemTaskType.QUICK_SCAN
        ..status = canBeRestarted
            ? api.SystemTaskStatus.FAILED
            : api.SystemTaskStatus.SUCCEEDED
        ..startedAt = '2026-06-11T00:00:00Z'
        ..endedAt = '2026-06-11T00:00:01Z'
        ..groupId = _groupId
        ..canBeRestarted = canBeRestarted,
    );

Widget _app(
  PermissionsModel permissions,
  api.Activity activity, {
  int listingGroupId = _groupId,
}) {
  return MultiProvider(
    providers: [
      ChangeNotifierProvider<PermissionsModel>.value(value: permissions),
      ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
      ChangeNotifierProvider<GroupModel>(create: (_) => GroupModel()),
    ],
    child: MaterialApp(
      // ListItemLead pins a 50px-wide leading box (color block + date) and the
      // trailing column has a fixed-height status chip; with the default test
      // font those overflow by a few px. The overflow is cosmetic and unrelated
      // to the rerun gate under test, so shrink text to keep the fixed-size
      // sub-widgets within their boxes.
      home: MediaQuery(
        data: const MediaQueryData(textScaler: TextScaler.linear(0.6)),
        child: Scaffold(
          body: GroupActivityListItem(
            key: _itemKey,
            activity: activity,
            groupId: listingGroupId,
          ),
        ),
      ),
    ),
  );
}

bool _slideEnabled(WidgetTester tester) => tester
    .widget<SlidableWidget>(find.descendant(
      of: find.byKey(_itemKey),
      matching: find.byType(SlidableWidget),
    ))
    .slideEnabled;

void main() {
  testWidgets(
    'rerun enabled: has group.activities.rerun and the activity is restartable',
    (tester) async {
      await tester.pumpWidget(_app(
        _permissions([api.Permission.groupPeriodActivitiesPeriodRerun]),
        _activity(canBeRestarted: true),
      ));
      await tester.pump();
      expect(_slideEnabled(tester), isTrue);
    },
  );

  testWidgets(
    'rerun enabled in all-groups context: permission keyed by the '
    "activity's own group, not the synthetic listing group",
    (tester) async {
      await tester.pumpWidget(_app(
        _permissions([api.Permission.groupPeriodActivitiesPeriodRerun]),
        _activity(canBeRestarted: true),
        // The listing group never appears in the permission map (the "all
        // groups" view); the gate must fall through to activity.groupId.
        listingGroupId: _allGroupId,
      ));
      await tester.pump();
      expect(_slideEnabled(tester), isTrue);
    },
  );

  testWidgets(
    'rerun disabled: lacks group.activities.rerun (member, read-only)',
    (tester) async {
      await tester.pumpWidget(_app(
        // Member of the group but without the rerun permission.
        _permissions([api.Permission.groupPeriodReceiptsPeriodRead]),
        _activity(canBeRestarted: true),
      ));
      await tester.pump();
      expect(_slideEnabled(tester), isFalse);
    },
  );

  testWidgets(
    'rerun disabled: activity is not restartable even with permission',
    (tester) async {
      await tester.pumpWidget(_app(
        _permissions([api.Permission.groupPeriodActivitiesPeriodRerun]),
        _activity(canBeRestarted: false),
      ));
      await tester.pump();
      expect(_slideEnabled(tester), isFalse);
    },
  );

  testWidgets(
    'rerun disabled after a successful rerun: the hasBeenRerun latch',
    (tester) async {
      // hasBeenRerun is widget-local state set after a successful rerun API
      // call -- the only way to exercise it is to actually perform the swipe
      // + tap with the API client stubbed out. Don't read (or restore) the
      // original OpenApiClient.client: its initializer constructs a real
      // Openapi with a relative '/api' baseUrl, which dio rejects off-Web.
      // This file is the only consumer of the client in this test isolate,
      // so installing the mock for the file's lifetime is safe.
      final mockClient = _MockOpenapi();
      final mockTaskApi = _MockSystemTaskApi();
      when(() => mockClient.getSystemTaskApi()).thenReturn(mockTaskApi);
      when(() => mockTaskApi.rerunActivity(id: any(named: 'id'))).thenAnswer(
        (_) async => Response<void>(
          requestOptions: RequestOptions(path: ''),
          statusCode: 200,
        ),
      );
      OpenApiClient.client = mockClient;

      await tester.pumpWidget(_app(
        _permissions([api.Permission.groupPeriodActivitiesPeriodRerun]),
        _activity(canBeRestarted: true),
      ));
      await tester.pump();
      expect(_slideEnabled(tester), isTrue);

      // Reveal the end action pane and tap "Rerun".
      await tester.drag(find.byKey(_itemKey), const Offset(-300, 0));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Rerun'));
      await tester.pumpAndSettle();

      verify(() => mockTaskApi.rerunActivity(id: 1)).called(1);
      expect(_slideEnabled(tester), isFalse,
          reason: 'a successful rerun must latch the gate off so the same '
              'activity cannot be queued twice from this list item');
    },
  );
}
