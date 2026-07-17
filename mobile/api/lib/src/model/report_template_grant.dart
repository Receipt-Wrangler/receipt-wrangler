//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_template_grant.g.dart';

/// ReportTemplateGrant
///
/// Properties:
/// * [reportTemplateId] 
/// * [permissions] - The actions a GROUP role may perform on this report template (read, generate, update, delete, duplicate).
@BuiltValue()
abstract class ReportTemplateGrant implements Built<ReportTemplateGrant, ReportTemplateGrantBuilder> {
  @BuiltValueField(wireName: r'reportTemplateId')
  int get reportTemplateId;

  /// The actions a GROUP role may perform on this report template (read, generate, update, delete, duplicate).
  @BuiltValueField(wireName: r'permissions')
  BuiltList<String> get permissions;

  ReportTemplateGrant._();

  factory ReportTemplateGrant([void updates(ReportTemplateGrantBuilder b)]) = _$ReportTemplateGrant;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportTemplateGrantBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportTemplateGrant> get serializer => _$ReportTemplateGrantSerializer();
}

class _$ReportTemplateGrantSerializer implements PrimitiveSerializer<ReportTemplateGrant> {
  @override
  final Iterable<Type> types = const [ReportTemplateGrant, _$ReportTemplateGrant];

  @override
  final String wireName = r'ReportTemplateGrant';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportTemplateGrant object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'reportTemplateId';
    yield serializers.serialize(
      object.reportTemplateId,
      specifiedType: const FullType(int),
    );
    yield r'permissions';
    yield serializers.serialize(
      object.permissions,
      specifiedType: const FullType(BuiltList, [FullType(String)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportTemplateGrant object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportTemplateGrantBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'reportTemplateId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.reportTemplateId = valueDes;
          break;
        case r'permissions':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(String)]),
          ) as BuiltList<String>;
          result.permissions.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportTemplateGrant deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportTemplateGrantBuilder();
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

