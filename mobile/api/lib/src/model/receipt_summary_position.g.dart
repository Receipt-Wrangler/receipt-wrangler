// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_summary_position.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ReceiptSummaryPosition _$TOP = const ReceiptSummaryPosition._('TOP');
const ReceiptSummaryPosition _$BOTTOM =
    const ReceiptSummaryPosition._('BOTTOM');

ReceiptSummaryPosition _$valueOf(String name) {
  switch (name) {
    case 'TOP':
      return _$TOP;
    case 'BOTTOM':
      return _$BOTTOM;
    default:
      return _$BOTTOM;
  }
}

final BuiltSet<ReceiptSummaryPosition> _$values =
    BuiltSet<ReceiptSummaryPosition>(const <ReceiptSummaryPosition>[
  _$TOP,
  _$BOTTOM,
]);

class _$ReceiptSummaryPositionMeta {
  const _$ReceiptSummaryPositionMeta();
  ReceiptSummaryPosition get TOP => _$TOP;
  ReceiptSummaryPosition get BOTTOM => _$BOTTOM;
  ReceiptSummaryPosition valueOf(String name) => _$valueOf(name);
  BuiltSet<ReceiptSummaryPosition> get values => _$values;
}

abstract class _$ReceiptSummaryPositionMixin {
  // ignore: non_constant_identifier_names
  _$ReceiptSummaryPositionMeta get ReceiptSummaryPosition =>
      const _$ReceiptSummaryPositionMeta();
}

Serializer<ReceiptSummaryPosition> _$receiptSummaryPositionSerializer =
    _$ReceiptSummaryPositionSerializer();

class _$ReceiptSummaryPositionSerializer
    implements PrimitiveSerializer<ReceiptSummaryPosition> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'TOP': 'TOP',
    'BOTTOM': 'BOTTOM',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'TOP': 'TOP',
    'BOTTOM': 'BOTTOM',
  };

  @override
  final Iterable<Type> types = const <Type>[ReceiptSummaryPosition];
  @override
  final String wireName = 'ReceiptSummaryPosition';

  @override
  Object serialize(Serializers serializers, ReceiptSummaryPosition object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ReceiptSummaryPosition deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ReceiptSummaryPosition.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
