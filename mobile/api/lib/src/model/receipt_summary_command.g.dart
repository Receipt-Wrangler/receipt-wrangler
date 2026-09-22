// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_summary_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptSummaryCommand extends ReceiptSummaryCommand {
  @override
  final ReceiptPagedRequestFilter? filter;
  @override
  final int? configurationGroupId;

  factory _$ReceiptSummaryCommand(
          [void Function(ReceiptSummaryCommandBuilder)? updates]) =>
      (ReceiptSummaryCommandBuilder()..update(updates))._build();

  _$ReceiptSummaryCommand._({this.filter, this.configurationGroupId})
      : super._();
  @override
  ReceiptSummaryCommand rebuild(
          void Function(ReceiptSummaryCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptSummaryCommandBuilder toBuilder() =>
      ReceiptSummaryCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptSummaryCommand &&
        filter == other.filter &&
        configurationGroupId == other.configurationGroupId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, filter.hashCode);
    _$hash = $jc(_$hash, configurationGroupId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptSummaryCommand')
          ..add('filter', filter)
          ..add('configurationGroupId', configurationGroupId))
        .toString();
  }
}

class ReceiptSummaryCommandBuilder
    implements Builder<ReceiptSummaryCommand, ReceiptSummaryCommandBuilder> {
  _$ReceiptSummaryCommand? _$v;

  ReceiptPagedRequestFilterBuilder? _filter;
  ReceiptPagedRequestFilterBuilder get filter =>
      _$this._filter ??= ReceiptPagedRequestFilterBuilder();
  set filter(ReceiptPagedRequestFilterBuilder? filter) =>
      _$this._filter = filter;

  int? _configurationGroupId;
  int? get configurationGroupId => _$this._configurationGroupId;
  set configurationGroupId(int? configurationGroupId) =>
      _$this._configurationGroupId = configurationGroupId;

  ReceiptSummaryCommandBuilder() {
    ReceiptSummaryCommand._defaults(this);
  }

  ReceiptSummaryCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _filter = $v.filter?.toBuilder();
      _configurationGroupId = $v.configurationGroupId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptSummaryCommand other) {
    _$v = other as _$ReceiptSummaryCommand;
  }

  @override
  void update(void Function(ReceiptSummaryCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptSummaryCommand build() => _build();

  _$ReceiptSummaryCommand _build() {
    _$ReceiptSummaryCommand _$result;
    try {
      _$result = _$v ??
          _$ReceiptSummaryCommand._(
            filter: _filter?.build(),
            configurationGroupId: configurationGroupId,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'filter';
        _filter?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptSummaryCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
