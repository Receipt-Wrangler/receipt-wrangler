//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_tag_command.g.dart';

/// Tag to relate receipts to
///
/// Properties:
/// * [id] - Tag id
/// * [name] - Tag name
/// * [description] - Tag description
@BuiltValue()
abstract class UpsertTagCommand implements Built<UpsertTagCommand, UpsertTagCommandBuilder> {
  /// Tag id
  @BuiltValueField(wireName: r'id')
  int? get id;

  /// Tag name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Tag description
  @BuiltValueField(wireName: r'description')
  String? get description;

  UpsertTagCommand._();

  factory UpsertTagCommand([void updates(UpsertTagCommandBuilder b)]) = _$UpsertTagCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertTagCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertTagCommand> get serializer => _$UpsertTagCommandSerializer();
}

class _$UpsertTagCommandSerializer implements PrimitiveSerializer<UpsertTagCommand> {
  @override
  final Iterable<Type> types = const [UpsertTagCommand, _$UpsertTagCommand];

  @override
  final String wireName = r'UpsertTagCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertTagCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.id != null) {
      yield r'id';
      yield serializers.serialize(
        object.id,
        specifiedType: const FullType(int),
      );
    }
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
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertTagCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertTagCommandBuilder result,
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
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertTagCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertTagCommandBuilder();
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

