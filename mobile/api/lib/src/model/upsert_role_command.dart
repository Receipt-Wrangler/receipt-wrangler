//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/permission_scope.dart';
import 'package:openapi/src/model/report_template_grant.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/permission.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_role_command.g.dart';

/// UpsertRoleCommand
///
/// Properties:
/// * [name] 
/// * [description] 
/// * [scope] 
/// * [permissions] 
/// * [categoryGrants] - Category ids to restrict a GROUP role's members to. Only valid on group roles; omit or leave empty for unrestricted access.
/// * [tagGrants] - Tag ids to restrict a GROUP role's members to. Only valid on group roles; omit or leave empty for unrestricted access.
/// * [paidByUserGrants] - User ids whose receipts a GROUP role's members may see (by the receipt's \"paid by\" user). Only valid on group roles; omit or leave empty (with includeOwnPaidReceipts false) for unrestricted access.
/// * [includeOwnPaidReceipts] - Whether to also let each member see receipts they paid for. Only valid on group roles.
/// * [seesAllMembers] - Whether a GROUP role exempts its members from member-presence isolation (they see, and are seen by, every member of an isolated group). Only valid on group roles.
/// * [skipDefaultGroupCreation] - Whether users created with this APP role skip the automatic personal \"My Receipts\" group (the virtual \"All\" group is always created). Only valid on app roles; applies at user-creation time only.
/// * [requiresIndividualCategoryGrants] - Whether this GROUP role requires per-member category assignment. When true, a member with no individual category grants sees NO categories instead of falling back to the role's set, so an unassigned member fails closed. Only valid on group roles.
/// * [requiresIndividualTagGrants] - Tag counterpart of requiresIndividualCategoryGrants. Only valid on group roles.
/// * [reportTemplateGrants] - Per-template action grants for a GROUP role, restricting which report templates its members may act on. Only valid on group roles; omit or leave empty for unrestricted access.
@BuiltValue()
abstract class UpsertRoleCommand implements Built<UpsertRoleCommand, UpsertRoleCommandBuilder> {
  @BuiltValueField(wireName: r'name')
  String get name;

  @BuiltValueField(wireName: r'description')
  String? get description;

  @BuiltValueField(wireName: r'scope')
  PermissionScope get scope;
  // enum scopeEnum {  APP,  GROUP,  };

  @BuiltValueField(wireName: r'permissions')
  BuiltList<Permission> get permissions;

  /// Category ids to restrict a GROUP role's members to. Only valid on group roles; omit or leave empty for unrestricted access.
  @BuiltValueField(wireName: r'categoryGrants')
  BuiltList<int>? get categoryGrants;

  /// Tag ids to restrict a GROUP role's members to. Only valid on group roles; omit or leave empty for unrestricted access.
  @BuiltValueField(wireName: r'tagGrants')
  BuiltList<int>? get tagGrants;

  /// User ids whose receipts a GROUP role's members may see (by the receipt's \"paid by\" user). Only valid on group roles; omit or leave empty (with includeOwnPaidReceipts false) for unrestricted access.
  @BuiltValueField(wireName: r'paidByUserGrants')
  BuiltList<int>? get paidByUserGrants;

  /// Whether to also let each member see receipts they paid for. Only valid on group roles.
  @BuiltValueField(wireName: r'includeOwnPaidReceipts')
  bool? get includeOwnPaidReceipts;

  /// Whether a GROUP role exempts its members from member-presence isolation (they see, and are seen by, every member of an isolated group). Only valid on group roles.
  @BuiltValueField(wireName: r'seesAllMembers')
  bool? get seesAllMembers;

  /// Whether users created with this APP role skip the automatic personal \"My Receipts\" group (the virtual \"All\" group is always created). Only valid on app roles; applies at user-creation time only.
  @BuiltValueField(wireName: r'skipDefaultGroupCreation')
  bool? get skipDefaultGroupCreation;

  /// Whether this GROUP role requires per-member category assignment. When true, a member with no individual category grants sees NO categories instead of falling back to the role's set, so an unassigned member fails closed. Only valid on group roles.
  @BuiltValueField(wireName: r'requiresIndividualCategoryGrants')
  bool? get requiresIndividualCategoryGrants;

  /// Tag counterpart of requiresIndividualCategoryGrants. Only valid on group roles.
  @BuiltValueField(wireName: r'requiresIndividualTagGrants')
  bool? get requiresIndividualTagGrants;

  /// Per-template action grants for a GROUP role, restricting which report templates its members may act on. Only valid on group roles; omit or leave empty for unrestricted access.
  @BuiltValueField(wireName: r'reportTemplateGrants')
  BuiltList<ReportTemplateGrant>? get reportTemplateGrants;

  UpsertRoleCommand._();

  factory UpsertRoleCommand([void updates(UpsertRoleCommandBuilder b)]) = _$UpsertRoleCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertRoleCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertRoleCommand> get serializer => _$UpsertRoleCommandSerializer();
}

class _$UpsertRoleCommandSerializer implements PrimitiveSerializer<UpsertRoleCommand> {
  @override
  final Iterable<Type> types = const [UpsertRoleCommand, _$UpsertRoleCommand];

  @override
  final String wireName = r'UpsertRoleCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertRoleCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.description != null) {
      yield r'description';
      yield serializers.serialize(
        object.description,
        specifiedType: const FullType(String),
      );
    }
    yield r'scope';
    yield serializers.serialize(
      object.scope,
      specifiedType: const FullType(PermissionScope),
    );
    yield r'permissions';
    yield serializers.serialize(
      object.permissions,
      specifiedType: const FullType(BuiltList, [FullType(Permission)]),
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
    if (object.paidByUserGrants != null) {
      yield r'paidByUserGrants';
      yield serializers.serialize(
        object.paidByUserGrants,
        specifiedType: const FullType(BuiltList, [FullType(int)]),
      );
    }
    if (object.includeOwnPaidReceipts != null) {
      yield r'includeOwnPaidReceipts';
      yield serializers.serialize(
        object.includeOwnPaidReceipts,
        specifiedType: const FullType(bool),
      );
    }
    if (object.seesAllMembers != null) {
      yield r'seesAllMembers';
      yield serializers.serialize(
        object.seesAllMembers,
        specifiedType: const FullType(bool),
      );
    }
    if (object.skipDefaultGroupCreation != null) {
      yield r'skipDefaultGroupCreation';
      yield serializers.serialize(
        object.skipDefaultGroupCreation,
        specifiedType: const FullType(bool),
      );
    }
    if (object.requiresIndividualCategoryGrants != null) {
      yield r'requiresIndividualCategoryGrants';
      yield serializers.serialize(
        object.requiresIndividualCategoryGrants,
        specifiedType: const FullType(bool),
      );
    }
    if (object.requiresIndividualTagGrants != null) {
      yield r'requiresIndividualTagGrants';
      yield serializers.serialize(
        object.requiresIndividualTagGrants,
        specifiedType: const FullType(bool),
      );
    }
    if (object.reportTemplateGrants != null) {
      yield r'reportTemplateGrants';
      yield serializers.serialize(
        object.reportTemplateGrants,
        specifiedType: const FullType(BuiltList, [FullType(ReportTemplateGrant)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertRoleCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertRoleCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'description':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.description = valueDes;
          break;
        case r'scope':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(PermissionScope),
          ) as PermissionScope;
          result.scope = valueDes;
          break;
        case r'permissions':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Permission)]),
          ) as BuiltList<Permission>;
          result.permissions.replace(valueDes);
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
        case r'paidByUserGrants':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.paidByUserGrants.replace(valueDes);
          break;
        case r'includeOwnPaidReceipts':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.includeOwnPaidReceipts = valueDes;
          break;
        case r'seesAllMembers':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.seesAllMembers = valueDes;
          break;
        case r'skipDefaultGroupCreation':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.skipDefaultGroupCreation = valueDes;
          break;
        case r'requiresIndividualCategoryGrants':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.requiresIndividualCategoryGrants = valueDes;
          break;
        case r'requiresIndividualTagGrants':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.requiresIndividualTagGrants = valueDes;
          break;
        case r'reportTemplateGrants':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReportTemplateGrant)]),
          ) as BuiltList<ReportTemplateGrant>;
          result.reportTemplateGrants.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertRoleCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertRoleCommandBuilder();
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

