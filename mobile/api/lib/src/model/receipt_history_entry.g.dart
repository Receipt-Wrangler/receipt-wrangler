// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_history_entry.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptHistoryEntry extends ReceiptHistoryEntry {
  @override
  final int systemTaskId;
  @override
  final String action;
  @override
  final String recordedAt;
  @override
  final ReceiptHistoryActor actor;
  @override
  final bool hasSnapshots;
  @override
  final String? errorMessage;
  @override
  final String? beforeState;
  @override
  final String? beforeSourceType;
  @override
  final String? beforeSourceRecordedAt;
  @override
  final BuiltList<ReceiptHistoryChange> changes;

  factory _$ReceiptHistoryEntry(
          [void Function(ReceiptHistoryEntryBuilder)? updates]) =>
      (ReceiptHistoryEntryBuilder()..update(updates))._build();

  _$ReceiptHistoryEntry._(
      {required this.systemTaskId,
      required this.action,
      required this.recordedAt,
      required this.actor,
      required this.hasSnapshots,
      this.errorMessage,
      this.beforeState,
      this.beforeSourceType,
      this.beforeSourceRecordedAt,
      required this.changes})
      : super._();
  @override
  ReceiptHistoryEntry rebuild(
          void Function(ReceiptHistoryEntryBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptHistoryEntryBuilder toBuilder() =>
      ReceiptHistoryEntryBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptHistoryEntry &&
        systemTaskId == other.systemTaskId &&
        action == other.action &&
        recordedAt == other.recordedAt &&
        actor == other.actor &&
        hasSnapshots == other.hasSnapshots &&
        errorMessage == other.errorMessage &&
        beforeState == other.beforeState &&
        beforeSourceType == other.beforeSourceType &&
        beforeSourceRecordedAt == other.beforeSourceRecordedAt &&
        changes == other.changes;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, systemTaskId.hashCode);
    _$hash = $jc(_$hash, action.hashCode);
    _$hash = $jc(_$hash, recordedAt.hashCode);
    _$hash = $jc(_$hash, actor.hashCode);
    _$hash = $jc(_$hash, hasSnapshots.hashCode);
    _$hash = $jc(_$hash, errorMessage.hashCode);
    _$hash = $jc(_$hash, beforeState.hashCode);
    _$hash = $jc(_$hash, beforeSourceType.hashCode);
    _$hash = $jc(_$hash, beforeSourceRecordedAt.hashCode);
    _$hash = $jc(_$hash, changes.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptHistoryEntry')
          ..add('systemTaskId', systemTaskId)
          ..add('action', action)
          ..add('recordedAt', recordedAt)
          ..add('actor', actor)
          ..add('hasSnapshots', hasSnapshots)
          ..add('errorMessage', errorMessage)
          ..add('beforeState', beforeState)
          ..add('beforeSourceType', beforeSourceType)
          ..add('beforeSourceRecordedAt', beforeSourceRecordedAt)
          ..add('changes', changes))
        .toString();
  }
}

class ReceiptHistoryEntryBuilder
    implements Builder<ReceiptHistoryEntry, ReceiptHistoryEntryBuilder> {
  _$ReceiptHistoryEntry? _$v;

  int? _systemTaskId;
  int? get systemTaskId => _$this._systemTaskId;
  set systemTaskId(int? systemTaskId) => _$this._systemTaskId = systemTaskId;

  String? _action;
  String? get action => _$this._action;
  set action(String? action) => _$this._action = action;

  String? _recordedAt;
  String? get recordedAt => _$this._recordedAt;
  set recordedAt(String? recordedAt) => _$this._recordedAt = recordedAt;

  ReceiptHistoryActorBuilder? _actor;
  ReceiptHistoryActorBuilder get actor =>
      _$this._actor ??= ReceiptHistoryActorBuilder();
  set actor(ReceiptHistoryActorBuilder? actor) => _$this._actor = actor;

  bool? _hasSnapshots;
  bool? get hasSnapshots => _$this._hasSnapshots;
  set hasSnapshots(bool? hasSnapshots) => _$this._hasSnapshots = hasSnapshots;

  String? _errorMessage;
  String? get errorMessage => _$this._errorMessage;
  set errorMessage(String? errorMessage) => _$this._errorMessage = errorMessage;

  String? _beforeState;
  String? get beforeState => _$this._beforeState;
  set beforeState(String? beforeState) => _$this._beforeState = beforeState;

  String? _beforeSourceType;
  String? get beforeSourceType => _$this._beforeSourceType;
  set beforeSourceType(String? beforeSourceType) =>
      _$this._beforeSourceType = beforeSourceType;

  String? _beforeSourceRecordedAt;
  String? get beforeSourceRecordedAt => _$this._beforeSourceRecordedAt;
  set beforeSourceRecordedAt(String? beforeSourceRecordedAt) =>
      _$this._beforeSourceRecordedAt = beforeSourceRecordedAt;

  ListBuilder<ReceiptHistoryChange>? _changes;
  ListBuilder<ReceiptHistoryChange> get changes =>
      _$this._changes ??= ListBuilder<ReceiptHistoryChange>();
  set changes(ListBuilder<ReceiptHistoryChange>? changes) =>
      _$this._changes = changes;

  ReceiptHistoryEntryBuilder() {
    ReceiptHistoryEntry._defaults(this);
  }

  ReceiptHistoryEntryBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _systemTaskId = $v.systemTaskId;
      _action = $v.action;
      _recordedAt = $v.recordedAt;
      _actor = $v.actor.toBuilder();
      _hasSnapshots = $v.hasSnapshots;
      _errorMessage = $v.errorMessage;
      _beforeState = $v.beforeState;
      _beforeSourceType = $v.beforeSourceType;
      _beforeSourceRecordedAt = $v.beforeSourceRecordedAt;
      _changes = $v.changes.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptHistoryEntry other) {
    _$v = other as _$ReceiptHistoryEntry;
  }

  @override
  void update(void Function(ReceiptHistoryEntryBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptHistoryEntry build() => _build();

  _$ReceiptHistoryEntry _build() {
    _$ReceiptHistoryEntry _$result;
    try {
      _$result = _$v ??
          _$ReceiptHistoryEntry._(
            systemTaskId: BuiltValueNullFieldError.checkNotNull(
                systemTaskId, r'ReceiptHistoryEntry', 'systemTaskId'),
            action: BuiltValueNullFieldError.checkNotNull(
                action, r'ReceiptHistoryEntry', 'action'),
            recordedAt: BuiltValueNullFieldError.checkNotNull(
                recordedAt, r'ReceiptHistoryEntry', 'recordedAt'),
            actor: actor.build(),
            hasSnapshots: BuiltValueNullFieldError.checkNotNull(
                hasSnapshots, r'ReceiptHistoryEntry', 'hasSnapshots'),
            errorMessage: errorMessage,
            beforeState: beforeState,
            beforeSourceType: beforeSourceType,
            beforeSourceRecordedAt: beforeSourceRecordedAt,
            changes: changes.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'actor';
        actor.build();

        _$failedField = 'changes';
        changes.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptHistoryEntry', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
