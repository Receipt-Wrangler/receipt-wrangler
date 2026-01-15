//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'pie_chart_data_point.g.dart';

/// PieChartDataPoint
///
/// Properties:
/// * [label] - Label for the pie chart slice
/// * [value] - Value for the pie chart slice
@BuiltValue()
abstract class PieChartDataPoint implements Built<PieChartDataPoint, PieChartDataPointBuilder> {
  /// Label for the pie chart slice
  @BuiltValueField(wireName: r'label')
  String get label;

  /// Value for the pie chart slice
  @BuiltValueField(wireName: r'value')
  double get value;

  PieChartDataPoint._();

  factory PieChartDataPoint([void updates(PieChartDataPointBuilder b)]) = _$PieChartDataPoint;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PieChartDataPointBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PieChartDataPoint> get serializer => _$PieChartDataPointSerializer();
}

class _$PieChartDataPointSerializer implements PrimitiveSerializer<PieChartDataPoint> {
  @override
  final Iterable<Type> types = const [PieChartDataPoint, _$PieChartDataPoint];

  @override
  final String wireName = r'PieChartDataPoint';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PieChartDataPoint object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'label';
    yield serializers.serialize(
      object.label,
      specifiedType: const FullType(String),
    );
    yield r'value';
    yield serializers.serialize(
      object.value,
      specifiedType: const FullType(double),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    PieChartDataPoint object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PieChartDataPointBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'label':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.label = valueDes;
          break;
        case r'value':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.value = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PieChartDataPoint deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PieChartDataPointBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

