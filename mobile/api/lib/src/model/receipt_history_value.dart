//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_history_value.g.dart';

/// A typed value for the client to format
///
/// Properties:
/// * [type] - How to format the value. One of TEXT, CURRENCY (decimal string), DATE (calendar day, yyyy-MM-dd), INSTANT (RFC 3339), USER (see userId), STATUS (a receipt status key), BOOLEAN (\"true\" / \"false\"). Plain string.
/// * [value] 
/// * [userId] 
@BuiltValue()
abstract class ReceiptHistoryValue implements Built<ReceiptHistoryValue, ReceiptHistoryValueBuilder> {
  /// How to format the value. One of TEXT, CURRENCY (decimal string), DATE (calendar day, yyyy-MM-dd), INSTANT (RFC 3339), USER (see userId), STATUS (a receipt status key), BOOLEAN (\"true\" / \"false\"). Plain string.
  @BuiltValueField(wireName: r'type')
  String get type;

  @BuiltValueField(wireName: r'value')
  String? get value;

  @BuiltValueField(wireName: r'userId')
  int? get userId;

  ReceiptHistoryValue._();

  factory ReceiptHistoryValue([void updates(ReceiptHistoryValueBuilder b)]) = _$ReceiptHistoryValue;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptHistoryValueBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptHistoryValue> get serializer => _$ReceiptHistoryValueSerializer();
}

class _$ReceiptHistoryValueSerializer implements PrimitiveSerializer<ReceiptHistoryValue> {
  @override
  final Iterable<Type> types = const [ReceiptHistoryValue, _$ReceiptHistoryValue];

  @override
  final String wireName = r'ReceiptHistoryValue';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptHistoryValue object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'type';
    yield serializers.serialize(
      object.type,
      specifiedType: const FullType(String),
    );
    if (object.value != null) {
      yield r'value';
      yield serializers.serialize(
        object.value,
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
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptHistoryValue object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptHistoryValueBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.type = valueDes;
          break;
        case r'value':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.value = valueDes;
          break;
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptHistoryValue deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptHistoryValueBuilder();
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

