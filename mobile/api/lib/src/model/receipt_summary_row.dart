//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_status.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/receipt_summary_custom_field_total.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary_row.g.dart';

/// ReceiptSummaryRow
///
/// Properties:
/// * [status] 
/// * [receiptCount] - How many receipts this row covers. On the overall row this matches the table's own total count for the same filter.
/// * [total] - Sum of the receipts' amounts
/// * [customFieldTotals] - One entry per configured currency custom field, in the configured order. Always present; empty when the group totals only the receipt amount.
@BuiltValue()
abstract class ReceiptSummaryRow implements Built<ReceiptSummaryRow, ReceiptSummaryRowBuilder> {
  @BuiltValueField(wireName: r'status')
  ReceiptStatus get status;
  // enum statusEnum {  OPEN,  NEEDS_ATTENTION,  RESOLVED,  DRAFT,  DECLINED,  ,  };

  /// How many receipts this row covers. On the overall row this matches the table's own total count for the same filter.
  @BuiltValueField(wireName: r'receiptCount')
  int get receiptCount;

  /// Sum of the receipts' amounts
  @BuiltValueField(wireName: r'total')
  String get total;

  /// One entry per configured currency custom field, in the configured order. Always present; empty when the group totals only the receipt amount.
  @BuiltValueField(wireName: r'customFieldTotals')
  BuiltList<ReceiptSummaryCustomFieldTotal> get customFieldTotals;

  ReceiptSummaryRow._();

  factory ReceiptSummaryRow([void updates(ReceiptSummaryRowBuilder b)]) = _$ReceiptSummaryRow;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptSummaryRowBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptSummaryRow> get serializer => _$ReceiptSummaryRowSerializer();
}

class _$ReceiptSummaryRowSerializer implements PrimitiveSerializer<ReceiptSummaryRow> {
  @override
  final Iterable<Type> types = const [ReceiptSummaryRow, _$ReceiptSummaryRow];

  @override
  final String wireName = r'ReceiptSummaryRow';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptSummaryRow object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(ReceiptStatus),
    );
    yield r'receiptCount';
    yield serializers.serialize(
      object.receiptCount,
      specifiedType: const FullType(int),
    );
    yield r'total';
    yield serializers.serialize(
      object.total,
      specifiedType: const FullType(String),
    );
    yield r'customFieldTotals';
    yield serializers.serialize(
      object.customFieldTotals,
      specifiedType: const FullType(BuiltList, [FullType(ReceiptSummaryCustomFieldTotal)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptSummaryRow object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptSummaryRowBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptStatus),
          ) as ReceiptStatus;
          result.status = valueDes;
          break;
        case r'receiptCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptCount = valueDes;
          break;
        case r'total':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.total = valueDes;
          break;
        case r'customFieldTotals':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReceiptSummaryCustomFieldTotal)]),
          ) as BuiltList<ReceiptSummaryCustomFieldTotal>;
          result.customFieldTotals.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptSummaryRow deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptSummaryRowBuilder();
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

