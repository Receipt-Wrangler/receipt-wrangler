//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/claims.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/token_pair.dart';
import 'package:openapi/src/model/user_role.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';
import 'package:one_of/any_of.dart';

part 'get_new_refresh_token200_response.g.dart';

/// GetNewRefreshToken200Response
///
/// Properties:
/// * [jwt] - JWT token
/// * [refreshToken] - Refresh token
/// * [userId] - User foreign key
/// * [userRole] - User's role
/// * [displayName] - Display name
/// * [defaultAvatarColor] - Default avatar color
/// * [username] - User's username used to login
/// * [iss] - Issuer
/// * [sub] - Subject
/// * [aud] - Audience
/// * [exp] - Expiration time
/// * [nbf] - Not before
/// * [iat] - Issued at
/// * [jti] - JWT ID
@BuiltValue()
abstract class GetNewRefreshToken200Response implements Built<GetNewRefreshToken200Response, GetNewRefreshToken200ResponseBuilder> {
  /// Any Of [Claims], [TokenPair]
  AnyOf get anyOf;

  GetNewRefreshToken200Response._();

  factory GetNewRefreshToken200Response([void updates(GetNewRefreshToken200ResponseBuilder b)]) = _$GetNewRefreshToken200Response;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GetNewRefreshToken200ResponseBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<GetNewRefreshToken200Response> get serializer => _$GetNewRefreshToken200ResponseSerializer();
}

class _$GetNewRefreshToken200ResponseSerializer implements PrimitiveSerializer<GetNewRefreshToken200Response> {
  @override
  final Iterable<Type> types = const [GetNewRefreshToken200Response, _$GetNewRefreshToken200Response];

  @override
  final String wireName = r'GetNewRefreshToken200Response';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GetNewRefreshToken200Response object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
  }

  @override
  Object serialize(
    Serializers serializers,
    GetNewRefreshToken200Response object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final anyOf = object.anyOf;
    return serializers.serialize(anyOf, specifiedType: FullType(AnyOf, anyOf.valueTypes.map((type) => FullType(type)).toList()))!;
  }

  @override
  GetNewRefreshToken200Response deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GetNewRefreshToken200ResponseBuilder();
    Object? anyOfDataSrc;
    final targetType = const FullType(AnyOf, [FullType(TokenPair), FullType(Claims), ]);
    anyOfDataSrc = serialized;
    result.anyOf = serializers.deserialize(anyOfDataSrc, specifiedType: targetType) as AnyOf;
    return result.build();
  }
}

