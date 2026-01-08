//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'system_task_type.g.dart';

class SystemTaskType extends EnumClass {

  @BuiltValueEnumConst(wireName: r'OCR_PROCESSING')
  static const SystemTaskType OCR_PROCESSING = _$OCR_PROCESSING;
  @BuiltValueEnumConst(wireName: r'CHAT_COMPLETION')
  static const SystemTaskType CHAT_COMPLETION = _$CHAT_COMPLETION;
  @BuiltValueEnumConst(wireName: r'MAGIC_FILL')
  static const SystemTaskType MAGIC_FILL = _$MAGIC_FILL;
  @BuiltValueEnumConst(wireName: r'QUICK_SCAN')
  static const SystemTaskType QUICK_SCAN = _$QUICK_SCAN;
  @BuiltValueEnumConst(wireName: r'EMAIL_READ')
  static const SystemTaskType EMAIL_READ = _$EMAIL_READ;
  @BuiltValueEnumConst(wireName: r'EMAIL_UPLOAD')
  static const SystemTaskType EMAIL_UPLOAD = _$EMAIL_UPLOAD;
  @BuiltValueEnumConst(wireName: r'SYSTEM_EMAIL_CONNECTIVITY_CHECK')
  static const SystemTaskType SYSTEM_EMAIL_CONNECTIVITY_CHECK = _$SYSTEM_EMAIL_CONNECTIVITY_CHECK;
  @BuiltValueEnumConst(wireName: r'RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK')
  static const SystemTaskType RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK = _$RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK;
  @BuiltValueEnumConst(wireName: r'RECEIPT_UPLOADED')
  static const SystemTaskType RECEIPT_UPLOADED = _$RECEIPT_UPLOADED;
  @BuiltValueEnumConst(wireName: r'RECEIPT_UPDATED')
  static const SystemTaskType RECEIPT_UPDATED = _$RECEIPT_UPDATED;
  @BuiltValueEnumConst(wireName: r'PROMPT_GENERATED')
  static const SystemTaskType PROMPT_GENERATED = _$PROMPT_GENERATED;
  @BuiltValueEnumConst(wireName: r'API_KEY_DELETED')
  static const SystemTaskType API_KEY_DELETED = _$API_KEY_DELETED;

  static Serializer<SystemTaskType> get serializer => _$systemTaskTypeSerializer;

  const SystemTaskType._(String name): super(name);

  static BuiltSet<SystemTaskType> get values => _$values;
  static SystemTaskType valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class SystemTaskTypeMixin = Object with _$SystemTaskTypeMixin;

