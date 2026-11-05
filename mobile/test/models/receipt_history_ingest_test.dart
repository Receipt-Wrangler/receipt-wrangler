import 'package:flutter_test/flutter_test.dart';

import '../helpers/receipt_history_test_helpers.dart';

/// Guards the receipt history contract at the API boundary.
///
/// Every code on it -- the entry's action and before-state, the actor type,
/// a change's field and kind, a value's type -- is a plain `String` in
/// `swagger.yml`, never an enum. A closed built_value `EnumClass` throws on a
/// value it does not know, which would fail the whole history payload on every
/// already-released build the first time the server learns a new change type:
/// the mechanism behind the two `Permission` login outages. This parses codes no
/// build knows through the real generated deserializer, so it fails if a regen
/// ever closes one of them.
void main() {
  test('codes this build does not know deserialize as plain strings', () {
    final history = historyFromJson(historyJson([
      entryJson(
        action: 'ARCHIVED',
        actor: {'type': 'ROBOT'},
        beforeState: 'GUESSED',
        changes: [
          {
            'field': 'SPLIT_RULE',
            'kind': 'REORDERED',
            'from': {'type': 'PERCENT', 'value': '10'},
            'to': {'type': 'PERCENT', 'value': '20'},
          },
        ],
      ),
    ], changeCount: 1, saveCount: 1, contributorCount: 1));

    final entry = history.entries.single;
    expect(entry.action, 'ARCHIVED');
    expect(entry.actor.type, 'ROBOT');
    expect(entry.beforeState, 'GUESSED');
    expect(entry.changes.single.field, 'SPLIT_RULE');
    expect(entry.changes.single.kind, 'REORDERED');
    expect(entry.changes.single.to?.type, 'PERCENT');
  });

  test('an empty history parses with its entries present and empty', () {
    final history = historyFromJson(historyJson([]));

    expect(history.entries, isEmpty);
  });
}
