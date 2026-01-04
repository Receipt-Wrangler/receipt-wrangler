//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/ai_type.dart';
import 'package:openapi/src/model/ocr_engine.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_receipt_processing_settings_command.g.dart';

/// UpsertReceiptProcessingSettingsCommand
///
/// Properties:
/// * [name] - Name of the settings
/// * [description] - Description of the settings
/// * [aiType] 
/// * [url] - URL for custom endpoints
/// * [key] - Key for endpoints that require authentication
/// * [model] - LLM model
/// * [isVisionModel] - Is vision model
/// * [ocrEngine] 
/// * [promptId] - Prompt foreign key
@BuiltValue()
abstract class UpsertReceiptProcessingSettingsCommand implements Built<UpsertReceiptProcessingSettingsCommand, UpsertReceiptProcessingSettingsCommandBuilder> {
  /// Name of the settings
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Description of the settings
  @BuiltValueField(wireName: r'description')
  String? get description;

  @BuiltValueField(wireName: r'aiType')
  AiType get aiType;
  // enum aiTypeEnum {  OPEN_AI_CUSTOM,  OPEN_AI,  GEMINI,  OLLAMA,  };

  /// URL for custom endpoints
  @BuiltValueField(wireName: r'url')
  String? get url;

  /// Key for endpoints that require authentication
  @BuiltValueField(wireName: r'key')
  String? get key;

  /// LLM model
  @BuiltValueField(wireName: r'model')
  String? get model;

  /// Is vision model
  @BuiltValueField(wireName: r'isVisionModel')
  bool? get isVisionModel;

  @BuiltValueField(wireName: r'ocrEngine')
  OcrEngine get ocrEngine;
  // enum ocrEngineEnum {  TESSERACT,  EASY_OCR,  };

  /// Prompt foreign key
  @BuiltValueField(wireName: r'promptId')
  int get promptId;

  UpsertReceiptProcessingSettingsCommand._();

  factory UpsertReceiptProcessingSettingsCommand([void updates(UpsertReceiptProcessingSettingsCommandBuilder b)]) = _$UpsertReceiptProcessingSettingsCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertReceiptProcessingSettingsCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertReceiptProcessingSettingsCommand> get serializer => _$UpsertReceiptProcessingSettingsCommandSerializer();
}

class _$UpsertReceiptProcessingSettingsCommandSerializer implements PrimitiveSerializer<UpsertReceiptProcessingSettingsCommand> {
  @override
  final Iterable<Type> types = const [UpsertReceiptProcessingSettingsCommand, _$UpsertReceiptProcessingSettingsCommand];

  @override
  final String wireName = r'UpsertReceiptProcessingSettingsCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertReceiptProcessingSettingsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.description != null) {
      yield r'description';
      yield serializers.serialize(
        object.description,
        specifiedType: const FullType(String),
      );
    }
    yield r'aiType';
    yield serializers.serialize(
      object.aiType,
      specifiedType: const FullType(AiType),
    );
    if (object.url != null) {
      yield r'url';
      yield serializers.serialize(
        object.url,
        specifiedType: const FullType(String),
      );
    }
    if (object.key != null) {
      yield r'key';
      yield serializers.serialize(
        object.key,
        specifiedType: const FullType(String),
      );
    }
    if (object.model != null) {
      yield r'model';
      yield serializers.serialize(
        object.model,
        specifiedType: const FullType(String),
      );
    }
    if (object.isVisionModel != null) {
      yield r'isVisionModel';
      yield serializers.serialize(
        object.isVisionModel,
        specifiedType: const FullType(bool),
      );
    }
    yield r'ocrEngine';
    yield serializers.serialize(
      object.ocrEngine,
      specifiedType: const FullType(OcrEngine),
    );
    yield r'promptId';
    yield serializers.serialize(
      object.promptId,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertReceiptProcessingSettingsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertReceiptProcessingSettingsCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'description':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.description = valueDes;
          break;
        case r'aiType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(AiType),
          ) as AiType;
          result.aiType = valueDes;
          break;
        case r'url':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.url = valueDes;
          break;
        case r'key':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.key = valueDes;
          break;
        case r'model':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.model = valueDes;
          break;
        case r'isVisionModel':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isVisionModel = valueDes;
          break;
        case r'ocrEngine':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(OcrEngine),
          ) as OcrEngine;
          result.ocrEngine = valueDes;
          break;
        case r'promptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.promptId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertReceiptProcessingSettingsCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertReceiptProcessingSettingsCommandBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

