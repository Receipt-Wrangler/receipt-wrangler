// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'chart_grouping.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ChartGrouping _$CATEGORIES = const ChartGrouping._('CATEGORIES');
const ChartGrouping _$TAGS = const ChartGrouping._('TAGS');
const ChartGrouping _$PAIDBY = const ChartGrouping._('PAIDBY');

ChartGrouping _$valueOf(String name) {
  switch (name) {
    case 'CATEGORIES':
      return _$CATEGORIES;
    case 'TAGS':
      return _$TAGS;
    case 'PAIDBY':
      return _$PAIDBY;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<ChartGrouping> _$values =
    new BuiltSet<ChartGrouping>(const <ChartGrouping>[
  _$CATEGORIES,
  _$TAGS,
  _$PAIDBY,
]);

class _$ChartGroupingMeta {
  const _$ChartGroupingMeta();
  ChartGrouping get CATEGORIES => _$CATEGORIES;
  ChartGrouping get TAGS => _$TAGS;
  ChartGrouping get PAIDBY => _$PAIDBY;
  ChartGrouping valueOf(String name) => _$valueOf(name);
  BuiltSet<ChartGrouping> get values => _$values;
}

abstract class _$ChartGroupingMixin {
  // ignore: non_constant_identifier_names
  _$ChartGroupingMeta get ChartGrouping => const _$ChartGroupingMeta();
}

Serializer<ChartGrouping> _$chartGroupingSerializer =
    new _$ChartGroupingSerializer();

class _$ChartGroupingSerializer implements PrimitiveSerializer<ChartGrouping> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'CATEGORIES': 'CATEGORIES',
    'TAGS': 'TAGS',
    'PAIDBY': 'PAIDBY',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'CATEGORIES': 'CATEGORIES',
    'TAGS': 'TAGS',
    'PAIDBY': 'PAIDBY',
  };

  @override
  final Iterable<Type> types = const <Type>[ChartGrouping];
  @override
  final String wireName = 'ChartGrouping';

  @override
  Object serialize(Serializers serializers, ChartGrouping object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ChartGrouping deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ChartGrouping.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
