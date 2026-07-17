//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_template_option.g.dart';

/// ReportTemplateOption
///
/// Properties:
/// * [id] 
/// * [name] 
/// * [groupIds] - The ids of the groups this template covers.
@BuiltValue()
abstract class ReportTemplateOption implements Built<ReportTemplateOption, ReportTemplateOptionBuilder> {
  @BuiltValueField(wireName: r'id')
  int get id;

  @BuiltValueField(wireName: r'name')
  String get name;

  /// The ids of the groups this template covers.
  @BuiltValueField(wireName: r'groupIds')
  BuiltList<int> get groupIds;

  ReportTemplateOption._();

  factory ReportTemplateOption([void updates(ReportTemplateOptionBuilder b)]) = _$ReportTemplateOption;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportTemplateOptionBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportTemplateOption> get serializer => _$ReportTemplateOptionSerializer();
}

class _$ReportTemplateOptionSerializer implements PrimitiveSerializer<ReportTemplateOption> {
  @override
  final Iterable<Type> types = const [ReportTemplateOption, _$ReportTemplateOption];

  @override
  final String wireName = r'ReportTemplateOption';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportTemplateOption object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'groupIds';
    yield serializers.serialize(
      object.groupIds,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportTemplateOption object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportTemplateOptionBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'groupIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.groupIds.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportTemplateOption deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportTemplateOptionBuilder();
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

