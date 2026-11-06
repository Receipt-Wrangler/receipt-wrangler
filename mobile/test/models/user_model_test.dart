import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

import '../helpers/receipt_form_test_helpers.dart' show buildUserView;

void main() {
  group('UserModel.getUserById', () {
    test('finds a user by their string id', () {
      final model = UserModel()
        ..setUsers([buildUserView(id: 1, displayName: 'Ana'), buildUserView(id: 2)]);

      expect(model.getUserById('1')?.displayName, 'Ana');
    });

    // History names users who were deleted or are hidden from the viewer, so a
    // miss is ordinary. It used to throw a StateError despite the nullable type.
    test('returns null for an id no visible user has', () {
      final model = UserModel()..setUsers([buildUserView(id: 1)]);

      expect(model.getUserById('9'), isNull);
    });
  });
}
