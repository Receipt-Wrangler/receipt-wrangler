//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'check_email_connectivity_command.g.dart';

/// CheckEmailConnectivityCommand
///
/// Properties:
/// * [id] - System email id
/// * [host] - IMAP host
/// * [port] - IMAP port
/// * [username] - IMAP username
/// * [password] - IMAP password
@BuiltValue()
abstract class CheckEmailConnectivityCommand implements Built<CheckEmailConnectivityCommand, CheckEmailConnectivityCommandBuilder> {
  /// System email id
  @BuiltValueField(wireName: r'id')
  int? get id;

  /// IMAP host
  @BuiltValueField(wireName: r'host')
  String? get host;

  /// IMAP port
  @BuiltValueField(wireName: r'port')
  int? get port;

  /// IMAP username
  @BuiltValueField(wireName: r'username')
  String? get username;

  /// IMAP password
  @BuiltValueField(wireName: r'password')
  String? get password;

  CheckEmailConnectivityCommand._();

  factory CheckEmailConnectivityCommand([void updates(CheckEmailConnectivityCommandBuilder b)]) = _$CheckEmailConnectivityCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CheckEmailConnectivityCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<CheckEmailConnectivityCommand> get serializer => _$CheckEmailConnectivityCommandSerializer();
}

class _$CheckEmailConnectivityCommandSerializer implements PrimitiveSerializer<CheckEmailConnectivityCommand> {
  @override
  final Iterable<Type> types = const [CheckEmailConnectivityCommand, _$CheckEmailConnectivityCommand];

  @override
  final String wireName = r'CheckEmailConnectivityCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CheckEmailConnectivityCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.id != null) {
      yield r'id';
      yield serializers.serialize(
        object.id,
        specifiedType: const FullType(int),
      );
    }
    if (object.host != null) {
      yield r'host';
      yield serializers.serialize(
        object.host,
        specifiedType: const FullType(String),
      );
    }
    if (object.port != null) {
      yield r'port';
      yield serializers.serialize(
        object.port,
        specifiedType: const FullType(int),
      );
    }
    if (object.username != null) {
      yield r'username';
      yield serializers.serialize(
        object.username,
        specifiedType: const FullType(String),
      );
    }
    if (object.password != null) {
      yield r'password';
      yield serializers.serialize(
        object.password,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    CheckEmailConnectivityCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CheckEmailConnectivityCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
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
            specifiedType: const FullType(int),
          ) as int;
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
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  CheckEmailConnectivityCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CheckEmailConnectivityCommandBuilder();
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

