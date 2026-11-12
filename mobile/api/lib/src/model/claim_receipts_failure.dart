//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_receipts_failure.g.dart';

/// ClaimReceiptsFailure
///
/// Properties:
/// * [receiptId] - Requested receipt id that was not added (or removed)
/// * [message] - Why not. A receipt that does not exist, belongs to another group, or is hidden from the caller reads the same, so this never reveals that a receipt exists.
@BuiltValue()
abstract class ClaimReceiptsFailure implements Built<ClaimReceiptsFailure, ClaimReceiptsFailureBuilder> {
  /// Requested receipt id that was not added (or removed)
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  /// Why not. A receipt that does not exist, belongs to another group, or is hidden from the caller reads the same, so this never reveals that a receipt exists.
  @BuiltValueField(wireName: r'message')
  String get message;

  ClaimReceiptsFailure._();

  factory ClaimReceiptsFailure([void updates(ClaimReceiptsFailureBuilder b)]) = _$ClaimReceiptsFailure;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimReceiptsFailureBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimReceiptsFailure> get serializer => _$ClaimReceiptsFailureSerializer();
}

class _$ClaimReceiptsFailureSerializer implements PrimitiveSerializer<ClaimReceiptsFailure> {
  @override
  final Iterable<Type> types = const [ClaimReceiptsFailure, _$ClaimReceiptsFailure];

  @override
  final String wireName = r'ClaimReceiptsFailure';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimReceiptsFailure object, {
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
    ClaimReceiptsFailure object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimReceiptsFailureBuilder result,
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
  ClaimReceiptsFailure deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimReceiptsFailureBuilder();
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

