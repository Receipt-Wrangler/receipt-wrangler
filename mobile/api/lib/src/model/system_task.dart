//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/base_model.dart';
import 'package:openapi/src/model/associated_entity_type.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/system_task_type.dart';
import 'package:openapi/src/model/system_task_status.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'system_task.g.dart';

/// SystemTask
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
/// * [type] 
/// * [status] 
/// * [startedAt] 
/// * [endedAt] 
/// * [associatedEntityId] 
/// * [associatedEntityType] 
/// * [ranByUserId] 
/// * [receiptId] 
/// * [groupId] 
/// * [resultDescription] 
/// * [childSystemTasks] 
@BuiltValue()
abstract class SystemTask implements BaseModel, Built<SystemTask, SystemTaskBuilder> {
  @BuiltValueField(wireName: r'associatedEntityId')
  int? get associatedEntityId;

  @BuiltValueField(wireName: r'endedAt')
  String? get endedAt;

  @BuiltValueField(wireName: r'associatedEntityType')
  AssociatedEntityType? get associatedEntityType;
  // enum associatedEntityTypeEnum {  NOOP_ENTITY_TYPE,  RECEIPT,  SYSTEM_EMAIL,  RECEIPT_PROCESSING_SETTINGS,  PROMPT,  };

  @BuiltValueField(wireName: r'groupId')
  int? get groupId;

  @BuiltValueField(wireName: r'childSystemTasks')
  BuiltList<SystemTask>? get childSystemTasks;

  @BuiltValueField(wireName: r'startedAt')
  String? get startedAt;

  @BuiltValueField(wireName: r'resultDescription')
  String? get resultDescription;

  @BuiltValueField(wireName: r'type')
  SystemTaskType? get type;
  // enum typeEnum {  OCR_PROCESSING,  CHAT_COMPLETION,  MAGIC_FILL,  QUICK_SCAN,  EMAIL_READ,  EMAIL_UPLOAD,  SYSTEM_EMAIL_CONNECTIVITY_CHECK,  RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK,  RECEIPT_UPLOADED,  RECEIPT_UPDATED,  PROMPT_GENERATED,  };

  @BuiltValueField(wireName: r'receiptId')
  int? get receiptId;

  @BuiltValueField(wireName: r'ranByUserId')
  int? get ranByUserId;

  @BuiltValueField(wireName: r'status')
  SystemTaskStatus? get status;
  // enum statusEnum {  SUCCEEDED,  FAILED,  };

  SystemTask._();

  factory SystemTask([void updates(SystemTaskBuilder b)]) = _$SystemTask;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(SystemTaskBuilder b) => b
      ..createdBy = 0
      ..createdByString = ''
      ..updatedAt = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<SystemTask> get serializer => _$SystemTaskSerializer();
}

class _$SystemTaskSerializer implements PrimitiveSerializer<SystemTask> {
  @override
  final Iterable<Type> types = const [SystemTask, _$SystemTask];

  @override
  final String wireName = r'SystemTask';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    SystemTask object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.groupId != null) {
      yield r'groupId';
      yield serializers.serialize(
        object.groupId,
        specifiedType: const FullType(int),
      );
    }
    if (object.startedAt != null) {
      yield r'startedAt';
      yield serializers.serialize(
        object.startedAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.type != null) {
      yield r'type';
      yield serializers.serialize(
        object.type,
        specifiedType: const FullType(SystemTaskType),
      );
    }
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    if (object.associatedEntityId != null) {
      yield r'associatedEntityId';
      yield serializers.serialize(
        object.associatedEntityId,
        specifiedType: const FullType(int),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    if (object.endedAt != null) {
      yield r'endedAt';
      yield serializers.serialize(
        object.endedAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.associatedEntityType != null) {
      yield r'associatedEntityType';
      yield serializers.serialize(
        object.associatedEntityType,
        specifiedType: const FullType(AssociatedEntityType),
      );
    }
    if (object.childSystemTasks != null) {
      yield r'childSystemTasks';
      yield serializers.serialize(
        object.childSystemTasks,
        specifiedType: const FullType(BuiltList, [FullType(SystemTask)]),
      );
    }
    if (object.resultDescription != null) {
      yield r'resultDescription';
      yield serializers.serialize(
        object.resultDescription,
        specifiedType: const FullType(String),
      );
    }
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.receiptId != null) {
      yield r'receiptId';
      yield serializers.serialize(
        object.receiptId,
        specifiedType: const FullType(int),
      );
    }
    if (object.ranByUserId != null) {
      yield r'ranByUserId';
      yield serializers.serialize(
        object.ranByUserId,
        specifiedType: const FullType(int),
      );
    }
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
      );
    }
    if (object.status != null) {
      yield r'status';
      yield serializers.serialize(
        object.status,
        specifiedType: const FullType(SystemTaskStatus),
      );
    }
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    SystemTask object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required SystemTaskBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'startedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.startedAt = valueDes;
          break;
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(SystemTaskType),
          ) as SystemTaskType;
          result.type = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'associatedEntityId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.associatedEntityId = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'endedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.endedAt = valueDes;
          break;
        case r'associatedEntityType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(AssociatedEntityType),
          ) as AssociatedEntityType;
          result.associatedEntityType = valueDes;
          break;
        case r'childSystemTasks':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(SystemTask)]),
          ) as BuiltList<SystemTask>;
          result.childSystemTasks.replace(valueDes);
          break;
        case r'resultDescription':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.resultDescription = valueDes;
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'ranByUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.ranByUserId = valueDes;
          break;
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(SystemTaskStatus),
          ) as SystemTaskStatus;
          result.status = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  SystemTask deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = SystemTaskBuilder();
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

