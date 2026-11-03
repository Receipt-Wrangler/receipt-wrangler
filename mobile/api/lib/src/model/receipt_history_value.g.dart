// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_history_value.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptHistoryValue extends ReceiptHistoryValue {
  @override
  final String type;
  @override
  final String? value;
  @override
  final int? userId;

  factory _$ReceiptHistoryValue(
          [void Function(ReceiptHistoryValueBuilder)? updates]) =>
      (ReceiptHistoryValueBuilder()..update(updates))._build();

  _$ReceiptHistoryValue._({required this.type, this.value, this.userId})
      : super._();
  @override
  ReceiptHistoryValue rebuild(
          void Function(ReceiptHistoryValueBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptHistoryValueBuilder toBuilder() =>
      ReceiptHistoryValueBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptHistoryValue &&
        type == other.type &&
        value == other.value &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, value.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptHistoryValue')
          ..add('type', type)
          ..add('value', value)
          ..add('userId', userId))
        .toString();
  }
}

class ReceiptHistoryValueBuilder
    implements Builder<ReceiptHistoryValue, ReceiptHistoryValueBuilder> {
  _$ReceiptHistoryValue? _$v;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  String? _value;
  String? get value => _$this._value;
  set value(String? value) => _$this._value = value;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  ReceiptHistoryValueBuilder() {
    ReceiptHistoryValue._defaults(this);
  }

  ReceiptHistoryValueBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _type = $v.type;
      _value = $v.value;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptHistoryValue other) {
    _$v = other as _$ReceiptHistoryValue;
  }

  @override
  void update(void Function(ReceiptHistoryValueBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptHistoryValue build() => _build();

  _$ReceiptHistoryValue _build() {
    final _$result = _$v ??
        _$ReceiptHistoryValue._(
          type: BuiltValueNullFieldError.checkNotNull(
              type, r'ReceiptHistoryValue', 'type'),
          value: value,
          userId: userId,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
