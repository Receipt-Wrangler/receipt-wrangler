//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_history_actor.g.dart';

/// Who or what made a history entry
///
/// Properties:
/// * [type] - One of USER, QUICK_SCAN, EMAIL, SYSTEM. Plain string.
/// * [userId] - The user, for USER; the scanning user (when visible), for QUICK_SCAN
@BuiltValue()
abstract class ReceiptHistoryActor implements Built<ReceiptHistoryActor, ReceiptHistoryActorBuilder> {
  /// One of USER, QUICK_SCAN, EMAIL, SYSTEM. Plain string.
  @BuiltValueField(wireName: r'type')
  String get type;

  /// The user, for USER; the scanning user (when visible), for QUICK_SCAN
  @BuiltValueField(wireName: r'userId')
  int? get userId;

  ReceiptHistoryActor._();

  factory ReceiptHistoryActor([void updates(ReceiptHistoryActorBuilder b)]) = _$ReceiptHistoryActor;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptHistoryActorBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptHistoryActor> get serializer => _$ReceiptHistoryActorSerializer();
}

class _$ReceiptHistoryActorSerializer implements PrimitiveSerializer<ReceiptHistoryActor> {
  @override
  final Iterable<Type> types = const [ReceiptHistoryActor, _$ReceiptHistoryActor];

  @override
  final String wireName = r'ReceiptHistoryActor';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptHistoryActor object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'type';
    yield serializers.serialize(
      object.type,
      specifiedType: const FullType(String),
    );
    if (object.userId != null) {
      yield r'userId';
      yield serializers.serialize(
        object.userId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptHistoryActor object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptHistoryActorBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.type = valueDes;
          break;
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptHistoryActor deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptHistoryActorBuilder();
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

