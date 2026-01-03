// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group_role.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const GroupRole _$OWNER = const GroupRole._('OWNER');
const GroupRole _$VIEWER = const GroupRole._('VIEWER');
const GroupRole _$EDITOR = const GroupRole._('EDITOR');

GroupRole _$valueOf(String name) {
  switch (name) {
    case 'OWNER':
      return _$OWNER;
    case 'VIEWER':
      return _$VIEWER;
    case 'EDITOR':
      return _$EDITOR;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<GroupRole> _$values = new BuiltSet<GroupRole>(const <GroupRole>[
  _$OWNER,
  _$VIEWER,
  _$EDITOR,
]);

class _$GroupRoleMeta {
  const _$GroupRoleMeta();
  GroupRole get OWNER => _$OWNER;
  GroupRole get VIEWER => _$VIEWER;
  GroupRole get EDITOR => _$EDITOR;
  GroupRole valueOf(String name) => _$valueOf(name);
  BuiltSet<GroupRole> get values => _$values;
}

abstract class _$GroupRoleMixin {
  // ignore: non_constant_identifier_names
  _$GroupRoleMeta get GroupRole => const _$GroupRoleMeta();
}

Serializer<GroupRole> _$groupRoleSerializer = new _$GroupRoleSerializer();

class _$GroupRoleSerializer implements PrimitiveSerializer<GroupRole> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'OWNER': 'OWNER',
    'VIEWER': 'VIEWER',
    'EDITOR': 'EDITOR',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'OWNER': 'OWNER',
    'VIEWER': 'VIEWER',
    'EDITOR': 'EDITOR',
  };

  @override
  final Iterable<Type> types = const <Type>[GroupRole];
  @override
  final String wireName = 'GroupRole';

  @override
  Object serialize(Serializers serializers, GroupRole object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  GroupRole deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      GroupRole.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
