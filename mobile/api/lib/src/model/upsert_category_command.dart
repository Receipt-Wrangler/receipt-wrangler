//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_category_command.g.dart';

/// UpsertCategoryCommand
///
/// Properties:
/// * [id] - Category id
/// * [name] - Category name
/// * [description] - Category description
@BuiltValue()
abstract class UpsertCategoryCommand implements Built<UpsertCategoryCommand, UpsertCategoryCommandBuilder> {
  /// Category id
  @BuiltValueField(wireName: r'id')
  int? get id;

  /// Category name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Category description
  @BuiltValueField(wireName: r'description')
  String? get description;

  UpsertCategoryCommand._();

  factory UpsertCategoryCommand([void updates(UpsertCategoryCommandBuilder b)]) = _$UpsertCategoryCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertCategoryCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertCategoryCommand> get serializer => _$UpsertCategoryCommandSerializer();
}

class _$UpsertCategoryCommandSerializer implements PrimitiveSerializer<UpsertCategoryCommand> {
  @override
  final Iterable<Type> types = const [UpsertCategoryCommand, _$UpsertCategoryCommand];

  @override
  final String wireName = r'UpsertCategoryCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertCategoryCommand object, {
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
    UpsertCategoryCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertCategoryCommandBuilder result,
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
  UpsertCategoryCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertCategoryCommandBuilder();
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

