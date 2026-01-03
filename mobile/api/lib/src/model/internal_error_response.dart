//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'internal_error_response.g.dart';

/// InternalErrorResponse
///
/// Properties:
/// * [errorMsg] - Error message
@BuiltValue()
abstract class InternalErrorResponse implements Built<InternalErrorResponse, InternalErrorResponseBuilder> {
  /// Error message
  @BuiltValueField(wireName: r'errorMsg')
  String get errorMsg;

  InternalErrorResponse._();

  factory InternalErrorResponse([void updates(InternalErrorResponseBuilder b)]) = _$InternalErrorResponse;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(InternalErrorResponseBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<InternalErrorResponse> get serializer => _$InternalErrorResponseSerializer();
}

class _$InternalErrorResponseSerializer implements PrimitiveSerializer<InternalErrorResponse> {
  @override
  final Iterable<Type> types = const [InternalErrorResponse, _$InternalErrorResponse];

  @override
  final String wireName = r'InternalErrorResponse';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    InternalErrorResponse object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'errorMsg';
    yield serializers.serialize(
      object.errorMsg,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    InternalErrorResponse object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required InternalErrorResponseBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'errorMsg':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.errorMsg = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  InternalErrorResponse deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = InternalErrorResponseBuilder();
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

