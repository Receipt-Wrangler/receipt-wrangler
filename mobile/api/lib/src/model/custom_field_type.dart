//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'custom_field_type.g.dart';

class CustomFieldType extends EnumClass {

  @BuiltValueEnumConst(wireName: r'TEXT')
  static const CustomFieldType TEXT = _$TEXT;
  @BuiltValueEnumConst(wireName: r'DATE')
  static const CustomFieldType DATE = _$DATE;
  @BuiltValueEnumConst(wireName: r'SELECT')
  static const CustomFieldType SELECT = _$SELECT;
  @BuiltValueEnumConst(wireName: r'CURRENCY')
  static const CustomFieldType CURRENCY = _$CURRENCY;
  @BuiltValueEnumConst(wireName: r'BOOLEAN')
  static const CustomFieldType BOOLEAN = _$BOOLEAN;

  static Serializer<CustomFieldType> get serializer => _$customFieldTypeSerializer;

  const CustomFieldType._(String name): super(name);

  static BuiltSet<CustomFieldType> get values => _$values;
  static CustomFieldType valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class CustomFieldTypeMixin = Object with _$CustomFieldTypeMixin;

