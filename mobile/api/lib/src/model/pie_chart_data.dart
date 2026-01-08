//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/pie_chart_data_point.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'pie_chart_data.g.dart';

/// PieChartData
///
/// Properties:
/// * [data] - Array of pie chart data points
@BuiltValue()
abstract class PieChartData implements Built<PieChartData, PieChartDataBuilder> {
  /// Array of pie chart data points
  @BuiltValueField(wireName: r'data')
  BuiltList<PieChartDataPoint> get data;

  PieChartData._();

  factory PieChartData([void updates(PieChartDataBuilder b)]) = _$PieChartData;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PieChartDataBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PieChartData> get serializer => _$PieChartDataSerializer();
}

class _$PieChartDataSerializer implements PrimitiveSerializer<PieChartData> {
  @override
  final Iterable<Type> types = const [PieChartData, _$PieChartData];

  @override
  final String wireName = r'PieChartData';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PieChartData object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'data';
    yield serializers.serialize(
      object.data,
      specifiedType: const FullType(BuiltList, [FullType(PieChartDataPoint)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    PieChartData object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PieChartDataBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'data':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(PieChartDataPoint)]),
          ) as BuiltList<PieChartDataPoint>;
          result.data.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PieChartData deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PieChartDataBuilder();
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

