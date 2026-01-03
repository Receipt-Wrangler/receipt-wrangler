//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/sort_direction.dart';
import 'package:openapi/src/model/paged_request_command.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'paged_activity_request_command.g.dart';

/// PagedActivityRequestCommand
///
/// Properties:
/// * [page] - Page number
/// * [pageSize] - Number of records per page
/// * [orderBy] - field to order on
/// * [sortDirection] 
/// * [groupIds] 
@BuiltValue()
abstract class PagedActivityRequestCommand implements PagedRequestCommand, Built<PagedActivityRequestCommand, PagedActivityRequestCommandBuilder> {
  @BuiltValueField(wireName: r'groupIds')
  BuiltList<int>? get groupIds;

  PagedActivityRequestCommand._();

  factory PagedActivityRequestCommand([void updates(PagedActivityRequestCommandBuilder b)]) = _$PagedActivityRequestCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PagedActivityRequestCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedActivityRequestCommand> get serializer => _$PagedActivityRequestCommandSerializer();
}

class _$PagedActivityRequestCommandSerializer implements PrimitiveSerializer<PagedActivityRequestCommand> {
  @override
  final Iterable<Type> types = const [PagedActivityRequestCommand, _$PagedActivityRequestCommand];

  @override
  final String wireName = r'PagedActivityRequestCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedActivityRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'pageSize';
    yield serializers.serialize(
      object.pageSize,
      specifiedType: const FullType(int),
    );
    if (object.orderBy != null) {
      yield r'orderBy';
      yield serializers.serialize(
        object.orderBy,
        specifiedType: const FullType(String),
      );
    }
    if (object.sortDirection != null) {
      yield r'sortDirection';
      yield serializers.serialize(
        object.sortDirection,
        specifiedType: const FullType(SortDirection),
      );
    }
    yield r'page';
    yield serializers.serialize(
      object.page,
      specifiedType: const FullType(int),
    );
    if (object.groupIds != null) {
      yield r'groupIds';
      yield serializers.serialize(
        object.groupIds,
        specifiedType: const FullType(BuiltList, [FullType(int)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedActivityRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PagedActivityRequestCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'pageSize':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.pageSize = valueDes;
          break;
        case r'orderBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.orderBy = valueDes;
          break;
        case r'sortDirection':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(SortDirection),
          ) as SortDirection;
          result.sortDirection = valueDes;
          break;
        case r'page':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.page = valueDes;
          break;
        case r'groupIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.groupIds.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PagedActivityRequestCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PagedActivityRequestCommandBuilder();
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

