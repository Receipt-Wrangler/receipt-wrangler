//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_system_email_command.g.dart';

/// UpsertSystemEmailCommand
///
/// Properties:
/// * [host] - IMAP host
/// * [port] - IMAP port
/// * [username] - IMAP username
/// * [password] - IMAP password
/// * [useStartTLS] - Whether to use STARTTLS
@BuiltValue()
abstract class UpsertSystemEmailCommand implements Built<UpsertSystemEmailCommand, UpsertSystemEmailCommandBuilder> {
  /// IMAP host
  @BuiltValueField(wireName: r'host')
  String get host;

  /// IMAP port
  @BuiltValueField(wireName: r'port')
  String get port;

  /// IMAP username
  @BuiltValueField(wireName: r'username')
  String get username;

  /// IMAP password
  @BuiltValueField(wireName: r'password')
  String get password;

  /// Whether to use STARTTLS
  @BuiltValueField(wireName: r'useStartTLS')
  bool? get useStartTLS;

  UpsertSystemEmailCommand._();

  factory UpsertSystemEmailCommand([void updates(UpsertSystemEmailCommandBuilder b)]) = _$UpsertSystemEmailCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertSystemEmailCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertSystemEmailCommand> get serializer => _$UpsertSystemEmailCommandSerializer();
}

class _$UpsertSystemEmailCommandSerializer implements PrimitiveSerializer<UpsertSystemEmailCommand> {
  @override
  final Iterable<Type> types = const [UpsertSystemEmailCommand, _$UpsertSystemEmailCommand];

  @override
  final String wireName = r'UpsertSystemEmailCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertSystemEmailCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'host';
    yield serializers.serialize(
      object.host,
      specifiedType: const FullType(String),
    );
    yield r'port';
    yield serializers.serialize(
      object.port,
      specifiedType: const FullType(String),
    );
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
    if (object.useStartTLS != null) {
      yield r'useStartTLS';
      yield serializers.serialize(
        object.useStartTLS,
        specifiedType: const FullType(bool),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertSystemEmailCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertSystemEmailCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'host':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.host = valueDes;
          break;
        case r'port':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.port = valueDes;
          break;
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
        case r'useStartTLS':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.useStartTLS = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertSystemEmailCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertSystemEmailCommandBuilder();
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

