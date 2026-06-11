// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_role_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertRoleCommand extends UpsertRoleCommand {
  @override
  final String name;
  @override
  final String? description;
  @override
  final PermissionScope scope;
  @override
  final BuiltList<Permission> permissions;

  factory _$UpsertRoleCommand(
          [void Function(UpsertRoleCommandBuilder)? updates]) =>
      (UpsertRoleCommandBuilder()..update(updates))._build();

  _$UpsertRoleCommand._(
      {required this.name,
      this.description,
      required this.scope,
      required this.permissions})
      : super._();
  @override
  UpsertRoleCommand rebuild(void Function(UpsertRoleCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertRoleCommandBuilder toBuilder() =>
      UpsertRoleCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertRoleCommand &&
        name == other.name &&
        description == other.description &&
        scope == other.scope &&
        permissions == other.permissions;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, scope.hashCode);
    _$hash = $jc(_$hash, permissions.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertRoleCommand')
          ..add('name', name)
          ..add('description', description)
          ..add('scope', scope)
          ..add('permissions', permissions))
        .toString();
  }
}

class UpsertRoleCommandBuilder
    implements Builder<UpsertRoleCommand, UpsertRoleCommandBuilder> {
  _$UpsertRoleCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  PermissionScope? _scope;
  PermissionScope? get scope => _$this._scope;
  set scope(PermissionScope? scope) => _$this._scope = scope;

  ListBuilder<Permission>? _permissions;
  ListBuilder<Permission> get permissions =>
      _$this._permissions ??= ListBuilder<Permission>();
  set permissions(ListBuilder<Permission>? permissions) =>
      _$this._permissions = permissions;

  UpsertRoleCommandBuilder() {
    UpsertRoleCommand._defaults(this);
  }

  UpsertRoleCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _description = $v.description;
      _scope = $v.scope;
      _permissions = $v.permissions.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertRoleCommand other) {
    _$v = other as _$UpsertRoleCommand;
  }

  @override
  void update(void Function(UpsertRoleCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertRoleCommand build() => _build();

  _$UpsertRoleCommand _build() {
    _$UpsertRoleCommand _$result;
    try {
      _$result = _$v ??
          _$UpsertRoleCommand._(
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertRoleCommand', 'name'),
            description: description,
            scope: BuiltValueNullFieldError.checkNotNull(
                scope, r'UpsertRoleCommand', 'scope'),
            permissions: permissions.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'permissions';
        permissions.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'UpsertRoleCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
