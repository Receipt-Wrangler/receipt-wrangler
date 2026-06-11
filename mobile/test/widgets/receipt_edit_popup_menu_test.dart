import 'package:built_collection/built_collection.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

PermissionsModel _modelWith(List<Permission> group7Permissions) {
  final model = PermissionsModel();
  model.setPermissions(
    BuiltList<Permission>(),
    BuiltMap<String, BuiltList<Permission>>({
      '7': BuiltList<Permission>(group7Permissions),
    }),
  );
  return model;
}

const _menuKey = ValueKey('popup-menu-under-test');

// formState is passed explicitly so the widget never falls back to reading the
// route (getFormStateFromContext), keeping the test free of a GoRouter.
Widget _wrap(PermissionsModel model, WranglerFormState formState) {
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
        ),
      ),
    ),
  );
}

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
  });
}
