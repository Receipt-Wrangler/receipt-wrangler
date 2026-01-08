//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/api_key_scope.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_api_key_command.g.dart';

/// UpsertApiKeyCommand
///
/// Properties:
/// * [name] - API key name
/// * [description] - API key description
/// * [scope] 
@BuiltValue()
abstract class UpsertApiKeyCommand implements Built<UpsertApiKeyCommand, UpsertApiKeyCommandBuilder> {
  /// API key name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// API key description
  @BuiltValueField(wireName: r'description')
  String? get description;

  @BuiltValueField(wireName: r'scope')
  ApiKeyScope get scope;
  // enum scopeEnum {  r,  w,  rw,  };

  UpsertApiKeyCommand._();

  factory UpsertApiKeyCommand([void updates(UpsertApiKeyCommandBuilder b)]) = _$UpsertApiKeyCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertApiKeyCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertApiKeyCommand> get serializer => _$UpsertApiKeyCommandSerializer();
}

class _$UpsertApiKeyCommandSerializer implements PrimitiveSerializer<UpsertApiKeyCommand> {
  @override
  final Iterable<Type> types = const [UpsertApiKeyCommand, _$UpsertApiKeyCommand];

  @override
  final String wireName = r'UpsertApiKeyCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertApiKeyCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.description != null) {
      yield r'description';
      yield serializers.serialize(
        object.description,
        specifiedType: const FullType(String),
      );
    }
    yield r'scope';
    yield serializers.serialize(
      object.scope,
      specifiedType: const FullType(ApiKeyScope),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertApiKeyCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertApiKeyCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'description':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.description = valueDes;
          break;
        case r'scope':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ApiKeyScope),
          ) as ApiKeyScope;
          result.scope = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertApiKeyCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertApiKeyCommandBuilder();
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

