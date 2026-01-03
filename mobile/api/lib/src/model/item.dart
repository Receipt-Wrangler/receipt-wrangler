//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/category.dart';
import 'package:openapi/src/model/tag.dart';
import 'package:openapi/src/model/item_status.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'item.g.dart';

/// Itemized item on a receipt
///
/// Properties:
/// * [isTaxed] - Is taxed (not used)
/// * [amount] - Amount the item costs
/// * [chargedToUserId] - User foreign key
/// * [createdAt] 
/// * [createdBy] 
/// * [id] 
/// * [name] - Item name
/// * [receiptId] - Receipt foreign key
/// * [status] 
/// * [categories] - Categories associated to the item
/// * [tags] - Tags associated to the item
/// * [updatedAt] 
@BuiltValue()
abstract class Item implements Built<Item, ItemBuilder> {
  /// Is taxed (not used)
  @BuiltValueField(wireName: r'IsTaxed')
  bool? get isTaxed;

  /// Amount the item costs
  @BuiltValueField(wireName: r'amount')
  String get amount;

  /// User foreign key
  @BuiltValueField(wireName: r'chargedToUserId')
  int get chargedToUserId;

  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  @BuiltValueField(wireName: r'id')
  int? get id;

  /// Item name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Receipt foreign key
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  @BuiltValueField(wireName: r'status')
  ItemStatus get status;
  // enum statusEnum {  OPEN,  RESOLVED,  DRAFT,  };

  /// Categories associated to the item
  @BuiltValueField(wireName: r'categories')
  BuiltList<Category>? get categories;

  /// Tags associated to the item
  @BuiltValueField(wireName: r'tags')
  BuiltList<Tag>? get tags;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  Item._();

  factory Item([void updates(ItemBuilder b)]) = _$Item;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ItemBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Item> get serializer => _$ItemSerializer();
}

class _$ItemSerializer implements PrimitiveSerializer<Item> {
  @override
  final Iterable<Type> types = const [Item, _$Item];

  @override
  final String wireName = r'Item';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Item object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.isTaxed != null) {
      yield r'IsTaxed';
      yield serializers.serialize(
        object.isTaxed,
        specifiedType: const FullType(bool),
      );
    }
    yield r'amount';
    yield serializers.serialize(
      object.amount,
      specifiedType: const FullType(String),
    );
    yield r'chargedToUserId';
    yield serializers.serialize(
      object.chargedToUserId,
      specifiedType: const FullType(int),
    );
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
    yield r'receiptId';
    yield serializers.serialize(
      object.receiptId,
      specifiedType: const FullType(int),
    );
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(ItemStatus),
    );
    if (object.categories != null) {
      yield r'categories';
      yield serializers.serialize(
        object.categories,
        specifiedType: const FullType(BuiltList, [FullType(Category)]),
      );
    }
    if (object.tags != null) {
      yield r'tags';
      yield serializers.serialize(
        object.tags,
        specifiedType: const FullType(BuiltList, [FullType(Tag)]),
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
    Item object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ItemBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'IsTaxed':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isTaxed = valueDes;
          break;
        case r'amount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.amount = valueDes;
          break;
        case r'chargedToUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.chargedToUserId = valueDes;
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
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ItemStatus),
          ) as ItemStatus;
          result.status = valueDes;
          break;
        case r'categories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Category)]),
          ) as BuiltList<Category>;
          result.categories.replace(valueDes);
          break;
        case r'tags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Tag)]),
          ) as BuiltList<Tag>;
          result.tags.replace(valueDes);
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
  Item deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ItemBuilder();
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

