//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/associated_group.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group_filter.g.dart';

/// GroupFilter
///
/// Properties:
/// * [associatedGroup] 
@BuiltValue()
abstract class GroupFilter implements Built<GroupFilter, GroupFilterBuilder> {
  @BuiltValueField(wireName: r'associatedGroup')
  AssociatedGroup? get associatedGroup;
  // enum associatedGroupEnum {  MINE,  ALL,  };

  GroupFilter._();

  factory GroupFilter([void updates(GroupFilterBuilder b)]) = _$GroupFilter;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GroupFilterBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<GroupFilter> get serializer => _$GroupFilterSerializer();
}

class _$GroupFilterSerializer implements PrimitiveSerializer<GroupFilter> {
  @override
  final Iterable<Type> types = const [GroupFilter, _$GroupFilter];

  @override
  final String wireName = r'GroupFilter';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GroupFilter object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.associatedGroup != null) {
      yield r'associatedGroup';
      yield serializers.serialize(
        object.associatedGroup,
        specifiedType: const FullType(AssociatedGroup),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    GroupFilter object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GroupFilterBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'associatedGroup':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(AssociatedGroup),
          ) as AssociatedGroup;
          result.associatedGroup = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  GroupFilter deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GroupFilterBuilder();
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

