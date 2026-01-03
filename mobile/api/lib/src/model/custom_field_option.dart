//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/base_model.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'custom_field_option.g.dart';

/// CustomFieldOption
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
/// * [value] - Custom Field Option value
/// * [customFieldId] - Custom Field Id
@BuiltValue()
abstract class CustomFieldOption implements BaseModel, Built<CustomFieldOption, CustomFieldOptionBuilder> {
  /// Custom Field Id
  @BuiltValueField(wireName: r'customFieldId')
  int get customFieldId;

  /// Custom Field Option value
  @BuiltValueField(wireName: r'value')
  String? get value;

  CustomFieldOption._();

  factory CustomFieldOption([void updates(CustomFieldOptionBuilder b)]) = _$CustomFieldOption;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CustomFieldOptionBuilder b) => b
      ..createdBy = 0
      ..createdByString = ''
      ..updatedAt = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<CustomFieldOption> get serializer => _$CustomFieldOptionSerializer();
}

class _$CustomFieldOptionSerializer implements PrimitiveSerializer<CustomFieldOption> {
  @override
  final Iterable<Type> types = const [CustomFieldOption, _$CustomFieldOption];

  @override
  final String wireName = r'CustomFieldOption';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CustomFieldOption object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    yield r'customFieldId';
    yield serializers.serialize(
      object.customFieldId,
      specifiedType: const FullType(int),
    );
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.value != null) {
      yield r'value';
      yield serializers.serialize(
        object.value,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
      );
    }
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    CustomFieldOption object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CustomFieldOptionBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'customFieldId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.customFieldId = valueDes;
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'value':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.value = valueDes;
          break;
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  CustomFieldOption deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CustomFieldOptionBuilder();
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

