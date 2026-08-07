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

part 'role.g.dart';

/// Role
///
/// Properties:
/// * [id] 
/// * [name] 
/// * [description] 
/// * [scope] 
/// * [isDefault] - Whether this role is the default for its scope — assigned to new accounts (APP) or to group creators (GROUP). Exactly one role per scope is the default.
/// * [isSystem] 
/// * [permissions] 
/// * [assignedCount] - Number of users or group members currently assigned this role
/// * [categoryGrants] - Category ids a GROUP role restricts its members to. Empty means unrestricted (members may use every category). Always empty for app roles.
/// * [tagGrants] - Tag ids a GROUP role restricts its members to. Empty means unrestricted (members may use every tag). Always empty for app roles.
/// * [paidByUserGrants] - User ids whose receipts a GROUP role lets its members see (by the receipt's \"paid by\" user). Empty with includeOwnPaidReceipts false means unrestricted (members see every payer's receipts). Always empty for app roles.
/// * [includeOwnPaidReceipts] - Whether a GROUP role lets each member see receipts they paid for. Part of the paid-by visibility filter; always false for app roles.
/// * [seesAllMembers] - Whether a GROUP role exempts its members from member-presence isolation (they see, and are seen by, every member of an isolated group). Always false for app roles.
/// * [skipDefaultGroupCreation] - Whether users created with this APP role skip the automatic personal \"My Receipts\" group. The virtual \"All\" group is always created. Applies at user-creation time only — changing it never adds or removes a group for an existing user. Always false for group roles.
/// * [requiresIndividualCategoryGrants] - Whether a GROUP role requires per-member category assignment. When true, a member holding this role with no individual category grants sees NO categories, rather than falling back to the role's set. Always false for app roles.
/// * [requiresIndividualTagGrants] - Tag counterpart of requiresIndividualCategoryGrants. Always false for app roles.
/// * [reportTemplateGrants] - Per-template action grants restricting which report templates a GROUP role's members may act on. Empty means unrestricted (every template the role's group access reaches). Always empty for app roles.
@BuiltValue()
abstract class Role implements Built<Role, RoleBuilder> {
  @BuiltValueField(wireName: r'id')
  int get id;

  @BuiltValueField(wireName: r'name')
  String get name;

  @BuiltValueField(wireName: r'description')
  String? get description;

  @BuiltValueField(wireName: r'scope')
  PermissionScope get scope;
  // enum scopeEnum {  APP,  GROUP,  };

  /// Whether this role is the default for its scope — assigned to new accounts (APP) or to group creators (GROUP). Exactly one role per scope is the default.
  @BuiltValueField(wireName: r'isDefault')
  bool get isDefault;

  @BuiltValueField(wireName: r'isSystem')
  bool get isSystem;

  @BuiltValueField(wireName: r'permissions')
  BuiltList<Permission> get permissions;

  /// Number of users or group members currently assigned this role
  @BuiltValueField(wireName: r'assignedCount')
  int? get assignedCount;

  /// Category ids a GROUP role restricts its members to. Empty means unrestricted (members may use every category). Always empty for app roles.
  @BuiltValueField(wireName: r'categoryGrants')
  BuiltList<int>? get categoryGrants;

  /// Tag ids a GROUP role restricts its members to. Empty means unrestricted (members may use every tag). Always empty for app roles.
  @BuiltValueField(wireName: r'tagGrants')
  BuiltList<int>? get tagGrants;

  /// User ids whose receipts a GROUP role lets its members see (by the receipt's \"paid by\" user). Empty with includeOwnPaidReceipts false means unrestricted (members see every payer's receipts). Always empty for app roles.
  @BuiltValueField(wireName: r'paidByUserGrants')
  BuiltList<int>? get paidByUserGrants;

  /// Whether a GROUP role lets each member see receipts they paid for. Part of the paid-by visibility filter; always false for app roles.
  @BuiltValueField(wireName: r'includeOwnPaidReceipts')
  bool? get includeOwnPaidReceipts;

  /// Whether a GROUP role exempts its members from member-presence isolation (they see, and are seen by, every member of an isolated group). Always false for app roles.
  @BuiltValueField(wireName: r'seesAllMembers')
  bool? get seesAllMembers;

  /// Whether users created with this APP role skip the automatic personal \"My Receipts\" group. The virtual \"All\" group is always created. Applies at user-creation time only — changing it never adds or removes a group for an existing user. Always false for group roles.
  @BuiltValueField(wireName: r'skipDefaultGroupCreation')
  bool? get skipDefaultGroupCreation;

  /// Whether a GROUP role requires per-member category assignment. When true, a member holding this role with no individual category grants sees NO categories, rather than falling back to the role's set. Always false for app roles.
  @BuiltValueField(wireName: r'requiresIndividualCategoryGrants')
  bool? get requiresIndividualCategoryGrants;

  /// Tag counterpart of requiresIndividualCategoryGrants. Always false for app roles.
  @BuiltValueField(wireName: r'requiresIndividualTagGrants')
  bool? get requiresIndividualTagGrants;

  /// Per-template action grants restricting which report templates a GROUP role's members may act on. Empty means unrestricted (every template the role's group access reaches). Always empty for app roles.
  @BuiltValueField(wireName: r'reportTemplateGrants')
  BuiltList<ReportTemplateGrant>? get reportTemplateGrants;

  Role._();

  factory Role([void updates(RoleBuilder b)]) = _$Role;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(RoleBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Role> get serializer => _$RoleSerializer();
}

class _$RoleSerializer implements PrimitiveSerializer<Role> {
  @override
  final Iterable<Type> types = const [Role, _$Role];

  @override
  final String wireName = r'Role';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Role object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
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
    yield r'isDefault';
    yield serializers.serialize(
      object.isDefault,
      specifiedType: const FullType(bool),
    );
    yield r'isSystem';
    yield serializers.serialize(
      object.isSystem,
      specifiedType: const FullType(bool),
    );
    yield r'permissions';
    yield serializers.serialize(
      object.permissions,
      specifiedType: const FullType(BuiltList, [FullType(Permission)]),
    );
    if (object.assignedCount != null) {
      yield r'assignedCount';
      yield serializers.serialize(
        object.assignedCount,
        specifiedType: const FullType(int),
      );
    }
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
    Role object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required RoleBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
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
        case r'isDefault':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isDefault = valueDes;
          break;
        case r'isSystem':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isSystem = valueDes;
          break;
        case r'permissions':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Permission)]),
          ) as BuiltList<Permission>;
          result.permissions.replace(valueDes);
          break;
        case r'assignedCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.assignedCount = valueDes;
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
  Role deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = RoleBuilder();
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

