// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_api_key_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertApiKeyCommand extends UpsertApiKeyCommand {
  @override
  final String name;
  @override
  final String? description;
  @override
  final ApiKeyScope scope;

  factory _$UpsertApiKeyCommand(
          [void Function(UpsertApiKeyCommandBuilder)? updates]) =>
      (new UpsertApiKeyCommandBuilder()..update(updates))._build();

  _$UpsertApiKeyCommand._(
      {required this.name, this.description, required this.scope})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(name, r'UpsertApiKeyCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        scope, r'UpsertApiKeyCommand', 'scope');
  }

  @override
  UpsertApiKeyCommand rebuild(
          void Function(UpsertApiKeyCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertApiKeyCommandBuilder toBuilder() =>
      new UpsertApiKeyCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertApiKeyCommand &&
        name == other.name &&
        description == other.description &&
        scope == other.scope;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, scope.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertApiKeyCommand')
          ..add('name', name)
          ..add('description', description)
          ..add('scope', scope))
        .toString();
  }
}

class UpsertApiKeyCommandBuilder
    implements Builder<UpsertApiKeyCommand, UpsertApiKeyCommandBuilder> {
  _$UpsertApiKeyCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  ApiKeyScope? _scope;
  ApiKeyScope? get scope => _$this._scope;
  set scope(ApiKeyScope? scope) => _$this._scope = scope;

  UpsertApiKeyCommandBuilder() {
    UpsertApiKeyCommand._defaults(this);
  }

  UpsertApiKeyCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _description = $v.description;
      _scope = $v.scope;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertApiKeyCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertApiKeyCommand;
  }

  @override
  void update(void Function(UpsertApiKeyCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertApiKeyCommand build() => _build();

  _$UpsertApiKeyCommand _build() {
    final _$result = _$v ??
        new _$UpsertApiKeyCommand._(
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertApiKeyCommand', 'name'),
            description: description,
            scope: BuiltValueNullFieldError.checkNotNull(
                scope, r'UpsertApiKeyCommand', 'scope'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
