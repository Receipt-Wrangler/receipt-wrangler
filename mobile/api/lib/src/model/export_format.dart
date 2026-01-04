//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'export_format.g.dart';

class ExportFormat extends EnumClass {

  @BuiltValueEnumConst(wireName: r'CSV')
  static const ExportFormat CSV = _$CSV;

  static Serializer<ExportFormat> get serializer => _$exportFormatSerializer;

  const ExportFormat._(String name): super(name);

  static BuiltSet<ExportFormat> get values => _$values;
  static ExportFormat valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ExportFormatMixin = Object with _$ExportFormatMixin;

