//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/sort_direction.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'paged_request_command.g.dart';

/// PagedRequestCommand
///
/// Properties:
/// * [page] - Page number
/// * [pageSize] - Number of records per page
/// * [orderBy] - field to order on
/// * [sortDirection] 
@BuiltValue(instantiable: false)
abstract class PagedRequestCommand  {
  /// Page number
  @BuiltValueField(wireName: r'page')
  int get page;

  /// Number of records per page
  @BuiltValueField(wireName: r'pageSize')
  int get pageSize;

  /// field to order on
  @BuiltValueField(wireName: r'orderBy')
  String? get orderBy;

  @BuiltValueField(wireName: r'sortDirection')
  SortDirection? get sortDirection;
  // enum sortDirectionEnum {  asc,  desc,  ,  };

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedRequestCommand> get serializer => _$PagedRequestCommandSerializer();
}

class _$PagedRequestCommandSerializer implements PrimitiveSerializer<PagedRequestCommand> {
  @override
  final Iterable<Type> types = const [PagedRequestCommand];

  @override
  final String wireName = r'PagedRequestCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'page';
    yield serializers.serialize(
      object.page,
      specifiedType: const FullType(int),
    );
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
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  @override
  PagedRequestCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return serializers.deserialize(serialized, specifiedType: FullType($PagedRequestCommand)) as $PagedRequestCommand;
  }
}

/// a concrete implementation of [PagedRequestCommand], since [PagedRequestCommand] is not instantiable
@BuiltValue(instantiable: true)
abstract class $PagedRequestCommand implements PagedRequestCommand, Built<$PagedRequestCommand, $PagedRequestCommandBuilder> {
  $PagedRequestCommand._();

  factory $PagedRequestCommand([void Function($PagedRequestCommandBuilder)? updates]) = _$$PagedRequestCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults($PagedRequestCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<$PagedRequestCommand> get serializer => _$$PagedRequestCommandSerializer();
}

class _$$PagedRequestCommandSerializer implements PrimitiveSerializer<$PagedRequestCommand> {
  @override
  final Iterable<Type> types = const [$PagedRequestCommand, _$$PagedRequestCommand];

  @override
  final String wireName = r'$PagedRequestCommand';

  @override
  Object serialize(
    Serializers serializers,
    $PagedRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return serializers.serialize(object, specifiedType: FullType(PagedRequestCommand))!;
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PagedRequestCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'page':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.page = valueDes;
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
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  $PagedRequestCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = $PagedRequestCommandBuilder();
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

