// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_period.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ReportPeriodPresetEnum _$reportPeriodPresetEnum_thisMonth =
    const ReportPeriodPresetEnum._('thisMonth');
const ReportPeriodPresetEnum _$reportPeriodPresetEnum_lastMonth =
    const ReportPeriodPresetEnum._('lastMonth');
const ReportPeriodPresetEnum _$reportPeriodPresetEnum_mtd =
    const ReportPeriodPresetEnum._('mtd');
const ReportPeriodPresetEnum _$reportPeriodPresetEnum_qtd =
    const ReportPeriodPresetEnum._('qtd');
const ReportPeriodPresetEnum _$reportPeriodPresetEnum_ytd =
    const ReportPeriodPresetEnum._('ytd');
const ReportPeriodPresetEnum _$reportPeriodPresetEnum_custom =
    const ReportPeriodPresetEnum._('custom');

ReportPeriodPresetEnum _$reportPeriodPresetEnumValueOf(String name) {
  switch (name) {
    case 'thisMonth':
      return _$reportPeriodPresetEnum_thisMonth;
    case 'lastMonth':
      return _$reportPeriodPresetEnum_lastMonth;
    case 'mtd':
      return _$reportPeriodPresetEnum_mtd;
    case 'qtd':
      return _$reportPeriodPresetEnum_qtd;
    case 'ytd':
      return _$reportPeriodPresetEnum_ytd;
    case 'custom':
      return _$reportPeriodPresetEnum_custom;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<ReportPeriodPresetEnum> _$reportPeriodPresetEnumValues =
    BuiltSet<ReportPeriodPresetEnum>(const <ReportPeriodPresetEnum>[
  _$reportPeriodPresetEnum_thisMonth,
  _$reportPeriodPresetEnum_lastMonth,
  _$reportPeriodPresetEnum_mtd,
  _$reportPeriodPresetEnum_qtd,
  _$reportPeriodPresetEnum_ytd,
  _$reportPeriodPresetEnum_custom,
]);

Serializer<ReportPeriodPresetEnum> _$reportPeriodPresetEnumSerializer =
    _$ReportPeriodPresetEnumSerializer();

class _$ReportPeriodPresetEnumSerializer
    implements PrimitiveSerializer<ReportPeriodPresetEnum> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'thisMonth': 'this_month',
    'lastMonth': 'last_month',
    'mtd': 'mtd',
    'qtd': 'qtd',
    'ytd': 'ytd',
    'custom': 'custom',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'this_month': 'thisMonth',
    'last_month': 'lastMonth',
    'mtd': 'mtd',
    'qtd': 'qtd',
    'ytd': 'ytd',
    'custom': 'custom',
  };

  @override
  final Iterable<Type> types = const <Type>[ReportPeriodPresetEnum];
  @override
  final String wireName = 'ReportPeriodPresetEnum';

  @override
  Object serialize(Serializers serializers, ReportPeriodPresetEnum object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ReportPeriodPresetEnum deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ReportPeriodPresetEnum.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

class _$ReportPeriod extends ReportPeriod {
  @override
  final ReportPeriodPresetEnum preset;
  @override
  final String? startDate;
  @override
  final String? endDate;

  factory _$ReportPeriod([void Function(ReportPeriodBuilder)? updates]) =>
      (ReportPeriodBuilder()..update(updates))._build();

  _$ReportPeriod._({required this.preset, this.startDate, this.endDate})
      : super._();
  @override
  ReportPeriod rebuild(void Function(ReportPeriodBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportPeriodBuilder toBuilder() => ReportPeriodBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportPeriod &&
        preset == other.preset &&
        startDate == other.startDate &&
        endDate == other.endDate;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, preset.hashCode);
    _$hash = $jc(_$hash, startDate.hashCode);
    _$hash = $jc(_$hash, endDate.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportPeriod')
          ..add('preset', preset)
          ..add('startDate', startDate)
          ..add('endDate', endDate))
        .toString();
  }
}

class ReportPeriodBuilder
    implements Builder<ReportPeriod, ReportPeriodBuilder> {
  _$ReportPeriod? _$v;

  ReportPeriodPresetEnum? _preset;
  ReportPeriodPresetEnum? get preset => _$this._preset;
  set preset(ReportPeriodPresetEnum? preset) => _$this._preset = preset;

  String? _startDate;
  String? get startDate => _$this._startDate;
  set startDate(String? startDate) => _$this._startDate = startDate;

  String? _endDate;
  String? get endDate => _$this._endDate;
  set endDate(String? endDate) => _$this._endDate = endDate;

  ReportPeriodBuilder() {
    ReportPeriod._defaults(this);
  }

  ReportPeriodBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _preset = $v.preset;
      _startDate = $v.startDate;
      _endDate = $v.endDate;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportPeriod other) {
    _$v = other as _$ReportPeriod;
  }

  @override
  void update(void Function(ReportPeriodBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportPeriod build() => _build();

  _$ReportPeriod _build() {
    final _$result = _$v ??
        _$ReportPeriod._(
          preset: BuiltValueNullFieldError.checkNotNull(
              preset, r'ReportPeriod', 'preset'),
          startDate: startDate,
          endDate: endDate,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
