//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/chart_grouping.dart';
import 'package:openapi/src/model/receipt_paged_request_filter.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'pie_chart_data_command.g.dart';

/// PieChartDataCommand
///
/// Properties:
/// * [chartGrouping] - What to group the pie chart by
/// * [filter] - Optional filter for receipts
@BuiltValue()
abstract class PieChartDataCommand implements Built<PieChartDataCommand, PieChartDataCommandBuilder> {
  /// What to group the pie chart by
  @BuiltValueField(wireName: r'chartGrouping')
  ChartGrouping get chartGrouping;
  // enum chartGroupingEnum {  CATEGORIES,  TAGS,  PAIDBY,  };

  /// Optional filter for receipts
  @BuiltValueField(wireName: r'filter')
  ReceiptPagedRequestFilter? get filter;

  PieChartDataCommand._();

  factory PieChartDataCommand([void updates(PieChartDataCommandBuilder b)]) = _$PieChartDataCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PieChartDataCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PieChartDataCommand> get serializer => _$PieChartDataCommandSerializer();
}

class _$PieChartDataCommandSerializer implements PrimitiveSerializer<PieChartDataCommand> {
  @override
  final Iterable<Type> types = const [PieChartDataCommand, _$PieChartDataCommand];

  @override
  final String wireName = r'PieChartDataCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PieChartDataCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'chartGrouping';
    yield serializers.serialize(
      object.chartGrouping,
      specifiedType: const FullType(ChartGrouping),
    );
    if (object.filter != null) {
      yield r'filter';
      yield serializers.serialize(
        object.filter,
        specifiedType: const FullType(ReceiptPagedRequestFilter),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    PieChartDataCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PieChartDataCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'chartGrouping':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ChartGrouping),
          ) as ChartGrouping;
          result.chartGrouping = valueDes;
          break;
        case r'filter':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptPagedRequestFilter),
          ) as ReceiptPagedRequestFilter;
          result.filter.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PieChartDataCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PieChartDataCommandBuilder();
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

