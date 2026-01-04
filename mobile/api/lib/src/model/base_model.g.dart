// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'base_model.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

abstract class BaseModelBuilder {
  void replace(BaseModel other);
  void update(void Function(BaseModelBuilder) updates);
  int? get id;
  set id(int? id);

  String? get createdAt;
  set createdAt(String? createdAt);

  int? get createdBy;
  set createdBy(int? createdBy);

  String? get createdByString;
  set createdByString(String? createdByString);

  String? get updatedAt;
  set updatedAt(String? updatedAt);
}

class _$$BaseModel extends $BaseModel {
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

  factory _$$BaseModel([void Function($BaseModelBuilder)? updates]) =>
      (new $BaseModelBuilder()..update(updates))._build();

  _$$BaseModel._(
      {required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'$BaseModel', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'$BaseModel', 'createdAt');
  }

  @override
  $BaseModel rebuild(void Function($BaseModelBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  $BaseModelBuilder toBuilder() => new $BaseModelBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is $BaseModel &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
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
    return (newBuiltValueToStringHelper(r'$BaseModel')
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class $BaseModelBuilder
    implements Builder<$BaseModel, $BaseModelBuilder>, BaseModelBuilder {
  _$$BaseModel? _$v;

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

  $BaseModelBuilder() {
    $BaseModel._defaults(this);
  }

  $BaseModelBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
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
  void replace(covariant $BaseModel other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$$BaseModel;
  }

  @override
  void update(void Function($BaseModelBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  $BaseModel build() => _build();

  _$$BaseModel _build() {
    final _$result = _$v ??
        new _$$BaseModel._(
            id: BuiltValueNullFieldError.checkNotNull(id, r'$BaseModel', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'$BaseModel', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
