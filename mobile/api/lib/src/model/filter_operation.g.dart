// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'filter_operation.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const FilterOperation _$CONTAINS = const FilterOperation._('CONTAINS');
const FilterOperation _$EQUALS = const FilterOperation._('EQUALS');
const FilterOperation _$GREATER_THAN = const FilterOperation._('GREATER_THAN');
const FilterOperation _$LESS_THAN = const FilterOperation._('LESS_THAN');
const FilterOperation _$BETWEEN = const FilterOperation._('BETWEEN');
const FilterOperation _$WITHIN_CURRENT_MONTH =
    const FilterOperation._('WITHIN_CURRENT_MONTH');
const FilterOperation _$empty = const FilterOperation._('empty');

FilterOperation _$valueOf(String name) {
  switch (name) {
    case 'CONTAINS':
      return _$CONTAINS;
    case 'EQUALS':
      return _$EQUALS;
    case 'GREATER_THAN':
      return _$GREATER_THAN;
    case 'LESS_THAN':
      return _$LESS_THAN;
    case 'BETWEEN':
      return _$BETWEEN;
    case 'WITHIN_CURRENT_MONTH':
      return _$WITHIN_CURRENT_MONTH;
    case 'empty':
      return _$empty;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<FilterOperation> _$values =
    new BuiltSet<FilterOperation>(const <FilterOperation>[
  _$CONTAINS,
  _$EQUALS,
  _$GREATER_THAN,
  _$LESS_THAN,
  _$BETWEEN,
  _$WITHIN_CURRENT_MONTH,
  _$empty,
]);

class _$FilterOperationMeta {
  const _$FilterOperationMeta();
  FilterOperation get CONTAINS => _$CONTAINS;
  FilterOperation get EQUALS => _$EQUALS;
  FilterOperation get GREATER_THAN => _$GREATER_THAN;
  FilterOperation get LESS_THAN => _$LESS_THAN;
  FilterOperation get BETWEEN => _$BETWEEN;
  FilterOperation get WITHIN_CURRENT_MONTH => _$WITHIN_CURRENT_MONTH;
  FilterOperation get empty => _$empty;
  FilterOperation valueOf(String name) => _$valueOf(name);
  BuiltSet<FilterOperation> get values => _$values;
}

abstract class _$FilterOperationMixin {
  // ignore: non_constant_identifier_names
  _$FilterOperationMeta get FilterOperation => const _$FilterOperationMeta();
}

Serializer<FilterOperation> _$filterOperationSerializer =
    new _$FilterOperationSerializer();

class _$FilterOperationSerializer
    implements PrimitiveSerializer<FilterOperation> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'CONTAINS': 'CONTAINS',
    'EQUALS': 'EQUALS',
    'GREATER_THAN': 'GREATER_THAN',
    'LESS_THAN': 'LESS_THAN',
    'BETWEEN': 'BETWEEN',
    'WITHIN_CURRENT_MONTH': 'WITHIN_CURRENT_MONTH',
    'empty': '',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'CONTAINS': 'CONTAINS',
    'EQUALS': 'EQUALS',
    'GREATER_THAN': 'GREATER_THAN',
    'LESS_THAN': 'LESS_THAN',
    'BETWEEN': 'BETWEEN',
    'WITHIN_CURRENT_MONTH': 'WITHIN_CURRENT_MONTH',
    '': 'empty',
  };

  @override
  final Iterable<Type> types = const <Type>[FilterOperation];
  @override
  final String wireName = 'FilterOperation';

  @override
  Object serialize(Serializers serializers, FilterOperation object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  FilterOperation deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      FilterOperation.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
