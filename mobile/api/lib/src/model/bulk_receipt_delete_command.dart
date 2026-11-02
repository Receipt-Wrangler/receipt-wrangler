//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bulk_receipt_delete_command.g.dart';

/// BulkReceiptDeleteCommand
///
/// Properties:
/// * [receiptIds] - Receipt ids to delete. Between 1 and 500 ids, with no duplicates
@BuiltValue()
abstract class BulkReceiptDeleteCommand implements Built<BulkReceiptDeleteCommand, BulkReceiptDeleteCommandBuilder> {
  /// Receipt ids to delete. Between 1 and 500 ids, with no duplicates
  @BuiltValueField(wireName: r'receiptIds')
  BuiltList<int> get receiptIds;

  BulkReceiptDeleteCommand._();

  factory BulkReceiptDeleteCommand([void updates(BulkReceiptDeleteCommandBuilder b)]) = _$BulkReceiptDeleteCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BulkReceiptDeleteCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BulkReceiptDeleteCommand> get serializer => _$BulkReceiptDeleteCommandSerializer();
}

class _$BulkReceiptDeleteCommandSerializer implements PrimitiveSerializer<BulkReceiptDeleteCommand> {
  @override
  final Iterable<Type> types = const [BulkReceiptDeleteCommand, _$BulkReceiptDeleteCommand];

  @override
  final String wireName = r'BulkReceiptDeleteCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BulkReceiptDeleteCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'receiptIds';
    yield serializers.serialize(
      object.receiptIds,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BulkReceiptDeleteCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BulkReceiptDeleteCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'receiptIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.receiptIds.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BulkReceiptDeleteCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BulkReceiptDeleteCommandBuilder();
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

