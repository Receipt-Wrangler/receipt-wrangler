import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:openapi/openapi.dart' as api;

import '../constants/colors.dart';
import 'receipts.dart' show titleCaseStatusName;

// Pure formatting for the receipt history screen. The server computes the diff
// and the redaction (api/CLAUDE.md -> "Receipt history"); this only turns its
// typed rows into text. Every code on the contract is a plain string, so each
// helper has a fallback for one this build does not know.
//
// The wording mirrors the desktop drawer's
// (desktop/src/receipts/receipt-history-drawer/receipt-history.utils.ts), so
// the two clients describe the same change the same way. Nothing ties them
// together mechanically; keep them in step by hand.

/// Change rows shown before "Show N more".
const historyCollapseAfter = 3;

/// A comment chip is clipped to this many characters.
const historyCommentClipLength = 60;

const unknownUserName = 'Unknown user';

/// A value that is present but empty, so its chip is never blank.
const emptyValueText = '(empty)';

const defaultActionVerb = 'updated this receipt';

typedef UserLookup = api.UserView? Function(int? userId);

String userDisplayName(api.UserView? user) {
  if (user == null) {
    return unknownUserName;
  }
  if (user.displayName.isNotEmpty) {
    return user.displayName;
  }
  return user.username.isNotEmpty ? user.username : unknownUserName;
}

/// `SOME_CODE` -> `Some Code`: the label for a code this build has no wording for.
String titleCaseCode(String? code) {
  return code == null || code.isEmpty ? '' : titleCaseStatusName(code);
}

const _actionVerbs = {
  'CREATED': 'created this receipt',
  'UPDATED': 'edited this receipt',
  'STATUS_CHANGED': 'changed the status',
  'RESOLVED': 'resolved this receipt',
  'IMAGES_COMMENTS_UPDATED': 'updated images and comments',
  'FAILED': 'tried to save changes',
};

String historyActionVerb(String? action) {
  return _actionVerbs[action] ?? defaultActionVerb;
}

// Named after the receipt form's own fields, so a row reads like the field it
// changed ("Receipt Date", not a bare "Date" beside "Resolved Date").
const _fieldLabels = {
  'NAME': 'Name',
  'AMOUNT': 'Amount',
  'DATE': 'Receipt Date',
  'GROUP': 'Group',
  'PAID_BY': 'Paid By',
  'STATUS': 'Status',
  'RESOLVED_DATE': 'Resolved Date',
  'CATEGORIES': 'Categories',
  'TAGS': 'Tags',
  'COMMENT': 'Comment',
  'IMAGE': 'Image',
  'EMAIL_INBOX': 'Inbox',
  'EMAIL_SENDER': 'Sender',
};

String _withSubject(String prefix, String? subject) {
  return subject == null || subject.isEmpty ? prefix : '$prefix: $subject';
}

String historyFieldLabel(api.ReceiptHistoryChange change, UserLookup findUser) {
  switch (change.field) {
    case 'CUSTOM_FIELD':
      return change.label?.isNotEmpty == true ? change.label! : 'Custom field';
    case 'ITEM':
      return _withSubject('Item', change.label);
    case 'SHARE':
      return _withSubject('Share', userDisplayName(findUser(change.userId)));
    case 'ITEM_CATEGORIES':
      return _withSubject('Item categories', change.label);
    case 'ITEM_TAGS':
      return _withSubject('Item tags', change.label);
    default:
      return _fieldLabels[change.field] ?? titleCaseCode(change.field);
  }
}

enum HistoryActorKind { user, quickScan, email, system }

/// Who made an entry, as the card's header shows them.
class HistoryActorDisplay {
  const HistoryActorDisplay({
    required this.kind,
    required this.name,
    this.user,
    this.icon,
    this.color,
  });

  final HistoryActorKind kind;
  final String name;

  /// The user behind the avatar, when known, so it shows their colour and initial.
  final api.UserView? user;

  /// Shown in place of an initial.
  final IconData? icon;

  /// The avatar's background; null for a known user, whose own colour applies.
  final Color? color;
}

HistoryActorDisplay historyActorDisplay(
  api.ReceiptHistoryActor actor,
  UserLookup findUser,
) {
  switch (actor.type) {
    case 'USER':
      final user = findUser(actor.userId);
      return user != null
          ? HistoryActorDisplay(
              kind: HistoryActorKind.user,
              name: userDisplayName(user),
              user: user,
            )
          : const HistoryActorDisplay(
              kind: HistoryActorKind.user,
              name: unknownUserName,
              icon: Icons.person,
              color: slate400,
            );
    case 'QUICK_SCAN':
      return const HistoryActorDisplay(
        kind: HistoryActorKind.quickScan,
        name: 'Quick Scan',
        icon: Icons.photo_camera,
        color: historyQuickScanActor,
      );
    case 'EMAIL':
      return const HistoryActorDisplay(
        kind: HistoryActorKind.email,
        name: 'Email Upload',
        icon: Icons.mail,
        color: slate500,
      );
    default:
      return const HistoryActorDisplay(
        kind: HistoryActorKind.system,
        name: 'System',
        icon: Icons.settings,
        color: slate400,
      );
  }
}

String _counted(int count, String singular, String plural) {
  return '$count ${count == 1 ? singular : plural}';
}

String historySummary(api.ReceiptHistory history) {
  return [
    _counted(history.changeCount, 'change', 'changes'),
    _counted(history.saveCount, 'save', 'saves'),
    _counted(history.contributorCount, 'contributor', 'contributors'),
  ].join(' · ');
}

/// How far an older entry's "before" can be trusted. COMPLETE, absent and
/// unknown states need no notice.
String? historyBeforeStateNotice(api.ReceiptHistoryEntry entry) {
  switch (entry.beforeState) {
    case 'REBUILT':
      final source = switch (entry.beforeSourceType) {
        'RECEIPT_UPDATED' => 'the previous update',
        'RECEIPT_UPLOADED' => 'the copy saved when it was created',
        _ => 'an earlier copy',
      };
      final date = _formatInstant(entry.beforeSourceRecordedAt, 'MMM d, yyyy');
      final dated = date == null ? '' : ' ($date)';
      return 'Recorded by an older version — compared against $source$dated. '
          'Bulk status changes can also show here.';
    case 'INCOMPLETE':
      return 'Recorded by an older version with no earlier complete copy. '
          'Some changes may be inaccurate.';
    default:
      return null;
  }
}

/// One day of entries, in the device's time zone.
class HistoryDay {
  const HistoryDay({required this.key, required this.label, required this.entries});

  /// `yyyy-MM-dd`, or [unknownDayKey].
  final String key;
  final String label;
  final List<api.ReceiptHistoryEntry> entries;
}

const unknownDayKey = 'unknown';

/// Groups entries by the device-local day they were recorded on, keeping the
/// server's order (newest first) across and within days. The label names the
/// year only outside the current one, as the design's "Tuesday, Oct 13" does.
///
/// Device-local, not the app time zone: mobile does not read
/// `AppData.timeZone` anywhere yet (mobile/CLAUDE.md -> "App time zone").
List<HistoryDay> groupHistoryByDay(
  Iterable<api.ReceiptHistoryEntry> entries, {
  DateTime? now,
}) {
  final currentYear = (now ?? DateTime.now()).year;
  final days = <String, HistoryDay>{};

  for (final entry in entries) {
    final recorded = DateTime.tryParse(entry.recordedAt)?.toLocal();
    final key = recorded == null
        ? unknownDayKey
        : DateFormat('yyyy-MM-dd').format(recorded);
    final day = days.putIfAbsent(key, () {
      final label = recorded == null
          ? 'Unknown date'
          : DateFormat(recorded.year == currentYear ? 'EEEE, MMM d' : 'EEEE, MMM d, y')
              .format(recorded);
      return HistoryDay(key: key, label: label, entries: []);
    });
    day.entries.add(entry);
  }

  return days.values.toList();
}

/// The time an entry was recorded, device-local: `4:52 PM`.
String historyTime(String recordedAt) {
  return _formatInstant(recordedAt, 'h:mm a') ?? '';
}

String? _formatInstant(String? value, String pattern) {
  final parsed = value == null ? null : DateTime.tryParse(value);
  return parsed == null ? null : DateFormat(pattern).format(parsed.toLocal());
}

/// A calendar day (`yyyy-MM-dd`) shown as that day, whatever the device zone.
///
/// Never through `toLocal()` (which `formatDate` in `utils/date.dart` does): a
/// stored day is midnight UTC, and converting it moves it to the previous day
/// anywhere west of UTC.
String? formatHistoryCalendarDate(String value) {
  final parsed = DateTime.tryParse(value);
  if (parsed == null) {
    return null;
  }
  return DateFormat('MMM d, yyyy')
      .format(DateTime(parsed.year, parsed.month, parsed.day));
}

/// How each value type is turned into text. The currency formatter needs a
/// `BuildContext`, so the caller injects it along with the user lookup.
class HistoryValueFormatters {
  const HistoryValueFormatters({required this.currency, required this.findUser});

  final String Function(String value) currency;
  final UserLookup findUser;
}

String formatHistoryValue(
  api.ReceiptHistoryValue value,
  HistoryValueFormatters formatters,
) {
  final raw = value.value ?? '';
  switch (value.type) {
    case 'TEXT':
      return raw;
    case 'CURRENCY':
      return raw.isEmpty ? '' : formatters.currency(raw);
    case 'DATE':
      return raw.isEmpty ? '' : (formatHistoryCalendarDate(raw) ?? raw);
    case 'INSTANT':
      return raw.isEmpty ? '' : (_formatInstant(raw, 'MMM d, yyyy, h:mm a') ?? raw);
    case 'USER':
      return userDisplayName(formatters.findUser(value.userId));
    case 'STATUS':
      return titleCaseCode(raw);
    case 'BOOLEAN':
      return raw == 'true' ? 'Yes' : (raw == 'false' ? 'No' : raw);
    default:
      return raw;
  }
}

/// One change as the card renders it.
class HistoryChangeRow {
  const HistoryChangeRow({
    required this.field,
    required this.kind,
    required this.label,
    this.from,
    this.to,
  });

  /// The raw codes.
  final String field;
  final String kind;

  final String label;
  final String? from;
  final String? to;

  /// The arrow sits between two values.
  bool get arrow => from != null && to != null;
}

String _clip(String text) {
  return text.length > historyCommentClipLength
      ? '${text.substring(0, historyCommentClipLength - 3)}…'
      : text;
}

String? _chipText(
  api.ReceiptHistoryChange change,
  api.ReceiptHistoryValue? value,
  HistoryValueFormatters formatters,
) {
  if (value == null) {
    return null;
  }

  final text = formatHistoryValue(value, formatters);
  if (change.field != 'COMMENT') {
    return text.isEmpty ? emptyValueText : text;
  }

  // A comment reads as its author and its clipped words: Maya Chen: “Took the bag fee off…”
  final author =
      change.userId != null ? '${userDisplayName(formatters.findUser(change.userId))}: ' : '';
  return '$author“${_clip(text)}”';
}

/// ADDED shows only the new value and REMOVED only the old; CHANGED (and a kind
/// this build does not know) shows whichever sides were sent.
HistoryChangeRow buildHistoryChangeRow(
  api.ReceiptHistoryChange change,
  HistoryValueFormatters formatters,
) {
  return HistoryChangeRow(
    field: change.field,
    kind: change.kind,
    label: historyFieldLabel(change, formatters.findUser),
    from: change.kind == 'ADDED' ? null : _chipText(change, change.from, formatters),
    to: change.kind == 'REMOVED' ? null : _chipText(change, change.to, formatters),
  );
}

/// A user's `defaultAvatarColor` (`#RRGGBB`), or null when absent or malformed.
Color? parseAvatarColor(String? hex) {
  if (hex == null) {
    return null;
  }
  final digits = hex.startsWith('#') ? hex.substring(1) : hex;
  if (digits.length != 6) {
    return null;
  }
  final value = int.tryParse(digits, radix: 16);
  return value == null ? null : Color(0xFF000000 | value);
}
