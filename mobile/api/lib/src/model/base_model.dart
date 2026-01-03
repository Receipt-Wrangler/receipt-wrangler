//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'base_model.g.dart';

/// BaseModel
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
@BuiltValue(instantiable: false)
abstract class BaseModel  {
  @BuiltValueField(wireName: r'id')
  int get id;

  @BuiltValueField(wireName: r'createdAt')
  String get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  /// Created by entity's name
  @BuiltValueField(wireName: r'createdByString')
  String? get createdByString;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  @BuiltValueSerializer(custom: true)
  static Serializer<BaseModel> get serializer => _$BaseModelSerializer();
}

class _$BaseModelSerializer implements PrimitiveSerializer<BaseModel> {
  @override
  final Iterable<Type> types = const [BaseModel];

  @override
  final String wireName = r'BaseModel';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BaseModel object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
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
    BaseModel object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  @override
  BaseModel deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return serializers.deserialize(serialized, specifiedType: FullType($BaseModel)) as $BaseModel;
  }
}

/// a concrete implementation of [BaseModel], since [BaseModel] is not instantiable
@BuiltValue(instantiable: true)
abstract class $BaseModel implements BaseModel, Built<$BaseModel, $BaseModelBuilder> {
  $BaseModel._();

  factory $BaseModel([void Function($BaseModelBuilder)? updates]) = _$$BaseModel;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults($BaseModelBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<$BaseModel> get serializer => _$$BaseModelSerializer();
}

class _$$BaseModelSerializer implements PrimitiveSerializer<$BaseModel> {
  @override
  final Iterable<Type> types = const [$BaseModel, _$$BaseModel];

  @override
  final String wireName = r'$BaseModel';

  @override
  Object serialize(
    Serializers serializers,
    $BaseModel object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return serializers.serialize(object, specifiedType: FullType(BaseModel))!;
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BaseModelBuilder result,
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
  $BaseModel deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = $BaseModelBuilder();
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

