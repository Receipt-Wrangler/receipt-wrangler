// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'custom_field_value.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CustomFieldValue extends CustomFieldValue {
  @override
  final String? dateValue;
  @override
  final String? stringValue;
  @override
  final int? selectValue;
  @override
  final String? currencyValue;
  @override
  final int customFieldId;
  @override
  final bool? booleanValue;
  @override
  final int receiptId;
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

  factory _$CustomFieldValue(
          [void Function(CustomFieldValueBuilder)? updates]) =>
      (new CustomFieldValueBuilder()..update(updates))._build();

  _$CustomFieldValue._(
      {this.dateValue,
      this.stringValue,
      this.selectValue,
      this.currencyValue,
      required this.customFieldId,
      this.booleanValue,
      required this.receiptId,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        customFieldId, r'CustomFieldValue', 'customFieldId');
    BuiltValueNullFieldError.checkNotNull(
        receiptId, r'CustomFieldValue', 'receiptId');
    BuiltValueNullFieldError.checkNotNull(id, r'CustomFieldValue', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'CustomFieldValue', 'createdAt');
  }

  @override
  CustomFieldValue rebuild(void Function(CustomFieldValueBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CustomFieldValueBuilder toBuilder() =>
      new CustomFieldValueBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CustomFieldValue &&
        dateValue == other.dateValue &&
        stringValue == other.stringValue &&
        selectValue == other.selectValue &&
        currencyValue == other.currencyValue &&
        customFieldId == other.customFieldId &&
        booleanValue == other.booleanValue &&
        receiptId == other.receiptId &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, dateValue.hashCode);
    _$hash = $jc(_$hash, stringValue.hashCode);
    _$hash = $jc(_$hash, selectValue.hashCode);
    _$hash = $jc(_$hash, currencyValue.hashCode);
    _$hash = $jc(_$hash, customFieldId.hashCode);
    _$hash = $jc(_$hash, booleanValue.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
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
    return (newBuiltValueToStringHelper(r'CustomFieldValue')
          ..add('dateValue', dateValue)
          ..add('stringValue', stringValue)
          ..add('selectValue', selectValue)
          ..add('currencyValue', currencyValue)
          ..add('customFieldId', customFieldId)
          ..add('booleanValue', booleanValue)
          ..add('receiptId', receiptId)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class CustomFieldValueBuilder
    implements
        Builder<CustomFieldValue, CustomFieldValueBuilder>,
        BaseModelBuilder {
  _$CustomFieldValue? _$v;

  String? _dateValue;
  String? get dateValue => _$this._dateValue;
  set dateValue(covariant String? dateValue) => _$this._dateValue = dateValue;

  String? _stringValue;
  String? get stringValue => _$this._stringValue;
  set stringValue(covariant String? stringValue) =>
      _$this._stringValue = stringValue;

  int? _selectValue;
  int? get selectValue => _$this._selectValue;
  set selectValue(covariant int? selectValue) =>
      _$this._selectValue = selectValue;

  String? _currencyValue;
  String? get currencyValue => _$this._currencyValue;
  set currencyValue(covariant String? currencyValue) =>
      _$this._currencyValue = currencyValue;

  int? _customFieldId;
  int? get customFieldId => _$this._customFieldId;
  set customFieldId(covariant int? customFieldId) =>
      _$this._customFieldId = customFieldId;

  bool? _booleanValue;
  bool? get booleanValue => _$this._booleanValue;
  set booleanValue(covariant bool? booleanValue) =>
      _$this._booleanValue = booleanValue;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(covariant int? receiptId) => _$this._receiptId = receiptId;

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

  CustomFieldValueBuilder() {
    CustomFieldValue._defaults(this);
  }

  CustomFieldValueBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _dateValue = $v.dateValue;
      _stringValue = $v.stringValue;
      _selectValue = $v.selectValue;
      _currencyValue = $v.currencyValue;
      _customFieldId = $v.customFieldId;
      _booleanValue = $v.booleanValue;
      _receiptId = $v.receiptId;
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
  void replace(covariant CustomFieldValue other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$CustomFieldValue;
  }

  @override
  void update(void Function(CustomFieldValueBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CustomFieldValue build() => _build();

  _$CustomFieldValue _build() {
    final _$result = _$v ??
        new _$CustomFieldValue._(
            dateValue: dateValue,
            stringValue: stringValue,
            selectValue: selectValue,
            currencyValue: currencyValue,
            customFieldId: BuiltValueNullFieldError.checkNotNull(
                customFieldId, r'CustomFieldValue', 'customFieldId'),
            booleanValue: booleanValue,
            receiptId: BuiltValueNullFieldError.checkNotNull(
                receiptId, r'CustomFieldValue', 'receiptId'),
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'CustomFieldValue', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'CustomFieldValue', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
