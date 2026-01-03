// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group_status.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const GroupStatus _$ACTIVE = const GroupStatus._('ACTIVE');
const GroupStatus _$ARCHIVED = const GroupStatus._('ARCHIVED');

GroupStatus _$valueOf(String name) {
  switch (name) {
    case 'ACTIVE':
      return _$ACTIVE;
    case 'ARCHIVED':
      return _$ARCHIVED;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<GroupStatus> _$values =
    new BuiltSet<GroupStatus>(const <GroupStatus>[
  _$ACTIVE,
  _$ARCHIVED,
]);

class _$GroupStatusMeta {
  const _$GroupStatusMeta();
  GroupStatus get ACTIVE => _$ACTIVE;
  GroupStatus get ARCHIVED => _$ARCHIVED;
  GroupStatus valueOf(String name) => _$valueOf(name);
  BuiltSet<GroupStatus> get values => _$values;
}

abstract class _$GroupStatusMixin {
  // ignore: non_constant_identifier_names
  _$GroupStatusMeta get GroupStatus => const _$GroupStatusMeta();
}

Serializer<GroupStatus> _$groupStatusSerializer = new _$GroupStatusSerializer();

class _$GroupStatusSerializer implements PrimitiveSerializer<GroupStatus> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'ACTIVE': 'ACTIVE',
    'ARCHIVED': 'ARCHIVED',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'ACTIVE': 'ACTIVE',
    'ARCHIVED': 'ARCHIVED',
  };

  @override
  final Iterable<Type> types = const <Type>[GroupStatus];
  @override
  final String wireName = 'GroupStatus';

  @override
  Object serialize(Serializers serializers, GroupStatus object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  GroupStatus deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      GroupStatus.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
