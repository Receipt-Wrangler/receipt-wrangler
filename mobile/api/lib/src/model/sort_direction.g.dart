// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'sort_direction.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const SortDirection _$asc = const SortDirection._('asc');
const SortDirection _$desc = const SortDirection._('desc');
const SortDirection _$empty = const SortDirection._('empty');

SortDirection _$valueOf(String name) {
  switch (name) {
    case 'asc':
      return _$asc;
    case 'desc':
      return _$desc;
    case 'empty':
      return _$empty;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<SortDirection> _$values =
    new BuiltSet<SortDirection>(const <SortDirection>[
  _$asc,
  _$desc,
  _$empty,
]);

class _$SortDirectionMeta {
  const _$SortDirectionMeta();
  SortDirection get asc => _$asc;
  SortDirection get desc => _$desc;
  SortDirection get empty => _$empty;
  SortDirection valueOf(String name) => _$valueOf(name);
  BuiltSet<SortDirection> get values => _$values;
}

abstract class _$SortDirectionMixin {
  // ignore: non_constant_identifier_names
  _$SortDirectionMeta get SortDirection => const _$SortDirectionMeta();
}

Serializer<SortDirection> _$sortDirectionSerializer =
    new _$SortDirectionSerializer();

class _$SortDirectionSerializer implements PrimitiveSerializer<SortDirection> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'asc': 'asc',
    'desc': 'desc',
    'empty': '',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'asc': 'asc',
    'desc': 'desc',
    '': 'empty',
  };

  @override
  final Iterable<Type> types = const <Type>[SortDirection];
  @override
  final String wireName = 'SortDirection';

  @override
  Object serialize(Serializers serializers, SortDirection object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  SortDirection deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      SortDirection.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
