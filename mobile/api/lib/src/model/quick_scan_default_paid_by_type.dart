//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'quick_scan_default_paid_by_type.g.dart';

class QuickScanDefaultPaidByType extends EnumClass {

  /// How the default paid-by is resolved for optional quick-scan paid-by
  @BuiltValueEnumConst(wireName: r'UPLOADER')
  static const QuickScanDefaultPaidByType UPLOADER = _$UPLOADER;
  /// How the default paid-by is resolved for optional quick-scan paid-by
  @BuiltValueEnumConst(wireName: r'USER')
  static const QuickScanDefaultPaidByType USER = _$USER;
  /// How the default paid-by is resolved for optional quick-scan paid-by
  @BuiltValueEnumConst(wireName: r'')
  static const QuickScanDefaultPaidByType empty = _$empty;

  static Serializer<QuickScanDefaultPaidByType> get serializer => _$quickScanDefaultPaidByTypeSerializer;

  const QuickScanDefaultPaidByType._(String name): super(name);

  static BuiltSet<QuickScanDefaultPaidByType> get values => _$values;
  static QuickScanDefaultPaidByType valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class QuickScanDefaultPaidByTypeMixin = Object with _$QuickScanDefaultPaidByTypeMixin;

