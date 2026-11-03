//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_history_snapshots.g.dart';

/// The receipt before and after one history entry, as JSON text, redacted to what the caller may see. A creation entry has no \"before\".
///
/// Properties:
/// * [before] 
/// * [after] 
@BuiltValue()
abstract class ReceiptHistorySnapshots implements Built<ReceiptHistorySnapshots, ReceiptHistorySnapshotsBuilder> {
  @BuiltValueField(wireName: r'before')
  String? get before;

  @BuiltValueField(wireName: r'after')
  String? get after;

  ReceiptHistorySnapshots._();

  factory ReceiptHistorySnapshots([void updates(ReceiptHistorySnapshotsBuilder b)]) = _$ReceiptHistorySnapshots;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptHistorySnapshotsBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptHistorySnapshots> get serializer => _$ReceiptHistorySnapshotsSerializer();
}

class _$ReceiptHistorySnapshotsSerializer implements PrimitiveSerializer<ReceiptHistorySnapshots> {
  @override
  final Iterable<Type> types = const [ReceiptHistorySnapshots, _$ReceiptHistorySnapshots];

  @override
  final String wireName = r'ReceiptHistorySnapshots';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptHistorySnapshots object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.before != null) {
      yield r'before';
      yield serializers.serialize(
        object.before,
        specifiedType: const FullType(String),
      );
    }
    if (object.after != null) {
      yield r'after';
      yield serializers.serialize(
        object.after,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptHistorySnapshots object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptHistorySnapshotsBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'before':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.before = valueDes;
          break;
        case r'after':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.after = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptHistorySnapshots deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptHistorySnapshotsBuilder();
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

