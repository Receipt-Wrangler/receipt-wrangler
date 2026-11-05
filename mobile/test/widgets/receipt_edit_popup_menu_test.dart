import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import '../helpers/permission_test_helpers.dart';

PermissionsModel _modelWith(List<Permission> group7Permissions) =>
    seededPermissions(group: {7: group7Permissions});

const _menuKey = ValueKey('popup-menu-under-test');

// formState is passed explicitly so the widget never falls back to reading the
// route (getFormStateFromContext), keeping the test free of a GoRouter.
Widget _wrap(
  PermissionsModel model,
  WranglerFormState formState, {
  List<PopupMenuEntry> alwaysShownChildren = const [],
}) {
  return ChangeNotifierProvider<PermissionsModel>.value(
    value: model,
    child: MaterialApp(
      home: Scaffold(
        body: ReceiptEditPopupMenu(
          key: _menuKey,
          groupId: 7,
          formState: formState,
          popupMenuChildren: const [
            PopupMenuItem<int>(value: 0, child: Text('Edit')),
          ],
          alwaysShownChildren: alwaysShownChildren,
        ),
      ),
    ),
  );
}

const _history = [PopupMenuItem<int>(value: 1, child: Text('History'))];

// House style: locate via the keyed widget under test, not a bare type lookup
// (a second PopupMenuButton elsewhere in the tree would silently match).
Finder _menuButton() => find.descendant(
      of: find.byKey(_menuKey),
      matching: find.byType(PopupMenuButton<dynamic>),
    );

void main() {
  group('ReceiptEditPopupMenu', () {
    testWidgets('shows the menu when the user can update receipts',
        (tester) async {
      await tester.pumpWidget(_wrap(
        _modelWith([Permission.groupPeriodReceiptsPeriodUpdate]),
        WranglerFormState.view,
      ));

      expect(_menuButton(), findsOneWidget);
    });

    testWidgets('hides the menu when the user cannot update receipts',
        (tester) async {
      await tester.pumpWidget(_wrap(
        _modelWith([Permission.groupPeriodReceiptsPeriodRead]),
        WranglerFormState.view,
      ));

      expect(_menuButton(), findsNothing);
    });

    testWidgets('always shows the menu in add form state (creating)',
        (tester) async {
      // No update permission, but add state forces the menu (you are creating).
      await tester.pumpWidget(_wrap(
        _modelWith([Permission.groupPeriodReceiptsPeriodRead]),
        WranglerFormState.add,
      ));

      expect(_menuButton(), findsOneWidget);
    });

    // The receipt view's History entry is gated by its caller, not by the edit
    // permission, so a read-only member who may see the history still gets it.
    testWidgets('shows only the always-shown entries to a user who cannot edit',
        (tester) async {
      await tester.pumpWidget(_wrap(
        _modelWith([Permission.groupPeriodReceiptsPeriodRead]),
        WranglerFormState.view,
        alwaysShownChildren: _history,
      ));

      await tester.tap(_menuButton());
      await tester.pumpAndSettle();

      expect(find.text('History'), findsOneWidget);
      expect(find.text('Edit'), findsNothing);
    });

    testWidgets('lists the edit entries before the always-shown ones',
        (tester) async {
      await tester.pumpWidget(_wrap(
        _modelWith([Permission.groupPeriodReceiptsPeriodUpdate]),
        WranglerFormState.view,
        alwaysShownChildren: _history,
      ));

      await tester.tap(_menuButton());
      await tester.pumpAndSettle();

      final edit = tester.getTopLeft(find.text('Edit'));
      final history = tester.getTopLeft(find.text('History'));
      expect(edit.dy, lessThan(history.dy));
    });
  });
}
