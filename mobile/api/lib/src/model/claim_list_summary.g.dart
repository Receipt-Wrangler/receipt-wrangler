// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_list_summary.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimListSummary extends ClaimListSummary {
  @override
  final String asOfDate;
  @override
  final int claimCount;
  @override
  final int awaitingCount;
  @override
  final String outstandingTotal;
  @override
  final int? oldestDaysWaiting;

  factory _$ClaimListSummary(
          [void Function(ClaimListSummaryBuilder)? updates]) =>
      (ClaimListSummaryBuilder()..update(updates))._build();

  _$ClaimListSummary._(
      {required this.asOfDate,
      required this.claimCount,
      required this.awaitingCount,
      required this.outstandingTotal,
      this.oldestDaysWaiting})
      : super._();
  @override
  ClaimListSummary rebuild(void Function(ClaimListSummaryBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimListSummaryBuilder toBuilder() =>
      ClaimListSummaryBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimListSummary &&
        asOfDate == other.asOfDate &&
        claimCount == other.claimCount &&
        awaitingCount == other.awaitingCount &&
        outstandingTotal == other.outstandingTotal &&
        oldestDaysWaiting == other.oldestDaysWaiting;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, asOfDate.hashCode);
    _$hash = $jc(_$hash, claimCount.hashCode);
    _$hash = $jc(_$hash, awaitingCount.hashCode);
    _$hash = $jc(_$hash, outstandingTotal.hashCode);
    _$hash = $jc(_$hash, oldestDaysWaiting.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ClaimListSummary')
          ..add('asOfDate', asOfDate)
          ..add('claimCount', claimCount)
          ..add('awaitingCount', awaitingCount)
          ..add('outstandingTotal', outstandingTotal)
          ..add('oldestDaysWaiting', oldestDaysWaiting))
        .toString();
  }
}

class ClaimListSummaryBuilder
    implements Builder<ClaimListSummary, ClaimListSummaryBuilder> {
  _$ClaimListSummary? _$v;

  String? _asOfDate;
  String? get asOfDate => _$this._asOfDate;
  set asOfDate(String? asOfDate) => _$this._asOfDate = asOfDate;

  int? _claimCount;
  int? get claimCount => _$this._claimCount;
  set claimCount(int? claimCount) => _$this._claimCount = claimCount;

  int? _awaitingCount;
  int? get awaitingCount => _$this._awaitingCount;
  set awaitingCount(int? awaitingCount) =>
      _$this._awaitingCount = awaitingCount;

  String? _outstandingTotal;
  String? get outstandingTotal => _$this._outstandingTotal;
  set outstandingTotal(String? outstandingTotal) =>
      _$this._outstandingTotal = outstandingTotal;

  int? _oldestDaysWaiting;
  int? get oldestDaysWaiting => _$this._oldestDaysWaiting;
  set oldestDaysWaiting(int? oldestDaysWaiting) =>
      _$this._oldestDaysWaiting = oldestDaysWaiting;

  ClaimListSummaryBuilder() {
    ClaimListSummary._defaults(this);
  }

  ClaimListSummaryBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _asOfDate = $v.asOfDate;
      _claimCount = $v.claimCount;
      _awaitingCount = $v.awaitingCount;
      _outstandingTotal = $v.outstandingTotal;
      _oldestDaysWaiting = $v.oldestDaysWaiting;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimListSummary other) {
    _$v = other as _$ClaimListSummary;
  }

  @override
  void update(void Function(ClaimListSummaryBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimListSummary build() => _build();

  _$ClaimListSummary _build() {
    final _$result = _$v ??
        _$ClaimListSummary._(
          asOfDate: BuiltValueNullFieldError.checkNotNull(
              asOfDate, r'ClaimListSummary', 'asOfDate'),
          claimCount: BuiltValueNullFieldError.checkNotNull(
              claimCount, r'ClaimListSummary', 'claimCount'),
          awaitingCount: BuiltValueNullFieldError.checkNotNull(
              awaitingCount, r'ClaimListSummary', 'awaitingCount'),
          outstandingTotal: BuiltValueNullFieldError.checkNotNull(
              outstandingTotal, r'ClaimListSummary', 'outstandingTotal'),
          oldestDaysWaiting: oldestDaysWaiting,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
