// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bulk_receipt_delete_failure.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BulkReceiptDeleteFailure extends BulkReceiptDeleteFailure {
  @override
  final int receiptId;
  @override
  final String message;

  factory _$BulkReceiptDeleteFailure(
          [void Function(BulkReceiptDeleteFailureBuilder)? updates]) =>
      (BulkReceiptDeleteFailureBuilder()..update(updates))._build();

  _$BulkReceiptDeleteFailure._({required this.receiptId, required this.message})
      : super._();
  @override
  BulkReceiptDeleteFailure rebuild(
          void Function(BulkReceiptDeleteFailureBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BulkReceiptDeleteFailureBuilder toBuilder() =>
      BulkReceiptDeleteFailureBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BulkReceiptDeleteFailure &&
        receiptId == other.receiptId &&
        message == other.message;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, message.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BulkReceiptDeleteFailure')
          ..add('receiptId', receiptId)
          ..add('message', message))
        .toString();
  }
}

class BulkReceiptDeleteFailureBuilder
    implements
        Builder<BulkReceiptDeleteFailure, BulkReceiptDeleteFailureBuilder> {
  _$BulkReceiptDeleteFailure? _$v;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  String? _message;
  String? get message => _$this._message;
  set message(String? message) => _$this._message = message;

  BulkReceiptDeleteFailureBuilder() {
    BulkReceiptDeleteFailure._defaults(this);
  }

  BulkReceiptDeleteFailureBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _receiptId = $v.receiptId;
      _message = $v.message;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BulkReceiptDeleteFailure other) {
    _$v = other as _$BulkReceiptDeleteFailure;
  }

  @override
  void update(void Function(BulkReceiptDeleteFailureBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BulkReceiptDeleteFailure build() => _build();

  _$BulkReceiptDeleteFailure _build() {
    final _$result = _$v ??
        _$BulkReceiptDeleteFailure._(
          receiptId: BuiltValueNullFieldError.checkNotNull(
              receiptId, r'BulkReceiptDeleteFailure', 'receiptId'),
          message: BuiltValueNullFieldError.checkNotNull(
              message, r'BulkReceiptDeleteFailure', 'message'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
