//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'queue_name.g.dart';

class QueueName extends EnumClass {

  @BuiltValueEnumConst(wireName: r'quick_scan')
  static const QueueName quickScan = _$quickScan;
  @BuiltValueEnumConst(wireName: r'email_polling')
  static const QueueName emailPolling = _$emailPolling;
  @BuiltValueEnumConst(wireName: r'email_receipt_processing')
  static const QueueName emailReceiptProcessing = _$emailReceiptProcessing;
  @BuiltValueEnumConst(wireName: r'email_receipt_image_cleanup')
  static const QueueName emailReceiptImageCleanup = _$emailReceiptImageCleanup;

  static Serializer<QueueName> get serializer => _$queueNameSerializer;

  const QueueName._(String name): super(name);

  static BuiltSet<QueueName> get values => _$values;
  static QueueName valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class QueueNameMixin = Object with _$QueueNameMixin;

