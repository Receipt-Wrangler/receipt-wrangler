//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'filter_operation.g.dart';

class FilterOperation extends EnumClass {

  @BuiltValueEnumConst(wireName: r'CONTAINS')
  static const FilterOperation CONTAINS = _$CONTAINS;
  @BuiltValueEnumConst(wireName: r'EQUALS')
  static const FilterOperation EQUALS = _$EQUALS;
  @BuiltValueEnumConst(wireName: r'GREATER_THAN')
  static const FilterOperation GREATER_THAN = _$GREATER_THAN;
  @BuiltValueEnumConst(wireName: r'LESS_THAN')
  static const FilterOperation LESS_THAN = _$LESS_THAN;
  @BuiltValueEnumConst(wireName: r'BETWEEN')
  static const FilterOperation BETWEEN = _$BETWEEN;
  @BuiltValueEnumConst(wireName: r'WITHIN_CURRENT_MONTH')
  static const FilterOperation WITHIN_CURRENT_MONTH = _$WITHIN_CURRENT_MONTH;
  @BuiltValueEnumConst(wireName: r'')
  static const FilterOperation empty = _$empty;

  static Serializer<FilterOperation> get serializer => _$filterOperationSerializer;

  const FilterOperation._(String name): super(name);

  static BuiltSet<FilterOperation> get values => _$values;
  static FilterOperation valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class FilterOperationMixin = Object with _$FilterOperationMixin;

