//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'currency_separator.g.dart';

class CurrencySeparator extends EnumClass {

  @BuiltValueEnumConst(wireName: r',')
  static const CurrencySeparator comma = _$comma;
  @BuiltValueEnumConst(wireName: r'.')
  static const CurrencySeparator period = _$period;

  static Serializer<CurrencySeparator> get serializer => _$currencySeparatorSerializer;

  const CurrencySeparator._(String name): super(name);

  static BuiltSet<CurrencySeparator> get values => _$values;
  static CurrencySeparator valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class CurrencySeparatorMixin = Object with _$CurrencySeparatorMixin;

