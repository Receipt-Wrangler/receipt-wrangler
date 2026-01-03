// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'prompt.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Prompt extends Prompt {
  @override
  final String name;
  @override
  final String? description;
  @override
  final String prompt;
  @override
  final int id;
  @override
  final String createdAt;
  @override
  final int? createdBy;
  @override
  final String? createdByString;
  @override
  final String? updatedAt;

  factory _$Prompt([void Function(PromptBuilder)? updates]) =>
      (new PromptBuilder()..update(updates))._build();

  _$Prompt._(
      {required this.name,
      this.description,
      required this.prompt,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(name, r'Prompt', 'name');
    BuiltValueNullFieldError.checkNotNull(prompt, r'Prompt', 'prompt');
    BuiltValueNullFieldError.checkNotNull(id, r'Prompt', 'id');
    BuiltValueNullFieldError.checkNotNull(createdAt, r'Prompt', 'createdAt');
  }

  @override
  Prompt rebuild(void Function(PromptBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PromptBuilder toBuilder() => new PromptBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Prompt &&
        name == other.name &&
        description == other.description &&
        prompt == other.prompt &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, prompt.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Prompt')
          ..add('name', name)
          ..add('description', description)
          ..add('prompt', prompt)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class PromptBuilder
    implements Builder<Prompt, PromptBuilder>, BaseModelBuilder {
  _$Prompt? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(covariant String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(covariant String? description) =>
      _$this._description = description;

  String? _prompt;
  String? get prompt => _$this._prompt;
  set prompt(covariant String? prompt) => _$this._prompt = prompt;

  int? _id;
  int? get id => _$this._id;
  set id(covariant int? id) => _$this._id = id;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(covariant String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(covariant int? createdBy) => _$this._createdBy = createdBy;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(covariant String? createdByString) =>
      _$this._createdByString = createdByString;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(covariant String? updatedAt) => _$this._updatedAt = updatedAt;

  PromptBuilder() {
    Prompt._defaults(this);
  }

  PromptBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _description = $v.description;
      _prompt = $v.prompt;
      _id = $v.id;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _createdByString = $v.createdByString;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant Prompt other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Prompt;
  }

  @override
  void update(void Function(PromptBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Prompt build() => _build();

  _$Prompt _build() {
    final _$result = _$v ??
        new _$Prompt._(
            name:
                BuiltValueNullFieldError.checkNotNull(name, r'Prompt', 'name'),
            description: description,
            prompt: BuiltValueNullFieldError.checkNotNull(
                prompt, r'Prompt', 'prompt'),
            id: BuiltValueNullFieldError.checkNotNull(id, r'Prompt', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'Prompt', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
