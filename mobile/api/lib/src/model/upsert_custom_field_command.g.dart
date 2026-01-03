// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_custom_field_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertCustomFieldCommand extends UpsertCustomFieldCommand {
  @override
  final String name;
  @override
  final CustomFieldType type;
  @override
  final String? description;
  @override
  final BuiltList<UpsertCustomFieldOptionCommand>? options;

  factory _$UpsertCustomFieldCommand(
          [void Function(UpsertCustomFieldCommandBuilder)? updates]) =>
      (new UpsertCustomFieldCommandBuilder()..update(updates))._build();

  _$UpsertCustomFieldCommand._(
      {required this.name, required this.type, this.description, this.options})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        name, r'UpsertCustomFieldCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        type, r'UpsertCustomFieldCommand', 'type');
  }

  @override
  UpsertCustomFieldCommand rebuild(
          void Function(UpsertCustomFieldCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertCustomFieldCommandBuilder toBuilder() =>
      new UpsertCustomFieldCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertCustomFieldCommand &&
        name == other.name &&
        type == other.type &&
        description == other.description &&
        options == other.options;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, options.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertCustomFieldCommand')
          ..add('name', name)
          ..add('type', type)
          ..add('description', description)
          ..add('options', options))
        .toString();
  }
}

class UpsertCustomFieldCommandBuilder
    implements
        Builder<UpsertCustomFieldCommand, UpsertCustomFieldCommandBuilder> {
  _$UpsertCustomFieldCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  CustomFieldType? _type;
  CustomFieldType? get type => _$this._type;
  set type(CustomFieldType? type) => _$this._type = type;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  ListBuilder<UpsertCustomFieldOptionCommand>? _options;
  ListBuilder<UpsertCustomFieldOptionCommand> get options =>
      _$this._options ??= new ListBuilder<UpsertCustomFieldOptionCommand>();
  set options(ListBuilder<UpsertCustomFieldOptionCommand>? options) =>
      _$this._options = options;

  UpsertCustomFieldCommandBuilder() {
    UpsertCustomFieldCommand._defaults(this);
  }

  UpsertCustomFieldCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _type = $v.type;
      _description = $v.description;
      _options = $v.options?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertCustomFieldCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertCustomFieldCommand;
  }

  @override
  void update(void Function(UpsertCustomFieldCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertCustomFieldCommand build() => _build();

  _$UpsertCustomFieldCommand _build() {
    _$UpsertCustomFieldCommand _$result;
    try {
      _$result = _$v ??
          new _$UpsertCustomFieldCommand._(
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'UpsertCustomFieldCommand', 'name'),
              type: BuiltValueNullFieldError.checkNotNull(
                  type, r'UpsertCustomFieldCommand', 'type'),
              description: description,
              options: _options?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'options';
        _options?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UpsertCustomFieldCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
