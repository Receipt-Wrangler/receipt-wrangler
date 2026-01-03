//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bulk_status_update_command.g.dart';

/// BulkStatusUpdateCommand
///
/// Properties:
/// * [comment] - Optional comment to leave on each receipt
/// * [status] - Status to update to
/// * [receiptIds] - Receipt ids to update
@BuiltValue()
abstract class BulkStatusUpdateCommand implements Built<BulkStatusUpdateCommand, BulkStatusUpdateCommandBuilder> {
  /// Optional comment to leave on each receipt
  @BuiltValueField(wireName: r'comment')
  String? get comment;

  /// Status to update to
  @BuiltValueField(wireName: r'status')
  String get status;

  /// Receipt ids to update
  @BuiltValueField(wireName: r'receiptIds')
  BuiltList<int> get receiptIds;

  BulkStatusUpdateCommand._();

  factory BulkStatusUpdateCommand([void updates(BulkStatusUpdateCommandBuilder b)]) = _$BulkStatusUpdateCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BulkStatusUpdateCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BulkStatusUpdateCommand> get serializer => _$BulkStatusUpdateCommandSerializer();
}

class _$BulkStatusUpdateCommandSerializer implements PrimitiveSerializer<BulkStatusUpdateCommand> {
  @override
  final Iterable<Type> types = const [BulkStatusUpdateCommand, _$BulkStatusUpdateCommand];

  @override
  final String wireName = r'BulkStatusUpdateCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BulkStatusUpdateCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.comment != null) {
      yield r'comment';
      yield serializers.serialize(
        object.comment,
        specifiedType: const FullType(String),
      );
    }
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(String),
    );
    yield r'receiptIds';
    yield serializers.serialize(
      object.receiptIds,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BulkStatusUpdateCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BulkStatusUpdateCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'comment':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.comment = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.status = valueDes;
          break;
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
  BulkStatusUpdateCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BulkStatusUpdateCommandBuilder();
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

