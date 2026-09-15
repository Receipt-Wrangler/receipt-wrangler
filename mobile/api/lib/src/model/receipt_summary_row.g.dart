// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_summary_row.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptSummaryRow extends ReceiptSummaryRow {
  @override
  final ReceiptStatus status;
  @override
  final int receiptCount;
  @override
  final String total;
  @override
  final BuiltList<ReceiptSummaryCustomFieldTotal> customFieldTotals;

  factory _$ReceiptSummaryRow(
          [void Function(ReceiptSummaryRowBuilder)? updates]) =>
      (ReceiptSummaryRowBuilder()..update(updates))._build();

  _$ReceiptSummaryRow._(
      {required this.status,
      required this.receiptCount,
      required this.total,
      required this.customFieldTotals})
      : super._();
  @override
  ReceiptSummaryRow rebuild(void Function(ReceiptSummaryRowBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptSummaryRowBuilder toBuilder() =>
      ReceiptSummaryRowBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptSummaryRow &&
        status == other.status &&
        receiptCount == other.receiptCount &&
        total == other.total &&
        customFieldTotals == other.customFieldTotals;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, receiptCount.hashCode);
    _$hash = $jc(_$hash, total.hashCode);
    _$hash = $jc(_$hash, customFieldTotals.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptSummaryRow')
          ..add('status', status)
          ..add('receiptCount', receiptCount)
          ..add('total', total)
          ..add('customFieldTotals', customFieldTotals))
        .toString();
  }
}

class ReceiptSummaryRowBuilder
    implements Builder<ReceiptSummaryRow, ReceiptSummaryRowBuilder> {
  _$ReceiptSummaryRow? _$v;

  ReceiptStatus? _status;
  ReceiptStatus? get status => _$this._status;
  set status(ReceiptStatus? status) => _$this._status = status;

  int? _receiptCount;
  int? get receiptCount => _$this._receiptCount;
  set receiptCount(int? receiptCount) => _$this._receiptCount = receiptCount;

  String? _total;
  String? get total => _$this._total;
  set total(String? total) => _$this._total = total;

  ListBuilder<ReceiptSummaryCustomFieldTotal>? _customFieldTotals;
  ListBuilder<ReceiptSummaryCustomFieldTotal> get customFieldTotals =>
      _$this._customFieldTotals ??=
          ListBuilder<ReceiptSummaryCustomFieldTotal>();
  set customFieldTotals(
          ListBuilder<ReceiptSummaryCustomFieldTotal>? customFieldTotals) =>
      _$this._customFieldTotals = customFieldTotals;

  ReceiptSummaryRowBuilder() {
    ReceiptSummaryRow._defaults(this);
  }

  ReceiptSummaryRowBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _status = $v.status;
      _receiptCount = $v.receiptCount;
      _total = $v.total;
      _customFieldTotals = $v.customFieldTotals.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptSummaryRow other) {
    _$v = other as _$ReceiptSummaryRow;
  }

  @override
  void update(void Function(ReceiptSummaryRowBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptSummaryRow build() => _build();

  _$ReceiptSummaryRow _build() {
    _$ReceiptSummaryRow _$result;
    try {
      _$result = _$v ??
          _$ReceiptSummaryRow._(
            status: BuiltValueNullFieldError.checkNotNull(
                status, r'ReceiptSummaryRow', 'status'),
            receiptCount: BuiltValueNullFieldError.checkNotNull(
                receiptCount, r'ReceiptSummaryRow', 'receiptCount'),
            total: BuiltValueNullFieldError.checkNotNull(
                total, r'ReceiptSummaryRow', 'total'),
            customFieldTotals: customFieldTotals.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'customFieldTotals';
        customFieldTotals.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptSummaryRow', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
