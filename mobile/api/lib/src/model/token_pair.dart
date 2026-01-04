//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'token_pair.g.dart';

/// TokenPair
///
/// Properties:
/// * [jwt] - JWT token
/// * [refreshToken] - Refresh token
@BuiltValue()
abstract class TokenPair implements Built<TokenPair, TokenPairBuilder> {
  /// JWT token
  @BuiltValueField(wireName: r'jwt')
  String get jwt;

  /// Refresh token
  @BuiltValueField(wireName: r'refreshToken')
  String get refreshToken;

  TokenPair._();

  factory TokenPair([void updates(TokenPairBuilder b)]) = _$TokenPair;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(TokenPairBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<TokenPair> get serializer => _$TokenPairSerializer();
}

class _$TokenPairSerializer implements PrimitiveSerializer<TokenPair> {
  @override
  final Iterable<Type> types = const [TokenPair, _$TokenPair];

  @override
  final String wireName = r'TokenPair';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    TokenPair object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'jwt';
    yield serializers.serialize(
      object.jwt,
      specifiedType: const FullType(String),
    );
    yield r'refreshToken';
    yield serializers.serialize(
      object.refreshToken,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    TokenPair object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required TokenPairBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'jwt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.jwt = valueDes;
          break;
        case r'refreshToken':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.refreshToken = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  TokenPair deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = TokenPairBuilder();
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

