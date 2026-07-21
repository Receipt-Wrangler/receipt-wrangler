// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_column.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ReportColumnKindEnum _$reportColumnKindEnum_dimension =
    const ReportColumnKindEnum._('dimension');
const ReportColumnKindEnum _$reportColumnKindEnum_aggregate =
    const ReportColumnKindEnum._('aggregate');
const ReportColumnKindEnum _$reportColumnKindEnum_formula =
    const ReportColumnKindEnum._('formula');

ReportColumnKindEnum _$reportColumnKindEnumValueOf(String name) {
  switch (name) {
    case 'dimension':
      return _$reportColumnKindEnum_dimension;
    case 'aggregate':
      return _$reportColumnKindEnum_aggregate;
    case 'formula':
      return _$reportColumnKindEnum_formula;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<ReportColumnKindEnum> _$reportColumnKindEnumValues =
    BuiltSet<ReportColumnKindEnum>(const <ReportColumnKindEnum>[
  _$reportColumnKindEnum_dimension,
  _$reportColumnKindEnum_aggregate,
  _$reportColumnKindEnum_formula,
]);

const ReportColumnAggFuncEnum _$reportColumnAggFuncEnum_SUM =
    const ReportColumnAggFuncEnum._('SUM');
const ReportColumnAggFuncEnum _$reportColumnAggFuncEnum_COUNT =
    const ReportColumnAggFuncEnum._('COUNT');
const ReportColumnAggFuncEnum _$reportColumnAggFuncEnum_AVG =
    const ReportColumnAggFuncEnum._('AVG');
const ReportColumnAggFuncEnum _$reportColumnAggFuncEnum_MIN =
    const ReportColumnAggFuncEnum._('MIN');
const ReportColumnAggFuncEnum _$reportColumnAggFuncEnum_MAX =
    const ReportColumnAggFuncEnum._('MAX');

ReportColumnAggFuncEnum _$reportColumnAggFuncEnumValueOf(String name) {
  switch (name) {
    case 'SUM':
      return _$reportColumnAggFuncEnum_SUM;
    case 'COUNT':
      return _$reportColumnAggFuncEnum_COUNT;
    case 'AVG':
      return _$reportColumnAggFuncEnum_AVG;
    case 'MIN':
      return _$reportColumnAggFuncEnum_MIN;
    case 'MAX':
      return _$reportColumnAggFuncEnum_MAX;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<ReportColumnAggFuncEnum> _$reportColumnAggFuncEnumValues =
    BuiltSet<ReportColumnAggFuncEnum>(const <ReportColumnAggFuncEnum>[
  _$reportColumnAggFuncEnum_SUM,
  _$reportColumnAggFuncEnum_COUNT,
  _$reportColumnAggFuncEnum_AVG,
  _$reportColumnAggFuncEnum_MIN,
  _$reportColumnAggFuncEnum_MAX,
]);

Serializer<ReportColumnKindEnum> _$reportColumnKindEnumSerializer =
    _$ReportColumnKindEnumSerializer();
Serializer<ReportColumnAggFuncEnum> _$reportColumnAggFuncEnumSerializer =
    _$ReportColumnAggFuncEnumSerializer();

class _$ReportColumnKindEnumSerializer
    implements PrimitiveSerializer<ReportColumnKindEnum> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'dimension': 'dimension',
    'aggregate': 'aggregate',
    'formula': 'formula',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'dimension': 'dimension',
    'aggregate': 'aggregate',
    'formula': 'formula',
  };

  @override
  final Iterable<Type> types = const <Type>[ReportColumnKindEnum];
  @override
  final String wireName = 'ReportColumnKindEnum';

  @override
  Object serialize(Serializers serializers, ReportColumnKindEnum object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ReportColumnKindEnum deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ReportColumnKindEnum.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

class _$ReportColumnAggFuncEnumSerializer
    implements PrimitiveSerializer<ReportColumnAggFuncEnum> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'SUM': 'SUM',
    'COUNT': 'COUNT',
    'AVG': 'AVG',
    'MIN': 'MIN',
    'MAX': 'MAX',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'SUM': 'SUM',
    'COUNT': 'COUNT',
    'AVG': 'AVG',
    'MIN': 'MIN',
    'MAX': 'MAX',
  };

  @override
  final Iterable<Type> types = const <Type>[ReportColumnAggFuncEnum];
  @override
  final String wireName = 'ReportColumnAggFuncEnum';

  @override
  Object serialize(Serializers serializers, ReportColumnAggFuncEnum object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ReportColumnAggFuncEnum deserialize(
          Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ReportColumnAggFuncEnum.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

class _$ReportColumn extends ReportColumn {
  @override
  final ReportColumnKindEnum kind;
  @override
  final String name;
  @override
  final String? label;
  @override
  final String? field;
  @override
  final ReportColumnAggFuncEnum? aggFunc;
  @override
  final String? measure;
  @override
  final String? expr;

  factory _$ReportColumn([void Function(ReportColumnBuilder)? updates]) =>
      (ReportColumnBuilder()..update(updates))._build();

  _$ReportColumn._(
      {required this.kind,
      required this.name,
      this.label,
      this.field,
      this.aggFunc,
      this.measure,
      this.expr})
      : super._();
  @override
  ReportColumn rebuild(void Function(ReportColumnBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportColumnBuilder toBuilder() => ReportColumnBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportColumn &&
        kind == other.kind &&
        name == other.name &&
        label == other.label &&
        field == other.field &&
        aggFunc == other.aggFunc &&
        measure == other.measure &&
        expr == other.expr;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, kind.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, label.hashCode);
    _$hash = $jc(_$hash, field.hashCode);
    _$hash = $jc(_$hash, aggFunc.hashCode);
    _$hash = $jc(_$hash, measure.hashCode);
    _$hash = $jc(_$hash, expr.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportColumn')
          ..add('kind', kind)
          ..add('name', name)
          ..add('label', label)
          ..add('field', field)
          ..add('aggFunc', aggFunc)
          ..add('measure', measure)
          ..add('expr', expr))
        .toString();
  }
}

class ReportColumnBuilder
    implements Builder<ReportColumn, ReportColumnBuilder> {
  _$ReportColumn? _$v;

  ReportColumnKindEnum? _kind;
  ReportColumnKindEnum? get kind => _$this._kind;
  set kind(ReportColumnKindEnum? kind) => _$this._kind = kind;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _label;
  String? get label => _$this._label;
  set label(String? label) => _$this._label = label;

  String? _field;
  String? get field => _$this._field;
  set field(String? field) => _$this._field = field;

  ReportColumnAggFuncEnum? _aggFunc;
  ReportColumnAggFuncEnum? get aggFunc => _$this._aggFunc;
  set aggFunc(ReportColumnAggFuncEnum? aggFunc) => _$this._aggFunc = aggFunc;

  String? _measure;
  String? get measure => _$this._measure;
  set measure(String? measure) => _$this._measure = measure;

  String? _expr;
  String? get expr => _$this._expr;
  set expr(String? expr) => _$this._expr = expr;

  ReportColumnBuilder() {
    ReportColumn._defaults(this);
  }

  ReportColumnBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _kind = $v.kind;
      _name = $v.name;
      _label = $v.label;
      _field = $v.field;
      _aggFunc = $v.aggFunc;
      _measure = $v.measure;
      _expr = $v.expr;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportColumn other) {
    _$v = other as _$ReportColumn;
  }

  @override
  void update(void Function(ReportColumnBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportColumn build() => _build();

  _$ReportColumn _build() {
    final _$result = _$v ??
        _$ReportColumn._(
          kind: BuiltValueNullFieldError.checkNotNull(
              kind, r'ReportColumn', 'kind'),
          name: BuiltValueNullFieldError.checkNotNull(
              name, r'ReportColumn', 'name'),
          label: label,
          field: field,
          aggFunc: aggFunc,
          measure: measure,
          expr: expr,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
