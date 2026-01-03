//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'update_profile_command.g.dart';

/// Command to update user's profile
///
/// Properties:
/// * [displayName] - User's displayName
/// * [defaultAvatarColor] - Color of default avatar
@BuiltValue()
abstract class UpdateProfileCommand implements Built<UpdateProfileCommand, UpdateProfileCommandBuilder> {
  /// User's displayName
  @BuiltValueField(wireName: r'displayName')
  String get displayName;

  /// Color of default avatar
  @BuiltValueField(wireName: r'defaultAvatarColor')
  String get defaultAvatarColor;

  UpdateProfileCommand._();

  factory UpdateProfileCommand([void updates(UpdateProfileCommandBuilder b)]) = _$UpdateProfileCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpdateProfileCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpdateProfileCommand> get serializer => _$UpdateProfileCommandSerializer();
}

class _$UpdateProfileCommandSerializer implements PrimitiveSerializer<UpdateProfileCommand> {
  @override
  final Iterable<Type> types = const [UpdateProfileCommand, _$UpdateProfileCommand];

  @override
  final String wireName = r'UpdateProfileCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpdateProfileCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'displayName';
    yield serializers.serialize(
      object.displayName,
      specifiedType: const FullType(String),
    );
    yield r'defaultAvatarColor';
    yield serializers.serialize(
      object.defaultAvatarColor,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UpdateProfileCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpdateProfileCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'displayName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.displayName = valueDes;
          break;
        case r'defaultAvatarColor':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.defaultAvatarColor = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpdateProfileCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpdateProfileCommandBuilder();
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

