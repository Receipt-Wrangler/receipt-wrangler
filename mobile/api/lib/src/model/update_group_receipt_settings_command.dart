//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'update_group_receipt_settings_command.g.dart';

/// UpdateGroupReceiptSettingsCommand
///
/// Properties:
/// * [hideImages] - Hide receipt images
/// * [hideReceiptCategories] - Hide receipt categories
/// * [hideReceiptTags] - Hide receipt tags
/// * [hideItemCategories] - Hide receipt item categories
/// * [hideItemTags] - Hide receipt item tags
/// * [hideComments] - Hide receipt comments
@BuiltValue()
abstract class UpdateGroupReceiptSettingsCommand implements Built<UpdateGroupReceiptSettingsCommand, UpdateGroupReceiptSettingsCommandBuilder> {
  /// Hide receipt images
  @BuiltValueField(wireName: r'hideImages')
  bool? get hideImages;

  /// Hide receipt categories
  @BuiltValueField(wireName: r'hideReceiptCategories')
  bool? get hideReceiptCategories;

  /// Hide receipt tags
  @BuiltValueField(wireName: r'hideReceiptTags')
  bool? get hideReceiptTags;

  /// Hide receipt item categories
  @BuiltValueField(wireName: r'hideItemCategories')
  bool? get hideItemCategories;

  /// Hide receipt item tags
  @BuiltValueField(wireName: r'hideItemTags')
  bool? get hideItemTags;

  /// Hide receipt comments
  @BuiltValueField(wireName: r'hideComments')
  bool? get hideComments;

  UpdateGroupReceiptSettingsCommand._();

  factory UpdateGroupReceiptSettingsCommand([void updates(UpdateGroupReceiptSettingsCommandBuilder b)]) = _$UpdateGroupReceiptSettingsCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpdateGroupReceiptSettingsCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpdateGroupReceiptSettingsCommand> get serializer => _$UpdateGroupReceiptSettingsCommandSerializer();
}

class _$UpdateGroupReceiptSettingsCommandSerializer implements PrimitiveSerializer<UpdateGroupReceiptSettingsCommand> {
  @override
  final Iterable<Type> types = const [UpdateGroupReceiptSettingsCommand, _$UpdateGroupReceiptSettingsCommand];

  @override
  final String wireName = r'UpdateGroupReceiptSettingsCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpdateGroupReceiptSettingsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.hideImages != null) {
      yield r'hideImages';
      yield serializers.serialize(
        object.hideImages,
        specifiedType: const FullType(bool),
      );
    }
    if (object.hideReceiptCategories != null) {
      yield r'hideReceiptCategories';
      yield serializers.serialize(
        object.hideReceiptCategories,
        specifiedType: const FullType(bool),
      );
    }
    if (object.hideReceiptTags != null) {
      yield r'hideReceiptTags';
      yield serializers.serialize(
        object.hideReceiptTags,
        specifiedType: const FullType(bool),
      );
    }
    if (object.hideItemCategories != null) {
      yield r'hideItemCategories';
      yield serializers.serialize(
        object.hideItemCategories,
        specifiedType: const FullType(bool),
      );
    }
    if (object.hideItemTags != null) {
      yield r'hideItemTags';
      yield serializers.serialize(
        object.hideItemTags,
        specifiedType: const FullType(bool),
      );
    }
    if (object.hideComments != null) {
      yield r'hideComments';
      yield serializers.serialize(
        object.hideComments,
        specifiedType: const FullType(bool),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpdateGroupReceiptSettingsCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpdateGroupReceiptSettingsCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'hideImages':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideImages = valueDes;
          break;
        case r'hideReceiptCategories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideReceiptCategories = valueDes;
          break;
        case r'hideReceiptTags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideReceiptTags = valueDes;
          break;
        case r'hideItemCategories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideItemCategories = valueDes;
          break;
        case r'hideItemTags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideItemTags = valueDes;
          break;
        case r'hideComments':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideComments = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpdateGroupReceiptSettingsCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpdateGroupReceiptSettingsCommandBuilder();
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

