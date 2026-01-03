//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/user_role.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'user_view.g.dart';

/// User in the system
///
/// Properties:
/// * [username] - User's username used to login
/// * [createdAt] 
/// * [createdBy] 
/// * [defaultAvatarColor] - Default avatar color
/// * [displayName] - Display name
/// * [id] 
/// * [isDummyUser] - Is dummy user
/// * [updatedAt] 
/// * [userRole] - User's role
@BuiltValue()
abstract class UserView implements Built<UserView, UserViewBuilder> {
  /// User's username used to login
  @BuiltValueField(wireName: r'username')
  String get username;

  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  /// Default avatar color
  @BuiltValueField(wireName: r'defaultAvatarColor')
  String? get defaultAvatarColor;

  /// Display name
  @BuiltValueField(wireName: r'displayName')
  String get displayName;

  @BuiltValueField(wireName: r'id')
  int get id;

  /// Is dummy user
  @BuiltValueField(wireName: r'isDummyUser')
  bool get isDummyUser;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// User's role
  @BuiltValueField(wireName: r'userRole')
  UserRole get userRole;
  // enum userRoleEnum {  ADMIN,  USER,  };

  UserView._();

  factory UserView([void updates(UserViewBuilder b)]) = _$UserView;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UserViewBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UserView> get serializer => _$UserViewSerializer();
}

class _$UserViewSerializer implements PrimitiveSerializer<UserView> {
  @override
  final Iterable<Type> types = const [UserView, _$UserView];

  @override
  final String wireName = r'UserView';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UserView object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'username';
    yield serializers.serialize(
      object.username,
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
    if (object.defaultAvatarColor != null) {
      yield r'defaultAvatarColor';
      yield serializers.serialize(
        object.defaultAvatarColor,
        specifiedType: const FullType(String),
      );
    }
    yield r'displayName';
    yield serializers.serialize(
      object.displayName,
      specifiedType: const FullType(String),
    );
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'isDummyUser';
    yield serializers.serialize(
      object.isDummyUser,
      specifiedType: const FullType(bool),
    );
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
    yield r'userRole';
    yield serializers.serialize(
      object.userRole,
      specifiedType: const FullType(UserRole),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UserView object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UserViewBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'username':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.username = valueDes;
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
        case r'defaultAvatarColor':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.defaultAvatarColor = valueDes;
          break;
        case r'displayName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.displayName = valueDes;
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'isDummyUser':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isDummyUser = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        case r'userRole':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(UserRole),
          ) as UserRole;
          result.userRole = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UserView deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UserViewBuilder();
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

