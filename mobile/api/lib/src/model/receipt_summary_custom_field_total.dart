//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary_custom_field_total.g.dart';

/// ReceiptSummaryCustomFieldTotal
///
/// Properties:
/// * [customFieldId] 
/// * [name] - The custom field's name, so the client needs no second lookup
/// * [total] - Sum of this field's currency values across the row's receipts
@BuiltValue()
abstract class ReceiptSummaryCustomFieldTotal implements Built<ReceiptSummaryCustomFieldTotal, ReceiptSummaryCustomFieldTotalBuilder> {
  @BuiltValueField(wireName: r'customFieldId')
  int get customFieldId;

  /// The custom field's name, so the client needs no second lookup
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Sum of this field's currency values across the row's receipts
  @BuiltValueField(wireName: r'total')
  String get total;

  ReceiptSummaryCustomFieldTotal._();

  factory ReceiptSummaryCustomFieldTotal([void updates(ReceiptSummaryCustomFieldTotalBuilder b)]) = _$ReceiptSummaryCustomFieldTotal;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptSummaryCustomFieldTotalBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptSummaryCustomFieldTotal> get serializer => _$ReceiptSummaryCustomFieldTotalSerializer();
}

class _$ReceiptSummaryCustomFieldTotalSerializer implements PrimitiveSerializer<ReceiptSummaryCustomFieldTotal> {
  @override
  final Iterable<Type> types = const [ReceiptSummaryCustomFieldTotal, _$ReceiptSummaryCustomFieldTotal];

  @override
  final String wireName = r'ReceiptSummaryCustomFieldTotal';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptSummaryCustomFieldTotal object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'customFieldId';
    yield serializers.serialize(
      object.customFieldId,
      specifiedType: const FullType(int),
    );
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'total';
    yield serializers.serialize(
      object.total,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptSummaryCustomFieldTotal object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptSummaryCustomFieldTotalBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'customFieldId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.customFieldId = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'total':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.total = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptSummaryCustomFieldTotal deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptSummaryCustomFieldTotalBuilder();
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

