//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/associated_api_keys.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'api_key_filter.g.dart';

/// ApiKeyFilter
///
/// Properties:
/// * [associatedApiKeys] 
@BuiltValue()
abstract class ApiKeyFilter implements Built<ApiKeyFilter, ApiKeyFilterBuilder> {
  @BuiltValueField(wireName: r'associatedApiKeys')
  AssociatedApiKeys? get associatedApiKeys;
  // enum associatedApiKeysEnum {  MINE,  ALL,  };

  ApiKeyFilter._();

  factory ApiKeyFilter([void updates(ApiKeyFilterBuilder b)]) = _$ApiKeyFilter;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ApiKeyFilterBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ApiKeyFilter> get serializer => _$ApiKeyFilterSerializer();
}

class _$ApiKeyFilterSerializer implements PrimitiveSerializer<ApiKeyFilter> {
  @override
  final Iterable<Type> types = const [ApiKeyFilter, _$ApiKeyFilter];

  @override
  final String wireName = r'ApiKeyFilter';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ApiKeyFilter object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.associatedApiKeys != null) {
      yield r'associatedApiKeys';
      yield serializers.serialize(
        object.associatedApiKeys,
        specifiedType: const FullType(AssociatedApiKeys),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ApiKeyFilter object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ApiKeyFilterBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'associatedApiKeys':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(AssociatedApiKeys),
          ) as AssociatedApiKeys;
          result.associatedApiKeys = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ApiKeyFilter deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ApiKeyFilterBuilder();
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

