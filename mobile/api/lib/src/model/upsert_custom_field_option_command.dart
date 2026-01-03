//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_custom_field_option_command.g.dart';

/// UpsertCustomFieldOptionCommand
///
/// Properties:
/// * [value] - Custom Field Option value
/// * [customFieldId] - Custom Field Id
@BuiltValue()
abstract class UpsertCustomFieldOptionCommand implements Built<UpsertCustomFieldOptionCommand, UpsertCustomFieldOptionCommandBuilder> {
  /// Custom Field Option value
  @BuiltValueField(wireName: r'value')
  String? get value;

  /// Custom Field Id
  @BuiltValueField(wireName: r'customFieldId')
  int get customFieldId;

  UpsertCustomFieldOptionCommand._();

  factory UpsertCustomFieldOptionCommand([void updates(UpsertCustomFieldOptionCommandBuilder b)]) = _$UpsertCustomFieldOptionCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertCustomFieldOptionCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertCustomFieldOptionCommand> get serializer => _$UpsertCustomFieldOptionCommandSerializer();
}

class _$UpsertCustomFieldOptionCommandSerializer implements PrimitiveSerializer<UpsertCustomFieldOptionCommand> {
  @override
  final Iterable<Type> types = const [UpsertCustomFieldOptionCommand, _$UpsertCustomFieldOptionCommand];

  @override
  final String wireName = r'UpsertCustomFieldOptionCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertCustomFieldOptionCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.value != null) {
      yield r'value';
      yield serializers.serialize(
        object.value,
        specifiedType: const FullType(String),
      );
    }
    yield r'customFieldId';
    yield serializers.serialize(
      object.customFieldId,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertCustomFieldOptionCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertCustomFieldOptionCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'value':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.value = valueDes;
          break;
        case r'customFieldId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.customFieldId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertCustomFieldOptionCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertCustomFieldOptionCommandBuilder();
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

