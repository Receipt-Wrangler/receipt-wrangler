//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/base_model.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group_receipt_settings.g.dart';

/// GroupReceiptSettings
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
/// * [groupId] - Group foreign key
/// * [hideImages] - Hide receipt images
/// * [hideReceiptCategories] - Hide receipt categories
/// * [hideReceiptTags] - Hide receipt tags
/// * [hideItemCategories] - Hide receipt item categories
/// * [hideItemTags] - Hide receipt item tags
/// * [hideComments] - Hide receipt comments
/// * [hideShareCategories] - Hide share categories
/// * [hideShareTags] - Hide share tags
@BuiltValue()
abstract class GroupReceiptSettings implements BaseModel, Built<GroupReceiptSettings, GroupReceiptSettingsBuilder> {
  /// Hide receipt tags
  @BuiltValueField(wireName: r'hideReceiptTags')
  bool? get hideReceiptTags;

  /// Hide receipt categories
  @BuiltValueField(wireName: r'hideReceiptCategories')
  bool? get hideReceiptCategories;

  /// Group foreign key
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  /// Hide receipt item categories
  @BuiltValueField(wireName: r'hideItemCategories')
  bool? get hideItemCategories;

  /// Hide share tags
  @BuiltValueField(wireName: r'hideShareTags')
  bool? get hideShareTags;

  /// Hide receipt images
  @BuiltValueField(wireName: r'hideImages')
  bool? get hideImages;

  /// Hide receipt item tags
  @BuiltValueField(wireName: r'hideItemTags')
  bool? get hideItemTags;

  /// Hide receipt comments
  @BuiltValueField(wireName: r'hideComments')
  bool? get hideComments;

  /// Hide share categories
  @BuiltValueField(wireName: r'hideShareCategories')
  bool? get hideShareCategories;

  GroupReceiptSettings._();

  factory GroupReceiptSettings([void updates(GroupReceiptSettingsBuilder b)]) = _$GroupReceiptSettings;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GroupReceiptSettingsBuilder b) => b
      ..createdBy = 0
      ..createdByString = ''
      ..updatedAt = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<GroupReceiptSettings> get serializer => _$GroupReceiptSettingsSerializer();
}

class _$GroupReceiptSettingsSerializer implements PrimitiveSerializer<GroupReceiptSettings> {
  @override
  final Iterable<Type> types = const [GroupReceiptSettings, _$GroupReceiptSettings];

  @override
  final String wireName = r'GroupReceiptSettings';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GroupReceiptSettings object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
    if (object.hideImages != null) {
      yield r'hideImages';
      yield serializers.serialize(
        object.hideImages,
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
    if (object.hideReceiptTags != null) {
      yield r'hideReceiptTags';
      yield serializers.serialize(
        object.hideReceiptTags,
        specifiedType: const FullType(bool),
      );
    }
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    if (object.hideReceiptCategories != null) {
      yield r'hideReceiptCategories';
      yield serializers.serialize(
        object.hideReceiptCategories,
        specifiedType: const FullType(bool),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    if (object.hideItemCategories != null) {
      yield r'hideItemCategories';
      yield serializers.serialize(
        object.hideItemCategories,
        specifiedType: const FullType(bool),
      );
    }
    if (object.hideShareTags != null) {
      yield r'hideShareTags';
      yield serializers.serialize(
        object.hideShareTags,
        specifiedType: const FullType(bool),
      );
    }
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.hideItemTags != null) {
      yield r'hideItemTags';
      yield serializers.serialize(
        object.hideItemTags,
        specifiedType: const FullType(bool),
      );
    }
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
      );
    }
    if (object.hideShareCategories != null) {
      yield r'hideShareCategories';
      yield serializers.serialize(
        object.hideShareCategories,
        specifiedType: const FullType(bool),
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
    GroupReceiptSettings object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GroupReceiptSettingsBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'hideImages':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideImages = valueDes;
          break;
        case r'hideComments':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideComments = valueDes;
          break;
        case r'hideReceiptTags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideReceiptTags = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'hideReceiptCategories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideReceiptCategories = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'hideItemCategories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideItemCategories = valueDes;
          break;
        case r'hideShareTags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideShareTags = valueDes;
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'hideItemTags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideItemTags = valueDes;
          break;
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
          break;
        case r'hideShareCategories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hideShareCategories = valueDes;
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
  GroupReceiptSettings deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GroupReceiptSettingsBuilder();
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

