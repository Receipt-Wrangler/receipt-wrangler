//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/widget.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'dashboard.g.dart';

/// Dashboard for a user
///
/// Properties:
/// * [createdAt] 
/// * [createdBy] 
/// * [id] 
/// * [name] - Dashboard name
/// * [groupId] - Group foreign key
/// * [userId] - User foreign key
/// * [updatedAt] 
/// * [widgets] - Widgets associated to dashboard
@BuiltValue()
abstract class Dashboard implements Built<Dashboard, DashboardBuilder> {
  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  @BuiltValueField(wireName: r'id')
  int get id;

  /// Dashboard name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Group foreign key
  @BuiltValueField(wireName: r'groupId')
  int? get groupId;

  /// User foreign key
  @BuiltValueField(wireName: r'userId')
  int get userId;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// Widgets associated to dashboard
  @BuiltValueField(wireName: r'widgets')
  BuiltList<Widget>? get widgets;

  Dashboard._();

  factory Dashboard([void updates(DashboardBuilder b)]) = _$Dashboard;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(DashboardBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Dashboard> get serializer => _$DashboardSerializer();
}

class _$DashboardSerializer implements PrimitiveSerializer<Dashboard> {
  @override
  final Iterable<Type> types = const [Dashboard, _$Dashboard];

  @override
  final String wireName = r'Dashboard';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Dashboard object, {
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
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.groupId != null) {
      yield r'groupId';
      yield serializers.serialize(
        object.groupId,
        specifiedType: const FullType(int),
      );
    }
    yield r'userId';
    yield serializers.serialize(
      object.userId,
      specifiedType: const FullType(int),
    );
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.widgets != null) {
      yield r'widgets';
      yield serializers.serialize(
        object.widgets,
        specifiedType: const FullType(BuiltList, [FullType(Widget)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    Dashboard object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required DashboardBuilder result,
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
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        case r'widgets':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Widget)]),
          ) as BuiltList<Widget>;
          result.widgets.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Dashboard deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = DashboardBuilder();
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

