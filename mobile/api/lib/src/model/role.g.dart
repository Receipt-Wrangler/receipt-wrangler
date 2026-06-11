// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'role.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Role extends Role {
  @override
  final int id;
  @override
  final String name;
  @override
  final String? description;
  @override
  final PermissionScope scope;
  @override
  final bool isDefault;
  @override
  final bool isSystem;
  @override
  final BuiltList<Permission> permissions;
  @override
  final int? assignedCount;

  factory _$Role([void Function(RoleBuilder)? updates]) =>
      (RoleBuilder()..update(updates))._build();

  _$Role._(
      {required this.id,
      required this.name,
      this.description,
      required this.scope,
      required this.isDefault,
      required this.isSystem,
      required this.permissions,
      this.assignedCount})
      : super._();
  @override
  Role rebuild(void Function(RoleBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  RoleBuilder toBuilder() => RoleBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Role &&
        id == other.id &&
        name == other.name &&
        description == other.description &&
        scope == other.scope &&
        isDefault == other.isDefault &&
        isSystem == other.isSystem &&
        permissions == other.permissions &&
        assignedCount == other.assignedCount;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, scope.hashCode);
    _$hash = $jc(_$hash, isDefault.hashCode);
    _$hash = $jc(_$hash, isSystem.hashCode);
    _$hash = $jc(_$hash, permissions.hashCode);
    _$hash = $jc(_$hash, assignedCount.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Role')
          ..add('id', id)
          ..add('name', name)
          ..add('description', description)
          ..add('scope', scope)
          ..add('isDefault', isDefault)
          ..add('isSystem', isSystem)
          ..add('permissions', permissions)
          ..add('assignedCount', assignedCount))
        .toString();
  }
}

class RoleBuilder implements Builder<Role, RoleBuilder> {
  _$Role? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  PermissionScope? _scope;
  PermissionScope? get scope => _$this._scope;
  set scope(PermissionScope? scope) => _$this._scope = scope;

  bool? _isDefault;
  bool? get isDefault => _$this._isDefault;
  set isDefault(bool? isDefault) => _$this._isDefault = isDefault;

  bool? _isSystem;
  bool? get isSystem => _$this._isSystem;
  set isSystem(bool? isSystem) => _$this._isSystem = isSystem;

  ListBuilder<Permission>? _permissions;
  ListBuilder<Permission> get permissions =>
      _$this._permissions ??= ListBuilder<Permission>();
  set permissions(ListBuilder<Permission>? permissions) =>
      _$this._permissions = permissions;

  int? _assignedCount;
  int? get assignedCount => _$this._assignedCount;
  set assignedCount(int? assignedCount) =>
      _$this._assignedCount = assignedCount;

  RoleBuilder() {
    Role._defaults(this);
  }

  RoleBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _name = $v.name;
      _description = $v.description;
      _scope = $v.scope;
      _isDefault = $v.isDefault;
      _isSystem = $v.isSystem;
      _permissions = $v.permissions.toBuilder();
      _assignedCount = $v.assignedCount;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Role other) {
    _$v = other as _$Role;
  }

  @override
  void update(void Function(RoleBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Role build() => _build();

  _$Role _build() {
    _$Role _$result;
    try {
      _$result = _$v ??
          _$Role._(
            id: BuiltValueNullFieldError.checkNotNull(id, r'Role', 'id'),
            name: BuiltValueNullFieldError.checkNotNull(name, r'Role', 'name'),
            description: description,
            scope:
                BuiltValueNullFieldError.checkNotNull(scope, r'Role', 'scope'),
            isDefault: BuiltValueNullFieldError.checkNotNull(
                isDefault, r'Role', 'isDefault'),
            isSystem: BuiltValueNullFieldError.checkNotNull(
                isSystem, r'Role', 'isSystem'),
            permissions: permissions.build(),
            assignedCount: assignedCount,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'permissions';
        permissions.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(r'Role', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
