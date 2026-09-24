//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary_position.g.dart';

class ReceiptSummaryPosition extends EnumClass {

  /// Where a group's receipt summary renders relative to its receipts list. A client that meets a value added later must degrade to BOTTOM, which is where the block rendered before the setting existed, rather than failing the whole payload.
  @BuiltValueEnumConst(wireName: r'TOP')
  static const ReceiptSummaryPosition TOP = _$TOP;
  /// Where a group's receipt summary renders relative to its receipts list. A client that meets a value added later must degrade to BOTTOM, which is where the block rendered before the setting existed, rather than failing the whole payload.
  // Applied by api/patches/apply-dart-dio-patches.sh -- do not hand-edit, and do not
  // drop it as generator noise. Without `fallback: true` the generated _$valueOf throws
  // on an unrecognized wire value, failing the WHOLE enclosing payload rather than the
  // one field. See mobile/CLAUDE.md for the two outages that came of it.
  @BuiltValueEnumConst(wireName: r'BOTTOM', fallback: true)
  static const ReceiptSummaryPosition BOTTOM = _$BOTTOM;

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

