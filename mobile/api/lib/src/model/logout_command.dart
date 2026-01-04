//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'logout_command.g.dart';

/// LogoutCommand
///
/// Properties:
/// * [refreshToken] - Refresh token
@BuiltValue()
abstract class LogoutCommand implements Built<LogoutCommand, LogoutCommandBuilder> {
  /// Refresh token
  @BuiltValueField(wireName: r'refreshToken')
  String get refreshToken;

  LogoutCommand._();

  factory LogoutCommand([void updates(LogoutCommandBuilder b)]) = _$LogoutCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(LogoutCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<LogoutCommand> get serializer => _$LogoutCommandSerializer();
}

class _$LogoutCommandSerializer implements PrimitiveSerializer<LogoutCommand> {
  @override
  final Iterable<Type> types = const [LogoutCommand, _$LogoutCommand];

  @override
  final String wireName = r'LogoutCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    LogoutCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'refreshToken';
    yield serializers.serialize(
      object.refreshToken,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    LogoutCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required LogoutCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
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
  LogoutCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = LogoutCommandBuilder();
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

