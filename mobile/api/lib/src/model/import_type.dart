//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'import_type.g.dart';

class ImportType extends EnumClass {

  @BuiltValueEnumConst(wireName: r'IMPORT_CONFIG')
  static const ImportType IMPORT_CONFIG = _$IMPORT_CONFIG;

  static Serializer<ImportType> get serializer => _$importTypeSerializer;

  const ImportType._(String name): super(name);

  static BuiltSet<ImportType> get values => _$values;
  static ImportType valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ImportTypeMixin = Object with _$ImportTypeMixin;

