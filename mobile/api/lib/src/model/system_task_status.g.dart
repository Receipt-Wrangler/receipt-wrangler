// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'system_task_status.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const SystemTaskStatus _$SUCCEEDED = const SystemTaskStatus._('SUCCEEDED');
const SystemTaskStatus _$FAILED = const SystemTaskStatus._('FAILED');

SystemTaskStatus _$valueOf(String name) {
  switch (name) {
    case 'SUCCEEDED':
      return _$SUCCEEDED;
    case 'FAILED':
      return _$FAILED;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<SystemTaskStatus> _$values =
    new BuiltSet<SystemTaskStatus>(const <SystemTaskStatus>[
  _$SUCCEEDED,
  _$FAILED,
]);

class _$SystemTaskStatusMeta {
  const _$SystemTaskStatusMeta();
  SystemTaskStatus get SUCCEEDED => _$SUCCEEDED;
  SystemTaskStatus get FAILED => _$FAILED;
  SystemTaskStatus valueOf(String name) => _$valueOf(name);
  BuiltSet<SystemTaskStatus> get values => _$values;
}

abstract class _$SystemTaskStatusMixin {
  // ignore: non_constant_identifier_names
  _$SystemTaskStatusMeta get SystemTaskStatus => const _$SystemTaskStatusMeta();
}

Serializer<SystemTaskStatus> _$systemTaskStatusSerializer =
    new _$SystemTaskStatusSerializer();

class _$SystemTaskStatusSerializer
    implements PrimitiveSerializer<SystemTaskStatus> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'SUCCEEDED': 'SUCCEEDED',
    'FAILED': 'FAILED',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'SUCCEEDED': 'SUCCEEDED',
    'FAILED': 'FAILED',
  };

  @override
  final Iterable<Type> types = const <Type>[SystemTaskStatus];
  @override
  final String wireName = 'SystemTaskStatus';

  @override
  Object serialize(Serializers serializers, SystemTaskStatus object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  SystemTaskStatus deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      SystemTaskStatus.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
