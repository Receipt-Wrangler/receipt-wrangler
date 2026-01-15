// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'associated_entity_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const AssociatedEntityType _$NOOP_ENTITY_TYPE =
    const AssociatedEntityType._('NOOP_ENTITY_TYPE');
const AssociatedEntityType _$RECEIPT = const AssociatedEntityType._('RECEIPT');
const AssociatedEntityType _$SYSTEM_EMAIL =
    const AssociatedEntityType._('SYSTEM_EMAIL');
const AssociatedEntityType _$RECEIPT_PROCESSING_SETTINGS =
    const AssociatedEntityType._('RECEIPT_PROCESSING_SETTINGS');
const AssociatedEntityType _$PROMPT = const AssociatedEntityType._('PROMPT');
const AssociatedEntityType _$API_KEY = const AssociatedEntityType._('API_KEY');

AssociatedEntityType _$valueOf(String name) {
  switch (name) {
    case 'NOOP_ENTITY_TYPE':
      return _$NOOP_ENTITY_TYPE;
    case 'RECEIPT':
      return _$RECEIPT;
    case 'SYSTEM_EMAIL':
      return _$SYSTEM_EMAIL;
    case 'RECEIPT_PROCESSING_SETTINGS':
      return _$RECEIPT_PROCESSING_SETTINGS;
    case 'PROMPT':
      return _$PROMPT;
    case 'API_KEY':
      return _$API_KEY;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<AssociatedEntityType> _$values =
    new BuiltSet<AssociatedEntityType>(const <AssociatedEntityType>[
  _$NOOP_ENTITY_TYPE,
  _$RECEIPT,
  _$SYSTEM_EMAIL,
  _$RECEIPT_PROCESSING_SETTINGS,
  _$PROMPT,
  _$API_KEY,
]);

class _$AssociatedEntityTypeMeta {
  const _$AssociatedEntityTypeMeta();
  AssociatedEntityType get NOOP_ENTITY_TYPE => _$NOOP_ENTITY_TYPE;
  AssociatedEntityType get RECEIPT => _$RECEIPT;
  AssociatedEntityType get SYSTEM_EMAIL => _$SYSTEM_EMAIL;
  AssociatedEntityType get RECEIPT_PROCESSING_SETTINGS =>
      _$RECEIPT_PROCESSING_SETTINGS;
  AssociatedEntityType get PROMPT => _$PROMPT;
  AssociatedEntityType get API_KEY => _$API_KEY;
  AssociatedEntityType valueOf(String name) => _$valueOf(name);
  BuiltSet<AssociatedEntityType> get values => _$values;
}

abstract class _$AssociatedEntityTypeMixin {
  // ignore: non_constant_identifier_names
  _$AssociatedEntityTypeMeta get AssociatedEntityType =>
      const _$AssociatedEntityTypeMeta();
}

Serializer<AssociatedEntityType> _$associatedEntityTypeSerializer =
    new _$AssociatedEntityTypeSerializer();

class _$AssociatedEntityTypeSerializer
    implements PrimitiveSerializer<AssociatedEntityType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'NOOP_ENTITY_TYPE': 'NOOP_ENTITY_TYPE',
    'RECEIPT': 'RECEIPT',
    'SYSTEM_EMAIL': 'SYSTEM_EMAIL',
    'RECEIPT_PROCESSING_SETTINGS': 'RECEIPT_PROCESSING_SETTINGS',
    'PROMPT': 'PROMPT',
    'API_KEY': 'API_KEY',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'NOOP_ENTITY_TYPE': 'NOOP_ENTITY_TYPE',
    'RECEIPT': 'RECEIPT',
    'SYSTEM_EMAIL': 'SYSTEM_EMAIL',
    'RECEIPT_PROCESSING_SETTINGS': 'RECEIPT_PROCESSING_SETTINGS',
    'PROMPT': 'PROMPT',
    'API_KEY': 'API_KEY',
  };

  @override
  final Iterable<Type> types = const <Type>[AssociatedEntityType];
  @override
  final String wireName = 'AssociatedEntityType';

  @override
  Object serialize(Serializers serializers, AssociatedEntityType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  AssociatedEntityType deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      AssociatedEntityType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
