// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'custom_field_option.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CustomFieldOption extends CustomFieldOption {
  @override
  final int customFieldId;
  @override
  final String? value;
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

  factory _$CustomFieldOption(
          [void Function(CustomFieldOptionBuilder)? updates]) =>
      (new CustomFieldOptionBuilder()..update(updates))._build();

  _$CustomFieldOption._(
      {required this.customFieldId,
      this.value,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        customFieldId, r'CustomFieldOption', 'customFieldId');
    BuiltValueNullFieldError.checkNotNull(id, r'CustomFieldOption', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'CustomFieldOption', 'createdAt');
  }

  @override
  CustomFieldOption rebuild(void Function(CustomFieldOptionBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CustomFieldOptionBuilder toBuilder() =>
      new CustomFieldOptionBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CustomFieldOption &&
        customFieldId == other.customFieldId &&
        value == other.value &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, customFieldId.hashCode);
    _$hash = $jc(_$hash, value.hashCode);
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
    return (newBuiltValueToStringHelper(r'CustomFieldOption')
          ..add('customFieldId', customFieldId)
          ..add('value', value)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class CustomFieldOptionBuilder
    implements
        Builder<CustomFieldOption, CustomFieldOptionBuilder>,
        BaseModelBuilder {
  _$CustomFieldOption? _$v;

  int? _customFieldId;
  int? get customFieldId => _$this._customFieldId;
  set customFieldId(covariant int? customFieldId) =>
      _$this._customFieldId = customFieldId;

  String? _value;
  String? get value => _$this._value;
  set value(covariant String? value) => _$this._value = value;

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

  CustomFieldOptionBuilder() {
    CustomFieldOption._defaults(this);
  }

  CustomFieldOptionBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _customFieldId = $v.customFieldId;
      _value = $v.value;
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
  void replace(covariant CustomFieldOption other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$CustomFieldOption;
  }

  @override
  void update(void Function(CustomFieldOptionBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CustomFieldOption build() => _build();

  _$CustomFieldOption _build() {
    final _$result = _$v ??
        new _$CustomFieldOption._(
            customFieldId: BuiltValueNullFieldError.checkNotNull(
                customFieldId, r'CustomFieldOption', 'customFieldId'),
            value: value,
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'CustomFieldOption', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'CustomFieldOption', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
