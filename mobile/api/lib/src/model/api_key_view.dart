//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'api_key_view.g.dart';

/// ApiKeyView
///
/// Properties:
/// * [id] - API key ID
/// * [createdAt] - Creation timestamp
/// * [updatedAt] - Last update timestamp
/// * [createdBy] - ID of the user who created this API key
/// * [createdByString] - String representation of the creator
/// * [name] - API key name
/// * [description] - API key description
/// * [userId] - ID of the user who owns this API key
/// * [scope] - API key scope/permissions
/// * [lastUsedAt] - When the API key was last used
@BuiltValue()
abstract class ApiKeyView implements Built<ApiKeyView, ApiKeyViewBuilder> {
  /// API key ID
  @BuiltValueField(wireName: r'id')
  String? get id;

  /// Creation timestamp
  @BuiltValueField(wireName: r'createdAt')
  DateTime? get createdAt;

  /// Last update timestamp
  @BuiltValueField(wireName: r'updatedAt')
  DateTime? get updatedAt;

  /// ID of the user who created this API key
  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  /// String representation of the creator
  @BuiltValueField(wireName: r'createdByString')
  String? get createdByString;

  /// API key name
  @BuiltValueField(wireName: r'name')
  String? get name;

  /// API key description
  @BuiltValueField(wireName: r'description')
  String? get description;

  /// ID of the user who owns this API key
  @BuiltValueField(wireName: r'userId')
  int? get userId;

  /// API key scope/permissions
  @BuiltValueField(wireName: r'scope')
  String? get scope;

  /// When the API key was last used
  @BuiltValueField(wireName: r'lastUsedAt')
  DateTime? get lastUsedAt;

  ApiKeyView._();

  factory ApiKeyView([void updates(ApiKeyViewBuilder b)]) = _$ApiKeyView;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ApiKeyViewBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ApiKeyView> get serializer => _$ApiKeyViewSerializer();
}

class _$ApiKeyViewSerializer implements PrimitiveSerializer<ApiKeyView> {
  @override
  final Iterable<Type> types = const [ApiKeyView, _$ApiKeyView];

  @override
  final String wireName = r'ApiKeyView';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ApiKeyView object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.id != null) {
      yield r'id';
      yield serializers.serialize(
        object.id,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(DateTime),
      );
    }
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
      );
    }
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(String),
      );
    }
    if (object.description != null) {
      yield r'description';
      yield serializers.serialize(
        object.description,
        specifiedType: const FullType(String),
      );
    }
    if (object.userId != null) {
      yield r'userId';
      yield serializers.serialize(
        object.userId,
        specifiedType: const FullType(int),
      );
    }
    if (object.scope != null) {
      yield r'scope';
      yield serializers.serialize(
        object.scope,
        specifiedType: const FullType(String),
      );
    }
    if (object.lastUsedAt != null) {
      yield r'lastUsedAt';
      yield serializers.serialize(
        object.lastUsedAt,
        specifiedType: const FullType(DateTime),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ApiKeyView object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ApiKeyViewBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.id = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.createdAt = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.updatedAt = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
          break;
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
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        case r'scope':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.scope = valueDes;
          break;
        case r'lastUsedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.lastUsedAt = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ApiKeyView deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ApiKeyViewBuilder();
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

