//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/json_object.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_paged_request_filter.g.dart';

/// ReceiptPagedRequestFilter
///
/// Properties:
/// * [date] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [amount] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [name] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [paidBy] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [categories] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [tags] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [status] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [resolvedDate] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
/// * [createdAt] - Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
@BuiltValue()
abstract class ReceiptPagedRequestFilter implements Built<ReceiptPagedRequestFilter, ReceiptPagedRequestFilterBuilder> {
  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'date')
  JsonObject? get date;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'amount')
  JsonObject? get amount;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'name')
  JsonObject? get name;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'paidBy')
  JsonObject? get paidBy;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'categories')
  JsonObject? get categories;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'tags')
  JsonObject? get tags;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'status')
  JsonObject? get status;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'resolvedDate')
  JsonObject? get resolvedDate;

  /// Contains two keys: operation of type FilterOperation and value which can a different type depending on the field.
  @BuiltValueField(wireName: r'createdAt')
  JsonObject? get createdAt;

  ReceiptPagedRequestFilter._();

  factory ReceiptPagedRequestFilter([void updates(ReceiptPagedRequestFilterBuilder b)]) = _$ReceiptPagedRequestFilter;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptPagedRequestFilterBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptPagedRequestFilter> get serializer => _$ReceiptPagedRequestFilterSerializer();
}

class _$ReceiptPagedRequestFilterSerializer implements PrimitiveSerializer<ReceiptPagedRequestFilter> {
  @override
  final Iterable<Type> types = const [ReceiptPagedRequestFilter, _$ReceiptPagedRequestFilter];

  @override
  final String wireName = r'ReceiptPagedRequestFilter';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptPagedRequestFilter object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.date != null) {
      yield r'date';
      yield serializers.serialize(
        object.date,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.amount != null) {
      yield r'amount';
      yield serializers.serialize(
        object.amount,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.paidBy != null) {
      yield r'paidBy';
      yield serializers.serialize(
        object.paidBy,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.categories != null) {
      yield r'categories';
      yield serializers.serialize(
        object.categories,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.tags != null) {
      yield r'tags';
      yield serializers.serialize(
        object.tags,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.status != null) {
      yield r'status';
      yield serializers.serialize(
        object.status,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.resolvedDate != null) {
      yield r'resolvedDate';
      yield serializers.serialize(
        object.resolvedDate,
        specifiedType: const FullType(JsonObject),
      );
    }
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(JsonObject),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptPagedRequestFilter object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptPagedRequestFilterBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'date':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.date = valueDes;
          break;
        case r'amount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.amount = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.name = valueDes;
          break;
        case r'paidBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.paidBy = valueDes;
          break;
        case r'categories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.categories = valueDes;
          break;
        case r'tags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.tags = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.status = valueDes;
          break;
        case r'resolvedDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.resolvedDate = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(JsonObject),
          ) as JsonObject;
          result.createdAt = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptPagedRequestFilter deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptPagedRequestFilterBuilder();
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

