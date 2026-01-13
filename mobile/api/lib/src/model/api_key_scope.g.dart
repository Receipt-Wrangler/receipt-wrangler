// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_key_scope.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ApiKeyScope _$r = const ApiKeyScope._('r');
const ApiKeyScope _$w = const ApiKeyScope._('w');
const ApiKeyScope _$rw = const ApiKeyScope._('rw');

ApiKeyScope _$valueOf(String name) {
  switch (name) {
    case 'r':
      return _$r;
    case 'w':
      return _$w;
    case 'rw':
      return _$rw;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<ApiKeyScope> _$values =
    new BuiltSet<ApiKeyScope>(const <ApiKeyScope>[
  _$r,
  _$w,
  _$rw,
]);

class _$ApiKeyScopeMeta {
  const _$ApiKeyScopeMeta();
  ApiKeyScope get r => _$r;
  ApiKeyScope get w => _$w;
  ApiKeyScope get rw => _$rw;
  ApiKeyScope valueOf(String name) => _$valueOf(name);
  BuiltSet<ApiKeyScope> get values => _$values;
}

abstract class _$ApiKeyScopeMixin {
  // ignore: non_constant_identifier_names
  _$ApiKeyScopeMeta get ApiKeyScope => const _$ApiKeyScopeMeta();
}

Serializer<ApiKeyScope> _$apiKeyScopeSerializer = new _$ApiKeyScopeSerializer();

class _$ApiKeyScopeSerializer implements PrimitiveSerializer<ApiKeyScope> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'r': 'r',
    'w': 'w',
    'rw': 'rw',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'r': 'r',
    'w': 'w',
    'rw': 'rw',
  };

  @override
  final Iterable<Type> types = const <Type>[ApiKeyScope];
  @override
  final String wireName = 'ApiKeyScope';

  @override
  Object serialize(Serializers serializers, ApiKeyScope object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ApiKeyScope deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ApiKeyScope.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
