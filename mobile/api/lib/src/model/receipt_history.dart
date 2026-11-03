//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/receipt_history_entry.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_history.g.dart';

/// A receipt's change history, newest entry first
///
/// Properties:
/// * [entries] 
/// * [changeCount] - Change rows across the successful entries other than the creation
/// * [saveCount] - Successful entries
/// * [contributorCount] - Distinct actors across all entries
@BuiltValue()
abstract class ReceiptHistory implements Built<ReceiptHistory, ReceiptHistoryBuilder> {
  @BuiltValueField(wireName: r'entries')
  BuiltList<ReceiptHistoryEntry> get entries;

  /// Change rows across the successful entries other than the creation
  @BuiltValueField(wireName: r'changeCount')
  int get changeCount;

  /// Successful entries
  @BuiltValueField(wireName: r'saveCount')
  int get saveCount;

  /// Distinct actors across all entries
  @BuiltValueField(wireName: r'contributorCount')
  int get contributorCount;

  ReceiptHistory._();

  factory ReceiptHistory([void updates(ReceiptHistoryBuilder b)]) = _$ReceiptHistory;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptHistoryBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptHistory> get serializer => _$ReceiptHistorySerializer();
}

class _$ReceiptHistorySerializer implements PrimitiveSerializer<ReceiptHistory> {
  @override
  final Iterable<Type> types = const [ReceiptHistory, _$ReceiptHistory];

  @override
  final String wireName = r'ReceiptHistory';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptHistory object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'entries';
    yield serializers.serialize(
      object.entries,
      specifiedType: const FullType(BuiltList, [FullType(ReceiptHistoryEntry)]),
    );
    yield r'changeCount';
    yield serializers.serialize(
      object.changeCount,
      specifiedType: const FullType(int),
    );
    yield r'saveCount';
    yield serializers.serialize(
      object.saveCount,
      specifiedType: const FullType(int),
    );
    yield r'contributorCount';
    yield serializers.serialize(
      object.contributorCount,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptHistory object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptHistoryBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'entries':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReceiptHistoryEntry)]),
          ) as BuiltList<ReceiptHistoryEntry>;
          result.entries.replace(valueDes);
          break;
        case r'changeCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.changeCount = valueDes;
          break;
        case r'saveCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.saveCount = valueDes;
          break;
        case r'contributorCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.contributorCount = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptHistory deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptHistoryBuilder();
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

