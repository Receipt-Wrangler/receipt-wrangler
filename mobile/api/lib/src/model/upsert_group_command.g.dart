// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_group_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertGroupCommand extends UpsertGroupCommand {
  @override
  final BuiltList<UpsertGroupMemberCommand> groupMembers;
  @override
  final bool? isDefault;
  @override
  final String name;
  @override
  final bool? isAllGroup;
  @override
  final GroupStatus status;

  factory _$UpsertGroupCommand(
          [void Function(UpsertGroupCommandBuilder)? updates]) =>
      (new UpsertGroupCommandBuilder()..update(updates))._build();

  _$UpsertGroupCommand._(
      {required this.groupMembers,
      this.isDefault,
      required this.name,
      this.isAllGroup,
      required this.status})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        groupMembers, r'UpsertGroupCommand', 'groupMembers');
    BuiltValueNullFieldError.checkNotNull(name, r'UpsertGroupCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        status, r'UpsertGroupCommand', 'status');
  }

  @override
  UpsertGroupCommand rebuild(
          void Function(UpsertGroupCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertGroupCommandBuilder toBuilder() =>
      new UpsertGroupCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertGroupCommand &&
        groupMembers == other.groupMembers &&
        isDefault == other.isDefault &&
        name == other.name &&
        isAllGroup == other.isAllGroup &&
        status == other.status;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, groupMembers.hashCode);
    _$hash = $jc(_$hash, isDefault.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, isAllGroup.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertGroupCommand')
          ..add('groupMembers', groupMembers)
          ..add('isDefault', isDefault)
          ..add('name', name)
          ..add('isAllGroup', isAllGroup)
          ..add('status', status))
        .toString();
  }
}

class UpsertGroupCommandBuilder
    implements Builder<UpsertGroupCommand, UpsertGroupCommandBuilder> {
  _$UpsertGroupCommand? _$v;

  ListBuilder<UpsertGroupMemberCommand>? _groupMembers;
  ListBuilder<UpsertGroupMemberCommand> get groupMembers =>
      _$this._groupMembers ??= new ListBuilder<UpsertGroupMemberCommand>();
  set groupMembers(ListBuilder<UpsertGroupMemberCommand>? groupMembers) =>
      _$this._groupMembers = groupMembers;

  bool? _isDefault;
  bool? get isDefault => _$this._isDefault;
  set isDefault(bool? isDefault) => _$this._isDefault = isDefault;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  bool? _isAllGroup;
  bool? get isAllGroup => _$this._isAllGroup;
  set isAllGroup(bool? isAllGroup) => _$this._isAllGroup = isAllGroup;

  GroupStatus? _status;
  GroupStatus? get status => _$this._status;
  set status(GroupStatus? status) => _$this._status = status;

  UpsertGroupCommandBuilder() {
    UpsertGroupCommand._defaults(this);
  }

  UpsertGroupCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _groupMembers = $v.groupMembers.toBuilder();
      _isDefault = $v.isDefault;
      _name = $v.name;
      _isAllGroup = $v.isAllGroup;
      _status = $v.status;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertGroupCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertGroupCommand;
  }

  @override
  void update(void Function(UpsertGroupCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertGroupCommand build() => _build();

  _$UpsertGroupCommand _build() {
    _$UpsertGroupCommand _$result;
    try {
      _$result = _$v ??
          new _$UpsertGroupCommand._(
              groupMembers: groupMembers.build(),
              isDefault: isDefault,
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'UpsertGroupCommand', 'name'),
              isAllGroup: isAllGroup,
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'UpsertGroupCommand', 'status'));
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'groupMembers';
        groupMembers.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UpsertGroupCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
