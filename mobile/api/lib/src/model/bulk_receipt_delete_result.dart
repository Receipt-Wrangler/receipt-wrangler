//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/bulk_receipt_delete_failure.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bulk_receipt_delete_result.g.dart';

/// BulkReceiptDeleteResult
///
/// Properties:
/// * [deletedReceiptIds] - Receipt ids that were deleted
/// * [failures] - One entry per requested receipt that was not deleted
@BuiltValue()
abstract class BulkReceiptDeleteResult implements Built<BulkReceiptDeleteResult, BulkReceiptDeleteResultBuilder> {
  /// Receipt ids that were deleted
  @BuiltValueField(wireName: r'deletedReceiptIds')
  BuiltList<int> get deletedReceiptIds;

  /// One entry per requested receipt that was not deleted
  @BuiltValueField(wireName: r'failures')
  BuiltList<BulkReceiptDeleteFailure> get failures;

  BulkReceiptDeleteResult._();

  factory BulkReceiptDeleteResult([void updates(BulkReceiptDeleteResultBuilder b)]) = _$BulkReceiptDeleteResult;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BulkReceiptDeleteResultBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BulkReceiptDeleteResult> get serializer => _$BulkReceiptDeleteResultSerializer();
}

class _$BulkReceiptDeleteResultSerializer implements PrimitiveSerializer<BulkReceiptDeleteResult> {
  @override
  final Iterable<Type> types = const [BulkReceiptDeleteResult, _$BulkReceiptDeleteResult];

  @override
  final String wireName = r'BulkReceiptDeleteResult';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BulkReceiptDeleteResult object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'deletedReceiptIds';
    yield serializers.serialize(
      object.deletedReceiptIds,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
    yield r'failures';
    yield serializers.serialize(
      object.failures,
      specifiedType: const FullType(BuiltList, [FullType(BulkReceiptDeleteFailure)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BulkReceiptDeleteResult object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BulkReceiptDeleteResultBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'deletedReceiptIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.deletedReceiptIds.replace(valueDes);
          break;
        case r'failures':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(BulkReceiptDeleteFailure)]),
          ) as BuiltList<BulkReceiptDeleteFailure>;
          result.failures.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BulkReceiptDeleteResult deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BulkReceiptDeleteResultBuilder();
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

