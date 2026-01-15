//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/widget_type.dart';
import 'package:built_value/json_object.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_widget_command.g.dart';

/// UpsertWidgetCommand
///
/// Properties:
/// * [name] - Widget name
/// * [widgetType] - Type of widget
/// * [configuration] - Configuration of widget
@BuiltValue()
abstract class UpsertWidgetCommand implements Built<UpsertWidgetCommand, UpsertWidgetCommandBuilder> {
  /// Widget name
  @BuiltValueField(wireName: r'name')
  String? get name;

  /// Type of widget
  @BuiltValueField(wireName: r'widgetType')
  WidgetType get widgetType;
  // enum widgetTypeEnum {  GROUP_SUMMARY,  FILTERED_RECEIPTS,  GROUP_ACTIVITY,  PIE_CHART,  };

  /// Configuration of widget
  @BuiltValueField(wireName: r'configuration')
  BuiltMap<String, JsonObject?>? get configuration;

  UpsertWidgetCommand._();

  factory UpsertWidgetCommand([void updates(UpsertWidgetCommandBuilder b)]) = _$UpsertWidgetCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertWidgetCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertWidgetCommand> get serializer => _$UpsertWidgetCommandSerializer();
}

class _$UpsertWidgetCommandSerializer implements PrimitiveSerializer<UpsertWidgetCommand> {
  @override
  final Iterable<Type> types = const [UpsertWidgetCommand, _$UpsertWidgetCommand];

  @override
  final String wireName = r'UpsertWidgetCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertWidgetCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(String),
      );
    }
    yield r'widgetType';
    yield serializers.serialize(
      object.widgetType,
      specifiedType: const FullType(WidgetType),
    );
    if (object.configuration != null) {
      yield r'configuration';
      yield serializers.serialize(
        object.configuration,
        specifiedType: const FullType(BuiltMap, [FullType(String), FullType.nullable(JsonObject)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertWidgetCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertWidgetCommandBuilder result,
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
        case r'widgetType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(WidgetType),
          ) as WidgetType;
          result.widgetType = valueDes;
          break;
        case r'configuration':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltMap, [FullType(String), FullType.nullable(JsonObject)]),
          ) as BuiltMap<String, JsonObject?>;
          result.configuration.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertWidgetCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertWidgetCommandBuilder();
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

