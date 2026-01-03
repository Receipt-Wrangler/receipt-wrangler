//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/base_model.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'custom_field_value.g.dart';

/// CustomFieldValue
///
/// Properties:
/// * [id] 
/// * [createdAt] 
/// * [createdBy] 
/// * [createdByString] - Created by entity's name
/// * [updatedAt] 
/// * [receiptId] - Receipt Id
/// * [customFieldId] - Custom Field ID
/// * [stringValue] - Custom Field String Value
/// * [dateValue] - Custom Field Date Value
/// * [selectValue] - Custom Field Select Value
/// * [currencyValue] - Custom Field Currency Value
/// * [booleanValue] - Custom Field Boolean Value
@BuiltValue()
abstract class CustomFieldValue implements BaseModel, Built<CustomFieldValue, CustomFieldValueBuilder> {
  /// Custom Field Date Value
  @BuiltValueField(wireName: r'dateValue')
  String? get dateValue;

  /// Custom Field String Value
  @BuiltValueField(wireName: r'stringValue')
  String? get stringValue;

  /// Custom Field Select Value
  @BuiltValueField(wireName: r'selectValue')
  int? get selectValue;

  /// Custom Field Currency Value
  @BuiltValueField(wireName: r'currencyValue')
  String? get currencyValue;

  /// Custom Field ID
  @BuiltValueField(wireName: r'customFieldId')
  int get customFieldId;

  /// Custom Field Boolean Value
  @BuiltValueField(wireName: r'booleanValue')
  bool? get booleanValue;

  /// Receipt Id
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  CustomFieldValue._();

  factory CustomFieldValue([void updates(CustomFieldValueBuilder b)]) = _$CustomFieldValue;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CustomFieldValueBuilder b) => b
      ..createdBy = 0
      ..createdByString = ''
      ..updatedAt = '';

  @BuiltValueSerializer(custom: true)
  static Serializer<CustomFieldValue> get serializer => _$CustomFieldValueSerializer();
}

class _$CustomFieldValueSerializer implements PrimitiveSerializer<CustomFieldValue> {
  @override
  final Iterable<Type> types = const [CustomFieldValue, _$CustomFieldValue];

  @override
  final String wireName = r'CustomFieldValue';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CustomFieldValue object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.dateValue != null) {
      yield r'dateValue';
      yield serializers.serialize(
        object.dateValue,
        specifiedType: const FullType(String),
      );
    }
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    if (object.stringValue != null) {
      yield r'stringValue';
      yield serializers.serialize(
        object.stringValue,
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
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
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
    yield r'customFieldId';
    yield serializers.serialize(
      object.customFieldId,
      specifiedType: const FullType(int),
    );
    if (object.booleanValue != null) {
      yield r'booleanValue';
      yield serializers.serialize(
        object.booleanValue,
        specifiedType: const FullType(bool),
      );
    }
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'receiptId';
    yield serializers.serialize(
      object.receiptId,
      specifiedType: const FullType(int),
    );
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
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
    CustomFieldValue object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CustomFieldValueBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'dateValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.dateValue = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'stringValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.stringValue = valueDes;
          break;
        case r'selectValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.selectValue = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'currencyValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.currencyValue = valueDes;
          break;
        case r'customFieldId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.customFieldId = valueDes;
          break;
        case r'booleanValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.booleanValue = valueDes;
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
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
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
  CustomFieldValue deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CustomFieldValueBuilder();
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

