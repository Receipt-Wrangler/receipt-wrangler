//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_status.g.dart';

class ReceiptStatus extends EnumClass {

  /// Status of a receipt
  @BuiltValueEnumConst(wireName: r'OPEN')
  static const ReceiptStatus OPEN = _$OPEN;
  /// Status of a receipt
  @BuiltValueEnumConst(wireName: r'NEEDS_ATTENTION')
  static const ReceiptStatus NEEDS_ATTENTION = _$NEEDS_ATTENTION;
  /// Status of a receipt
  @BuiltValueEnumConst(wireName: r'RESOLVED')
  static const ReceiptStatus RESOLVED = _$RESOLVED;
  /// Status of a receipt
  @BuiltValueEnumConst(wireName: r'DRAFT')
  static const ReceiptStatus DRAFT = _$DRAFT;
  /// Status of a receipt
  @BuiltValueEnumConst(wireName: r'')
  static const ReceiptStatus empty = _$empty;

  static Serializer<ReceiptStatus> get serializer => _$receiptStatusSerializer;

  const ReceiptStatus._(String name): super(name);

  static BuiltSet<ReceiptStatus> get values => _$values;
  static ReceiptStatus valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ReceiptStatusMixin = Object with _$ReceiptStatusMixin;

