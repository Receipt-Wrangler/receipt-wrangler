//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'reset_password_command.g.dart';

/// Command to reset user's password profile
///
/// Properties:
/// * [password] - User's new password
@BuiltValue()
abstract class ResetPasswordCommand implements Built<ResetPasswordCommand, ResetPasswordCommandBuilder> {
  /// User's new password
  @BuiltValueField(wireName: r'password')
  String get password;

  ResetPasswordCommand._();

  factory ResetPasswordCommand([void updates(ResetPasswordCommandBuilder b)]) = _$ResetPasswordCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ResetPasswordCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ResetPasswordCommand> get serializer => _$ResetPasswordCommandSerializer();
}

class _$ResetPasswordCommandSerializer implements PrimitiveSerializer<ResetPasswordCommand> {
  @override
  final Iterable<Type> types = const [ResetPasswordCommand, _$ResetPasswordCommand];

  @override
  final String wireName = r'ResetPasswordCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ResetPasswordCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'password';
    yield serializers.serialize(
      object.password,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ResetPasswordCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ResetPasswordCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
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
  ResetPasswordCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ResetPasswordCommandBuilder();
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

