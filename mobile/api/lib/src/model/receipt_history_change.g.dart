// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_history_change.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptHistoryChange extends ReceiptHistoryChange {
  @override
  final String field;
  @override
  final String kind;
  @override
  final String? label;
  @override
  final int? userId;
  @override
  final ReceiptHistoryValue? from;
  @override
  final ReceiptHistoryValue? to;

  factory _$ReceiptHistoryChange(
          [void Function(ReceiptHistoryChangeBuilder)? updates]) =>
      (ReceiptHistoryChangeBuilder()..update(updates))._build();

  _$ReceiptHistoryChange._(
      {required this.field,
      required this.kind,
      this.label,
      this.userId,
      this.from,
      this.to})
      : super._();
  @override
  ReceiptHistoryChange rebuild(
          void Function(ReceiptHistoryChangeBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptHistoryChangeBuilder toBuilder() =>
      ReceiptHistoryChangeBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptHistoryChange &&
        field == other.field &&
        kind == other.kind &&
        label == other.label &&
        userId == other.userId &&
        from == other.from &&
        to == other.to;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, field.hashCode);
    _$hash = $jc(_$hash, kind.hashCode);
    _$hash = $jc(_$hash, label.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jc(_$hash, from.hashCode);
    _$hash = $jc(_$hash, to.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptHistoryChange')
          ..add('field', field)
          ..add('kind', kind)
          ..add('label', label)
          ..add('userId', userId)
          ..add('from', from)
          ..add('to', to))
        .toString();
  }
}

class ReceiptHistoryChangeBuilder
    implements Builder<ReceiptHistoryChange, ReceiptHistoryChangeBuilder> {
  _$ReceiptHistoryChange? _$v;

  String? _field;
  String? get field => _$this._field;
  set field(String? field) => _$this._field = field;

  String? _kind;
  String? get kind => _$this._kind;
  set kind(String? kind) => _$this._kind = kind;

  String? _label;
  String? get label => _$this._label;
  set label(String? label) => _$this._label = label;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  ReceiptHistoryValueBuilder? _from;
  ReceiptHistoryValueBuilder get from =>
      _$this._from ??= ReceiptHistoryValueBuilder();
  set from(ReceiptHistoryValueBuilder? from) => _$this._from = from;

  ReceiptHistoryValueBuilder? _to;
  ReceiptHistoryValueBuilder get to =>
      _$this._to ??= ReceiptHistoryValueBuilder();
  set to(ReceiptHistoryValueBuilder? to) => _$this._to = to;

  ReceiptHistoryChangeBuilder() {
    ReceiptHistoryChange._defaults(this);
  }

  ReceiptHistoryChangeBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _field = $v.field;
      _kind = $v.kind;
      _label = $v.label;
      _userId = $v.userId;
      _from = $v.from?.toBuilder();
      _to = $v.to?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptHistoryChange other) {
    _$v = other as _$ReceiptHistoryChange;
  }

  @override
  void update(void Function(ReceiptHistoryChangeBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptHistoryChange build() => _build();

  _$ReceiptHistoryChange _build() {
    _$ReceiptHistoryChange _$result;
    try {
      _$result = _$v ??
          _$ReceiptHistoryChange._(
            field: BuiltValueNullFieldError.checkNotNull(
                field, r'ReceiptHistoryChange', 'field'),
            kind: BuiltValueNullFieldError.checkNotNull(
                kind, r'ReceiptHistoryChange', 'kind'),
            label: label,
            userId: userId,
            from: _from?.build(),
            to: _to?.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'from';
        _from?.build();
        _$failedField = 'to';
        _to?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptHistoryChange', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
