//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'update_group_member_grants_command.g.dart';

/// Per-member category/tag assignment. The membership being edited is taken from the URL, never this body.
///
/// Properties:
/// * [categoryGrants] - Category ids to assign to this member. Every id must sit within the ceiling set by the member's group role, or the request is rejected with 400. An empty array clears the member's category restriction, handing them back to their role's set.
/// * [tagGrants] - Tag counterpart of categoryGrants.
@BuiltValue()
abstract class UpdateGroupMemberGrantsCommand implements Built<UpdateGroupMemberGrantsCommand, UpdateGroupMemberGrantsCommandBuilder> {
  /// Category ids to assign to this member. Every id must sit within the ceiling set by the member's group role, or the request is rejected with 400. An empty array clears the member's category restriction, handing them back to their role's set.
  @BuiltValueField(wireName: r'categoryGrants')
  BuiltList<int>? get categoryGrants;

  /// Tag counterpart of categoryGrants.
  @BuiltValueField(wireName: r'tagGrants')
  BuiltList<int>? get tagGrants;

  UpdateGroupMemberGrantsCommand._();

  factory UpdateGroupMemberGrantsCommand([void updates(UpdateGroupMemberGrantsCommandBuilder b)]) = _$UpdateGroupMemberGrantsCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpdateGroupMemberGrantsCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpdateGroupMemberGrantsCommand> get serializer => _$UpdateGroupMemberGrantsCommandSerializer();
}

class _$UpdateGroupMemberGrantsCommandSerializer implements PrimitiveSerializer<UpdateGroupMemberGrantsCommand> {
  @override
  final Iterable<Type> types = const [UpdateGroupMemberGrantsCommand, _$UpdateGroupMemberGrantsCommand];

  @override
  final String wireName = r'UpdateGroupMemberGrantsCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpdateGroupMemberGrantsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.categoryGrants != null) {
      yield r'categoryGrants';
      yield serializers.serialize(
        object.categoryGrants,
        specifiedType: const FullType(BuiltList, [FullType(int)]),
      );
    }
    if (object.tagGrants != null) {
      yield r'tagGrants';
      yield serializers.serialize(
        object.tagGrants,
        specifiedType: const FullType(BuiltList, [FullType(int)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpdateGroupMemberGrantsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpdateGroupMemberGrantsCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'categoryGrants':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.categoryGrants.replace(valueDes);
          break;
        case r'tagGrants':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.tagGrants.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpdateGroupMemberGrantsCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpdateGroupMemberGrantsCommandBuilder();
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

