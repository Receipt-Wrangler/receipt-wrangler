// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'category_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CategoryView extends CategoryView {
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

  factory _$CategoryView([void Function(CategoryViewBuilder)? updates]) =>
      (new CategoryViewBuilder()..update(updates))._build();

  _$CategoryView._(
      {this.createdAt,
      this.createdBy,
      required this.id,
      required this.name,
      this.description,
      this.updatedAt,
      required this.numberOfReceipts})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'CategoryView', 'id');
    BuiltValueNullFieldError.checkNotNull(name, r'CategoryView', 'name');
    BuiltValueNullFieldError.checkNotNull(
        numberOfReceipts, r'CategoryView', 'numberOfReceipts');
  }

  @override
  CategoryView rebuild(void Function(CategoryViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CategoryViewBuilder toBuilder() => new CategoryViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CategoryView &&
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
    return (newBuiltValueToStringHelper(r'CategoryView')
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

class CategoryViewBuilder
    implements Builder<CategoryView, CategoryViewBuilder> {
  _$CategoryView? _$v;

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

  CategoryViewBuilder() {
    CategoryView._defaults(this);
  }

  CategoryViewBuilder get _$this {
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
  void replace(CategoryView other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$CategoryView;
  }

  @override
  void update(void Function(CategoryViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CategoryView build() => _build();

  _$CategoryView _build() {
    final _$result = _$v ??
        new _$CategoryView._(
            createdAt: createdAt,
            createdBy: createdBy,
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'CategoryView', 'id'),
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'CategoryView', 'name'),
            description: description,
            updatedAt: updatedAt,
            numberOfReceipts: BuiltValueNullFieldError.checkNotNull(
                numberOfReceipts, r'CategoryView', 'numberOfReceipts'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
