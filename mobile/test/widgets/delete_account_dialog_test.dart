import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/profile/widgets/delete_account_dialog.dart';

// The delete-account dialog masks its password the same way the login form
// does, and offers the same eye to reveal it -- the icon names the ACTION (an
// open eye while masked, crossed-out while revealed), matching the desktop
// visibility eye. Both of the app's password fields must behave identically,
// so what auth_form_test asserts about the login field is asserted here too.
void main() {
  const toggleKey = ValueKey('password-visibility-toggle');

  /// Opens the dialog through its public entry point -- the widget itself is
  /// private, so `showDeleteAccountDialog` is the only way in, and it is also
  /// how the profile screen reaches it.
  Future<void> pumpDialog(WidgetTester tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Builder(
          builder: (context) => Scaffold(
            body: Center(
              child: ElevatedButton(
                onPressed: () => showDeleteAccountDialog(context),
                child: const Text('open'),
              ),
            ),
          ),
        ),
      ),
    );

    await tester.tap(find.text('open'));
    await tester.pumpAndSettle();
  }

  bool isObscured(WidgetTester tester) =>
      tester.widget<TextField>(find.byType(TextField)).obscureText;

  Future<void> tapToggle(WidgetTester tester) async {
    await tester.tap(find.byKey(toggleKey));
    await tester.pump();
  }

  testWidgets('the password starts masked, offering to show it',
      (tester) async {
    await pumpDialog(tester);

    expect(isObscured(tester), isTrue);
    expect(find.byKey(toggleKey), findsOneWidget);
    expect(
      tester.widget<IconButton>(find.byKey(toggleKey)).tooltip,
      'Show Password',
    );
    expect(
      find.descendant(
          of: find.byKey(toggleKey), matching: find.byIcon(Icons.visibility)),
      findsOneWidget,
    );
  });

  testWidgets('tapping the eye reveals the password and offers to hide it',
      (tester) async {
    await pumpDialog(tester);

    await tapToggle(tester);

    expect(isObscured(tester), isFalse);
    expect(
      tester.widget<IconButton>(find.byKey(toggleKey)).tooltip,
      'Hide Password',
    );
    expect(
      find.descendant(
          of: find.byKey(toggleKey),
          matching: find.byIcon(Icons.visibility_off)),
      findsOneWidget,
    );
  });

  testWidgets('tapping it again re-masks', (tester) async {
    await pumpDialog(tester);

    await tapToggle(tester);
    await tapToggle(tester);

    expect(isObscured(tester), isTrue);
    expect(
      tester.widget<IconButton>(find.byKey(toggleKey)).tooltip,
      'Show Password',
    );
  });

  testWidgets('revealing does not submit the dialog', (tester) async {
    await pumpDialog(tester);

    await tester.enterText(find.byType(TextField), 'MySecretPass123');
    await tester.pump();

    await tapToggle(tester);

    // Still open, still holding what was typed -- the eye is not a submit.
    expect(find.text('Delete Account'), findsWidgets);
    expect(find.text('MySecretPass123'), findsOneWidget);
    expect(isObscured(tester), isFalse);
  });
}
