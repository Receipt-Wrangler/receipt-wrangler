// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bulk_user_delete_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BulkUserDeleteCommand extends BulkUserDeleteCommand {
  @override
  final BuiltList<String> userIds;

  factory _$BulkUserDeleteCommand(
          [void Function(BulkUserDeleteCommandBuilder)? updates]) =>
      (new BulkUserDeleteCommandBuilder()..update(updates))._build();

  _$BulkUserDeleteCommand._({required this.userIds}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        userIds, r'BulkUserDeleteCommand', 'userIds');
  }

  @override
  BulkUserDeleteCommand rebuild(
          void Function(BulkUserDeleteCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BulkUserDeleteCommandBuilder toBuilder() =>
      new BulkUserDeleteCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BulkUserDeleteCommand && userIds == other.userIds;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, userIds.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BulkUserDeleteCommand')
          ..add('userIds', userIds))
        .toString();
  }
}

class BulkUserDeleteCommandBuilder
    implements Builder<BulkUserDeleteCommand, BulkUserDeleteCommandBuilder> {
  _$BulkUserDeleteCommand? _$v;

  ListBuilder<String>? _userIds;
  ListBuilder<String> get userIds =>
      _$this._userIds ??= new ListBuilder<String>();
  set userIds(ListBuilder<String>? userIds) => _$this._userIds = userIds;

  BulkUserDeleteCommandBuilder() {
    BulkUserDeleteCommand._defaults(this);
  }

  BulkUserDeleteCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _userIds = $v.userIds.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BulkUserDeleteCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$BulkUserDeleteCommand;
  }

  @override
  void update(void Function(BulkUserDeleteCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BulkUserDeleteCommand build() => _build();

  _$BulkUserDeleteCommand _build() {
    _$BulkUserDeleteCommand _$result;
    try {
      _$result = _$v ?? new _$BulkUserDeleteCommand._(userIds: userIds.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'userIds';
        userIds.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'BulkUserDeleteCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
