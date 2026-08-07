//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group_member.g.dart';

/// Group member
///
/// Properties:
/// * [createdAt] 
/// * [groupId] - Group compound primary key
/// * [groupRoleId] - Id of the modern group role assigned to the member
/// * [updatedAt] 
/// * [userId] - User compound primary key
/// * [categoryGrants] - Category ids this individual member may see, narrowing WITHIN whatever their group role allows (the two layers intersect). Empty means the member adds no narrowing of their own: they fall back to the role's set, EXCEPT when the role sets requiresIndividualCategoryGrants, in which case an unassigned member sees no categories at all (fail closed). Read only here — write via PUT /group/{groupId}/member/{userId}/grants.
/// * [tagGrants] - Tag counterpart of categoryGrants. Restricted independently of categories.
@BuiltValue()
abstract class GroupMember implements Built<GroupMember, GroupMemberBuilder> {
  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  /// Group compound primary key
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  /// Id of the modern group role assigned to the member
  @BuiltValueField(wireName: r'groupRoleId')
  int? get groupRoleId;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// User compound primary key
  @BuiltValueField(wireName: r'userId')
  int get userId;

  /// Category ids this individual member may see, narrowing WITHIN whatever their group role allows (the two layers intersect). Empty means the member adds no narrowing of their own: they fall back to the role's set, EXCEPT when the role sets requiresIndividualCategoryGrants, in which case an unassigned member sees no categories at all (fail closed). Read only here — write via PUT /group/{groupId}/member/{userId}/grants.
  @BuiltValueField(wireName: r'categoryGrants')
  BuiltList<int>? get categoryGrants;

  /// Tag counterpart of categoryGrants. Restricted independently of categories.
  @BuiltValueField(wireName: r'tagGrants')
  BuiltList<int>? get tagGrants;

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
    if (object.groupRoleId != null) {
      yield r'groupRoleId';
      yield serializers.serialize(
        object.groupRoleId,
        specifiedType: const FullType(int),
      );
    }
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
        case r'groupRoleId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupRoleId = valueDes;
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

