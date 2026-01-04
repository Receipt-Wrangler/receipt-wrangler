// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'tag_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$TagView extends TagView {
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final int id;
  @override
  final String name;
  @override
  final String? description;
  @override
  final String? updatedAt;
  @override
  final int numberOfReceipts;

  factory _$TagView([void Function(TagViewBuilder)? updates]) =>
      (new TagViewBuilder()..update(updates))._build();

  _$TagView._(
      {this.createdAt,
      this.createdBy,
      required this.id,
      required this.name,
      this.description,
      this.updatedAt,
      required this.numberOfReceipts})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'TagView', 'id');
    BuiltValueNullFieldError.checkNotNull(name, r'TagView', 'name');
    BuiltValueNullFieldError.checkNotNull(
        numberOfReceipts, r'TagView', 'numberOfReceipts');
  }

  @override
  TagView rebuild(void Function(TagViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  TagViewBuilder toBuilder() => new TagViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is TagView &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        id == other.id &&
        name == other.name &&
        description == other.description &&
        updatedAt == other.updatedAt &&
        numberOfReceipts == other.numberOfReceipts;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, numberOfReceipts.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'TagView')
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('id', id)
          ..add('name', name)
          ..add('description', description)
          ..add('updatedAt', updatedAt)
          ..add('numberOfReceipts', numberOfReceipts))
        .toString();
  }
}

class TagViewBuilder implements Builder<TagView, TagViewBuilder> {
  _$TagView? _$v;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _numberOfReceipts;
  int? get numberOfReceipts => _$this._numberOfReceipts;
  set numberOfReceipts(int? numberOfReceipts) =>
      _$this._numberOfReceipts = numberOfReceipts;

  TagViewBuilder() {
    TagView._defaults(this);
  }

  TagViewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _id = $v.id;
      _name = $v.name;
      _description = $v.description;
      _updatedAt = $v.updatedAt;
      _numberOfReceipts = $v.numberOfReceipts;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(TagView other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$TagView;
  }

  @override
  void update(void Function(TagViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  TagView build() => _build();

  _$TagView _build() {
    final _$result = _$v ??
        new _$TagView._(
            createdAt: createdAt,
            createdBy: createdBy,
            id: BuiltValueNullFieldError.checkNotNull(id, r'TagView', 'id'),
            name:
                BuiltValueNullFieldError.checkNotNull(name, r'TagView', 'name'),
            description: description,
            updatedAt: updatedAt,
            numberOfReceipts: BuiltValueNullFieldError.checkNotNull(
                numberOfReceipts, r'TagView', 'numberOfReceipts'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
