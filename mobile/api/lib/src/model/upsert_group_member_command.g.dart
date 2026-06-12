// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_group_member_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertGroupMemberCommand extends UpsertGroupMemberCommand {
  @override
  final int groupId;
  @override
  final int? groupRoleId;
  @override
  final int userId;

  factory _$UpsertGroupMemberCommand(
          [void Function(UpsertGroupMemberCommandBuilder)? updates]) =>
      (UpsertGroupMemberCommandBuilder()..update(updates))._build();

  _$UpsertGroupMemberCommand._(
      {required this.groupId, this.groupRoleId, required this.userId})
      : super._();
  @override
  UpsertGroupMemberCommand rebuild(
          void Function(UpsertGroupMemberCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertGroupMemberCommandBuilder toBuilder() =>
      UpsertGroupMemberCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertGroupMemberCommand &&
        groupId == other.groupId &&
        groupRoleId == other.groupRoleId &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, groupRoleId.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertGroupMemberCommand')
          ..add('groupId', groupId)
          ..add('groupRoleId', groupRoleId)
          ..add('userId', userId))
        .toString();
  }
}

class UpsertGroupMemberCommandBuilder
    implements
        Builder<UpsertGroupMemberCommand, UpsertGroupMemberCommandBuilder> {
  _$UpsertGroupMemberCommand? _$v;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _groupRoleId;
  int? get groupRoleId => _$this._groupRoleId;
  set groupRoleId(int? groupRoleId) => _$this._groupRoleId = groupRoleId;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  UpsertGroupMemberCommandBuilder() {
    UpsertGroupMemberCommand._defaults(this);
  }

  UpsertGroupMemberCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _groupId = $v.groupId;
      _groupRoleId = $v.groupRoleId;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertGroupMemberCommand other) {
    _$v = other as _$UpsertGroupMemberCommand;
  }

  @override
  void update(void Function(UpsertGroupMemberCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertGroupMemberCommand build() => _build();

  _$UpsertGroupMemberCommand _build() {
    final _$result = _$v ??
        _$UpsertGroupMemberCommand._(
          groupId: BuiltValueNullFieldError.checkNotNull(
              groupId, r'UpsertGroupMemberCommand', 'groupId'),
          groupRoleId: groupRoleId,
          userId: BuiltValueNullFieldError.checkNotNull(
              userId, r'UpsertGroupMemberCommand', 'userId'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
