// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_receipts_failure.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimReceiptsFailure extends ClaimReceiptsFailure {
  @override
  final int receiptId;
  @override
  final String message;

  factory _$ClaimReceiptsFailure(
          [void Function(ClaimReceiptsFailureBuilder)? updates]) =>
      (ClaimReceiptsFailureBuilder()..update(updates))._build();

  _$ClaimReceiptsFailure._({required this.receiptId, required this.message})
      : super._();
  @override
  ClaimReceiptsFailure rebuild(
          void Function(ClaimReceiptsFailureBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimReceiptsFailureBuilder toBuilder() =>
      ClaimReceiptsFailureBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimReceiptsFailure &&
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
    return (newBuiltValueToStringHelper(r'ClaimReceiptsFailure')
          ..add('receiptId', receiptId)
          ..add('message', message))
        .toString();
  }
}

class ClaimReceiptsFailureBuilder
    implements Builder<ClaimReceiptsFailure, ClaimReceiptsFailureBuilder> {
  _$ClaimReceiptsFailure? _$v;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  String? _message;
  String? get message => _$this._message;
  set message(String? message) => _$this._message = message;

  ClaimReceiptsFailureBuilder() {
    ClaimReceiptsFailure._defaults(this);
  }

  ClaimReceiptsFailureBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _receiptId = $v.receiptId;
      _message = $v.message;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimReceiptsFailure other) {
    _$v = other as _$ClaimReceiptsFailure;
  }

  @override
  void update(void Function(ClaimReceiptsFailureBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimReceiptsFailure build() => _build();

  _$ClaimReceiptsFailure _build() {
    final _$result = _$v ??
        _$ClaimReceiptsFailure._(
          receiptId: BuiltValueNullFieldError.checkNotNull(
              receiptId, r'ClaimReceiptsFailure', 'receiptId'),
          message: BuiltValueNullFieldError.checkNotNull(
              message, r'ClaimReceiptsFailure', 'message'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
