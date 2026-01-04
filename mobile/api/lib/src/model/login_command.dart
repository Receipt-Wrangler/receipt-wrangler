//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'login_command.g.dart';

/// LoginCommand
///
/// Properties:
/// * [username] - User's username
/// * [password] - User's password
@BuiltValue()
abstract class LoginCommand implements Built<LoginCommand, LoginCommandBuilder> {
  /// User's username
  @BuiltValueField(wireName: r'username')
  String get username;

  /// User's password
  @BuiltValueField(wireName: r'password')
  String get password;

  LoginCommand._();

  factory LoginCommand([void updates(LoginCommandBuilder b)]) = _$LoginCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(LoginCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<LoginCommand> get serializer => _$LoginCommandSerializer();
}

class _$LoginCommandSerializer implements PrimitiveSerializer<LoginCommand> {
  @override
  final Iterable<Type> types = const [LoginCommand, _$LoginCommand];

  @override
  final String wireName = r'LoginCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    LoginCommand object, {
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
  }

  @override
  Object serialize(
    Serializers serializers,
    LoginCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required LoginCommandBuilder result,
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
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  LoginCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = LoginCommandBuilder();
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

