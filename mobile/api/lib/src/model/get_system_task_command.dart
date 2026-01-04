//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/associated_entity_type.dart';
import 'package:openapi/src/model/sort_direction.dart';
import 'package:openapi/src/model/paged_request_command.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'get_system_task_command.g.dart';

/// GetSystemTaskCommand
///
/// Properties:
/// * [associatedEntityId] - Associated entity id
/// * [associatedEntityType] 
/// * [page] - Page number
/// * [pageSize] - Number of records per page
/// * [orderBy] - field to order on
/// * [sortDirection] 
@BuiltValue()
abstract class GetSystemTaskCommand implements PagedRequestCommand, Built<GetSystemTaskCommand, GetSystemTaskCommandBuilder> {
  /// Associated entity id
  @BuiltValueField(wireName: r'associatedEntityId')
  int? get associatedEntityId;

  @BuiltValueField(wireName: r'associatedEntityType')
  AssociatedEntityType? get associatedEntityType;
  // enum associatedEntityTypeEnum {  NOOP_ENTITY_TYPE,  RECEIPT,  SYSTEM_EMAIL,  RECEIPT_PROCESSING_SETTINGS,  PROMPT,  };

  GetSystemTaskCommand._();

  factory GetSystemTaskCommand([void updates(GetSystemTaskCommandBuilder b)]) = _$GetSystemTaskCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GetSystemTaskCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<GetSystemTaskCommand> get serializer => _$GetSystemTaskCommandSerializer();
}

class _$GetSystemTaskCommandSerializer implements PrimitiveSerializer<GetSystemTaskCommand> {
  @override
  final Iterable<Type> types = const [GetSystemTaskCommand, _$GetSystemTaskCommand];

  @override
  final String wireName = r'GetSystemTaskCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GetSystemTaskCommand object, {
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
    if (object.associatedEntityId != null) {
      yield r'associatedEntityId';
      yield serializers.serialize(
        object.associatedEntityId,
        specifiedType: const FullType(int),
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
    if (object.associatedEntityType != null) {
      yield r'associatedEntityType';
      yield serializers.serialize(
        object.associatedEntityType,
        specifiedType: const FullType(AssociatedEntityType),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    GetSystemTaskCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GetSystemTaskCommandBuilder result,
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
        case r'associatedEntityId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.associatedEntityId = valueDes;
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
        case r'associatedEntityType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(AssociatedEntityType),
          ) as AssociatedEntityType;
          result.associatedEntityType = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  GetSystemTaskCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GetSystemTaskCommandBuilder();
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

