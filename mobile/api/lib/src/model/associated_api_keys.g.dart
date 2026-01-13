// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'associated_api_keys.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const AssociatedApiKeys _$MINE = const AssociatedApiKeys._('MINE');
const AssociatedApiKeys _$ALL = const AssociatedApiKeys._('ALL');

AssociatedApiKeys _$valueOf(String name) {
  switch (name) {
    case 'MINE':
      return _$MINE;
    case 'ALL':
      return _$ALL;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<AssociatedApiKeys> _$values =
    new BuiltSet<AssociatedApiKeys>(const <AssociatedApiKeys>[
  _$MINE,
  _$ALL,
]);

class _$AssociatedApiKeysMeta {
  const _$AssociatedApiKeysMeta();
  AssociatedApiKeys get MINE => _$MINE;
  AssociatedApiKeys get ALL => _$ALL;
  AssociatedApiKeys valueOf(String name) => _$valueOf(name);
  BuiltSet<AssociatedApiKeys> get values => _$values;
}

abstract class _$AssociatedApiKeysMixin {
  // ignore: non_constant_identifier_names
  _$AssociatedApiKeysMeta get AssociatedApiKeys =>
      const _$AssociatedApiKeysMeta();
}

Serializer<AssociatedApiKeys> _$associatedApiKeysSerializer =
    new _$AssociatedApiKeysSerializer();

class _$AssociatedApiKeysSerializer
    implements PrimitiveSerializer<AssociatedApiKeys> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'MINE': 'MINE',
    'ALL': 'ALL',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'MINE': 'MINE',
    'ALL': 'ALL',
  };

  @override
  final Iterable<Type> types = const <Type>[AssociatedApiKeys];
  @override
  final String wireName = 'AssociatedApiKeys';

  @override
  Object serialize(Serializers serializers, AssociatedApiKeys object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  AssociatedApiKeys deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      AssociatedApiKeys.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
