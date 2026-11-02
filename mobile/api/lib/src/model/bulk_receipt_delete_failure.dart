//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bulk_receipt_delete_failure.g.dart';

/// BulkReceiptDeleteFailure
///
/// Properties:
/// * [receiptId] - Requested receipt id that was not deleted
/// * [message] - Why it was not deleted. A missing receipt and one the caller may not delete read the same, so this never reveals that a receipt exists
@BuiltValue()
abstract class BulkReceiptDeleteFailure implements Built<BulkReceiptDeleteFailure, BulkReceiptDeleteFailureBuilder> {
  /// Requested receipt id that was not deleted
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  /// Why it was not deleted. A missing receipt and one the caller may not delete read the same, so this never reveals that a receipt exists
  @BuiltValueField(wireName: r'message')
  String get message;

  BulkReceiptDeleteFailure._();

  factory BulkReceiptDeleteFailure([void updates(BulkReceiptDeleteFailureBuilder b)]) = _$BulkReceiptDeleteFailure;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BulkReceiptDeleteFailureBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BulkReceiptDeleteFailure> get serializer => _$BulkReceiptDeleteFailureSerializer();
}

class _$BulkReceiptDeleteFailureSerializer implements PrimitiveSerializer<BulkReceiptDeleteFailure> {
  @override
  final Iterable<Type> types = const [BulkReceiptDeleteFailure, _$BulkReceiptDeleteFailure];

  @override
  final String wireName = r'BulkReceiptDeleteFailure';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BulkReceiptDeleteFailure object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'receiptId';
    yield serializers.serialize(
      object.receiptId,
      specifiedType: const FullType(int),
    );
    yield r'message';
    yield serializers.serialize(
      object.message,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BulkReceiptDeleteFailure object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BulkReceiptDeleteFailureBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'message':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.message = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BulkReceiptDeleteFailure deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BulkReceiptDeleteFailureBuilder();
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

