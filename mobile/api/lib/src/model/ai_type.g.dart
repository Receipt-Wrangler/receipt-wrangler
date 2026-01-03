// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'ai_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const AiType _$OPEN_AI_CUSTOM = const AiType._('OPEN_AI_CUSTOM');
const AiType _$OPEN_AI = const AiType._('OPEN_AI');
const AiType _$GEMINI = const AiType._('GEMINI');
const AiType _$OLLAMA = const AiType._('OLLAMA');

AiType _$valueOf(String name) {
  switch (name) {
    case 'OPEN_AI_CUSTOM':
      return _$OPEN_AI_CUSTOM;
    case 'OPEN_AI':
      return _$OPEN_AI;
    case 'GEMINI':
      return _$GEMINI;
    case 'OLLAMA':
      return _$OLLAMA;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<AiType> _$values = new BuiltSet<AiType>(const <AiType>[
  _$OPEN_AI_CUSTOM,
  _$OPEN_AI,
  _$GEMINI,
  _$OLLAMA,
]);

class _$AiTypeMeta {
  const _$AiTypeMeta();
  AiType get OPEN_AI_CUSTOM => _$OPEN_AI_CUSTOM;
  AiType get OPEN_AI => _$OPEN_AI;
  AiType get GEMINI => _$GEMINI;
  AiType get OLLAMA => _$OLLAMA;
  AiType valueOf(String name) => _$valueOf(name);
  BuiltSet<AiType> get values => _$values;
}

abstract class _$AiTypeMixin {
  // ignore: non_constant_identifier_names
  _$AiTypeMeta get AiType => const _$AiTypeMeta();
}

Serializer<AiType> _$aiTypeSerializer = new _$AiTypeSerializer();

class _$AiTypeSerializer implements PrimitiveSerializer<AiType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'OPEN_AI_CUSTOM': 'OPEN_AI_CUSTOM',
    'OPEN_AI': 'OPEN_AI',
    'GEMINI': 'GEMINI',
    'OLLAMA': 'OLLAMA',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'OPEN_AI_CUSTOM': 'OPEN_AI_CUSTOM',
    'OPEN_AI': 'OPEN_AI',
    'GEMINI': 'GEMINI',
    'OLLAMA': 'OLLAMA',
  };

  @override
  final Iterable<Type> types = const <Type>[AiType];
  @override
  final String wireName = 'AiType';

  @override
  Object serialize(Serializers serializers, AiType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  AiType deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      AiType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
