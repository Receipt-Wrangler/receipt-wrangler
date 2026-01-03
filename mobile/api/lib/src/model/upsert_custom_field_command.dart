//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/upsert_custom_field_option_command.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/custom_field_type.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_custom_field_command.g.dart';

/// UpsertCustomFieldCommand
///
/// Properties:
/// * [name] - Custom Field name
/// * [type] 
/// * [description] - Custom Field description
/// * [options] 
@BuiltValue()
abstract class UpsertCustomFieldCommand implements Built<UpsertCustomFieldCommand, UpsertCustomFieldCommandBuilder> {
  /// Custom Field name
  @BuiltValueField(wireName: r'name')
  String get name;

  @BuiltValueField(wireName: r'type')
  CustomFieldType get type;
  // enum typeEnum {  TEXT,  DATE,  SELECT,  CURRENCY,  BOOLEAN,  };

  /// Custom Field description
  @BuiltValueField(wireName: r'description')
  String? get description;

  @BuiltValueField(wireName: r'options')
  BuiltList<UpsertCustomFieldOptionCommand>? get options;

  UpsertCustomFieldCommand._();

  factory UpsertCustomFieldCommand([void updates(UpsertCustomFieldCommandBuilder b)]) = _$UpsertCustomFieldCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertCustomFieldCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertCustomFieldCommand> get serializer => _$UpsertCustomFieldCommandSerializer();
}

class _$UpsertCustomFieldCommandSerializer implements PrimitiveSerializer<UpsertCustomFieldCommand> {
  @override
  final Iterable<Type> types = const [UpsertCustomFieldCommand, _$UpsertCustomFieldCommand];

  @override
  final String wireName = r'UpsertCustomFieldCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertCustomFieldCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'type';
    yield serializers.serialize(
      object.type,
      specifiedType: const FullType(CustomFieldType),
    );
    if (object.description != null) {
      yield r'description';
      yield serializers.serialize(
        object.description,
        specifiedType: const FullType(String),
      );
    }
    if (object.options != null) {
      yield r'options';
      yield serializers.serialize(
        object.options,
        specifiedType: const FullType(BuiltList, [FullType(UpsertCustomFieldOptionCommand)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertCustomFieldCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertCustomFieldCommandBuilder result,
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
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(CustomFieldType),
          ) as CustomFieldType;
          result.type = valueDes;
          break;
        case r'description':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.description = valueDes;
          break;
        case r'options':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(UpsertCustomFieldOptionCommand)]),
          ) as BuiltList<UpsertCustomFieldOptionCommand>;
          result.options.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertCustomFieldCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertCustomFieldCommandBuilder();
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

