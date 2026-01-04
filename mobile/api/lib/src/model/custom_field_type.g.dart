// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'custom_field_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const CustomFieldType _$TEXT = const CustomFieldType._('TEXT');
const CustomFieldType _$DATE = const CustomFieldType._('DATE');
const CustomFieldType _$SELECT = const CustomFieldType._('SELECT');
const CustomFieldType _$CURRENCY = const CustomFieldType._('CURRENCY');
const CustomFieldType _$BOOLEAN = const CustomFieldType._('BOOLEAN');

CustomFieldType _$valueOf(String name) {
  switch (name) {
    case 'TEXT':
      return _$TEXT;
    case 'DATE':
      return _$DATE;
    case 'SELECT':
      return _$SELECT;
    case 'CURRENCY':
      return _$CURRENCY;
    case 'BOOLEAN':
      return _$BOOLEAN;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<CustomFieldType> _$values =
    new BuiltSet<CustomFieldType>(const <CustomFieldType>[
  _$TEXT,
  _$DATE,
  _$SELECT,
  _$CURRENCY,
  _$BOOLEAN,
]);

class _$CustomFieldTypeMeta {
  const _$CustomFieldTypeMeta();
  CustomFieldType get TEXT => _$TEXT;
  CustomFieldType get DATE => _$DATE;
  CustomFieldType get SELECT => _$SELECT;
  CustomFieldType get CURRENCY => _$CURRENCY;
  CustomFieldType get BOOLEAN => _$BOOLEAN;
  CustomFieldType valueOf(String name) => _$valueOf(name);
  BuiltSet<CustomFieldType> get values => _$values;
}

abstract class _$CustomFieldTypeMixin {
  // ignore: non_constant_identifier_names
  _$CustomFieldTypeMeta get CustomFieldType => const _$CustomFieldTypeMeta();
}

Serializer<CustomFieldType> _$customFieldTypeSerializer =
    new _$CustomFieldTypeSerializer();

class _$CustomFieldTypeSerializer
    implements PrimitiveSerializer<CustomFieldType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'TEXT': 'TEXT',
    'DATE': 'DATE',
    'SELECT': 'SELECT',
    'CURRENCY': 'CURRENCY',
    'BOOLEAN': 'BOOLEAN',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'TEXT': 'TEXT',
    'DATE': 'DATE',
    'SELECT': 'SELECT',
    'CURRENCY': 'CURRENCY',
    'BOOLEAN': 'BOOLEAN',
  };

  @override
  final Iterable<Type> types = const <Type>[CustomFieldType];
  @override
  final String wireName = 'CustomFieldType';

  @override
  Object serialize(Serializers serializers, CustomFieldType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  CustomFieldType deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      CustomFieldType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
