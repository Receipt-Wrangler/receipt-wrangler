//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_custom_field_value_command.g.dart';

/// UpsertCustomFieldValueCommand
///
/// Properties:
/// * [receiptId] - Receipt Id
/// * [customFieldId] - Custom Field ID
/// * [stringValue] - Custom Field String Value
/// * [dateValue] - Custom Field Date Value
/// * [selectValue] - Custom Field Select Value
/// * [currencyValue] - Custom Field Currency Value
/// * [booleanValue] - Custom Field Boolean Value
@BuiltValue()
abstract class UpsertCustomFieldValueCommand implements Built<UpsertCustomFieldValueCommand, UpsertCustomFieldValueCommandBuilder> {
  /// Receipt Id
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  /// Custom Field ID
  @BuiltValueField(wireName: r'customFieldId')
  int get customFieldId;

  /// Custom Field String Value
  @BuiltValueField(wireName: r'stringValue')
  String? get stringValue;

  /// Custom Field Date Value
  @BuiltValueField(wireName: r'dateValue')
  String? get dateValue;

  /// Custom Field Select Value
  @BuiltValueField(wireName: r'selectValue')
  int? get selectValue;

  /// Custom Field Currency Value
  @BuiltValueField(wireName: r'currencyValue')
  String? get currencyValue;

  /// Custom Field Boolean Value
  @BuiltValueField(wireName: r'booleanValue')
  bool? get booleanValue;

  UpsertCustomFieldValueCommand._();

  factory UpsertCustomFieldValueCommand([void updates(UpsertCustomFieldValueCommandBuilder b)]) = _$UpsertCustomFieldValueCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertCustomFieldValueCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertCustomFieldValueCommand> get serializer => _$UpsertCustomFieldValueCommandSerializer();
}

class _$UpsertCustomFieldValueCommandSerializer implements PrimitiveSerializer<UpsertCustomFieldValueCommand> {
  @override
  final Iterable<Type> types = const [UpsertCustomFieldValueCommand, _$UpsertCustomFieldValueCommand];

  @override
  final String wireName = r'UpsertCustomFieldValueCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertCustomFieldValueCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'receiptId';
    yield serializers.serialize(
      object.receiptId,
      specifiedType: const FullType(int),
    );
    yield r'customFieldId';
    yield serializers.serialize(
      object.customFieldId,
      specifiedType: const FullType(int),
    );
    if (object.stringValue != null) {
      yield r'stringValue';
      yield serializers.serialize(
        object.stringValue,
        specifiedType: const FullType(String),
      );
    }
    if (object.dateValue != null) {
      yield r'dateValue';
      yield serializers.serialize(
        object.dateValue,
        specifiedType: const FullType(String),
      );
    }
    if (object.selectValue != null) {
      yield r'selectValue';
      yield serializers.serialize(
        object.selectValue,
        specifiedType: const FullType(int),
      );
    }
    if (object.currencyValue != null) {
      yield r'currencyValue';
      yield serializers.serialize(
        object.currencyValue,
        specifiedType: const FullType(String),
      );
    }
    if (object.booleanValue != null) {
      yield r'booleanValue';
      yield serializers.serialize(
        object.booleanValue,
        specifiedType: const FullType(bool),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertCustomFieldValueCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertCustomFieldValueCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'customFieldId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.customFieldId = valueDes;
          break;
        case r'stringValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.stringValue = valueDes;
          break;
        case r'dateValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.dateValue = valueDes;
          break;
        case r'selectValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.selectValue = valueDes;
          break;
        case r'currencyValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.currencyValue = valueDes;
          break;
        case r'booleanValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.booleanValue = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertCustomFieldValueCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertCustomFieldValueCommandBuilder();
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

