//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary_position.g.dart';

class ReceiptSummaryPosition extends EnumClass {

  /// Where a group's receipt summary renders relative to its receipts list. The empty member exists so an already-released client tolerates a value added later rather than failing the whole payload; the server never sends it.
  @BuiltValueEnumConst(wireName: r'TOP')
  static const ReceiptSummaryPosition TOP = _$TOP;
  /// Where a group's receipt summary renders relative to its receipts list. The empty member exists so an already-released client tolerates a value added later rather than failing the whole payload; the server never sends it.
  @BuiltValueEnumConst(wireName: r'BOTTOM')
  static const ReceiptSummaryPosition BOTTOM = _$BOTTOM;
  /// Where a group's receipt summary renders relative to its receipts list. The empty member exists so an already-released client tolerates a value added later rather than failing the whole payload; the server never sends it.
  // HAND-PATCH (re-apply after every regen -- see mobile/CLAUDE.md "Known
  // dart-dio default-value regressions"): `fallback: true` makes the generated
  // `_$valueOf` return this member for an unrecognized wire value instead of
  // throwing ArgumentError, which would fail the ENTIRE enclosing payload.
  // GroupReceiptSettings rides on AppData, so that payload is LOGIN -- adding a
  // third position later would brick every already-released build at the sign-in
  // screen, which is exactly the two Permission outages. The openapi-generator
  // does not emit `fallback`, and nothing fails to compile without it.
  @BuiltValueEnumConst(wireName: r'', fallback: true)
  static const ReceiptSummaryPosition empty = _$empty;

  static Serializer<ReceiptSummaryPosition> get serializer => _$receiptSummaryPositionSerializer;

  const ReceiptSummaryPosition._(String name): super(name);

  static BuiltSet<ReceiptSummaryPosition> get values => _$values;
  static ReceiptSummaryPosition valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ReceiptSummaryPositionMixin = Object with _$ReceiptSummaryPositionMixin;

