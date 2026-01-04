//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'ai_type.g.dart';

class AiType extends EnumClass {

  @BuiltValueEnumConst(wireName: r'OPEN_AI_CUSTOM')
  static const AiType OPEN_AI_CUSTOM = _$OPEN_AI_CUSTOM;
  @BuiltValueEnumConst(wireName: r'OPEN_AI')
  static const AiType OPEN_AI = _$OPEN_AI;
  @BuiltValueEnumConst(wireName: r'GEMINI')
  static const AiType GEMINI = _$GEMINI;
  @BuiltValueEnumConst(wireName: r'OLLAMA')
  static const AiType OLLAMA = _$OLLAMA;

  static Serializer<AiType> get serializer => _$aiTypeSerializer;

  const AiType._(String name): super(name);

  static BuiltSet<AiType> get values => _$values;
  static AiType valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class AiTypeMixin = Object with _$AiTypeMixin;

