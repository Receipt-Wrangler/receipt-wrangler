// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_receipts_result.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimReceiptsResult extends ClaimReceiptsResult {
  @override
  final BuiltList<int> receiptIds;
  @override
  final BuiltList<ClaimReceiptsFailure> failures;

  factory _$ClaimReceiptsResult(
          [void Function(ClaimReceiptsResultBuilder)? updates]) =>
      (ClaimReceiptsResultBuilder()..update(updates))._build();

  _$ClaimReceiptsResult._({required this.receiptIds, required this.failures})
      : super._();
  @override
  ClaimReceiptsResult rebuild(
          void Function(ClaimReceiptsResultBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimReceiptsResultBuilder toBuilder() =>
      ClaimReceiptsResultBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimReceiptsResult &&
        receiptIds == other.receiptIds &&
        failures == other.failures;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, receiptIds.hashCode);
    _$hash = $jc(_$hash, failures.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ClaimReceiptsResult')
          ..add('receiptIds', receiptIds)
          ..add('failures', failures))
        .toString();
  }
}

class ClaimReceiptsResultBuilder
    implements Builder<ClaimReceiptsResult, ClaimReceiptsResultBuilder> {
  _$ClaimReceiptsResult? _$v;

  ListBuilder<int>? _receiptIds;
  ListBuilder<int> get receiptIds => _$this._receiptIds ??= ListBuilder<int>();
  set receiptIds(ListBuilder<int>? receiptIds) =>
      _$this._receiptIds = receiptIds;

  ListBuilder<ClaimReceiptsFailure>? _failures;
  ListBuilder<ClaimReceiptsFailure> get failures =>
      _$this._failures ??= ListBuilder<ClaimReceiptsFailure>();
  set failures(ListBuilder<ClaimReceiptsFailure>? failures) =>
      _$this._failures = failures;

  ClaimReceiptsResultBuilder() {
    ClaimReceiptsResult._defaults(this);
  }

  ClaimReceiptsResultBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _receiptIds = $v.receiptIds.toBuilder();
      _failures = $v.failures.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimReceiptsResult other) {
    _$v = other as _$ClaimReceiptsResult;
  }

  @override
  void update(void Function(ClaimReceiptsResultBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimReceiptsResult build() => _build();

  _$ClaimReceiptsResult _build() {
    _$ClaimReceiptsResult _$result;
    try {
      _$result = _$v ??
          _$ClaimReceiptsResult._(
            receiptIds: receiptIds.build(),
            failures: failures.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'receiptIds';
        receiptIds.build();
        _$failedField = 'failures';
        failures.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ClaimReceiptsResult', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
