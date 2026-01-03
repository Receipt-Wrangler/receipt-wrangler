// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'item_status.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ItemStatus _$OPEN = const ItemStatus._('OPEN');
const ItemStatus _$RESOLVED = const ItemStatus._('RESOLVED');
const ItemStatus _$DRAFT = const ItemStatus._('DRAFT');

ItemStatus _$valueOf(String name) {
  switch (name) {
    case 'OPEN':
      return _$OPEN;
    case 'RESOLVED':
      return _$RESOLVED;
    case 'DRAFT':
      return _$DRAFT;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<ItemStatus> _$values =
    new BuiltSet<ItemStatus>(const <ItemStatus>[
  _$OPEN,
  _$RESOLVED,
  _$DRAFT,
]);

class _$ItemStatusMeta {
  const _$ItemStatusMeta();
  ItemStatus get OPEN => _$OPEN;
  ItemStatus get RESOLVED => _$RESOLVED;
  ItemStatus get DRAFT => _$DRAFT;
  ItemStatus valueOf(String name) => _$valueOf(name);
  BuiltSet<ItemStatus> get values => _$values;
}

abstract class _$ItemStatusMixin {
  // ignore: non_constant_identifier_names
  _$ItemStatusMeta get ItemStatus => const _$ItemStatusMeta();
}

Serializer<ItemStatus> _$itemStatusSerializer = new _$ItemStatusSerializer();

class _$ItemStatusSerializer implements PrimitiveSerializer<ItemStatus> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'OPEN': 'OPEN',
    'RESOLVED': 'RESOLVED',
    'DRAFT': 'DRAFT',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'OPEN': 'OPEN',
    'RESOLVED': 'RESOLVED',
    'DRAFT': 'DRAFT',
  };

  @override
  final Iterable<Type> types = const <Type>[ItemStatus];
  @override
  final String wireName = 'ItemStatus';

  @override
  Object serialize(Serializers serializers, ItemStatus object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ItemStatus deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ItemStatus.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
