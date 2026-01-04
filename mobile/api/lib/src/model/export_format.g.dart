// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'export_format.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ExportFormat _$CSV = const ExportFormat._('CSV');

ExportFormat _$valueOf(String name) {
  switch (name) {
    case 'CSV':
      return _$CSV;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<ExportFormat> _$values =
    new BuiltSet<ExportFormat>(const <ExportFormat>[
  _$CSV,
]);

class _$ExportFormatMeta {
  const _$ExportFormatMeta();
  ExportFormat get CSV => _$CSV;
  ExportFormat valueOf(String name) => _$valueOf(name);
  BuiltSet<ExportFormat> get values => _$values;
}

abstract class _$ExportFormatMixin {
  // ignore: non_constant_identifier_names
  _$ExportFormatMeta get ExportFormat => const _$ExportFormatMeta();
}

Serializer<ExportFormat> _$exportFormatSerializer =
    new _$ExportFormatSerializer();

class _$ExportFormatSerializer implements PrimitiveSerializer<ExportFormat> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'CSV': 'CSV',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'CSV': 'CSV',
  };

  @override
  final Iterable<Type> types = const <Type>[ExportFormat];
  @override
  final String wireName = 'ExportFormat';

  @override
  Object serialize(Serializers serializers, ExportFormat object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ExportFormat deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ExportFormat.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
