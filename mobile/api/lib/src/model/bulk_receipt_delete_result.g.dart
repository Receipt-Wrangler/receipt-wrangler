// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bulk_receipt_delete_result.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BulkReceiptDeleteResult extends BulkReceiptDeleteResult {
  @override
  final BuiltList<int> deletedReceiptIds;
  @override
  final BuiltList<BulkReceiptDeleteFailure> failures;

  factory _$BulkReceiptDeleteResult(
          [void Function(BulkReceiptDeleteResultBuilder)? updates]) =>
      (BulkReceiptDeleteResultBuilder()..update(updates))._build();

  _$BulkReceiptDeleteResult._(
      {required this.deletedReceiptIds, required this.failures})
      : super._();
  @override
  BulkReceiptDeleteResult rebuild(
          void Function(BulkReceiptDeleteResultBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BulkReceiptDeleteResultBuilder toBuilder() =>
      BulkReceiptDeleteResultBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BulkReceiptDeleteResult &&
        deletedReceiptIds == other.deletedReceiptIds &&
        failures == other.failures;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, deletedReceiptIds.hashCode);
    _$hash = $jc(_$hash, failures.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BulkReceiptDeleteResult')
          ..add('deletedReceiptIds', deletedReceiptIds)
          ..add('failures', failures))
        .toString();
  }
}

class BulkReceiptDeleteResultBuilder
    implements
        Builder<BulkReceiptDeleteResult, BulkReceiptDeleteResultBuilder> {
  _$BulkReceiptDeleteResult? _$v;

  ListBuilder<int>? _deletedReceiptIds;
  ListBuilder<int> get deletedReceiptIds =>
      _$this._deletedReceiptIds ??= ListBuilder<int>();
  set deletedReceiptIds(ListBuilder<int>? deletedReceiptIds) =>
      _$this._deletedReceiptIds = deletedReceiptIds;

  ListBuilder<BulkReceiptDeleteFailure>? _failures;
  ListBuilder<BulkReceiptDeleteFailure> get failures =>
      _$this._failures ??= ListBuilder<BulkReceiptDeleteFailure>();
  set failures(ListBuilder<BulkReceiptDeleteFailure>? failures) =>
      _$this._failures = failures;

  BulkReceiptDeleteResultBuilder() {
    BulkReceiptDeleteResult._defaults(this);
  }

  BulkReceiptDeleteResultBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _deletedReceiptIds = $v.deletedReceiptIds.toBuilder();
      _failures = $v.failures.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BulkReceiptDeleteResult other) {
    _$v = other as _$BulkReceiptDeleteResult;
  }

  @override
  void update(void Function(BulkReceiptDeleteResultBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BulkReceiptDeleteResult build() => _build();

  _$BulkReceiptDeleteResult _build() {
    _$BulkReceiptDeleteResult _$result;
    try {
      _$result = _$v ??
          _$BulkReceiptDeleteResult._(
            deletedReceiptIds: deletedReceiptIds.build(),
            failures: failures.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'deletedReceiptIds';
        deletedReceiptIds.build();
        _$failedField = 'failures';
        failures.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'BulkReceiptDeleteResult', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
