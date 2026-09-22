// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_summary_custom_field_total.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptSummaryCustomFieldTotal extends ReceiptSummaryCustomFieldTotal {
  @override
  final int customFieldId;
  @override
  final String name;
  @override
  final String total;

  factory _$ReceiptSummaryCustomFieldTotal(
          [void Function(ReceiptSummaryCustomFieldTotalBuilder)? updates]) =>
      (ReceiptSummaryCustomFieldTotalBuilder()..update(updates))._build();

  _$ReceiptSummaryCustomFieldTotal._(
      {required this.customFieldId, required this.name, required this.total})
      : super._();
  @override
  ReceiptSummaryCustomFieldTotal rebuild(
          void Function(ReceiptSummaryCustomFieldTotalBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptSummaryCustomFieldTotalBuilder toBuilder() =>
      ReceiptSummaryCustomFieldTotalBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptSummaryCustomFieldTotal &&
        customFieldId == other.customFieldId &&
        name == other.name &&
        total == other.total;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, customFieldId.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, total.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptSummaryCustomFieldTotal')
          ..add('customFieldId', customFieldId)
          ..add('name', name)
          ..add('total', total))
        .toString();
  }
}

class ReceiptSummaryCustomFieldTotalBuilder
    implements
        Builder<ReceiptSummaryCustomFieldTotal,
            ReceiptSummaryCustomFieldTotalBuilder> {
  _$ReceiptSummaryCustomFieldTotal? _$v;

  int? _customFieldId;
  int? get customFieldId => _$this._customFieldId;
  set customFieldId(int? customFieldId) =>
      _$this._customFieldId = customFieldId;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _total;
  String? get total => _$this._total;
  set total(String? total) => _$this._total = total;

  ReceiptSummaryCustomFieldTotalBuilder() {
    ReceiptSummaryCustomFieldTotal._defaults(this);
  }

  ReceiptSummaryCustomFieldTotalBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _customFieldId = $v.customFieldId;
      _name = $v.name;
      _total = $v.total;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptSummaryCustomFieldTotal other) {
    _$v = other as _$ReceiptSummaryCustomFieldTotal;
  }

  @override
  void update(void Function(ReceiptSummaryCustomFieldTotalBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptSummaryCustomFieldTotal build() => _build();

  _$ReceiptSummaryCustomFieldTotal _build() {
    final _$result = _$v ??
        _$ReceiptSummaryCustomFieldTotal._(
          customFieldId: BuiltValueNullFieldError.checkNotNull(customFieldId,
              r'ReceiptSummaryCustomFieldTotal', 'customFieldId'),
          name: BuiltValueNullFieldError.checkNotNull(
              name, r'ReceiptSummaryCustomFieldTotal', 'name'),
          total: BuiltValueNullFieldError.checkNotNull(
              total, r'ReceiptSummaryCustomFieldTotal', 'total'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
