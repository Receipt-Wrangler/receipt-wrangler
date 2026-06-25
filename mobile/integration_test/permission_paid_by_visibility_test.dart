import 'dart:io' show Platform;

import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';

/// Permission coverage for group-role **paid-by visibility**. A member whose
/// group role restricts paid-by visibility to "their own receipts" must see only
/// the receipts they paid for: the backend filters the receipts list, so the
/// other-payer receipt never reaches the app. Mirrors the list assertion in the
/// desktop `paid-by-visibility.spec.ts` (the API-403 / direct-nav cases there are
/// desktop-only — mobile opens receipts by tapping the list, so a filtered-out
/// receipt is simply unreachable from the UI).
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
    'paid-by-own member sees only their own receipt in the group list',
    (tester) async {
      final paidBy = await provisionPaidByOwnMember();
      await loginAs(
        tester,
        username: paidBy.fixture.username,
        password: paidBy.fixture.password,
      );

      // Opening the group's receipt list waits for the member's own receipt,
      // proving the list has loaded.
      await openGroupReceipts(
        tester,
        paidBy.fixture.groupName!,
        paidBy.ownReceiptName,
      );

      // The own receipt is visible; the admin-paid receipt is filtered out by
      // the role's paid-by grant server-side, so it never appears in the list.
      expect(find.text(paidBy.ownReceiptName), findsOneWidget);
      expect(find.text(paidBy.hiddenReceiptName), findsNothing);
    },
  );
}
