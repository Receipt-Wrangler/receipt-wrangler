//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_period.g.dart';

/// ReportPeriod
///
/// Properties:
/// * [preset] 
/// * [startDate] - Start date (YYYY-MM-DD), read only when preset is custom
/// * [endDate] - End date (YYYY-MM-DD), read only when preset is custom
@BuiltValue()
abstract class ReportPeriod implements Built<ReportPeriod, ReportPeriodBuilder> {
  @BuiltValueField(wireName: r'preset')
  ReportPeriodPresetEnum get preset;
  // enum presetEnum {  this_month,  last_month,  mtd,  qtd,  ytd,  custom,  };

  /// Start date (YYYY-MM-DD), read only when preset is custom
  @BuiltValueField(wireName: r'startDate')
  String? get startDate;

  /// End date (YYYY-MM-DD), read only when preset is custom
  @BuiltValueField(wireName: r'endDate')
  String? get endDate;

  ReportPeriod._();

  factory ReportPeriod([void updates(ReportPeriodBuilder b)]) = _$ReportPeriod;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportPeriodBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportPeriod> get serializer => _$ReportPeriodSerializer();
}

class _$ReportPeriodSerializer implements PrimitiveSerializer<ReportPeriod> {
  @override
  final Iterable<Type> types = const [ReportPeriod, _$ReportPeriod];

  @override
  final String wireName = r'ReportPeriod';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportPeriod object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'preset';
    yield serializers.serialize(
      object.preset,
      specifiedType: const FullType(ReportPeriodPresetEnum),
    );
    if (object.startDate != null) {
      yield r'startDate';
      yield serializers.serialize(
        object.startDate,
        specifiedType: const FullType(String),
      );
    }
    if (object.endDate != null) {
      yield r'endDate';
      yield serializers.serialize(
        object.endDate,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportPeriod object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportPeriodBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'preset':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportPeriodPresetEnum),
          ) as ReportPeriodPresetEnum;
          result.preset = valueDes;
          break;
        case r'startDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.startDate = valueDes;
          break;
        case r'endDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.endDate = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportPeriod deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportPeriodBuilder();
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

class ReportPeriodPresetEnum extends EnumClass {

  @BuiltValueEnumConst(wireName: r'this_month')
  static const ReportPeriodPresetEnum thisMonth = _$reportPeriodPresetEnum_thisMonth;
  @BuiltValueEnumConst(wireName: r'last_month')
  static const ReportPeriodPresetEnum lastMonth = _$reportPeriodPresetEnum_lastMonth;
  @BuiltValueEnumConst(wireName: r'mtd')
  static const ReportPeriodPresetEnum mtd = _$reportPeriodPresetEnum_mtd;
  @BuiltValueEnumConst(wireName: r'qtd')
  static const ReportPeriodPresetEnum qtd = _$reportPeriodPresetEnum_qtd;
  @BuiltValueEnumConst(wireName: r'ytd')
  static const ReportPeriodPresetEnum ytd = _$reportPeriodPresetEnum_ytd;
  @BuiltValueEnumConst(wireName: r'custom')
  static const ReportPeriodPresetEnum custom = _$reportPeriodPresetEnum_custom;

  static Serializer<ReportPeriodPresetEnum> get serializer => _$reportPeriodPresetEnumSerializer;

  const ReportPeriodPresetEnum._(String name): super(name);

  static BuiltSet<ReportPeriodPresetEnum> get values => _$reportPeriodPresetEnumValues;
  static ReportPeriodPresetEnum valueOf(String name) => _$reportPeriodPresetEnumValueOf(name);
}

