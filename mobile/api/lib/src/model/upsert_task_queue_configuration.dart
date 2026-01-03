//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/queue_name.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_task_queue_configuration.g.dart';

/// UpsertTaskQueueConfiguration
///
/// Properties:
/// * [name] 
/// * [priority] - Queue priority
@BuiltValue()
abstract class UpsertTaskQueueConfiguration implements Built<UpsertTaskQueueConfiguration, UpsertTaskQueueConfigurationBuilder> {
  @BuiltValueField(wireName: r'name')
  QueueName? get name;
  // enum nameEnum {  quick_scan,  email_polling,  email_receipt_processing,  email_receipt_image_cleanup,  };

  /// Queue priority
  @BuiltValueField(wireName: r'priority')
  int? get priority;

  UpsertTaskQueueConfiguration._();

  factory UpsertTaskQueueConfiguration([void updates(UpsertTaskQueueConfigurationBuilder b)]) = _$UpsertTaskQueueConfiguration;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertTaskQueueConfigurationBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertTaskQueueConfiguration> get serializer => _$UpsertTaskQueueConfigurationSerializer();
}

class _$UpsertTaskQueueConfigurationSerializer implements PrimitiveSerializer<UpsertTaskQueueConfiguration> {
  @override
  final Iterable<Type> types = const [UpsertTaskQueueConfiguration, _$UpsertTaskQueueConfiguration];

  @override
  final String wireName = r'UpsertTaskQueueConfiguration';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertTaskQueueConfiguration object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(QueueName),
      );
    }
    if (object.priority != null) {
      yield r'priority';
      yield serializers.serialize(
        object.priority,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertTaskQueueConfiguration object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertTaskQueueConfigurationBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(QueueName),
          ) as QueueName;
          result.name = valueDes;
          break;
        case r'priority':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.priority = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertTaskQueueConfiguration deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertTaskQueueConfigurationBuilder();
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

