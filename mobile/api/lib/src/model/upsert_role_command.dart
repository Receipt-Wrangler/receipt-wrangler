//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/permission_scope.dart';
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

