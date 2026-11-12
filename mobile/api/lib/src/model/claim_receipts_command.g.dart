// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_receipts_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimReceiptsCommand extends ClaimReceiptsCommand {
  @override
  final BuiltList<int> receiptIds;

  factory _$ClaimReceiptsCommand(
          [void Function(ClaimReceiptsCommandBuilder)? updates]) =>
      (ClaimReceiptsCommandBuilder()..update(updates))._build();

  _$ClaimReceiptsCommand._({required this.receiptIds}) : super._();
  @override
  ClaimReceiptsCommand rebuild(
          void Function(ClaimReceiptsCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimReceiptsCommandBuilder toBuilder() =>
      ClaimReceiptsCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimReceiptsCommand && receiptIds == other.receiptIds;
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
    return (newBuiltValueToStringHelper(r'ClaimReceiptsCommand')
          ..add('receiptIds', receiptIds))
        .toString();
  }
}

class ClaimReceiptsCommandBuilder
    implements Builder<ClaimReceiptsCommand, ClaimReceiptsCommandBuilder> {
  _$ClaimReceiptsCommand? _$v;

  ListBuilder<int>? _receiptIds;
  ListBuilder<int> get receiptIds => _$this._receiptIds ??= ListBuilder<int>();
  set receiptIds(ListBuilder<int>? receiptIds) =>
      _$this._receiptIds = receiptIds;

  ClaimReceiptsCommandBuilder() {
    ClaimReceiptsCommand._defaults(this);
  }

  ClaimReceiptsCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _receiptIds = $v.receiptIds.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimReceiptsCommand other) {
    _$v = other as _$ClaimReceiptsCommand;
  }

  @override
  void update(void Function(ClaimReceiptsCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimReceiptsCommand build() => _build();

  _$ClaimReceiptsCommand _build() {
    _$ClaimReceiptsCommand _$result;
    try {
      _$result = _$v ??
          _$ClaimReceiptsCommand._(
            receiptIds: receiptIds.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'receiptIds';
        receiptIds.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ClaimReceiptsCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
