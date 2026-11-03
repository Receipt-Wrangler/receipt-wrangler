// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_history_snapshots.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptHistorySnapshots extends ReceiptHistorySnapshots {
  @override
  final String? before;
  @override
  final String? after;

  factory _$ReceiptHistorySnapshots(
          [void Function(ReceiptHistorySnapshotsBuilder)? updates]) =>
      (ReceiptHistorySnapshotsBuilder()..update(updates))._build();

  _$ReceiptHistorySnapshots._({this.before, this.after}) : super._();
  @override
  ReceiptHistorySnapshots rebuild(
          void Function(ReceiptHistorySnapshotsBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptHistorySnapshotsBuilder toBuilder() =>
      ReceiptHistorySnapshotsBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptHistorySnapshots &&
        before == other.before &&
        after == other.after;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, before.hashCode);
    _$hash = $jc(_$hash, after.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptHistorySnapshots')
          ..add('before', before)
          ..add('after', after))
        .toString();
  }
}

class ReceiptHistorySnapshotsBuilder
    implements
        Builder<ReceiptHistorySnapshots, ReceiptHistorySnapshotsBuilder> {
  _$ReceiptHistorySnapshots? _$v;

  String? _before;
  String? get before => _$this._before;
  set before(String? before) => _$this._before = before;

  String? _after;
  String? get after => _$this._after;
  set after(String? after) => _$this._after = after;

  ReceiptHistorySnapshotsBuilder() {
    ReceiptHistorySnapshots._defaults(this);
  }

  ReceiptHistorySnapshotsBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _before = $v.before;
      _after = $v.after;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptHistorySnapshots other) {
    _$v = other as _$ReceiptHistorySnapshots;
  }

  @override
  void update(void Function(ReceiptHistorySnapshotsBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptHistorySnapshots build() => _build();

  _$ReceiptHistorySnapshots _build() {
    final _$result = _$v ??
        _$ReceiptHistorySnapshots._(
          before: before,
          after: after,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
