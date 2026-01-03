//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/sort_direction.dart';
import 'package:openapi/src/model/paged_request_command.dart';
import 'package:openapi/src/model/group_filter.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'paged_group_request_command.g.dart';

/// PagedGroupRequestCommand
///
/// Properties:
/// * [page] - Page number
/// * [pageSize] - Number of records per page
/// * [orderBy] - field to order on
/// * [sortDirection] 
/// * [filter] 
@BuiltValue()
abstract class PagedGroupRequestCommand implements PagedRequestCommand, Built<PagedGroupRequestCommand, PagedGroupRequestCommandBuilder> {
  @BuiltValueField(wireName: r'filter')
  GroupFilter? get filter;

  PagedGroupRequestCommand._();

  factory PagedGroupRequestCommand([void updates(PagedGroupRequestCommandBuilder b)]) = _$PagedGroupRequestCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PagedGroupRequestCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedGroupRequestCommand> get serializer => _$PagedGroupRequestCommandSerializer();
}

class _$PagedGroupRequestCommandSerializer implements PrimitiveSerializer<PagedGroupRequestCommand> {
  @override
  final Iterable<Type> types = const [PagedGroupRequestCommand, _$PagedGroupRequestCommand];

  @override
  final String wireName = r'PagedGroupRequestCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedGroupRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.filter != null) {
      yield r'filter';
      yield serializers.serialize(
        object.filter,
        specifiedType: const FullType(GroupFilter),
      );
    }
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
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedGroupRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PagedGroupRequestCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'filter':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(GroupFilter),
          ) as GroupFilter;
          result.filter.replace(valueDes);
          break;
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
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PagedGroupRequestCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PagedGroupRequestCommandBuilder();
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

