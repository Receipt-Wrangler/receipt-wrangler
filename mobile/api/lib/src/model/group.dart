//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/group_receipt_settings.dart';
import 'package:openapi/src/model/group_settings.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/group_member.dart';
import 'package:openapi/src/model/group_status.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group.g.dart';

/// Group in the system
///
/// Properties:
/// * [createdAt] 
/// * [createdBy] 
/// * [groupSettings] 
/// * [groupReceiptSettings] 
/// * [groupMembers] - Members of the group
/// * [id] 
/// * [isDefault] - Is default group (not used yet)
/// * [name] - Name of the group
/// * [isAllGroup] - Is all group for user
/// * [status] 
/// * [updatedAt] 
@BuiltValue()
abstract class Group implements Built<Group, GroupBuilder> {
  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  @BuiltValueField(wireName: r'groupSettings')
  GroupSettings? get groupSettings;

  @BuiltValueField(wireName: r'groupReceiptSettings')
  GroupReceiptSettings get groupReceiptSettings;

  /// Members of the group
  @BuiltValueField(wireName: r'groupMembers')
  BuiltList<GroupMember> get groupMembers;

  @BuiltValueField(wireName: r'id')
  int get id;

  /// Is default group (not used yet)
  @BuiltValueField(wireName: r'isDefault')
  bool? get isDefault;

  /// Name of the group
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Is all group for user
  @BuiltValueField(wireName: r'isAllGroup')
  bool get isAllGroup;

  @BuiltValueField(wireName: r'status')
  GroupStatus get status;
  // enum statusEnum {  ACTIVE,  ARCHIVED,  };

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  Group._();

  factory Group([void updates(GroupBuilder b)]) = _$Group;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GroupBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Group> get serializer => _$GroupSerializer();
}

class _$GroupSerializer implements PrimitiveSerializer<Group> {
  @override
  final Iterable<Type> types = const [Group, _$Group];

  @override
  final String wireName = r'Group';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Group object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    if (object.groupSettings != null) {
      yield r'groupSettings';
      yield serializers.serialize(
        object.groupSettings,
        specifiedType: const FullType(GroupSettings),
      );
    }
    yield r'groupReceiptSettings';
    yield serializers.serialize(
      object.groupReceiptSettings,
      specifiedType: const FullType(GroupReceiptSettings),
    );
    yield r'groupMembers';
    yield serializers.serialize(
      object.groupMembers,
      specifiedType: const FullType(BuiltList, [FullType(GroupMember)]),
    );
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.isDefault != null) {
      yield r'isDefault';
      yield serializers.serialize(
        object.isDefault,
        specifiedType: const FullType(bool),
      );
    }
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'isAllGroup';
    yield serializers.serialize(
      object.isAllGroup,
      specifiedType: const FullType(bool),
    );
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(GroupStatus),
    );
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    Group object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GroupBuilder result,
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
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'groupSettings':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(GroupSettings),
          ) as GroupSettings;
          result.groupSettings.replace(valueDes);
          break;
        case r'groupReceiptSettings':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(GroupReceiptSettings),
          ) as GroupReceiptSettings;
          result.groupReceiptSettings.replace(valueDes);
          break;
        case r'groupMembers':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(GroupMember)]),
          ) as BuiltList<GroupMember>;
          result.groupMembers.replace(valueDes);
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'isDefault':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isDefault = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'isAllGroup':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isAllGroup = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(GroupStatus),
          ) as GroupStatus;
          result.status = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Group deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GroupBuilder();
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

