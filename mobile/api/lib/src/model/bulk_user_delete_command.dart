//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bulk_user_delete_command.g.dart';

/// BulkUserDeleteCommand
///
/// Properties:
/// * [userIds] - User IDs to delete
@BuiltValue()
abstract class BulkUserDeleteCommand implements Built<BulkUserDeleteCommand, BulkUserDeleteCommandBuilder> {
  /// User IDs to delete
  @BuiltValueField(wireName: r'userIds')
  BuiltList<String> get userIds;

  BulkUserDeleteCommand._();

  factory BulkUserDeleteCommand([void updates(BulkUserDeleteCommandBuilder b)]) = _$BulkUserDeleteCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BulkUserDeleteCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BulkUserDeleteCommand> get serializer => _$BulkUserDeleteCommandSerializer();
}

class _$BulkUserDeleteCommandSerializer implements PrimitiveSerializer<BulkUserDeleteCommand> {
  @override
  final Iterable<Type> types = const [BulkUserDeleteCommand, _$BulkUserDeleteCommand];

  @override
  final String wireName = r'BulkUserDeleteCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BulkUserDeleteCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'userIds';
    yield serializers.serialize(
      object.userIds,
      specifiedType: const FullType(BuiltList, [FullType(String)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BulkUserDeleteCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BulkUserDeleteCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'userIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(String)]),
          ) as BuiltList<String>;
          result.userIds.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BulkUserDeleteCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BulkUserDeleteCommandBuilder();
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

