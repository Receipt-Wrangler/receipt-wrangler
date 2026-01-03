//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/user_role.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'sign_up_command.g.dart';

/// SignUpCommand
///
/// Properties:
/// * [username] - User's username
/// * [password] - User's password
/// * [displayName] - User's displayname
/// * [isDummyUser] - Whether the user is a dummy user
/// * [userRole] - User's role
@BuiltValue()
abstract class SignUpCommand implements Built<SignUpCommand, SignUpCommandBuilder> {
  /// User's username
  @BuiltValueField(wireName: r'username')
  String get username;

  /// User's password
  @BuiltValueField(wireName: r'password')
  String get password;

  /// User's displayname
  @BuiltValueField(wireName: r'displayName')
  String? get displayName;

  /// Whether the user is a dummy user
  @BuiltValueField(wireName: r'isDummyUser')
  bool? get isDummyUser;

  /// User's role
  @BuiltValueField(wireName: r'userRole')
  UserRole? get userRole;
  // enum userRoleEnum {  ADMIN,  USER,  };

  SignUpCommand._();

  factory SignUpCommand([void updates(SignUpCommandBuilder b)]) = _$SignUpCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(SignUpCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<SignUpCommand> get serializer => _$SignUpCommandSerializer();
}

class _$SignUpCommandSerializer implements PrimitiveSerializer<SignUpCommand> {
  @override
  final Iterable<Type> types = const [SignUpCommand, _$SignUpCommand];

  @override
  final String wireName = r'SignUpCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    SignUpCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'username';
    yield serializers.serialize(
      object.username,
      specifiedType: const FullType(String),
    );
    yield r'password';
    yield serializers.serialize(
      object.password,
      specifiedType: const FullType(String),
    );
    if (object.displayName != null) {
      yield r'displayName';
      yield serializers.serialize(
        object.displayName,
        specifiedType: const FullType(String),
      );
    }
    if (object.isDummyUser != null) {
      yield r'isDummyUser';
      yield serializers.serialize(
        object.isDummyUser,
        specifiedType: const FullType(bool),
      );
    }
    if (object.userRole != null) {
      yield r'userRole';
      yield serializers.serialize(
        object.userRole,
        specifiedType: const FullType(UserRole),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    SignUpCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required SignUpCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'username':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.username = valueDes;
          break;
        case r'password':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.password = valueDes;
          break;
        case r'displayName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.displayName = valueDes;
          break;
        case r'isDummyUser':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isDummyUser = valueDes;
          break;
        case r'userRole':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(UserRole),
          ) as UserRole;
          result.userRole = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  SignUpCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = SignUpCommandBuilder();
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

