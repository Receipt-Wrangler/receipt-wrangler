//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';
import 'package:openapi/src/model/user_role.dart';

part 'claims.g.dart';

/// Claims
///
/// Properties:
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
abstract class Claims implements Built<Claims, ClaimsBuilder> {
  /// User foreign key
  @BuiltValueField(wireName: r'userId')
  int get userId;

  /// User's role
  @BuiltValueField(wireName: r'userRole')
  UserRole get userRole;
  // enum userRoleEnum {  ADMIN,  USER,  };

  /// Display name
  @BuiltValueField(wireName: r'displayName')
  String get displayName;

  /// Default avatar color
  @BuiltValueField(wireName: r'defaultAvatarColor')
  String get defaultAvatarColor;

  /// User's username used to login
  @BuiltValueField(wireName: r'username')
  String get username;

  /// Issuer
  @BuiltValueField(wireName: r'iss')
  String get iss;

  /// Subject
  @BuiltValueField(wireName: r'sub')
  String? get sub;

  /// Audience
  @BuiltValueField(wireName: r'aud')
  BuiltList<String>? get aud;

  /// Expiration time
  @BuiltValueField(wireName: r'exp')
  int get exp;

  /// Not before
  @BuiltValueField(wireName: r'nbf')
  int? get nbf;

  /// Issued at
  @BuiltValueField(wireName: r'iat')
  int? get iat;

  /// JWT ID
  @BuiltValueField(wireName: r'jti')
  String? get jti;

  Claims._();

  factory Claims([void updates(ClaimsBuilder b)]) = _$Claims;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimsBuilder b) => b
    ..userId = 0
    ..userRole = UserRole.USER
    ..displayName = ''
    ..defaultAvatarColor = ''
    ..username = ''
    ..iss = ''
    ..sub = ''
    ..aud = ListBuilder()
    ..exp = 0
    ..nbf = 0
    ..iat = 0
    ..jti = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<Claims> get serializer => _$ClaimsSerializer();
}

class _$ClaimsSerializer implements PrimitiveSerializer<Claims> {
  @override
  final Iterable<Type> types = const [Claims, _$Claims];

  @override
  final String wireName = r'Claims';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Claims object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'userId';
    yield serializers.serialize(
      object.userId,
      specifiedType: const FullType(int),
    );
    yield r'userRole';
    yield serializers.serialize(
      object.userRole,
      specifiedType: const FullType(UserRole),
    );
    yield r'displayName';
    yield serializers.serialize(
      object.displayName,
      specifiedType: const FullType(String),
    );
    yield r'defaultAvatarColor';
    yield serializers.serialize(
      object.defaultAvatarColor,
      specifiedType: const FullType(String),
    );
    yield r'username';
    yield serializers.serialize(
      object.username,
      specifiedType: const FullType(String),
    );
    yield r'iss';
    yield serializers.serialize(
      object.iss,
      specifiedType: const FullType(String),
    );
    if (object.sub != null) {
      yield r'sub';
      yield serializers.serialize(
        object.sub,
        specifiedType: const FullType(String),
      );
    }
    if (object.aud != null) {
      yield r'aud';
      yield serializers.serialize(
        object.aud,
        specifiedType: const FullType(BuiltList, [FullType(String)]),
      );
    }
    yield r'exp';
    yield serializers.serialize(
      object.exp,
      specifiedType: const FullType(int),
    );
    if (object.nbf != null) {
      yield r'nbf';
      yield serializers.serialize(
        object.nbf,
        specifiedType: const FullType(int),
      );
    }
    if (object.iat != null) {
      yield r'iat';
      yield serializers.serialize(
        object.iat,
        specifiedType: const FullType(int),
      );
    }
    if (object.jti != null) {
      yield r'jti';
      yield serializers.serialize(
        object.jti,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    Claims object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object,
            specifiedType: specifiedType)
        .toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimsBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        case r'userRole':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(UserRole),
          ) as UserRole;
          result.userRole = valueDes;
          break;
        case r'displayName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.displayName = valueDes;
          break;
        case r'defaultAvatarColor':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.defaultAvatarColor = valueDes;
          break;
        case r'username':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.username = valueDes;
          break;
        case r'iss':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.iss = valueDes;
          break;
        case r'sub':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.sub = valueDes;
          break;
        case r'aud':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(String)]),
          ) as BuiltList<String>;
          result.aud.replace(valueDes);
          break;
        case r'exp':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.exp = valueDes;
          break;
        case r'nbf':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.nbf = valueDes;
          break;
        case r'iat':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.iat = valueDes;
          break;
        case r'jti':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.jti = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Claims deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimsBuilder();
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
