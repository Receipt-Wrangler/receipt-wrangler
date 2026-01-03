// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_tag_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertTagCommand extends UpsertTagCommand {
  @override
  final int? id;
  @override
  final String name;
  @override
  final String? description;

  factory _$UpsertTagCommand(
          [void Function(UpsertTagCommandBuilder)? updates]) =>
      (new UpsertTagCommandBuilder()..update(updates))._build();

  _$UpsertTagCommand._({this.id, required this.name, this.description})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(name, r'UpsertTagCommand', 'name');
  }

  @override
  UpsertTagCommand rebuild(void Function(UpsertTagCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertTagCommandBuilder toBuilder() =>
      new UpsertTagCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertTagCommand &&
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
    return (newBuiltValueToStringHelper(r'UpsertTagCommand')
          ..add('id', id)
          ..add('name', name)
          ..add('description', description))
        .toString();
  }
}

class UpsertTagCommandBuilder
    implements Builder<UpsertTagCommand, UpsertTagCommandBuilder> {
  _$UpsertTagCommand? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  UpsertTagCommandBuilder() {
    UpsertTagCommand._defaults(this);
  }

  UpsertTagCommandBuilder get _$this {
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
  void replace(UpsertTagCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertTagCommand;
  }

  @override
  void update(void Function(UpsertTagCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertTagCommand build() => _build();

  _$UpsertTagCommand _build() {
    final _$result = _$v ??
        new _$UpsertTagCommand._(
            id: id,
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertTagCommand', 'name'),
            description: description);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
