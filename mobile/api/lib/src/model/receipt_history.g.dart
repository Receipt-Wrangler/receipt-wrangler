// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_history.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptHistory extends ReceiptHistory {
  @override
  final BuiltList<ReceiptHistoryEntry> entries;
  @override
  final int changeCount;
  @override
  final int saveCount;
  @override
  final int contributorCount;

  factory _$ReceiptHistory([void Function(ReceiptHistoryBuilder)? updates]) =>
      (ReceiptHistoryBuilder()..update(updates))._build();

  _$ReceiptHistory._(
      {required this.entries,
      required this.changeCount,
      required this.saveCount,
      required this.contributorCount})
      : super._();
  @override
  ReceiptHistory rebuild(void Function(ReceiptHistoryBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptHistoryBuilder toBuilder() => ReceiptHistoryBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptHistory &&
        entries == other.entries &&
        changeCount == other.changeCount &&
        saveCount == other.saveCount &&
        contributorCount == other.contributorCount;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, entries.hashCode);
    _$hash = $jc(_$hash, changeCount.hashCode);
    _$hash = $jc(_$hash, saveCount.hashCode);
    _$hash = $jc(_$hash, contributorCount.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptHistory')
          ..add('entries', entries)
          ..add('changeCount', changeCount)
          ..add('saveCount', saveCount)
          ..add('contributorCount', contributorCount))
        .toString();
  }
}

class ReceiptHistoryBuilder
    implements Builder<ReceiptHistory, ReceiptHistoryBuilder> {
  _$ReceiptHistory? _$v;

  ListBuilder<ReceiptHistoryEntry>? _entries;
  ListBuilder<ReceiptHistoryEntry> get entries =>
      _$this._entries ??= ListBuilder<ReceiptHistoryEntry>();
  set entries(ListBuilder<ReceiptHistoryEntry>? entries) =>
      _$this._entries = entries;

  int? _changeCount;
  int? get changeCount => _$this._changeCount;
  set changeCount(int? changeCount) => _$this._changeCount = changeCount;

  int? _saveCount;
  int? get saveCount => _$this._saveCount;
  set saveCount(int? saveCount) => _$this._saveCount = saveCount;

  int? _contributorCount;
  int? get contributorCount => _$this._contributorCount;
  set contributorCount(int? contributorCount) =>
      _$this._contributorCount = contributorCount;

  ReceiptHistoryBuilder() {
    ReceiptHistory._defaults(this);
  }

  ReceiptHistoryBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _entries = $v.entries.toBuilder();
      _changeCount = $v.changeCount;
      _saveCount = $v.saveCount;
      _contributorCount = $v.contributorCount;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptHistory other) {
    _$v = other as _$ReceiptHistory;
  }

  @override
  void update(void Function(ReceiptHistoryBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptHistory build() => _build();

  _$ReceiptHistory _build() {
    _$ReceiptHistory _$result;
    try {
      _$result = _$v ??
          _$ReceiptHistory._(
            entries: entries.build(),
            changeCount: BuiltValueNullFieldError.checkNotNull(
                changeCount, r'ReceiptHistory', 'changeCount'),
            saveCount: BuiltValueNullFieldError.checkNotNull(
                saveCount, r'ReceiptHistory', 'saveCount'),
            contributorCount: BuiltValueNullFieldError.checkNotNull(
                contributorCount, r'ReceiptHistory', 'contributorCount'),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'entries';
        entries.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptHistory', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
