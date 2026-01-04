//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'category_view.g.dart';

/// Category to relate receipts to
///
/// Properties:
/// * [createdAt] 
/// * [createdBy] 
/// * [id] 
/// * [name] - Name of the category
/// * [description] - Description of the category
/// * [updatedAt] 
/// * [numberOfReceipts] - Number of receipts associated with this category
@BuiltValue()
abstract class CategoryView implements Built<CategoryView, CategoryViewBuilder> {
  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  @BuiltValueField(wireName: r'id')
  int get id;

  /// Name of the category
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Description of the category
  @BuiltValueField(wireName: r'description')
  String? get description;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// Number of receipts associated with this category
  @BuiltValueField(wireName: r'numberOfReceipts')
  int get numberOfReceipts;

  CategoryView._();

  factory CategoryView([void updates(CategoryViewBuilder b)]) = _$CategoryView;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CategoryViewBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<CategoryView> get serializer => _$CategoryViewSerializer();
}

class _$CategoryViewSerializer implements PrimitiveSerializer<CategoryView> {
  @override
  final Iterable<Type> types = const [CategoryView, _$CategoryView];

  @override
  final String wireName = r'CategoryView';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CategoryView object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
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
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
    yield r'numberOfReceipts';
    yield serializers.serialize(
      object.numberOfReceipts,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    CategoryView object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CategoryViewBuilder result,
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
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        case r'numberOfReceipts':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.numberOfReceipts = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  CategoryView deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CategoryViewBuilder();
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

