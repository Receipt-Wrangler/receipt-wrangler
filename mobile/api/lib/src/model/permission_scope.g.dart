// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'permission_scope.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const PermissionScope _$APP = const PermissionScope._('APP');
const PermissionScope _$GROUP = const PermissionScope._('GROUP');

PermissionScope _$valueOf(String name) {
  switch (name) {
    case 'APP':
      return _$APP;
    case 'GROUP':
      return _$GROUP;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<PermissionScope> _$values =
    BuiltSet<PermissionScope>(const <PermissionScope>[
  _$APP,
  _$GROUP,
]);

class _$PermissionScopeMeta {
  const _$PermissionScopeMeta();
  PermissionScope get APP => _$APP;
  PermissionScope get GROUP => _$GROUP;
  PermissionScope valueOf(String name) => _$valueOf(name);
  BuiltSet<PermissionScope> get values => _$values;
}

abstract class _$PermissionScopeMixin {
  // ignore: non_constant_identifier_names
  _$PermissionScopeMeta get PermissionScope => const _$PermissionScopeMeta();
}

Serializer<PermissionScope> _$permissionScopeSerializer =
    _$PermissionScopeSerializer();

class _$PermissionScopeSerializer
    implements PrimitiveSerializer<PermissionScope> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'APP': 'APP',
    'GROUP': 'GROUP',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'APP': 'APP',
    'GROUP': 'GROUP',
  };

  @override
  final Iterable<Type> types = const <Type>[PermissionScope];
  @override
  final String wireName = 'PermissionScope';

  @override
  Object serialize(Serializers serializers, PermissionScope object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  PermissionScope deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      PermissionScope.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
