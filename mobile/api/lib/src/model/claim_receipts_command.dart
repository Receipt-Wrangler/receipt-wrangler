//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_receipts_command.g.dart';

/// ClaimReceiptsCommand
///
/// Properties:
/// * [receiptIds] - Receipt ids to add or remove. Between 1 and 500 ids, with no duplicates
@BuiltValue()
abstract class ClaimReceiptsCommand implements Built<ClaimReceiptsCommand, ClaimReceiptsCommandBuilder> {
  /// Receipt ids to add or remove. Between 1 and 500 ids, with no duplicates
  @BuiltValueField(wireName: r'receiptIds')
  BuiltList<int> get receiptIds;

  ClaimReceiptsCommand._();

  factory ClaimReceiptsCommand([void updates(ClaimReceiptsCommandBuilder b)]) = _$ClaimReceiptsCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimReceiptsCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimReceiptsCommand> get serializer => _$ClaimReceiptsCommandSerializer();
}

class _$ClaimReceiptsCommandSerializer implements PrimitiveSerializer<ClaimReceiptsCommand> {
  @override
  final Iterable<Type> types = const [ClaimReceiptsCommand, _$ClaimReceiptsCommand];

  @override
  final String wireName = r'ClaimReceiptsCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimReceiptsCommand object, {
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
    ClaimReceiptsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimReceiptsCommandBuilder result,
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
  ClaimReceiptsCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimReceiptsCommandBuilder();
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

