import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';

void main() {
  group('AuthModel.pendingServerUrl', () {
    test('defaults to null', () {
      expect(AuthModel().pendingServerUrl, isNull);
    });

    test('setPendingServerUrl stores the value and notifies', () {
      final model = AuthModel();
      var notified = 0;
      model.addListener(() => notified++);

      model.setPendingServerUrl('https://demo.receiptwrangler.io/api');

      expect(model.pendingServerUrl, 'https://demo.receiptwrangler.io/api');
      expect(notified, 1);
    });

    test('clearPendingServerUrl clears and notifies when a value is set', () {
      final model = AuthModel();
      model.setPendingServerUrl('https://demo.receiptwrangler.io/api');

      var notified = 0;
      model.addListener(() => notified++);
      model.clearPendingServerUrl();

      expect(model.pendingServerUrl, isNull);
      expect(notified, 1);
    });

    test('clearPendingServerUrl is a no-op (no notify) when already null', () {
      final model = AuthModel();
      var notified = 0;
      model.addListener(() => notified++);

      model.clearPendingServerUrl();

      expect(model.pendingServerUrl, isNull);
      expect(notified, 0);
    });
  });
}
