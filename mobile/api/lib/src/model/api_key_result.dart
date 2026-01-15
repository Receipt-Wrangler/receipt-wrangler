//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'api_key_result.g.dart';

/// ApiKeyResult
///
/// Properties:
/// * [key] - The generated API key
@BuiltValue()
abstract class ApiKeyResult implements Built<ApiKeyResult, ApiKeyResultBuilder> {
  /// The generated API key
  @BuiltValueField(wireName: r'key')
  String get key;

  ApiKeyResult._();

  factory ApiKeyResult([void updates(ApiKeyResultBuilder b)]) = _$ApiKeyResult;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ApiKeyResultBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ApiKeyResult> get serializer => _$ApiKeyResultSerializer();
}

class _$ApiKeyResultSerializer implements PrimitiveSerializer<ApiKeyResult> {
  @override
  final Iterable<Type> types = const [ApiKeyResult, _$ApiKeyResult];

  @override
  final String wireName = r'ApiKeyResult';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ApiKeyResult object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'key';
    yield serializers.serialize(
      object.key,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ApiKeyResult object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ApiKeyResultBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'key':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.key = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ApiKeyResult deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ApiKeyResultBuilder();
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

