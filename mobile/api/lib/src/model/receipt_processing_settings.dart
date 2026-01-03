//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/base_model.dart';
import 'package:openapi/src/model/ai_type.dart';
import 'package:openapi/src/model/prompt.dart';
import 'package:openapi/src/model/ocr_engine.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_processing_settings.g.dart';

/// ReceiptProcessingSettings
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
/// * [name] - Name of the settings
/// * [description] - Description of the settings
/// * [aiType] 
/// * [url] - URL for custom endpoints
/// * [key] - Key for endpoints that require authentication
/// * [model] - LLM model
/// * [isVisionModel] - Is vision model
/// * [ocrEngine] 
/// * [prompt] 
/// * [promptId] - Prompt foreign key
@BuiltValue()
abstract class ReceiptProcessingSettings implements BaseModel, Built<ReceiptProcessingSettings, ReceiptProcessingSettingsBuilder> {
  @BuiltValueField(wireName: r'ocrEngine')
  OcrEngine? get ocrEngine;
  // enum ocrEngineEnum {  TESSERACT,  EASY_OCR,  };

  /// Is vision model
  @BuiltValueField(wireName: r'isVisionModel')
  bool? get isVisionModel;

  @BuiltValueField(wireName: r'aiType')
  AiType? get aiType;
  // enum aiTypeEnum {  OPEN_AI_CUSTOM,  OPEN_AI,  GEMINI,  OLLAMA,  };

  /// Prompt foreign key
  @BuiltValueField(wireName: r'promptId')
  int? get promptId;

  /// Name of the settings
  @BuiltValueField(wireName: r'name')
  String? get name;

  /// Description of the settings
  @BuiltValueField(wireName: r'description')
  String? get description;

  /// LLM model
  @BuiltValueField(wireName: r'model')
  String? get model;

  @BuiltValueField(wireName: r'prompt')
  Prompt? get prompt;

  /// URL for custom endpoints
  @BuiltValueField(wireName: r'url')
  String? get url;

  /// Key for endpoints that require authentication
  @BuiltValueField(wireName: r'key')
  String? get key;

  ReceiptProcessingSettings._();

  factory ReceiptProcessingSettings([void updates(ReceiptProcessingSettingsBuilder b)]) = _$ReceiptProcessingSettings;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptProcessingSettingsBuilder b) => b
      ..createdBy = 0
      ..createdByString = ''
      ..updatedAt = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptProcessingSettings> get serializer => _$ReceiptProcessingSettingsSerializer();
}

class _$ReceiptProcessingSettingsSerializer implements PrimitiveSerializer<ReceiptProcessingSettings> {
  @override
  final Iterable<Type> types = const [ReceiptProcessingSettings, _$ReceiptProcessingSettings];

  @override
  final String wireName = r'ReceiptProcessingSettings';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptProcessingSettings object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.ocrEngine != null) {
      yield r'ocrEngine';
      yield serializers.serialize(
        object.ocrEngine,
        specifiedType: const FullType(OcrEngine),
      );
    }
    if (object.isVisionModel != null) {
      yield r'isVisionModel';
      yield serializers.serialize(
        object.isVisionModel,
        specifiedType: const FullType(bool),
      );
    }
    if (object.description != null) {
      yield r'description';
      yield serializers.serialize(
        object.description,
        specifiedType: const FullType(String),
      );
    }
    if (object.url != null) {
      yield r'url';
      yield serializers.serialize(
        object.url,
        specifiedType: const FullType(String),
      );
    }
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    if (object.aiType != null) {
      yield r'aiType';
      yield serializers.serialize(
        object.aiType,
        specifiedType: const FullType(AiType),
      );
    }
    if (object.promptId != null) {
      yield r'promptId';
      yield serializers.serialize(
        object.promptId,
        specifiedType: const FullType(int),
      );
    }
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
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
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.prompt != null) {
      yield r'prompt';
      yield serializers.serialize(
        object.prompt,
        specifiedType: const FullType(Prompt),
      );
    }
    if (object.key != null) {
      yield r'key';
      yield serializers.serialize(
        object.key,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
      );
    }
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptProcessingSettings object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptProcessingSettingsBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'ocrEngine':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(OcrEngine),
          ) as OcrEngine;
          result.ocrEngine = valueDes;
          break;
        case r'isVisionModel':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isVisionModel = valueDes;
          break;
        case r'description':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.description = valueDes;
          break;
        case r'url':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.url = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'aiType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(AiType),
          ) as AiType;
          result.aiType = valueDes;
          break;
        case r'promptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.promptId = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'model':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.model = valueDes;
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'prompt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(Prompt),
          ) as Prompt;
          result.prompt.replace(valueDes);
          break;
        case r'key':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.key = valueDes;
          break;
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptProcessingSettings deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptProcessingSettingsBuilder();
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

