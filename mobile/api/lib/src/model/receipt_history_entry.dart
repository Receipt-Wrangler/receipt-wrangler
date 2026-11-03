//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_history_change.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/receipt_history_actor.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_history_entry.g.dart';

/// One recorded event in a receipt's history
///
/// Properties:
/// * [systemTaskId] - The system task the entry was built from
/// * [action] - What the entry did. One of CREATED, UPDATED, STATUS_CHANGED, RESOLVED, IMAGES_COMMENTS_UPDATED, FAILED. A plain string, not an enum, so a client must tolerate values it does not know.
/// * [recordedAt] - When the entry was recorded (RFC 3339 instant)
/// * [actor] 
/// * [hasSnapshots] - Whether getReceiptHistorySnapshots has receipt copies for this entry
/// * [errorMessage] - The error a FAILED entry recorded
/// * [beforeState] - How far the compared \"before\" can be trusted, on update entries. One of COMPLETE, REBUILT (an older row compared against an earlier complete copy) or INCOMPLETE (an older row with no such copy). Plain string.
/// * [beforeSourceType] - On a REBUILT entry, the system task type the \"before\" copy came from
/// * [beforeSourceRecordedAt] - On a REBUILT entry, when the \"before\" copy was recorded (RFC 3339 instant)
/// * [changes] 
@BuiltValue()
abstract class ReceiptHistoryEntry implements Built<ReceiptHistoryEntry, ReceiptHistoryEntryBuilder> {
  /// The system task the entry was built from
  @BuiltValueField(wireName: r'systemTaskId')
  int get systemTaskId;

  /// What the entry did. One of CREATED, UPDATED, STATUS_CHANGED, RESOLVED, IMAGES_COMMENTS_UPDATED, FAILED. A plain string, not an enum, so a client must tolerate values it does not know.
  @BuiltValueField(wireName: r'action')
  String get action;

  /// When the entry was recorded (RFC 3339 instant)
  @BuiltValueField(wireName: r'recordedAt')
  String get recordedAt;

  @BuiltValueField(wireName: r'actor')
  ReceiptHistoryActor get actor;

  /// Whether getReceiptHistorySnapshots has receipt copies for this entry
  @BuiltValueField(wireName: r'hasSnapshots')
  bool get hasSnapshots;

  /// The error a FAILED entry recorded
  @BuiltValueField(wireName: r'errorMessage')
  String? get errorMessage;

  /// How far the compared \"before\" can be trusted, on update entries. One of COMPLETE, REBUILT (an older row compared against an earlier complete copy) or INCOMPLETE (an older row with no such copy). Plain string.
  @BuiltValueField(wireName: r'beforeState')
  String? get beforeState;

  /// On a REBUILT entry, the system task type the \"before\" copy came from
  @BuiltValueField(wireName: r'beforeSourceType')
  String? get beforeSourceType;

  /// On a REBUILT entry, when the \"before\" copy was recorded (RFC 3339 instant)
  @BuiltValueField(wireName: r'beforeSourceRecordedAt')
  String? get beforeSourceRecordedAt;

  @BuiltValueField(wireName: r'changes')
  BuiltList<ReceiptHistoryChange> get changes;

  ReceiptHistoryEntry._();

  factory ReceiptHistoryEntry([void updates(ReceiptHistoryEntryBuilder b)]) = _$ReceiptHistoryEntry;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptHistoryEntryBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptHistoryEntry> get serializer => _$ReceiptHistoryEntrySerializer();
}

class _$ReceiptHistoryEntrySerializer implements PrimitiveSerializer<ReceiptHistoryEntry> {
  @override
  final Iterable<Type> types = const [ReceiptHistoryEntry, _$ReceiptHistoryEntry];

  @override
  final String wireName = r'ReceiptHistoryEntry';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptHistoryEntry object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'systemTaskId';
    yield serializers.serialize(
      object.systemTaskId,
      specifiedType: const FullType(int),
    );
    yield r'action';
    yield serializers.serialize(
      object.action,
      specifiedType: const FullType(String),
    );
    yield r'recordedAt';
    yield serializers.serialize(
      object.recordedAt,
      specifiedType: const FullType(String),
    );
    yield r'actor';
    yield serializers.serialize(
      object.actor,
      specifiedType: const FullType(ReceiptHistoryActor),
    );
    yield r'hasSnapshots';
    yield serializers.serialize(
      object.hasSnapshots,
      specifiedType: const FullType(bool),
    );
    if (object.errorMessage != null) {
      yield r'errorMessage';
      yield serializers.serialize(
        object.errorMessage,
        specifiedType: const FullType(String),
      );
    }
    if (object.beforeState != null) {
      yield r'beforeState';
      yield serializers.serialize(
        object.beforeState,
        specifiedType: const FullType(String),
      );
    }
    if (object.beforeSourceType != null) {
      yield r'beforeSourceType';
      yield serializers.serialize(
        object.beforeSourceType,
        specifiedType: const FullType(String),
      );
    }
    if (object.beforeSourceRecordedAt != null) {
      yield r'beforeSourceRecordedAt';
      yield serializers.serialize(
        object.beforeSourceRecordedAt,
        specifiedType: const FullType(String),
      );
    }
    yield r'changes';
    yield serializers.serialize(
      object.changes,
      specifiedType: const FullType(BuiltList, [FullType(ReceiptHistoryChange)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptHistoryEntry object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptHistoryEntryBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'systemTaskId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.systemTaskId = valueDes;
          break;
        case r'action':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.action = valueDes;
          break;
        case r'recordedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.recordedAt = valueDes;
          break;
        case r'actor':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptHistoryActor),
          ) as ReceiptHistoryActor;
          result.actor.replace(valueDes);
          break;
        case r'hasSnapshots':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hasSnapshots = valueDes;
          break;
        case r'errorMessage':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.errorMessage = valueDes;
          break;
        case r'beforeState':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.beforeState = valueDes;
          break;
        case r'beforeSourceType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.beforeSourceType = valueDes;
          break;
        case r'beforeSourceRecordedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.beforeSourceRecordedAt = valueDes;
          break;
        case r'changes':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReceiptHistoryChange)]),
          ) as BuiltList<ReceiptHistoryChange>;
          result.changes.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptHistoryEntry deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptHistoryEntryBuilder();
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

