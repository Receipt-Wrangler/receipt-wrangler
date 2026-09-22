//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/receipt_summary_row.dart';
import 'package:openapi/src/model/receipt_summary_position.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary.g.dart';

/// ReceiptSummary
///
/// Properties:
/// * [enabled] - Whether the configuration group has the summary turned on. False comes back at a normal 200 with zeroed rows, so a client with stale group settings renders nothing rather than surfacing an error.
/// * [configurationGroupId] - The group whose settings produced this breakdown
/// * [position] 
/// * [overall] 
/// * [statuses] - One row per configured status, in ReceiptStatus declaration order. A configured status matching no receipt is still present, with zeroed figures. Always present; empty when no status is configured.
@BuiltValue()
abstract class ReceiptSummary implements Built<ReceiptSummary, ReceiptSummaryBuilder> {
  /// Whether the configuration group has the summary turned on. False comes back at a normal 200 with zeroed rows, so a client with stale group settings renders nothing rather than surfacing an error.
  @BuiltValueField(wireName: r'enabled')
  bool get enabled;

  /// The group whose settings produced this breakdown
  @BuiltValueField(wireName: r'configurationGroupId')
  int get configurationGroupId;

  @BuiltValueField(wireName: r'position')
  ReceiptSummaryPosition get position;
  // enum positionEnum {  TOP,  BOTTOM,  };

  @BuiltValueField(wireName: r'overall')
  ReceiptSummaryRow get overall;

  /// One row per configured status, in ReceiptStatus declaration order. A configured status matching no receipt is still present, with zeroed figures. Always present; empty when no status is configured.
  @BuiltValueField(wireName: r'statuses')
  BuiltList<ReceiptSummaryRow> get statuses;

  ReceiptSummary._();

  factory ReceiptSummary([void updates(ReceiptSummaryBuilder b)]) = _$ReceiptSummary;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptSummaryBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptSummary> get serializer => _$ReceiptSummarySerializer();
}

class _$ReceiptSummarySerializer implements PrimitiveSerializer<ReceiptSummary> {
  @override
  final Iterable<Type> types = const [ReceiptSummary, _$ReceiptSummary];

  @override
  final String wireName = r'ReceiptSummary';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptSummary object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'enabled';
    yield serializers.serialize(
      object.enabled,
      specifiedType: const FullType(bool),
    );
    yield r'configurationGroupId';
    yield serializers.serialize(
      object.configurationGroupId,
      specifiedType: const FullType(int),
    );
    yield r'position';
    yield serializers.serialize(
      object.position,
      specifiedType: const FullType(ReceiptSummaryPosition),
    );
    yield r'overall';
    yield serializers.serialize(
      object.overall,
      specifiedType: const FullType(ReceiptSummaryRow),
    );
    yield r'statuses';
    yield serializers.serialize(
      object.statuses,
      specifiedType: const FullType(BuiltList, [FullType(ReceiptSummaryRow)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptSummary object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptSummaryBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'enabled':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.enabled = valueDes;
          break;
        case r'configurationGroupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.configurationGroupId = valueDes;
          break;
        case r'position':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptSummaryPosition),
          ) as ReceiptSummaryPosition;
          result.position = valueDes;
          break;
        case r'overall':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptSummaryRow),
          ) as ReceiptSummaryRow;
          result.overall.replace(valueDes);
          break;
        case r'statuses':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReceiptSummaryRow)]),
          ) as BuiltList<ReceiptSummaryRow>;
          result.statuses.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptSummary deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptSummaryBuilder();
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

