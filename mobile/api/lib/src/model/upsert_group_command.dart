//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/upsert_group_member_command.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/group_status.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_group_command.g.dart';

/// UpsertGroupCommand
///
/// Properties:
/// * [groupMembers] - Members of the group
/// * [isDefault] - Is default group (not used yet)
/// * [name] - Name of the group
/// * [isAllGroup] - Is all group for user
/// * [status] 
@BuiltValue()
abstract class UpsertGroupCommand implements Built<UpsertGroupCommand, UpsertGroupCommandBuilder> {
  /// Members of the group
  @BuiltValueField(wireName: r'groupMembers')
  BuiltList<UpsertGroupMemberCommand> get groupMembers;

  /// Is default group (not used yet)
  @BuiltValueField(wireName: r'isDefault')
  bool? get isDefault;

  /// Name of the group
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Is all group for user
  @BuiltValueField(wireName: r'isAllGroup')
  bool? get isAllGroup;

  @BuiltValueField(wireName: r'status')
  GroupStatus get status;
  // enum statusEnum {  ACTIVE,  ARCHIVED,  };

  UpsertGroupCommand._();

  factory UpsertGroupCommand([void updates(UpsertGroupCommandBuilder b)]) = _$UpsertGroupCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertGroupCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertGroupCommand> get serializer => _$UpsertGroupCommandSerializer();
}

class _$UpsertGroupCommandSerializer implements PrimitiveSerializer<UpsertGroupCommand> {
  @override
  final Iterable<Type> types = const [UpsertGroupCommand, _$UpsertGroupCommand];

  @override
  final String wireName = r'UpsertGroupCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertGroupCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'groupMembers';
    yield serializers.serialize(
      object.groupMembers,
      specifiedType: const FullType(BuiltList, [FullType(UpsertGroupMemberCommand)]),
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
    if (object.isAllGroup != null) {
      yield r'isAllGroup';
      yield serializers.serialize(
        object.isAllGroup,
        specifiedType: const FullType(bool),
      );
    }
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(GroupStatus),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertGroupCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertGroupCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'groupMembers':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(UpsertGroupMemberCommand)]),
          ) as BuiltList<UpsertGroupMemberCommand>;
          result.groupMembers.replace(valueDes);
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
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertGroupCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertGroupCommandBuilder();
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

