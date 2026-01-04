//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group_settings_white_list_email.g.dart';

/// GroupSettingsWhiteListEmail
///
/// Properties:
/// * [id] - Group settings email id
/// * [groupSettingsId] - Group settings foreign key
/// * [email] - Email to match
/// * [createdAt] 
/// * [createdBy] 
/// * [updatedAt] 
@BuiltValue()
abstract class GroupSettingsWhiteListEmail implements Built<GroupSettingsWhiteListEmail, GroupSettingsWhiteListEmailBuilder> {
  /// Group settings email id
  @BuiltValueField(wireName: r'id')
  int get id;

  /// Group settings foreign key
  @BuiltValueField(wireName: r'groupSettingsId')
  int get groupSettingsId;

  /// Email to match
  @BuiltValueField(wireName: r'email')
  String get email;

  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  GroupSettingsWhiteListEmail._();

  factory GroupSettingsWhiteListEmail([void updates(GroupSettingsWhiteListEmailBuilder b)]) = _$GroupSettingsWhiteListEmail;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GroupSettingsWhiteListEmailBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<GroupSettingsWhiteListEmail> get serializer => _$GroupSettingsWhiteListEmailSerializer();
}

class _$GroupSettingsWhiteListEmailSerializer implements PrimitiveSerializer<GroupSettingsWhiteListEmail> {
  @override
  final Iterable<Type> types = const [GroupSettingsWhiteListEmail, _$GroupSettingsWhiteListEmail];

  @override
  final String wireName = r'GroupSettingsWhiteListEmail';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GroupSettingsWhiteListEmail object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'groupSettingsId';
    yield serializers.serialize(
      object.groupSettingsId,
      specifiedType: const FullType(int),
    );
    yield r'email';
    yield serializers.serialize(
      object.email,
      specifiedType: const FullType(String),
    );
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
    GroupSettingsWhiteListEmail object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GroupSettingsWhiteListEmailBuilder result,
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
        case r'groupSettingsId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupSettingsId = valueDes;
          break;
        case r'email':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.email = valueDes;
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
  GroupSettingsWhiteListEmail deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GroupSettingsWhiteListEmailBuilder();
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

