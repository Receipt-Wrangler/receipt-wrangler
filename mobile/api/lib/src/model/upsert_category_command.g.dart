// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_category_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertCategoryCommand extends UpsertCategoryCommand {
  @override
  final int? id;
  @override
  final String name;
  @override
  final String? description;

  factory _$UpsertCategoryCommand(
          [void Function(UpsertCategoryCommandBuilder)? updates]) =>
      (new UpsertCategoryCommandBuilder()..update(updates))._build();

  _$UpsertCategoryCommand._({this.id, required this.name, this.description})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        name, r'UpsertCategoryCommand', 'name');
  }

  @override
  UpsertCategoryCommand rebuild(
          void Function(UpsertCategoryCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertCategoryCommandBuilder toBuilder() =>
      new UpsertCategoryCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertCategoryCommand &&
        id == other.id &&
        name == other.name &&
        description == other.description;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertCategoryCommand')
          ..add('id', id)
          ..add('name', name)
          ..add('description', description))
        .toString();
  }
}

class UpsertCategoryCommandBuilder
    implements Builder<UpsertCategoryCommand, UpsertCategoryCommandBuilder> {
  _$UpsertCategoryCommand? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  UpsertCategoryCommandBuilder() {
    UpsertCategoryCommand._defaults(this);
  }

  UpsertCategoryCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _name = $v.name;
      _description = $v.description;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertCategoryCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertCategoryCommand;
  }

  @override
  void update(void Function(UpsertCategoryCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertCategoryCommand build() => _build();

  _$UpsertCategoryCommand _build() {
    final _$result = _$v ??
        new _$UpsertCategoryCommand._(
            id: id,
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertCategoryCommand', 'name'),
            description: description);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
