// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'file_data_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$FileDataView extends FileDataView {
  @override
  final String encodedImage;
  @override
  final String name;
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

  factory _$FileDataView([void Function(FileDataViewBuilder)? updates]) =>
      (new FileDataViewBuilder()..update(updates))._build();

  _$FileDataView._(
      {required this.encodedImage,
      required this.name,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        encodedImage, r'FileDataView', 'encodedImage');
    BuiltValueNullFieldError.checkNotNull(name, r'FileDataView', 'name');
    BuiltValueNullFieldError.checkNotNull(id, r'FileDataView', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'FileDataView', 'createdAt');
  }

  @override
  FileDataView rebuild(void Function(FileDataViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  FileDataViewBuilder toBuilder() => new FileDataViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is FileDataView &&
        encodedImage == other.encodedImage &&
        name == other.name &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, encodedImage.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
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
    return (newBuiltValueToStringHelper(r'FileDataView')
          ..add('encodedImage', encodedImage)
          ..add('name', name)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class FileDataViewBuilder
    implements Builder<FileDataView, FileDataViewBuilder>, BaseModelBuilder {
  _$FileDataView? _$v;

  String? _encodedImage;
  String? get encodedImage => _$this._encodedImage;
  set encodedImage(covariant String? encodedImage) =>
      _$this._encodedImage = encodedImage;

  String? _name;
  String? get name => _$this._name;
  set name(covariant String? name) => _$this._name = name;

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

  FileDataViewBuilder() {
    FileDataView._defaults(this);
  }

  FileDataViewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _encodedImage = $v.encodedImage;
      _name = $v.name;
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
  void replace(covariant FileDataView other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$FileDataView;
  }

  @override
  void update(void Function(FileDataViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  FileDataView build() => _build();

  _$FileDataView _build() {
    final _$result = _$v ??
        new _$FileDataView._(
            encodedImage: BuiltValueNullFieldError.checkNotNull(
                encodedImage, r'FileDataView', 'encodedImage'),
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'FileDataView', 'name'),
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'FileDataView', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'FileDataView', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
