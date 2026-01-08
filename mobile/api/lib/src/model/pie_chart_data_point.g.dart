// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'pie_chart_data_point.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PieChartDataPoint extends PieChartDataPoint {
  @override
  final String label;
  @override
  final double value;

  factory _$PieChartDataPoint(
          [void Function(PieChartDataPointBuilder)? updates]) =>
      (new PieChartDataPointBuilder()..update(updates))._build();

  _$PieChartDataPoint._({required this.label, required this.value}) : super._() {
    BuiltValueNullFieldError.checkNotNull(label, r'PieChartDataPoint', 'label');
    BuiltValueNullFieldError.checkNotNull(value, r'PieChartDataPoint', 'value');
  }

  @override
  PieChartDataPoint rebuild(void Function(PieChartDataPointBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PieChartDataPointBuilder toBuilder() =>
      new PieChartDataPointBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PieChartDataPoint &&
        label == other.label &&
        value == other.value;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, label.hashCode);
    _$hash = $jc(_$hash, value.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PieChartDataPoint')
          ..add('label', label)
          ..add('value', value))
        .toString();
  }
}

class PieChartDataPointBuilder
    implements Builder<PieChartDataPoint, PieChartDataPointBuilder> {
  _$PieChartDataPoint? _$v;

  String? _label;
  String? get label => _$this._label;
  set label(String? label) => _$this._label = label;

  double? _value;
  double? get value => _$this._value;
  set value(double? value) => _$this._value = value;

  PieChartDataPointBuilder() {
    PieChartDataPoint._defaults(this);
  }

  PieChartDataPointBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _label = $v.label;
      _value = $v.value;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PieChartDataPoint other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PieChartDataPoint;
  }

  @override
  void update(void Function(PieChartDataPointBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PieChartDataPoint build() => _build();

  _$PieChartDataPoint _build() {
    final _$result = _$v ??
        new _$PieChartDataPoint._(
            label: BuiltValueNullFieldError.checkNotNull(
                label, r'PieChartDataPoint', 'label'),
            value: BuiltValueNullFieldError.checkNotNull(
                value, r'PieChartDataPoint', 'value'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
