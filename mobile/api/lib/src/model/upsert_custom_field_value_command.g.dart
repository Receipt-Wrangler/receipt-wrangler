// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_custom_field_value_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertCustomFieldValueCommand extends UpsertCustomFieldValueCommand {
  @override
  final int receiptId;
  @override
  final int customFieldId;
  @override
  final String? stringValue;
  @override
  final String? dateValue;
  @override
  final int? selectValue;
  @override
  final String? currencyValue;
  @override
  final bool? booleanValue;

  factory _$UpsertCustomFieldValueCommand(
          [void Function(UpsertCustomFieldValueCommandBuilder)? updates]) =>
      (new UpsertCustomFieldValueCommandBuilder()..update(updates))._build();

  _$UpsertCustomFieldValueCommand._(
      {required this.receiptId,
      required this.customFieldId,
      this.stringValue,
      this.dateValue,
      this.selectValue,
      this.currencyValue,
      this.booleanValue})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        receiptId, r'UpsertCustomFieldValueCommand', 'receiptId');
    BuiltValueNullFieldError.checkNotNull(
        customFieldId, r'UpsertCustomFieldValueCommand', 'customFieldId');
  }

  @override
  UpsertCustomFieldValueCommand rebuild(
          void Function(UpsertCustomFieldValueCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertCustomFieldValueCommandBuilder toBuilder() =>
      new UpsertCustomFieldValueCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertCustomFieldValueCommand &&
        receiptId == other.receiptId &&
        customFieldId == other.customFieldId &&
        stringValue == other.stringValue &&
        dateValue == other.dateValue &&
        selectValue == other.selectValue &&
        currencyValue == other.currencyValue &&
        booleanValue == other.booleanValue;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, customFieldId.hashCode);
    _$hash = $jc(_$hash, stringValue.hashCode);
    _$hash = $jc(_$hash, dateValue.hashCode);
    _$hash = $jc(_$hash, selectValue.hashCode);
    _$hash = $jc(_$hash, currencyValue.hashCode);
    _$hash = $jc(_$hash, booleanValue.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertCustomFieldValueCommand')
          ..add('receiptId', receiptId)
          ..add('customFieldId', customFieldId)
          ..add('stringValue', stringValue)
          ..add('dateValue', dateValue)
          ..add('selectValue', selectValue)
          ..add('currencyValue', currencyValue)
          ..add('booleanValue', booleanValue))
        .toString();
  }
}

class UpsertCustomFieldValueCommandBuilder
    implements
        Builder<UpsertCustomFieldValueCommand,
            UpsertCustomFieldValueCommandBuilder> {
  _$UpsertCustomFieldValueCommand? _$v;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  int? _customFieldId;
  int? get customFieldId => _$this._customFieldId;
  set customFieldId(int? customFieldId) =>
      _$this._customFieldId = customFieldId;

  String? _stringValue;
  String? get stringValue => _$this._stringValue;
  set stringValue(String? stringValue) => _$this._stringValue = stringValue;

  String? _dateValue;
  String? get dateValue => _$this._dateValue;
  set dateValue(String? dateValue) => _$this._dateValue = dateValue;

  int? _selectValue;
  int? get selectValue => _$this._selectValue;
  set selectValue(int? selectValue) => _$this._selectValue = selectValue;

  String? _currencyValue;
  String? get currencyValue => _$this._currencyValue;
  set currencyValue(String? currencyValue) =>
      _$this._currencyValue = currencyValue;

  bool? _booleanValue;
  bool? get booleanValue => _$this._booleanValue;
  set booleanValue(bool? booleanValue) => _$this._booleanValue = booleanValue;

  UpsertCustomFieldValueCommandBuilder() {
    UpsertCustomFieldValueCommand._defaults(this);
  }

  UpsertCustomFieldValueCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _receiptId = $v.receiptId;
      _customFieldId = $v.customFieldId;
      _stringValue = $v.stringValue;
      _dateValue = $v.dateValue;
      _selectValue = $v.selectValue;
      _currencyValue = $v.currencyValue;
      _booleanValue = $v.booleanValue;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertCustomFieldValueCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertCustomFieldValueCommand;
  }

  @override
  void update(void Function(UpsertCustomFieldValueCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertCustomFieldValueCommand build() => _build();

  _$UpsertCustomFieldValueCommand _build() {
    final _$result = _$v ??
        new _$UpsertCustomFieldValueCommand._(
            receiptId: BuiltValueNullFieldError.checkNotNull(
                receiptId, r'UpsertCustomFieldValueCommand', 'receiptId'),
            customFieldId: BuiltValueNullFieldError.checkNotNull(customFieldId,
                r'UpsertCustomFieldValueCommand', 'customFieldId'),
            stringValue: stringValue,
            dateValue: dateValue,
            selectValue: selectValue,
            currencyValue: currencyValue,
            booleanValue: booleanValue);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
