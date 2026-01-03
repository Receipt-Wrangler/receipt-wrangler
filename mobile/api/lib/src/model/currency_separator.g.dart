// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'currency_separator.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const CurrencySeparator _$comma = const CurrencySeparator._('comma');
const CurrencySeparator _$period = const CurrencySeparator._('period');

CurrencySeparator _$valueOf(String name) {
  switch (name) {
    case 'comma':
      return _$comma;
    case 'period':
      return _$period;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<CurrencySeparator> _$values =
    new BuiltSet<CurrencySeparator>(const <CurrencySeparator>[
  _$comma,
  _$period,
]);

class _$CurrencySeparatorMeta {
  const _$CurrencySeparatorMeta();
  CurrencySeparator get comma => _$comma;
  CurrencySeparator get period => _$period;
  CurrencySeparator valueOf(String name) => _$valueOf(name);
  BuiltSet<CurrencySeparator> get values => _$values;
}

abstract class _$CurrencySeparatorMixin {
  // ignore: non_constant_identifier_names
  _$CurrencySeparatorMeta get CurrencySeparator =>
      const _$CurrencySeparatorMeta();
}

Serializer<CurrencySeparator> _$currencySeparatorSerializer =
    new _$CurrencySeparatorSerializer();

class _$CurrencySeparatorSerializer
    implements PrimitiveSerializer<CurrencySeparator> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'comma': ',',
    'period': '.',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    ',': 'comma',
    '.': 'period',
  };

  @override
  final Iterable<Type> types = const <Type>[CurrencySeparator];
  @override
  final String wireName = 'CurrencySeparator';

  @override
  Object serialize(Serializers serializers, CurrencySeparator object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  CurrencySeparator deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      CurrencySeparator.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
