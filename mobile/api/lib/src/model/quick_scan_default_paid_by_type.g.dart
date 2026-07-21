// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'quick_scan_default_paid_by_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const QuickScanDefaultPaidByType _$UPLOADER =
    const QuickScanDefaultPaidByType._('UPLOADER');
const QuickScanDefaultPaidByType _$USER =
    const QuickScanDefaultPaidByType._('USER');
const QuickScanDefaultPaidByType _$empty =
    const QuickScanDefaultPaidByType._('empty');

QuickScanDefaultPaidByType _$valueOf(String name) {
  switch (name) {
    case 'UPLOADER':
      return _$UPLOADER;
    case 'USER':
      return _$USER;
    case 'empty':
      return _$empty;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<QuickScanDefaultPaidByType> _$values =
    BuiltSet<QuickScanDefaultPaidByType>(const <QuickScanDefaultPaidByType>[
  _$UPLOADER,
  _$USER,
  _$empty,
]);

class _$QuickScanDefaultPaidByTypeMeta {
  const _$QuickScanDefaultPaidByTypeMeta();
  QuickScanDefaultPaidByType get UPLOADER => _$UPLOADER;
  QuickScanDefaultPaidByType get USER => _$USER;
  QuickScanDefaultPaidByType get empty => _$empty;
  QuickScanDefaultPaidByType valueOf(String name) => _$valueOf(name);
  BuiltSet<QuickScanDefaultPaidByType> get values => _$values;
}

abstract class _$QuickScanDefaultPaidByTypeMixin {
  // ignore: non_constant_identifier_names
  _$QuickScanDefaultPaidByTypeMeta get QuickScanDefaultPaidByType =>
      const _$QuickScanDefaultPaidByTypeMeta();
}

Serializer<QuickScanDefaultPaidByType> _$quickScanDefaultPaidByTypeSerializer =
    _$QuickScanDefaultPaidByTypeSerializer();

class _$QuickScanDefaultPaidByTypeSerializer
    implements PrimitiveSerializer<QuickScanDefaultPaidByType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'UPLOADER': 'UPLOADER',
    'USER': 'USER',
    'empty': '',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'UPLOADER': 'UPLOADER',
    'USER': 'USER',
    '': 'empty',
  };

  @override
  final Iterable<Type> types = const <Type>[QuickScanDefaultPaidByType];
  @override
  final String wireName = 'QuickScanDefaultPaidByType';

  @override
  Object serialize(Serializers serializers, QuickScanDefaultPaidByType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  QuickScanDefaultPaidByType deserialize(
          Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      QuickScanDefaultPaidByType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
