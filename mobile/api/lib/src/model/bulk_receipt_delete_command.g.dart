// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bulk_receipt_delete_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BulkReceiptDeleteCommand extends BulkReceiptDeleteCommand {
  @override
  final BuiltList<int> receiptIds;

  factory _$BulkReceiptDeleteCommand(
          [void Function(BulkReceiptDeleteCommandBuilder)? updates]) =>
      (BulkReceiptDeleteCommandBuilder()..update(updates))._build();

  _$BulkReceiptDeleteCommand._({required this.receiptIds}) : super._();
  @override
  BulkReceiptDeleteCommand rebuild(
          void Function(BulkReceiptDeleteCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BulkReceiptDeleteCommandBuilder toBuilder() =>
      BulkReceiptDeleteCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BulkReceiptDeleteCommand && receiptIds == other.receiptIds;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, receiptIds.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BulkReceiptDeleteCommand')
          ..add('receiptIds', receiptIds))
        .toString();
  }
}

class BulkReceiptDeleteCommandBuilder
    implements
        Builder<BulkReceiptDeleteCommand, BulkReceiptDeleteCommandBuilder> {
  _$BulkReceiptDeleteCommand? _$v;

  ListBuilder<int>? _receiptIds;
  ListBuilder<int> get receiptIds => _$this._receiptIds ??= ListBuilder<int>();
  set receiptIds(ListBuilder<int>? receiptIds) =>
      _$this._receiptIds = receiptIds;

  BulkReceiptDeleteCommandBuilder() {
    BulkReceiptDeleteCommand._defaults(this);
  }

  BulkReceiptDeleteCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _receiptIds = $v.receiptIds.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BulkReceiptDeleteCommand other) {
    _$v = other as _$BulkReceiptDeleteCommand;
  }

  @override
  void update(void Function(BulkReceiptDeleteCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BulkReceiptDeleteCommand build() => _build();

  _$BulkReceiptDeleteCommand _build() {
    _$BulkReceiptDeleteCommand _$result;
    try {
      _$result = _$v ??
          _$BulkReceiptDeleteCommand._(
            receiptIds: receiptIds.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'receiptIds';
        receiptIds.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'BulkReceiptDeleteCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
