//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'currency_symbol_position.g.dart';

class CurrencySymbolPosition extends EnumClass {

  @BuiltValueEnumConst(wireName: r'START')
  static const CurrencySymbolPosition START = _$START;
  @BuiltValueEnumConst(wireName: r'END')
  static const CurrencySymbolPosition END = _$END;

  static Serializer<CurrencySymbolPosition> get serializer => _$currencySymbolPositionSerializer;

  const CurrencySymbolPosition._(String name): super(name);

  static BuiltSet<CurrencySymbolPosition> get values => _$values;
  static CurrencySymbolPosition valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class CurrencySymbolPositionMixin = Object with _$CurrencySymbolPositionMixin;

