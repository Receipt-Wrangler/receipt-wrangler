// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'pie_chart_data.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PieChartData extends PieChartData {
  @override
  final BuiltList<PieChartDataPoint> data;

  factory _$PieChartData([void Function(PieChartDataBuilder)? updates]) =>
      (new PieChartDataBuilder()..update(updates))._build();

  _$PieChartData._({required this.data}) : super._() {
    BuiltValueNullFieldError.checkNotNull(data, r'PieChartData', 'data');
  }

  @override
  PieChartData rebuild(void Function(PieChartDataBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PieChartDataBuilder toBuilder() => new PieChartDataBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PieChartData && data == other.data;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, data.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PieChartData')..add('data', data))
        .toString();
  }
}

class PieChartDataBuilder
    implements Builder<PieChartData, PieChartDataBuilder> {
  _$PieChartData? _$v;

  ListBuilder<PieChartDataPoint>? _data;
  ListBuilder<PieChartDataPoint> get data =>
      _$this._data ??= new ListBuilder<PieChartDataPoint>();
  set data(ListBuilder<PieChartDataPoint>? data) => _$this._data = data;

  PieChartDataBuilder() {
    PieChartData._defaults(this);
  }

  PieChartDataBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _data = $v.data.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PieChartData other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PieChartData;
  }

  @override
  void update(void Function(PieChartDataBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PieChartData build() => _build();

  _$PieChartData _build() {
    _$PieChartData _$result;
    try {
      _$result = _$v ?? new _$PieChartData._(data: data.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'data';
        data.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'PieChartData', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
