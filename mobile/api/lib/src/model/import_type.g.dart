// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'import_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ImportType _$IMPORT_CONFIG = const ImportType._('IMPORT_CONFIG');

ImportType _$valueOf(String name) {
  switch (name) {
    case 'IMPORT_CONFIG':
      return _$IMPORT_CONFIG;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<ImportType> _$values =
    new BuiltSet<ImportType>(const <ImportType>[
  _$IMPORT_CONFIG,
]);

class _$ImportTypeMeta {
  const _$ImportTypeMeta();
  ImportType get IMPORT_CONFIG => _$IMPORT_CONFIG;
  ImportType valueOf(String name) => _$valueOf(name);
  BuiltSet<ImportType> get values => _$values;
}

abstract class _$ImportTypeMixin {
  // ignore: non_constant_identifier_names
  _$ImportTypeMeta get ImportType => const _$ImportTypeMeta();
}

Serializer<ImportType> _$importTypeSerializer = new _$ImportTypeSerializer();

class _$ImportTypeSerializer implements PrimitiveSerializer<ImportType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'IMPORT_CONFIG': 'IMPORT_CONFIG',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'IMPORT_CONFIG': 'IMPORT_CONFIG',
  };

  @override
  final Iterable<Type> types = const <Type>[ImportType];
  @override
  final String wireName = 'ImportType';

  @override
  Object serialize(Serializers serializers, ImportType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ImportType deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ImportType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
