//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'sort_direction.g.dart';

class SortDirection extends EnumClass {

  @BuiltValueEnumConst(wireName: r'asc')
  static const SortDirection asc = _$asc;
  @BuiltValueEnumConst(wireName: r'desc')
  static const SortDirection desc = _$desc;
  @BuiltValueEnumConst(wireName: r'')
  static const SortDirection empty = _$empty;

  static Serializer<SortDirection> get serializer => _$sortDirectionSerializer;

  const SortDirection._(String name): super(name);

  static BuiltSet<SortDirection> get values => _$values;
  static SortDirection valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class SortDirectionMixin = Object with _$SortDirectionMixin;

