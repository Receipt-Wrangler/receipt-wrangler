// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'currency_symbol_position.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const CurrencySymbolPosition _$START = const CurrencySymbolPosition._('START');
const CurrencySymbolPosition _$END = const CurrencySymbolPosition._('END');

CurrencySymbolPosition _$valueOf(String name) {
  switch (name) {
    case 'START':
      return _$START;
    case 'END':
      return _$END;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<CurrencySymbolPosition> _$values =
    new BuiltSet<CurrencySymbolPosition>(const <CurrencySymbolPosition>[
  _$START,
  _$END,
]);

class _$CurrencySymbolPositionMeta {
  const _$CurrencySymbolPositionMeta();
  CurrencySymbolPosition get START => _$START;
  CurrencySymbolPosition get END => _$END;
  CurrencySymbolPosition valueOf(String name) => _$valueOf(name);
  BuiltSet<CurrencySymbolPosition> get values => _$values;
}

abstract class _$CurrencySymbolPositionMixin {
  // ignore: non_constant_identifier_names
  _$CurrencySymbolPositionMeta get CurrencySymbolPosition =>
      const _$CurrencySymbolPositionMeta();
}

Serializer<CurrencySymbolPosition> _$currencySymbolPositionSerializer =
    new _$CurrencySymbolPositionSerializer();

class _$CurrencySymbolPositionSerializer
    implements PrimitiveSerializer<CurrencySymbolPosition> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'START': 'START',
    'END': 'END',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'START': 'START',
    'END': 'END',
  };

  @override
  final Iterable<Type> types = const <Type>[CurrencySymbolPosition];
  @override
  final String wireName = 'CurrencySymbolPosition';

  @override
  Object serialize(Serializers serializers, CurrencySymbolPosition object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  CurrencySymbolPosition deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      CurrencySymbolPosition.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
