//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/widget_type.dart';
import 'package:built_value/json_object.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'widget.g.dart';

/// Widget related to a user's dashboard
///
/// Properties:
/// * [createdAt] 
/// * [createdBy] 
/// * [id] 
/// * [name] - Widget name
/// * [dashboardId] - Dashboard foreign key
/// * [updatedAt] 
/// * [widgetType] - Type of widget
/// * [configuration] - Configuration of widget
@BuiltValue()
abstract class Widget implements Built<Widget, WidgetBuilder> {
  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  @BuiltValueField(wireName: r'id')
  int get id;

  /// Widget name
  @BuiltValueField(wireName: r'name')
  String? get name;

  /// Dashboard foreign key
  @BuiltValueField(wireName: r'dashboardId')
  int get dashboardId;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// Type of widget
  @BuiltValueField(wireName: r'widgetType')
  WidgetType? get widgetType;
  // enum widgetTypeEnum {  GROUP_SUMMARY,  FILTERED_RECEIPTS,  GROUP_ACTIVITY,  PIE_CHART,  };

  /// Configuration of widget
  @BuiltValueField(wireName: r'configuration')
  BuiltMap<String, JsonObject?>? get configuration;

  Widget._();

  factory Widget([void updates(WidgetBuilder b)]) = _$Widget;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(WidgetBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Widget> get serializer => _$WidgetSerializer();
}

class _$WidgetSerializer implements PrimitiveSerializer<Widget> {
  @override
  final Iterable<Type> types = const [Widget, _$Widget];

  @override
  final String wireName = r'Widget';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Widget object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(String),
      );
    }
    yield r'dashboardId';
    yield serializers.serialize(
      object.dashboardId,
      specifiedType: const FullType(int),
    );
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.widgetType != null) {
      yield r'widgetType';
      yield serializers.serialize(
        object.widgetType,
        specifiedType: const FullType(WidgetType),
      );
    }
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
    Widget object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required WidgetBuilder result,
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
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
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
        case r'dashboardId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.dashboardId = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
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
  Widget deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = WidgetBuilder();
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

