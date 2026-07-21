//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_detail.g.dart';

/// ReportDetail
///
/// Properties:
/// * [mode] 
/// * [by] - Dimension field key to aggregate leaf rows by (aggregate mode)
@BuiltValue()
abstract class ReportDetail implements Built<ReportDetail, ReportDetailBuilder> {
  @BuiltValueField(wireName: r'mode')
  ReportDetailModeEnum get mode;
  // enum modeEnum {  records,  aggregate,  };

  /// Dimension field key to aggregate leaf rows by (aggregate mode)
  @BuiltValueField(wireName: r'by')
  String? get by;

  ReportDetail._();

  factory ReportDetail([void updates(ReportDetailBuilder b)]) = _$ReportDetail;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportDetailBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportDetail> get serializer => _$ReportDetailSerializer();
}

class _$ReportDetailSerializer implements PrimitiveSerializer<ReportDetail> {
  @override
  final Iterable<Type> types = const [ReportDetail, _$ReportDetail];

  @override
  final String wireName = r'ReportDetail';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportDetail object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'mode';
    yield serializers.serialize(
      object.mode,
      specifiedType: const FullType(ReportDetailModeEnum),
    );
    if (object.by != null) {
      yield r'by';
      yield serializers.serialize(
        object.by,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportDetail object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportDetailBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'mode':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportDetailModeEnum),
          ) as ReportDetailModeEnum;
          result.mode = valueDes;
          break;
        case r'by':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.by = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportDetail deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportDetailBuilder();
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

class ReportDetailModeEnum extends EnumClass {

  @BuiltValueEnumConst(wireName: r'records')
  static const ReportDetailModeEnum records = _$reportDetailModeEnum_records;
  @BuiltValueEnumConst(wireName: r'aggregate')
  static const ReportDetailModeEnum aggregate = _$reportDetailModeEnum_aggregate;

  static Serializer<ReportDetailModeEnum> get serializer => _$reportDetailModeEnumSerializer;

  const ReportDetailModeEnum._(String name): super(name);

  static BuiltSet<ReportDetailModeEnum> get values => _$reportDetailModeEnumValues;
  static ReportDetailModeEnum valueOf(String name) => _$reportDetailModeEnumValueOf(name);
}

