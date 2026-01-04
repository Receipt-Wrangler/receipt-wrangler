//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/ai_type.dart';
import 'package:openapi/src/model/ocr_engine.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'check_receipt_processing_settings_connectivity_command.g.dart';

/// CheckReceiptProcessingSettingsConnectivityCommand
///
/// Properties:
/// * [id] - Receipt processing settings id
/// * [name] - Name of the settings
/// * [aiType] 
/// * [url] - URL for custom endpoints
/// * [key] - Key for endpoints that require authentication
/// * [model] - LLM model
/// * [numWorkers] - Number of workers to use
/// * [ocrEngine] 
/// * [promptId] - Prompt foreign key
@BuiltValue()
abstract class CheckReceiptProcessingSettingsConnectivityCommand implements Built<CheckReceiptProcessingSettingsConnectivityCommand, CheckReceiptProcessingSettingsConnectivityCommandBuilder> {
  /// Receipt processing settings id
  @BuiltValueField(wireName: r'id')
  int? get id;

  /// Name of the settings
  @BuiltValueField(wireName: r'name')
  String? get name;

  @BuiltValueField(wireName: r'aiType')
  AiType? get aiType;
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

  /// Number of workers to use
  @BuiltValueField(wireName: r'numWorkers')
  int? get numWorkers;

  @BuiltValueField(wireName: r'ocrEngine')
  OcrEngine? get ocrEngine;
  // enum ocrEngineEnum {  TESSERACT,  EASY_OCR,  };

  /// Prompt foreign key
  @BuiltValueField(wireName: r'promptId')
  int? get promptId;

  CheckReceiptProcessingSettingsConnectivityCommand._();

  factory CheckReceiptProcessingSettingsConnectivityCommand([void updates(CheckReceiptProcessingSettingsConnectivityCommandBuilder b)]) = _$CheckReceiptProcessingSettingsConnectivityCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CheckReceiptProcessingSettingsConnectivityCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<CheckReceiptProcessingSettingsConnectivityCommand> get serializer => _$CheckReceiptProcessingSettingsConnectivityCommandSerializer();
}

class _$CheckReceiptProcessingSettingsConnectivityCommandSerializer implements PrimitiveSerializer<CheckReceiptProcessingSettingsConnectivityCommand> {
  @override
  final Iterable<Type> types = const [CheckReceiptProcessingSettingsConnectivityCommand, _$CheckReceiptProcessingSettingsConnectivityCommand];

  @override
  final String wireName = r'CheckReceiptProcessingSettingsConnectivityCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CheckReceiptProcessingSettingsConnectivityCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.id != null) {
      yield r'id';
      yield serializers.serialize(
        object.id,
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
    if (object.aiType != null) {
      yield r'aiType';
      yield serializers.serialize(
        object.aiType,
        specifiedType: const FullType(AiType),
      );
    }
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
    if (object.numWorkers != null) {
      yield r'numWorkers';
      yield serializers.serialize(
        object.numWorkers,
        specifiedType: const FullType(int),
      );
    }
    if (object.ocrEngine != null) {
      yield r'ocrEngine';
      yield serializers.serialize(
        object.ocrEngine,
        specifiedType: const FullType(OcrEngine),
      );
    }
    if (object.promptId != null) {
      yield r'promptId';
      yield serializers.serialize(
        object.promptId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    CheckReceiptProcessingSettingsConnectivityCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CheckReceiptProcessingSettingsConnectivityCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
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
        case r'numWorkers':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.numWorkers = valueDes;
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
  CheckReceiptProcessingSettingsConnectivityCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CheckReceiptProcessingSettingsConnectivityCommandBuilder();
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

