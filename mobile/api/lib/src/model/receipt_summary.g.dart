// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_summary.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptSummary extends ReceiptSummary {
  @override
  final bool enabled;
  @override
  final int configurationGroupId;
  @override
  final ReceiptSummaryRow overall;
  @override
  final BuiltList<ReceiptSummaryRow> statuses;

  factory _$ReceiptSummary([void Function(ReceiptSummaryBuilder)? updates]) =>
      (ReceiptSummaryBuilder()..update(updates))._build();

  _$ReceiptSummary._(
      {required this.enabled,
      required this.configurationGroupId,
      required this.overall,
      required this.statuses})
      : super._();
  @override
  ReceiptSummary rebuild(void Function(ReceiptSummaryBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptSummaryBuilder toBuilder() => ReceiptSummaryBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptSummary &&
        enabled == other.enabled &&
        configurationGroupId == other.configurationGroupId &&
        overall == other.overall &&
        statuses == other.statuses;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, enabled.hashCode);
    _$hash = $jc(_$hash, configurationGroupId.hashCode);
    _$hash = $jc(_$hash, overall.hashCode);
    _$hash = $jc(_$hash, statuses.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptSummary')
          ..add('enabled', enabled)
          ..add('configurationGroupId', configurationGroupId)
          ..add('overall', overall)
          ..add('statuses', statuses))
        .toString();
  }
}

class ReceiptSummaryBuilder
    implements Builder<ReceiptSummary, ReceiptSummaryBuilder> {
  _$ReceiptSummary? _$v;

  bool? _enabled;
  bool? get enabled => _$this._enabled;
  set enabled(bool? enabled) => _$this._enabled = enabled;

  int? _configurationGroupId;
  int? get configurationGroupId => _$this._configurationGroupId;
  set configurationGroupId(int? configurationGroupId) =>
      _$this._configurationGroupId = configurationGroupId;

  ReceiptSummaryRowBuilder? _overall;
  ReceiptSummaryRowBuilder get overall =>
      _$this._overall ??= ReceiptSummaryRowBuilder();
  set overall(ReceiptSummaryRowBuilder? overall) => _$this._overall = overall;

  ListBuilder<ReceiptSummaryRow>? _statuses;
  ListBuilder<ReceiptSummaryRow> get statuses =>
      _$this._statuses ??= ListBuilder<ReceiptSummaryRow>();
  set statuses(ListBuilder<ReceiptSummaryRow>? statuses) =>
      _$this._statuses = statuses;

  ReceiptSummaryBuilder() {
    ReceiptSummary._defaults(this);
  }

  ReceiptSummaryBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _enabled = $v.enabled;
      _configurationGroupId = $v.configurationGroupId;
      _overall = $v.overall.toBuilder();
      _statuses = $v.statuses.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptSummary other) {
    _$v = other as _$ReceiptSummary;
  }

  @override
  void update(void Function(ReceiptSummaryBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptSummary build() => _build();

  _$ReceiptSummary _build() {
    _$ReceiptSummary _$result;
    try {
      _$result = _$v ??
          _$ReceiptSummary._(
            enabled: BuiltValueNullFieldError.checkNotNull(
                enabled, r'ReceiptSummary', 'enabled'),
            configurationGroupId: BuiltValueNullFieldError.checkNotNull(
                configurationGroupId,
                r'ReceiptSummary',
                'configurationGroupId'),
            overall: overall.build(),
            statuses: statuses.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'overall';
        overall.build();
        _$failedField = 'statuses';
        statuses.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptSummary', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
