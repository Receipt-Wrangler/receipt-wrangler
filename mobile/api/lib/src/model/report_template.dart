//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/base_model.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/report_request_command.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_template.g.dart';

/// ReportTemplate
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
/// * [name] - The template name (mirrors the saved report's name).
/// * [configuration] 
/// * [configurationVersion] - Schema version the stored configuration was written under.
/// * [allowedActions] - The actions the requesting user may perform on this template (read, generate, update, delete, duplicate), resolved per user and populated only on the list response. Drives the row action buttons.
@BuiltValue()
abstract class ReportTemplate implements BaseModel, Built<ReportTemplate, ReportTemplateBuilder> {
  /// Schema version the stored configuration was written under.
  @BuiltValueField(wireName: r'configurationVersion')
  int get configurationVersion;

  @BuiltValueField(wireName: r'configuration')
  ReportRequestCommand get configuration;

  /// The template name (mirrors the saved report's name).
  @BuiltValueField(wireName: r'name')
  String get name;

  /// The actions the requesting user may perform on this template (read, generate, update, delete, duplicate), resolved per user and populated only on the list response. Drives the row action buttons.
  @BuiltValueField(wireName: r'allowedActions')
  BuiltList<String>? get allowedActions;

  ReportTemplate._();

  factory ReportTemplate([void updates(ReportTemplateBuilder b)]) = _$ReportTemplate;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportTemplateBuilder b) => b
      ..createdBy = 0
      ..createdByString = ''
      ..updatedAt = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportTemplate> get serializer => _$ReportTemplateSerializer();
}

class _$ReportTemplateSerializer implements PrimitiveSerializer<ReportTemplate> {
  @override
  final Iterable<Type> types = const [ReportTemplate, _$ReportTemplate];

  @override
  final String wireName = r'ReportTemplate';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportTemplate object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    yield r'configurationVersion';
    yield serializers.serialize(
      object.configurationVersion,
      specifiedType: const FullType(int),
    );
    yield r'configuration';
    yield serializers.serialize(
      object.configuration,
      specifiedType: const FullType(ReportRequestCommand),
    );
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.allowedActions != null) {
      yield r'allowedActions';
      yield serializers.serialize(
        object.allowedActions,
        specifiedType: const FullType(BuiltList, [FullType(String)]),
      );
    }
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
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
    ReportTemplate object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportTemplateBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'configurationVersion':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.configurationVersion = valueDes;
          break;
        case r'configuration':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportRequestCommand),
          ) as ReportRequestCommand;
          result.configuration.replace(valueDes);
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'allowedActions':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(String)]),
          ) as BuiltList<String>;
          result.allowedActions.replace(valueDes);
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
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
  ReportTemplate deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportTemplateBuilder();
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

