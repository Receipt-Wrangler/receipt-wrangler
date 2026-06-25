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
  @override
  final BuiltList<int>? categoryGrants;
  @override
  final BuiltList<int>? tagGrants;
  @override
  final BuiltList<int>? paidByUserGrants;
  @override
  final bool? includeOwnPaidReceipts;

  factory _$UpsertRoleCommand(
          [void Function(UpsertRoleCommandBuilder)? updates]) =>
      (UpsertRoleCommandBuilder()..update(updates))._build();

  _$UpsertRoleCommand._(
      {required this.name,
      this.description,
      required this.scope,
      required this.permissions,
      this.categoryGrants,
      this.tagGrants,
      this.paidByUserGrants,
      this.includeOwnPaidReceipts})
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
        permissions == other.permissions &&
        categoryGrants == other.categoryGrants &&
        tagGrants == other.tagGrants &&
        paidByUserGrants == other.paidByUserGrants &&
        includeOwnPaidReceipts == other.includeOwnPaidReceipts;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, scope.hashCode);
    _$hash = $jc(_$hash, permissions.hashCode);
    _$hash = $jc(_$hash, categoryGrants.hashCode);
    _$hash = $jc(_$hash, tagGrants.hashCode);
    _$hash = $jc(_$hash, paidByUserGrants.hashCode);
    _$hash = $jc(_$hash, includeOwnPaidReceipts.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertRoleCommand')
          ..add('name', name)
          ..add('description', description)
          ..add('scope', scope)
          ..add('permissions', permissions)
          ..add('categoryGrants', categoryGrants)
          ..add('tagGrants', tagGrants)
          ..add('paidByUserGrants', paidByUserGrants)
          ..add('includeOwnPaidReceipts', includeOwnPaidReceipts))
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

  ListBuilder<int>? _categoryGrants;
  ListBuilder<int> get categoryGrants =>
      _$this._categoryGrants ??= ListBuilder<int>();
  set categoryGrants(ListBuilder<int>? categoryGrants) =>
      _$this._categoryGrants = categoryGrants;

  ListBuilder<int>? _tagGrants;
  ListBuilder<int> get tagGrants => _$this._tagGrants ??= ListBuilder<int>();
  set tagGrants(ListBuilder<int>? tagGrants) => _$this._tagGrants = tagGrants;

  ListBuilder<int>? _paidByUserGrants;
  ListBuilder<int> get paidByUserGrants =>
      _$this._paidByUserGrants ??= ListBuilder<int>();
  set paidByUserGrants(ListBuilder<int>? paidByUserGrants) =>
      _$this._paidByUserGrants = paidByUserGrants;

  bool? _includeOwnPaidReceipts;
  bool? get includeOwnPaidReceipts => _$this._includeOwnPaidReceipts;
  set includeOwnPaidReceipts(bool? includeOwnPaidReceipts) =>
      _$this._includeOwnPaidReceipts = includeOwnPaidReceipts;

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
      _categoryGrants = $v.categoryGrants?.toBuilder();
      _tagGrants = $v.tagGrants?.toBuilder();
      _paidByUserGrants = $v.paidByUserGrants?.toBuilder();
      _includeOwnPaidReceipts = $v.includeOwnPaidReceipts;
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
            categoryGrants: _categoryGrants?.build(),
            tagGrants: _tagGrants?.build(),
            paidByUserGrants: _paidByUserGrants?.build(),
            includeOwnPaidReceipts: includeOwnPaidReceipts,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'permissions';
        permissions.build();
        _$failedField = 'categoryGrants';
        _categoryGrants?.build();
        _$failedField = 'tagGrants';
        _tagGrants?.build();
        _$failedField = 'paidByUserGrants';
        _paidByUserGrants?.build();
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
