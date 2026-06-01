//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/permission_scope.dart';
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
/// * [isSystem] 
/// * [permissions] 
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

  @BuiltValueField(wireName: r'isSystem')
  bool get isSystem;

  @BuiltValueField(wireName: r'permissions')
  BuiltList<Permission> get permissions;

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

