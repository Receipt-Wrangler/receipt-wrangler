import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';

Tag _tag(int id, String name) =>
    Tag((b) => b
      ..id = id
      ..name = name);

void main() {
  group('TagModel.tagsForGroup', () {
    test('returns the per-group catalog for a known group', () {
      final model = TagModel();
      final urgent = _tag(1, 'Urgent');
      model.setGroupTags({
        7: [urgent],
      });

      expect(model.tagsForGroup(7), [urgent]);
    });

    test('isolates catalogs per group', () {
      final model = TagModel();
      final urgent = _tag(1, 'Urgent');
      final personal = _tag(2, 'Personal');
      model.setGroupTags({
        7: [urgent],
        9: [personal],
      });

      expect(model.tagsForGroup(7), [urgent]);
      expect(model.tagsForGroup(9), [personal]);
    });

    test('returns an empty list for an unknown group', () {
      final model = TagModel();
      model.setGroupTags({
        7: [_tag(1, 'Urgent')],
      });

      expect(model.tagsForGroup(99), isEmpty);
    });

    test('returns an empty list before hydration', () {
      expect(TagModel().tagsForGroup(7), isEmpty);
    });

    test('setGroupTags notifies listeners', () {
      final model = TagModel();
      var notified = 0;
      model.addListener(() => notified++);
      model.setGroupTags({});
      expect(notified, 1);
    });
  });
}
