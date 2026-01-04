//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/upsert_category_command.dart';
import 'package:openapi/src/model/item_status.dart';
import 'package:openapi/src/model/upsert_tag_command.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_item_command.g.dart';

/// UpsertItemCommand
///
/// Properties:
/// * [amount] - Amount the item costs
/// * [chargedToUserId] - User foreign key
/// * [name] - Item name
/// * [receiptId] - Receipt foreign key
/// * [status] 
/// * [categories] - Categories associated to item
/// * [tags] - Tags associated to item
@BuiltValue()
abstract class UpsertItemCommand implements Built<UpsertItemCommand, UpsertItemCommandBuilder> {
  /// Amount the item costs
  @BuiltValueField(wireName: r'amount')
  String get amount;

  /// User foreign key
  @BuiltValueField(wireName: r'chargedToUserId')
  int get chargedToUserId;

  /// Item name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Receipt foreign key
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  @BuiltValueField(wireName: r'status')
  ItemStatus get status;
  // enum statusEnum {  OPEN,  RESOLVED,  DRAFT,  };

  /// Categories associated to item
  @BuiltValueField(wireName: r'categories')
  BuiltList<UpsertCategoryCommand>? get categories;

  /// Tags associated to item
  @BuiltValueField(wireName: r'tags')
  BuiltList<UpsertTagCommand>? get tags;

  UpsertItemCommand._();

  factory UpsertItemCommand([void updates(UpsertItemCommandBuilder b)]) = _$UpsertItemCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertItemCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertItemCommand> get serializer => _$UpsertItemCommandSerializer();
}

class _$UpsertItemCommandSerializer implements PrimitiveSerializer<UpsertItemCommand> {
  @override
  final Iterable<Type> types = const [UpsertItemCommand, _$UpsertItemCommand];

  @override
  final String wireName = r'UpsertItemCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertItemCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
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
        specifiedType: const FullType(BuiltList, [FullType(UpsertCategoryCommand)]),
      );
    }
    if (object.tags != null) {
      yield r'tags';
      yield serializers.serialize(
        object.tags,
        specifiedType: const FullType(BuiltList, [FullType(UpsertTagCommand)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertItemCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertItemCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
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
            specifiedType: const FullType(BuiltList, [FullType(UpsertCategoryCommand)]),
          ) as BuiltList<UpsertCategoryCommand>;
          result.categories.replace(valueDes);
          break;
        case r'tags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(UpsertTagCommand)]),
          ) as BuiltList<UpsertTagCommand>;
          result.tags.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertItemCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertItemCommandBuilder();
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

