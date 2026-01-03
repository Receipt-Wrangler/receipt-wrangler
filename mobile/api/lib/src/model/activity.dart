//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/system_task_type.dart';
import 'package:openapi/src/model/system_task_status.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'activity.g.dart';

/// Activity
///
/// Properties:
/// * [id] 
/// * [type] 
/// * [status] 
/// * [startedAt] 
/// * [endedAt] 
/// * [ranByUserId] 
/// * [receiptId] 
/// * [groupId] 
/// * [canBeRestarted] 
@BuiltValue()
abstract class Activity implements Built<Activity, ActivityBuilder> {
  @BuiltValueField(wireName: r'id')
  int get id;

  @BuiltValueField(wireName: r'type')
  SystemTaskType get type;
  // enum typeEnum {  OCR_PROCESSING,  CHAT_COMPLETION,  MAGIC_FILL,  QUICK_SCAN,  EMAIL_READ,  EMAIL_UPLOAD,  SYSTEM_EMAIL_CONNECTIVITY_CHECK,  RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK,  RECEIPT_UPLOADED,  RECEIPT_UPDATED,  PROMPT_GENERATED,  };

  @BuiltValueField(wireName: r'status')
  SystemTaskStatus get status;
  // enum statusEnum {  SUCCEEDED,  FAILED,  };

  @BuiltValueField(wireName: r'startedAt')
  String get startedAt;

  @BuiltValueField(wireName: r'endedAt')
  String get endedAt;

  @BuiltValueField(wireName: r'ranByUserId')
  int? get ranByUserId;

  @BuiltValueField(wireName: r'receiptId')
  int? get receiptId;

  @BuiltValueField(wireName: r'groupId')
  int? get groupId;

  @BuiltValueField(wireName: r'canBeRestarted')
  bool? get canBeRestarted;

  Activity._();

  factory Activity([void updates(ActivityBuilder b)]) = _$Activity;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ActivityBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Activity> get serializer => _$ActivitySerializer();
}

class _$ActivitySerializer implements PrimitiveSerializer<Activity> {
  @override
  final Iterable<Type> types = const [Activity, _$Activity];

  @override
  final String wireName = r'Activity';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Activity object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'type';
    yield serializers.serialize(
      object.type,
      specifiedType: const FullType(SystemTaskType),
    );
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(SystemTaskStatus),
    );
    yield r'startedAt';
    yield serializers.serialize(
      object.startedAt,
      specifiedType: const FullType(String),
    );
    yield r'endedAt';
    yield serializers.serialize(
      object.endedAt,
      specifiedType: const FullType(String),
    );
    if (object.ranByUserId != null) {
      yield r'ranByUserId';
      yield serializers.serialize(
        object.ranByUserId,
        specifiedType: const FullType(int),
      );
    }
    if (object.receiptId != null) {
      yield r'receiptId';
      yield serializers.serialize(
        object.receiptId,
        specifiedType: const FullType(int),
      );
    }
    if (object.groupId != null) {
      yield r'groupId';
      yield serializers.serialize(
        object.groupId,
        specifiedType: const FullType(int),
      );
    }
    if (object.canBeRestarted != null) {
      yield r'canBeRestarted';
      yield serializers.serialize(
        object.canBeRestarted,
        specifiedType: const FullType(bool),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    Activity object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ActivityBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(SystemTaskType),
          ) as SystemTaskType;
          result.type = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(SystemTaskStatus),
          ) as SystemTaskStatus;
          result.status = valueDes;
          break;
        case r'startedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.startedAt = valueDes;
          break;
        case r'endedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.endedAt = valueDes;
          break;
        case r'ranByUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.ranByUserId = valueDes;
          break;
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'canBeRestarted':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.canBeRestarted = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Activity deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ActivityBuilder();
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

