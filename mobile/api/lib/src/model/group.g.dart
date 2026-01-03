// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Group extends Group {
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final GroupSettings? groupSettings;
  @override
  final GroupReceiptSettings groupReceiptSettings;
  @override
  final BuiltList<GroupMember> groupMembers;
  @override
  final int id;
  @override
  final bool? isDefault;
  @override
  final String name;
  @override
  final bool isAllGroup;
  @override
  final GroupStatus status;
  @override
  final String? updatedAt;

  factory _$Group([void Function(GroupBuilder)? updates]) =>
      (new GroupBuilder()..update(updates))._build();

  _$Group._(
      {this.createdAt,
      this.createdBy,
      this.groupSettings,
      required this.groupReceiptSettings,
      required this.groupMembers,
      required this.id,
      this.isDefault,
      required this.name,
      required this.isAllGroup,
      required this.status,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        groupReceiptSettings, r'Group', 'groupReceiptSettings');
    BuiltValueNullFieldError.checkNotNull(
        groupMembers, r'Group', 'groupMembers');
    BuiltValueNullFieldError.checkNotNull(id, r'Group', 'id');
    BuiltValueNullFieldError.checkNotNull(name, r'Group', 'name');
    BuiltValueNullFieldError.checkNotNull(isAllGroup, r'Group', 'isAllGroup');
    BuiltValueNullFieldError.checkNotNull(status, r'Group', 'status');
  }

  @override
  Group rebuild(void Function(GroupBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GroupBuilder toBuilder() => new GroupBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Group &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        groupSettings == other.groupSettings &&
        groupReceiptSettings == other.groupReceiptSettings &&
        groupMembers == other.groupMembers &&
        id == other.id &&
        isDefault == other.isDefault &&
        name == other.name &&
        isAllGroup == other.isAllGroup &&
        status == other.status &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, groupSettings.hashCode);
    _$hash = $jc(_$hash, groupReceiptSettings.hashCode);
    _$hash = $jc(_$hash, groupMembers.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, isDefault.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, isAllGroup.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Group')
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('groupSettings', groupSettings)
          ..add('groupReceiptSettings', groupReceiptSettings)
          ..add('groupMembers', groupMembers)
          ..add('id', id)
          ..add('isDefault', isDefault)
          ..add('name', name)
          ..add('isAllGroup', isAllGroup)
          ..add('status', status)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class GroupBuilder implements Builder<Group, GroupBuilder> {
  _$Group? _$v;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  GroupSettingsBuilder? _groupSettings;
  GroupSettingsBuilder get groupSettings =>
      _$this._groupSettings ??= new GroupSettingsBuilder();
  set groupSettings(GroupSettingsBuilder? groupSettings) =>
      _$this._groupSettings = groupSettings;

  GroupReceiptSettingsBuilder? _groupReceiptSettings;
  GroupReceiptSettingsBuilder get groupReceiptSettings =>
      _$this._groupReceiptSettings ??= new GroupReceiptSettingsBuilder();
  set groupReceiptSettings(GroupReceiptSettingsBuilder? groupReceiptSettings) =>
      _$this._groupReceiptSettings = groupReceiptSettings;

  ListBuilder<GroupMember>? _groupMembers;
  ListBuilder<GroupMember> get groupMembers =>
      _$this._groupMembers ??= new ListBuilder<GroupMember>();
  set groupMembers(ListBuilder<GroupMember>? groupMembers) =>
      _$this._groupMembers = groupMembers;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

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

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  GroupBuilder() {
    Group._defaults(this);
  }

  GroupBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _groupSettings = $v.groupSettings?.toBuilder();
      _groupReceiptSettings = $v.groupReceiptSettings.toBuilder();
      _groupMembers = $v.groupMembers.toBuilder();
      _id = $v.id;
      _isDefault = $v.isDefault;
      _name = $v.name;
      _isAllGroup = $v.isAllGroup;
      _status = $v.status;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Group other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Group;
  }

  @override
  void update(void Function(GroupBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Group build() => _build();

  _$Group _build() {
    _$Group _$result;
    try {
      _$result = _$v ??
          new _$Group._(
              createdAt: createdAt,
              createdBy: createdBy,
              groupSettings: _groupSettings?.build(),
              groupReceiptSettings: groupReceiptSettings.build(),
              groupMembers: groupMembers.build(),
              id: BuiltValueNullFieldError.checkNotNull(id, r'Group', 'id'),
              isDefault: isDefault,
              name:
                  BuiltValueNullFieldError.checkNotNull(name, r'Group', 'name'),
              isAllGroup: BuiltValueNullFieldError.checkNotNull(
                  isAllGroup, r'Group', 'isAllGroup'),
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'Group', 'status'),
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'groupSettings';
        _groupSettings?.build();
        _$failedField = 'groupReceiptSettings';
        groupReceiptSettings.build();
        _$failedField = 'groupMembers';
        groupMembers.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Group', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
