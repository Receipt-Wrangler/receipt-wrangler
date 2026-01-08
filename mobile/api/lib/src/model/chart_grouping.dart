//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'chart_grouping.g.dart';

class ChartGrouping extends EnumClass {

  @BuiltValueEnumConst(wireName: r'CATEGORIES')
  static const ChartGrouping CATEGORIES = _$CATEGORIES;
  @BuiltValueEnumConst(wireName: r'TAGS')
  static const ChartGrouping TAGS = _$TAGS;
  @BuiltValueEnumConst(wireName: r'PAIDBY')
  static const ChartGrouping PAIDBY = _$PAIDBY;

  static Serializer<ChartGrouping> get serializer => _$chartGroupingSerializer;

  const ChartGrouping._(String name): super(name);

  static BuiltSet<ChartGrouping> get values => _$values;
  static ChartGrouping valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ChartGroupingMixin = Object with _$ChartGroupingMixin;

