// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'queue_name.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const QueueName _$quickScan = const QueueName._('quickScan');
const QueueName _$emailPolling = const QueueName._('emailPolling');
const QueueName _$emailReceiptProcessing =
    const QueueName._('emailReceiptProcessing');
const QueueName _$emailReceiptImageCleanup =
    const QueueName._('emailReceiptImageCleanup');

QueueName _$valueOf(String name) {
  switch (name) {
    case 'quickScan':
      return _$quickScan;
    case 'emailPolling':
      return _$emailPolling;
    case 'emailReceiptProcessing':
      return _$emailReceiptProcessing;
    case 'emailReceiptImageCleanup':
      return _$emailReceiptImageCleanup;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<QueueName> _$values = new BuiltSet<QueueName>(const <QueueName>[
  _$quickScan,
  _$emailPolling,
  _$emailReceiptProcessing,
  _$emailReceiptImageCleanup,
]);

class _$QueueNameMeta {
  const _$QueueNameMeta();
  QueueName get quickScan => _$quickScan;
  QueueName get emailPolling => _$emailPolling;
  QueueName get emailReceiptProcessing => _$emailReceiptProcessing;
  QueueName get emailReceiptImageCleanup => _$emailReceiptImageCleanup;
  QueueName valueOf(String name) => _$valueOf(name);
  BuiltSet<QueueName> get values => _$values;
}

abstract class _$QueueNameMixin {
  // ignore: non_constant_identifier_names
  _$QueueNameMeta get QueueName => const _$QueueNameMeta();
}

Serializer<QueueName> _$queueNameSerializer = new _$QueueNameSerializer();

class _$QueueNameSerializer implements PrimitiveSerializer<QueueName> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'quickScan': 'quick_scan',
    'emailPolling': 'email_polling',
    'emailReceiptProcessing': 'email_receipt_processing',
    'emailReceiptImageCleanup': 'email_receipt_image_cleanup',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'quick_scan': 'quickScan',
    'email_polling': 'emailPolling',
    'email_receipt_processing': 'emailReceiptProcessing',
    'email_receipt_image_cleanup': 'emailReceiptImageCleanup',
  };

  @override
  final Iterable<Type> types = const <Type>[QueueName];
  @override
  final String wireName = 'QueueName';

  @override
  Object serialize(Serializers serializers, QueueName object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  QueueName deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      QueueName.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
