// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_history_actor.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptHistoryActor extends ReceiptHistoryActor {
  @override
  final String type;
  @override
  final int? userId;

  factory _$ReceiptHistoryActor(
          [void Function(ReceiptHistoryActorBuilder)? updates]) =>
      (ReceiptHistoryActorBuilder()..update(updates))._build();

  _$ReceiptHistoryActor._({required this.type, this.userId}) : super._();
  @override
  ReceiptHistoryActor rebuild(
          void Function(ReceiptHistoryActorBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptHistoryActorBuilder toBuilder() =>
      ReceiptHistoryActorBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptHistoryActor &&
        type == other.type &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptHistoryActor')
          ..add('type', type)
          ..add('userId', userId))
        .toString();
  }
}

class ReceiptHistoryActorBuilder
    implements Builder<ReceiptHistoryActor, ReceiptHistoryActorBuilder> {
  _$ReceiptHistoryActor? _$v;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  ReceiptHistoryActorBuilder() {
    ReceiptHistoryActor._defaults(this);
  }

  ReceiptHistoryActorBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _type = $v.type;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptHistoryActor other) {
    _$v = other as _$ReceiptHistoryActor;
  }

  @override
  void update(void Function(ReceiptHistoryActorBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptHistoryActor build() => _build();

  _$ReceiptHistoryActor _build() {
    final _$result = _$v ??
        _$ReceiptHistoryActor._(
          type: BuiltValueNullFieldError.checkNotNull(
              type, r'ReceiptHistoryActor', 'type'),
          userId: userId,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
