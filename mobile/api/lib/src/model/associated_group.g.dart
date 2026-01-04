// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'associated_group.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const AssociatedGroup _$MINE = const AssociatedGroup._('MINE');
const AssociatedGroup _$ALL = const AssociatedGroup._('ALL');

AssociatedGroup _$valueOf(String name) {
  switch (name) {
    case 'MINE':
      return _$MINE;
    case 'ALL':
      return _$ALL;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<AssociatedGroup> _$values =
    new BuiltSet<AssociatedGroup>(const <AssociatedGroup>[
  _$MINE,
  _$ALL,
]);

class _$AssociatedGroupMeta {
  const _$AssociatedGroupMeta();
  AssociatedGroup get MINE => _$MINE;
  AssociatedGroup get ALL => _$ALL;
  AssociatedGroup valueOf(String name) => _$valueOf(name);
  BuiltSet<AssociatedGroup> get values => _$values;
}

abstract class _$AssociatedGroupMixin {
  // ignore: non_constant_identifier_names
  _$AssociatedGroupMeta get AssociatedGroup => const _$AssociatedGroupMeta();
}

Serializer<AssociatedGroup> _$associatedGroupSerializer =
    new _$AssociatedGroupSerializer();

class _$AssociatedGroupSerializer
    implements PrimitiveSerializer<AssociatedGroup> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'MINE': 'MINE',
    'ALL': 'ALL',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'MINE': 'MINE',
    'ALL': 'ALL',
  };

  @override
  final Iterable<Type> types = const <Type>[AssociatedGroup];
  @override
  final String wireName = 'AssociatedGroup';

  @override
  Object serialize(Serializers serializers, AssociatedGroup object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  AssociatedGroup deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      AssociatedGroup.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
