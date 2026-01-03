// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'system_task_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const SystemTaskType _$OCR_PROCESSING =
    const SystemTaskType._('OCR_PROCESSING');
const SystemTaskType _$CHAT_COMPLETION =
    const SystemTaskType._('CHAT_COMPLETION');
const SystemTaskType _$MAGIC_FILL = const SystemTaskType._('MAGIC_FILL');
const SystemTaskType _$QUICK_SCAN = const SystemTaskType._('QUICK_SCAN');
const SystemTaskType _$EMAIL_READ = const SystemTaskType._('EMAIL_READ');
const SystemTaskType _$EMAIL_UPLOAD = const SystemTaskType._('EMAIL_UPLOAD');
const SystemTaskType _$SYSTEM_EMAIL_CONNECTIVITY_CHECK =
    const SystemTaskType._('SYSTEM_EMAIL_CONNECTIVITY_CHECK');
const SystemTaskType _$RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK =
    const SystemTaskType._('RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK');
const SystemTaskType _$RECEIPT_UPLOADED =
    const SystemTaskType._('RECEIPT_UPLOADED');
const SystemTaskType _$RECEIPT_UPDATED =
    const SystemTaskType._('RECEIPT_UPDATED');
const SystemTaskType _$PROMPT_GENERATED =
    const SystemTaskType._('PROMPT_GENERATED');

SystemTaskType _$valueOf(String name) {
  switch (name) {
    case 'OCR_PROCESSING':
      return _$OCR_PROCESSING;
    case 'CHAT_COMPLETION':
      return _$CHAT_COMPLETION;
    case 'MAGIC_FILL':
      return _$MAGIC_FILL;
    case 'QUICK_SCAN':
      return _$QUICK_SCAN;
    case 'EMAIL_READ':
      return _$EMAIL_READ;
    case 'EMAIL_UPLOAD':
      return _$EMAIL_UPLOAD;
    case 'SYSTEM_EMAIL_CONNECTIVITY_CHECK':
      return _$SYSTEM_EMAIL_CONNECTIVITY_CHECK;
    case 'RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK':
      return _$RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK;
    case 'RECEIPT_UPLOADED':
      return _$RECEIPT_UPLOADED;
    case 'RECEIPT_UPDATED':
      return _$RECEIPT_UPDATED;
    case 'PROMPT_GENERATED':
      return _$PROMPT_GENERATED;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<SystemTaskType> _$values =
    new BuiltSet<SystemTaskType>(const <SystemTaskType>[
  _$OCR_PROCESSING,
  _$CHAT_COMPLETION,
  _$MAGIC_FILL,
  _$QUICK_SCAN,
  _$EMAIL_READ,
  _$EMAIL_UPLOAD,
  _$SYSTEM_EMAIL_CONNECTIVITY_CHECK,
  _$RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK,
  _$RECEIPT_UPLOADED,
  _$RECEIPT_UPDATED,
  _$PROMPT_GENERATED,
]);

class _$SystemTaskTypeMeta {
  const _$SystemTaskTypeMeta();
  SystemTaskType get OCR_PROCESSING => _$OCR_PROCESSING;
  SystemTaskType get CHAT_COMPLETION => _$CHAT_COMPLETION;
  SystemTaskType get MAGIC_FILL => _$MAGIC_FILL;
  SystemTaskType get QUICK_SCAN => _$QUICK_SCAN;
  SystemTaskType get EMAIL_READ => _$EMAIL_READ;
  SystemTaskType get EMAIL_UPLOAD => _$EMAIL_UPLOAD;
  SystemTaskType get SYSTEM_EMAIL_CONNECTIVITY_CHECK =>
      _$SYSTEM_EMAIL_CONNECTIVITY_CHECK;
  SystemTaskType get RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK =>
      _$RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK;
  SystemTaskType get RECEIPT_UPLOADED => _$RECEIPT_UPLOADED;
  SystemTaskType get RECEIPT_UPDATED => _$RECEIPT_UPDATED;
  SystemTaskType get PROMPT_GENERATED => _$PROMPT_GENERATED;
  SystemTaskType valueOf(String name) => _$valueOf(name);
  BuiltSet<SystemTaskType> get values => _$values;
}

abstract class _$SystemTaskTypeMixin {
  // ignore: non_constant_identifier_names
  _$SystemTaskTypeMeta get SystemTaskType => const _$SystemTaskTypeMeta();
}

Serializer<SystemTaskType> _$systemTaskTypeSerializer =
    new _$SystemTaskTypeSerializer();

class _$SystemTaskTypeSerializer
    implements PrimitiveSerializer<SystemTaskType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'OCR_PROCESSING': 'OCR_PROCESSING',
    'CHAT_COMPLETION': 'CHAT_COMPLETION',
    'MAGIC_FILL': 'MAGIC_FILL',
    'QUICK_SCAN': 'QUICK_SCAN',
    'EMAIL_READ': 'EMAIL_READ',
    'EMAIL_UPLOAD': 'EMAIL_UPLOAD',
    'SYSTEM_EMAIL_CONNECTIVITY_CHECK': 'SYSTEM_EMAIL_CONNECTIVITY_CHECK',
    'RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK':
        'RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK',
    'RECEIPT_UPLOADED': 'RECEIPT_UPLOADED',
    'RECEIPT_UPDATED': 'RECEIPT_UPDATED',
    'PROMPT_GENERATED': 'PROMPT_GENERATED',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'OCR_PROCESSING': 'OCR_PROCESSING',
    'CHAT_COMPLETION': 'CHAT_COMPLETION',
    'MAGIC_FILL': 'MAGIC_FILL',
    'QUICK_SCAN': 'QUICK_SCAN',
    'EMAIL_READ': 'EMAIL_READ',
    'EMAIL_UPLOAD': 'EMAIL_UPLOAD',
    'SYSTEM_EMAIL_CONNECTIVITY_CHECK': 'SYSTEM_EMAIL_CONNECTIVITY_CHECK',
    'RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK':
        'RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK',
    'RECEIPT_UPLOADED': 'RECEIPT_UPLOADED',
    'RECEIPT_UPDATED': 'RECEIPT_UPDATED',
    'PROMPT_GENERATED': 'PROMPT_GENERATED',
  };

  @override
  final Iterable<Type> types = const <Type>[SystemTaskType];
  @override
  final String wireName = 'SystemTaskType';

  @override
  Object serialize(Serializers serializers, SystemTaskType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  SystemTaskType deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      SystemTaskType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
