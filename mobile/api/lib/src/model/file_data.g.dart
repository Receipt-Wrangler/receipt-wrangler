// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'file_data.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$FileData extends FileData {
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final String? fileType;
  @override
  final int id;
  @override
  final BuiltList<int>? imageData;
  @override
  final String? name;
  @override
  final int receiptId;
  @override
  final int? size;
  @override
  final String? updatedAt;

  factory _$FileData([void Function(FileDataBuilder)? updates]) =>
      (new FileDataBuilder()..update(updates))._build();

  _$FileData._(
      {this.createdAt,
      this.createdBy,
      this.fileType,
      required this.id,
      this.imageData,
      this.name,
      required this.receiptId,
      this.size,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'FileData', 'id');
    BuiltValueNullFieldError.checkNotNull(receiptId, r'FileData', 'receiptId');
  }

  @override
  FileData rebuild(void Function(FileDataBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  FileDataBuilder toBuilder() => new FileDataBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is FileData &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        fileType == other.fileType &&
        id == other.id &&
        imageData == other.imageData &&
        name == other.name &&
        receiptId == other.receiptId &&
        size == other.size &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, fileType.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, imageData.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, size.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'FileData')
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('fileType', fileType)
          ..add('id', id)
          ..add('imageData', imageData)
          ..add('name', name)
          ..add('receiptId', receiptId)
          ..add('size', size)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class FileDataBuilder implements Builder<FileData, FileDataBuilder> {
  _$FileData? _$v;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  String? _fileType;
  String? get fileType => _$this._fileType;
  set fileType(String? fileType) => _$this._fileType = fileType;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  ListBuilder<int>? _imageData;
  ListBuilder<int> get imageData =>
      _$this._imageData ??= new ListBuilder<int>();
  set imageData(ListBuilder<int>? imageData) => _$this._imageData = imageData;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  int? _size;
  int? get size => _$this._size;
  set size(int? size) => _$this._size = size;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  FileDataBuilder() {
    FileData._defaults(this);
  }

  FileDataBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _fileType = $v.fileType;
      _id = $v.id;
      _imageData = $v.imageData?.toBuilder();
      _name = $v.name;
      _receiptId = $v.receiptId;
      _size = $v.size;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(FileData other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$FileData;
  }

  @override
  void update(void Function(FileDataBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  FileData build() => _build();

  _$FileData _build() {
    _$FileData _$result;
    try {
      _$result = _$v ??
          new _$FileData._(
              createdAt: createdAt,
              createdBy: createdBy,
              fileType: fileType,
              id: BuiltValueNullFieldError.checkNotNull(id, r'FileData', 'id'),
              imageData: _imageData?.build(),
              name: name,
              receiptId: BuiltValueNullFieldError.checkNotNull(
                  receiptId, r'FileData', 'receiptId'),
              size: size,
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'imageData';
        _imageData?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'FileData', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
