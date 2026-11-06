import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/utils/receipt_history.dart';

import '../helpers/receipt_form_test_helpers.dart' show buildUserView;
import '../helpers/receipt_history_test_helpers.dart';

final _maya = buildUserView(id: 2, displayName: 'Maya Chen');

api.UserView? _findUser(int? id) => id == 2 ? _maya : null;

final _formatters = HistoryValueFormatters(
  currency: (value) => '\$$value',
  findUser: _findUser,
);

api.ReceiptHistoryValue _value(String type, {String? value, int? userId}) =>
    api.ReceiptHistoryValue((b) => b
      ..type = type
      ..value = value
      ..userId = userId);

void main() {
  group('historyActionVerb', () {
    test('words every action the server sends', () {
      expect(historyActionVerb('CREATED'), 'created this receipt');
      expect(historyActionVerb('UPDATED'), 'edited this receipt');
      expect(historyActionVerb('STATUS_CHANGED'), 'changed the status');
      expect(historyActionVerb('RESOLVED'), 'resolved this receipt');
      expect(historyActionVerb('IMAGES_COMMENTS_UPDATED'), 'updated images and comments');
      expect(historyActionVerb('FAILED'), 'tried to save changes');
    });

    test('falls back for an action this build does not know', () {
      expect(historyActionVerb('ARCHIVED'), defaultActionVerb);
      expect(historyActionVerb(null), defaultActionVerb);
    });
  });

  group('historyFieldLabel', () {
    String label(Map<String, Object?> json) =>
        historyFieldLabel(changeFromJson(json), _findUser);

    test('names the fixed fields after the receipt form', () {
      expect(label({'field': 'DATE', 'kind': 'CHANGED'}), 'Receipt Date');
      expect(label({'field': 'PAID_BY', 'kind': 'CHANGED'}), 'Paid By');
      expect(label({'field': 'EMAIL_INBOX', 'kind': 'ADDED'}), 'Inbox');
      expect(label({'field': 'EMAIL_SENDER', 'kind': 'ADDED'}), 'Sender');
    });

    test('names what an item, share or custom field change is about', () {
      expect(label({'field': 'ITEM', 'kind': 'ADDED', 'label': 'Bag Fee'}), 'Item: Bag Fee');
      expect(label({'field': 'SHARE', 'kind': 'ADDED', 'userId': 2}), 'Share: Maya Chen');
      expect(label({'field': 'SHARE', 'kind': 'ADDED', 'userId': 9}), 'Share: $unknownUserName');
      expect(label({'field': 'CUSTOM_FIELD', 'kind': 'CHANGED', 'label': 'Mileage'}), 'Mileage');
      expect(
        label({'field': 'ITEM_CATEGORIES', 'kind': 'ADDED', 'label': 'Produce'}),
        'Item categories: Produce',
      );
    });

    test('title-cases a field this build does not know', () {
      expect(label({'field': 'SPLIT_RULE', 'kind': 'CHANGED'}), 'Split Rule');
    });
  });

  group('historyActorDisplay', () {
    HistoryActorDisplay actor(Map<String, Object?> json) => historyActorDisplay(
          api.standardSerializers
              .deserializeWith(api.ReceiptHistoryActor.serializer, json)!,
          _findUser,
        );

    test('a known user shows as themselves, in their own colour', () {
      final display = actor({'type': 'USER', 'userId': 2});
      expect(display.name, 'Maya Chen');
      expect(display.user, _maya);
      expect(display.icon, isNull);
      expect(display.color, isNull);
    });

    test('a user the viewer cannot see shows as unknown, never throwing', () {
      final display = actor({'type': 'USER', 'userId': 9});
      expect(display.name, unknownUserName);
      expect(display.icon, Icons.person);
    });

    test('the non-user actors read as the design names them', () {
      expect(actor({'type': 'QUICK_SCAN'}).name, 'Quick Scan');
      expect(actor({'type': 'QUICK_SCAN'}).icon, Icons.photo_camera);
      expect(actor({'type': 'EMAIL'}).name, 'Email Upload');
      expect(actor({'type': 'EMAIL'}).icon, Icons.mail);
      expect(actor({'type': 'SYSTEM'}).name, 'System');
      expect(actor({'type': 'ROBOT'}).name, 'System');
    });
  });

  test('historySummary pluralises each count', () {
    expect(
      historySummary(historyFromJson(
          historyJson([], changeCount: 9, saveCount: 4, contributorCount: 2))),
      '9 changes · 4 saves · 2 contributors',
    );
    expect(
      historySummary(historyFromJson(
          historyJson([], changeCount: 1, saveCount: 1, contributorCount: 1))),
      '1 change · 1 save · 1 contributor',
    );
  });

  group('historyBeforeStateNotice', () {
    test('names the copy a rebuilt entry compared against', () {
      final notice = historyBeforeStateNotice(entryFromJson(entryJson(
        beforeState: 'REBUILT',
        beforeSourceType: 'RECEIPT_UPDATED',
        beforeSourceRecordedAt: '2026-10-16T12:00:00Z',
      )))!;
      expect(notice, startsWith('Recorded by an older version — compared against the previous update ('));
      expect(notice, endsWith('Bulk status changes can also show here.'));

      final created = historyBeforeStateNotice(entryFromJson(
          entryJson(beforeState: 'REBUILT', beforeSourceType: 'RECEIPT_UPLOADED')))!;
      expect(created, contains('the copy saved when it was created.'));
    });

    test('warns about an incomplete entry', () {
      expect(
        historyBeforeStateNotice(entryFromJson(entryJson(beforeState: 'INCOMPLETE'))),
        'Recorded by an older version with no earlier complete copy. '
        'Some changes may be inaccurate.',
      );
    });

    test('says nothing for a complete, absent or unknown state', () {
      expect(historyBeforeStateNotice(entryFromJson(entryJson(beforeState: 'COMPLETE'))), isNull);
      expect(historyBeforeStateNotice(entryFromJson(entryJson())), isNull);
      expect(historyBeforeStateNotice(entryFromJson(entryJson(beforeState: 'GUESSED'))), isNull);
    });
  });

  group('groupHistoryByDay', () {
    // Local times sent the way the server sends them (UTC), so the expected
    // days hold in whatever zone the suite runs in -- grouping is device-local.
    String at(DateTime local) => local.toUtc().toIso8601String();
    final entries = [
      entryJson(systemTaskId: 3, recordedAt: at(DateTime(2026, 10, 13, 12))),
      entryJson(systemTaskId: 2, recordedAt: at(DateTime(2026, 10, 13, 9))),
      entryJson(systemTaskId: 1, recordedAt: at(DateTime(2025, 12, 30, 12))),
      entryJson(systemTaskId: 0, recordedAt: 'not a date'),
    ].map(entryFromJson);

    test('keeps the server order within and across days', () {
      final days = groupHistoryByDay(entries, now: DateTime(2026, 11, 5));

      expect(days.map((day) => day.entries.map((e) => e.systemTaskId).toList()), [
        [3, 2],
        [1],
        [0],
      ]);
    });

    test('names the year only outside the current one', () {
      final days = groupHistoryByDay(entries, now: DateTime(2026, 11, 5));

      expect(days.map((day) => day.label), [
        'Tuesday, Oct 13',
        'Tuesday, Dec 30, 2025',
        'Unknown date',
      ]);
      expect(days.last.key, unknownDayKey);
    });
  });

  group('formatHistoryValue', () {
    String format(String type, {String? value, int? userId}) =>
        formatHistoryValue(_value(type, value: value, userId: userId), _formatters);

    test('formats each value type the server sends', () {
      expect(format('TEXT', value: 'Costco'), 'Costco');
      expect(format('CURRENCY', value: '12.00'), '\$12.00');
      expect(format('USER', userId: 2), 'Maya Chen');
      expect(format('USER', userId: 9), unknownUserName);
      expect(format('STATUS', value: 'NEEDS_ATTENTION'), 'Needs Attention');
      expect(format('BOOLEAN', value: 'true'), 'Yes');
      expect(format('BOOLEAN', value: 'false'), 'No');
      expect(format('INSTANT', value: '2026-10-20T16:52:31Z'), contains('2026'));
    });

    // A calendar day is stored as midnight UTC. Converting it to local time
    // would show the previous day anywhere west of UTC.
    test('shows a calendar day as that day, whatever the device zone', () {
      expect(format('DATE', value: '2026-10-11'), 'Oct 11, 2026');
      expect(format('DATE', value: '2026-10-11T00:00:00Z'), 'Oct 11, 2026');
    });

    test('shows the raw value for a type this build does not know', () {
      expect(format('COLOUR', value: '#ff0000'), '#ff0000');
    });
  });

  group('buildHistoryChangeRow', () {
    HistoryChangeRow row(Map<String, Object?> json) =>
        buildHistoryChangeRow(changeFromJson(json), _formatters);

    test('a change shows both sides with the arrow between', () {
      final changed = row(textChangeJson('NAME', from: 'COSTCO WHSE', to: 'Costco Wholesale'));
      expect(changed.label, 'Name');
      expect(changed.from, 'COSTCO WHSE');
      expect(changed.to, 'Costco Wholesale');
      expect(changed.arrow, isTrue);
    });

    test('an addition shows only the new side, a removal only the old', () {
      final added = row(textChangeJson('TAGS', to: 'Bulk'));
      expect(added.from, isNull);
      expect(added.to, 'Bulk');
      expect(added.arrow, isFalse);

      final removed = row(textChangeJson('IMAGE', from: 'costco.jpg'));
      expect(removed.from, 'costco.jpg');
      expect(removed.to, isNull);
    });

    test('an empty value still reads as something', () {
      expect(row(textChangeJson('NAME', from: '', to: 'Lunch')).from, emptyValueText);
    });

    test('a comment names its author and clips long text', () {
      final long = 'x' * 80;
      final comment = row({
        'field': 'COMMENT',
        'kind': 'ADDED',
        'userId': 2,
        'to': {'type': 'TEXT', 'value': long},
      });
      expect(comment.to, 'Maya Chen: “${'x' * (historyCommentClipLength - 3)}…”');
    });
  });

  test('parseAvatarColor reads a hex colour and ignores anything else', () {
    expect(parseAvatarColor('#a78bfa'), const Color(0xFFA78BFA));
    expect(parseAvatarColor('a78bfa'), const Color(0xFFA78BFA));
    expect(parseAvatarColor('#abc'), isNull);
    expect(parseAvatarColor(null), isNull);
  });
}
