import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for NotificationsApi
void main() {
  final instance = Openapi().getNotificationsApi();

  group(NotificationsApi, () {
    // Delete all notifications for user
    //
    // This deletes all notifications for a user
    //
    //Future deleteAllNotificationsForUser() async
    test('test deleteAllNotificationsForUser', () async {
      // TODO
    });

    // Delete notification by id
    //
    // This deletes a notification by id
    //
    //Future deleteNotificationById(int notificationId) async
    test('test deleteNotificationById', () async {
      // TODO
    });

    // Notification count
    //
    // This will get the notification count for the currently logged in user
    //
    //Future<int> getNotificationCount() async
    test('test getNotificationCount', () async {
      // TODO
    });

    // Get all user notifications
    //
    // This will get all the notifications for the currently logged in user
    //
    //Future<BuiltList<Notification>> getNotificationsForuser() async
    test('test getNotificationsForuser', () async {
      // TODO
    });

  });
}
