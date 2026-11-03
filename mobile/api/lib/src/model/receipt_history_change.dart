//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_history_value.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_history_change.g.dart';

/// One field-level change within a history entry
///
/// Properties:
/// * [field] - Which field changed. One of NAME, AMOUNT, DATE, GROUP, PAID_BY, STATUS, RESOLVED_DATE, CATEGORIES, TAGS, CUSTOM_FIELD, ITEM, SHARE, ITEM_CATEGORIES, ITEM_TAGS, COMMENT, IMAGE, EMAIL_INBOX, EMAIL_SENDER. Plain string.
/// * [kind] - One of CHANGED, ADDED, REMOVED. Plain string.
/// * [label] - What the change is about, where the field alone does not say: the item's name (ITEM, ITEM_CATEGORIES, ITEM_TAGS) or the custom field's name (CUSTOM_FIELD)
/// * [userId] - The share's charged-to user (SHARE) or the comment's author (COMMENT)
/// * [from] 
/// * [to] 
@BuiltValue()
abstract class ReceiptHistoryChange implements Built<ReceiptHistoryChange, ReceiptHistoryChangeBuilder> {
  /// Which field changed. One of NAME, AMOUNT, DATE, GROUP, PAID_BY, STATUS, RESOLVED_DATE, CATEGORIES, TAGS, CUSTOM_FIELD, ITEM, SHARE, ITEM_CATEGORIES, ITEM_TAGS, COMMENT, IMAGE, EMAIL_INBOX, EMAIL_SENDER. Plain string.
  @BuiltValueField(wireName: r'field')
  String get field;

  /// One of CHANGED, ADDED, REMOVED. Plain string.
  @BuiltValueField(wireName: r'kind')
  String get kind;

  /// What the change is about, where the field alone does not say: the item's name (ITEM, ITEM_CATEGORIES, ITEM_TAGS) or the custom field's name (CUSTOM_FIELD)
  @BuiltValueField(wireName: r'label')
  String? get label;

  /// The share's charged-to user (SHARE) or the comment's author (COMMENT)
  @BuiltValueField(wireName: r'userId')
  int? get userId;

  @BuiltValueField(wireName: r'from')
  ReceiptHistoryValue? get from;

  @BuiltValueField(wireName: r'to')
  ReceiptHistoryValue? get to;

  ReceiptHistoryChange._();

  factory ReceiptHistoryChange([void updates(ReceiptHistoryChangeBuilder b)]) = _$ReceiptHistoryChange;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptHistoryChangeBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptHistoryChange> get serializer => _$ReceiptHistoryChangeSerializer();
}

class _$ReceiptHistoryChangeSerializer implements PrimitiveSerializer<ReceiptHistoryChange> {
  @override
  final Iterable<Type> types = const [ReceiptHistoryChange, _$ReceiptHistoryChange];

  @override
  final String wireName = r'ReceiptHistoryChange';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptHistoryChange object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'field';
    yield serializers.serialize(
      object.field,
      specifiedType: const FullType(String),
    );
    yield r'kind';
    yield serializers.serialize(
      object.kind,
      specifiedType: const FullType(String),
    );
    if (object.label != null) {
      yield r'label';
      yield serializers.serialize(
        object.label,
        specifiedType: const FullType(String),
      );
    }
    if (object.userId != null) {
      yield r'userId';
      yield serializers.serialize(
        object.userId,
        specifiedType: const FullType(int),
      );
    }
    if (object.from != null) {
      yield r'from';
      yield serializers.serialize(
        object.from,
        specifiedType: const FullType(ReceiptHistoryValue),
      );
    }
    if (object.to != null) {
      yield r'to';
      yield serializers.serialize(
        object.to,
        specifiedType: const FullType(ReceiptHistoryValue),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptHistoryChange object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptHistoryChangeBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'field':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.field = valueDes;
          break;
        case r'kind':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.kind = valueDes;
          break;
        case r'label':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.label = valueDes;
          break;
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        case r'from':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptHistoryValue),
          ) as ReceiptHistoryValue;
          result.from.replace(valueDes);
          break;
        case r'to':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptHistoryValue),
          ) as ReceiptHistoryValue;
          result.to.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptHistoryChange deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptHistoryChangeBuilder();
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

