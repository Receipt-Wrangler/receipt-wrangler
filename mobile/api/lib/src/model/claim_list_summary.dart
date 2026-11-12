//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_list_summary.g.dart';

/// Figures for an \"outstanding and ageing\" header, over the claims in the list
///
/// Properties:
/// * [asOfDate] - Today in the app time zone (yyyy-MM-dd), the day daysWaiting is counted to
/// * [claimCount] - Number of claims in the list
/// * [awaitingCount] - Number of claims with awaitingSettlement true
/// * [outstandingTotal] - Sum of outstandingAmount over the claims with awaitingSettlement true
/// * [oldestDaysWaiting] - The largest daysWaiting in the list; absent when no claim has one
@BuiltValue()
abstract class ClaimListSummary implements Built<ClaimListSummary, ClaimListSummaryBuilder> {
  /// Today in the app time zone (yyyy-MM-dd), the day daysWaiting is counted to
  @BuiltValueField(wireName: r'asOfDate')
  String get asOfDate;

  /// Number of claims in the list
  @BuiltValueField(wireName: r'claimCount')
  int get claimCount;

  /// Number of claims with awaitingSettlement true
  @BuiltValueField(wireName: r'awaitingCount')
  int get awaitingCount;

  /// Sum of outstandingAmount over the claims with awaitingSettlement true
  @BuiltValueField(wireName: r'outstandingTotal')
  String get outstandingTotal;

  /// The largest daysWaiting in the list; absent when no claim has one
  @BuiltValueField(wireName: r'oldestDaysWaiting')
  int? get oldestDaysWaiting;

  ClaimListSummary._();

  factory ClaimListSummary([void updates(ClaimListSummaryBuilder b)]) = _$ClaimListSummary;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimListSummaryBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimListSummary> get serializer => _$ClaimListSummarySerializer();
}

class _$ClaimListSummarySerializer implements PrimitiveSerializer<ClaimListSummary> {
  @override
  final Iterable<Type> types = const [ClaimListSummary, _$ClaimListSummary];

  @override
  final String wireName = r'ClaimListSummary';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimListSummary object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'asOfDate';
    yield serializers.serialize(
      object.asOfDate,
      specifiedType: const FullType(String),
    );
    yield r'claimCount';
    yield serializers.serialize(
      object.claimCount,
      specifiedType: const FullType(int),
    );
    yield r'awaitingCount';
    yield serializers.serialize(
      object.awaitingCount,
      specifiedType: const FullType(int),
    );
    yield r'outstandingTotal';
    yield serializers.serialize(
      object.outstandingTotal,
      specifiedType: const FullType(String),
    );
    if (object.oldestDaysWaiting != null) {
      yield r'oldestDaysWaiting';
      yield serializers.serialize(
        object.oldestDaysWaiting,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ClaimListSummary object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimListSummaryBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'asOfDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.asOfDate = valueDes;
          break;
        case r'claimCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.claimCount = valueDes;
          break;
        case r'awaitingCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.awaitingCount = valueDes;
          break;
        case r'outstandingTotal':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.outstandingTotal = valueDes;
          break;
        case r'oldestDaysWaiting':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.oldestDaysWaiting = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ClaimListSummary deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimListSummaryBuilder();
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

