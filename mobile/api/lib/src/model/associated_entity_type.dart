//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'associated_entity_type.g.dart';

class AssociatedEntityType extends EnumClass {

  @BuiltValueEnumConst(wireName: r'NOOP_ENTITY_TYPE')
  static const AssociatedEntityType NOOP_ENTITY_TYPE = _$NOOP_ENTITY_TYPE;
  @BuiltValueEnumConst(wireName: r'RECEIPT')
  static const AssociatedEntityType RECEIPT = _$RECEIPT;
  @BuiltValueEnumConst(wireName: r'SYSTEM_EMAIL')
  static const AssociatedEntityType SYSTEM_EMAIL = _$SYSTEM_EMAIL;
  @BuiltValueEnumConst(wireName: r'RECEIPT_PROCESSING_SETTINGS')
  static const AssociatedEntityType RECEIPT_PROCESSING_SETTINGS = _$RECEIPT_PROCESSING_SETTINGS;
  @BuiltValueEnumConst(wireName: r'PROMPT')
  static const AssociatedEntityType PROMPT = _$PROMPT;

  static Serializer<AssociatedEntityType> get serializer => _$associatedEntityTypeSerializer;

  const AssociatedEntityType._(String name): super(name);

  static BuiltSet<AssociatedEntityType> get values => _$values;
  static AssociatedEntityType valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class AssociatedEntityTypeMixin = Object with _$AssociatedEntityTypeMixin;

