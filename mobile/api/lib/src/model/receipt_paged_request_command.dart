//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_paged_request_filter.dart';
import 'package:openapi/src/model/sort_direction.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_paged_request_command.g.dart';

/// ReceiptPagedRequestCommand
///
/// Properties:
/// * [page] - Page number
/// * [pageSize] - Number of records per page
/// * [orderBy] - field to order on
/// * [sortDirection] 
/// * [filter] 
/// * [fullReceipts] - Whether to include all receipt associations (receiptItems, comments, customFields, imageFiles, etc.)
@BuiltValue()
abstract class ReceiptPagedRequestCommand implements Built<ReceiptPagedRequestCommand, ReceiptPagedRequestCommandBuilder> {
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

  @BuiltValueField(wireName: r'filter')
  ReceiptPagedRequestFilter? get filter;

  /// Whether to include all receipt associations (receiptItems, comments, customFields, imageFiles, etc.)
  @BuiltValueField(wireName: r'fullReceipts')
  bool? get fullReceipts;

  ReceiptPagedRequestCommand._();

  factory ReceiptPagedRequestCommand([void updates(ReceiptPagedRequestCommandBuilder b)]) = _$ReceiptPagedRequestCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptPagedRequestCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptPagedRequestCommand> get serializer => _$ReceiptPagedRequestCommandSerializer();
}

class _$ReceiptPagedRequestCommandSerializer implements PrimitiveSerializer<ReceiptPagedRequestCommand> {
  @override
  final Iterable<Type> types = const [ReceiptPagedRequestCommand, _$ReceiptPagedRequestCommand];

  @override
  final String wireName = r'ReceiptPagedRequestCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptPagedRequestCommand object, {
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
    if (object.filter != null) {
      yield r'filter';
      yield serializers.serialize(
        object.filter,
        specifiedType: const FullType(ReceiptPagedRequestFilter),
      );
    }
    if (object.fullReceipts != null) {
      yield r'fullReceipts';
      yield serializers.serialize(
        object.fullReceipts,
        specifiedType: const FullType(bool),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptPagedRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptPagedRequestCommandBuilder result,
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
        case r'filter':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptPagedRequestFilter),
          ) as ReceiptPagedRequestFilter;
          result.filter.replace(valueDes);
          break;
        case r'fullReceipts':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.fullReceipts = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptPagedRequestCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptPagedRequestCommandBuilder();
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

