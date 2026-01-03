//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/group_role.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group_member.g.dart';

/// Group member
///
/// Properties:
/// * [createdAt] 
/// * [groupId] - Group compound primary key
/// * [groupRole] 
/// * [updatedAt] 
/// * [userId] - User compound primary key
@BuiltValue()
abstract class GroupMember implements Built<GroupMember, GroupMemberBuilder> {
  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  /// Group compound primary key
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  @BuiltValueField(wireName: r'groupRole')
  GroupRole get groupRole;
  // enum groupRoleEnum {  OWNER,  VIEWER,  EDITOR,  };

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// User compound primary key
  @BuiltValueField(wireName: r'userId')
  int get userId;

  GroupMember._();

  factory GroupMember([void updates(GroupMemberBuilder b)]) = _$GroupMember;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GroupMemberBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<GroupMember> get serializer => _$GroupMemberSerializer();
}

class _$GroupMemberSerializer implements PrimitiveSerializer<GroupMember> {
  @override
  final Iterable<Type> types = const [GroupMember, _$GroupMember];

  @override
  final String wireName = r'GroupMember';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GroupMember object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(String),
      );
    }
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
    yield r'groupRole';
    yield serializers.serialize(
      object.groupRole,
      specifiedType: const FullType(GroupRole),
    );
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
    yield r'userId';
    yield serializers.serialize(
      object.userId,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    GroupMember object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GroupMemberBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'groupRole':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(GroupRole),
          ) as GroupRole;
          result.groupRole = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
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
  GroupMember deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GroupMemberBuilder();
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

