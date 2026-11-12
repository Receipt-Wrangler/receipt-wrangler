//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/claim_receipts_failure.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_receipts_result.g.dart';

/// ClaimReceiptsResult
///
/// Properties:
/// * [receiptIds] - Requested receipt ids that were added (or removed), in request order
/// * [failures] - One entry per requested receipt that was not added (or removed)
@BuiltValue()
abstract class ClaimReceiptsResult implements Built<ClaimReceiptsResult, ClaimReceiptsResultBuilder> {
  /// Requested receipt ids that were added (or removed), in request order
  @BuiltValueField(wireName: r'receiptIds')
  BuiltList<int> get receiptIds;

  /// One entry per requested receipt that was not added (or removed)
  @BuiltValueField(wireName: r'failures')
  BuiltList<ClaimReceiptsFailure> get failures;

  ClaimReceiptsResult._();

  factory ClaimReceiptsResult([void updates(ClaimReceiptsResultBuilder b)]) = _$ClaimReceiptsResult;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimReceiptsResultBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimReceiptsResult> get serializer => _$ClaimReceiptsResultSerializer();
}

class _$ClaimReceiptsResultSerializer implements PrimitiveSerializer<ClaimReceiptsResult> {
  @override
  final Iterable<Type> types = const [ClaimReceiptsResult, _$ClaimReceiptsResult];

  @override
  final String wireName = r'ClaimReceiptsResult';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimReceiptsResult object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'receiptIds';
    yield serializers.serialize(
      object.receiptIds,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
    yield r'failures';
    yield serializers.serialize(
      object.failures,
      specifiedType: const FullType(BuiltList, [FullType(ClaimReceiptsFailure)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ClaimReceiptsResult object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimReceiptsResultBuilder result,
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
        case r'failures':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ClaimReceiptsFailure)]),
          ) as BuiltList<ClaimReceiptsFailure>;
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
  ClaimReceiptsResult deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimReceiptsResultBuilder();
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

