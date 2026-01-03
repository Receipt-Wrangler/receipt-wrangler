// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_prompt_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertPromptCommand extends UpsertPromptCommand {
  @override
  final String name;
  @override
  final String? description;
  @override
  final String prompt;

  factory _$UpsertPromptCommand(
          [void Function(UpsertPromptCommandBuilder)? updates]) =>
      (new UpsertPromptCommandBuilder()..update(updates))._build();

  _$UpsertPromptCommand._(
      {required this.name, this.description, required this.prompt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(name, r'UpsertPromptCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        prompt, r'UpsertPromptCommand', 'prompt');
  }

  @override
  UpsertPromptCommand rebuild(
          void Function(UpsertPromptCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertPromptCommandBuilder toBuilder() =>
      new UpsertPromptCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertPromptCommand &&
        name == other.name &&
        description == other.description &&
        prompt == other.prompt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, prompt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertPromptCommand')
          ..add('name', name)
          ..add('description', description)
          ..add('prompt', prompt))
        .toString();
  }
}

class UpsertPromptCommandBuilder
    implements Builder<UpsertPromptCommand, UpsertPromptCommandBuilder> {
  _$UpsertPromptCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  String? _prompt;
  String? get prompt => _$this._prompt;
  set prompt(String? prompt) => _$this._prompt = prompt;

  UpsertPromptCommandBuilder() {
    UpsertPromptCommand._defaults(this);
  }

  UpsertPromptCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _description = $v.description;
      _prompt = $v.prompt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertPromptCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertPromptCommand;
  }

  @override
  void update(void Function(UpsertPromptCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertPromptCommand build() => _build();

  _$UpsertPromptCommand _build() {
    final _$result = _$v ??
        new _$UpsertPromptCommand._(
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertPromptCommand', 'name'),
            description: description,
            prompt: BuiltValueNullFieldError.checkNotNull(
                prompt, r'UpsertPromptCommand', 'prompt'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
