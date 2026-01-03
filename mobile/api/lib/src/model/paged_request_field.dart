//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/paged_request_field_value.dart';
import 'package:openapi/src/model/filter_operation.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'paged_request_field.g.dart';

/// PagedRequestField
///
/// Properties:
/// * [operation] - Filter operation
/// * [value] 
@BuiltValue()
abstract class PagedRequestField implements Built<PagedRequestField, PagedRequestFieldBuilder> {
  /// Filter operation
  @BuiltValueField(wireName: r'operation')
  FilterOperation? get operation;
  // enum operationEnum {  CONTAINS,  EQUALS,  GREATER_THAN,  LESS_THAN,  ,  };

  @BuiltValueField(wireName: r'value')
  PagedRequestFieldValue? get value;

  PagedRequestField._();

  factory PagedRequestField([void updates(PagedRequestFieldBuilder b)]) = _$PagedRequestField;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PagedRequestFieldBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedRequestField> get serializer => _$PagedRequestFieldSerializer();
}

class _$PagedRequestFieldSerializer implements PrimitiveSerializer<PagedRequestField> {
  @override
  final Iterable<Type> types = const [PagedRequestField, _$PagedRequestField];

  @override
  final String wireName = r'PagedRequestField';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedRequestField object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.operation != null) {
      yield r'operation';
      yield serializers.serialize(
        object.operation,
        specifiedType: const FullType(FilterOperation),
      );
    }
    if (object.value != null) {
      yield r'value';
      yield serializers.serialize(
        object.value,
        specifiedType: const FullType(PagedRequestFieldValue),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedRequestField object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PagedRequestFieldBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'operation':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(FilterOperation),
          ) as FilterOperation;
          result.operation = valueDes;
          break;
        case r'value':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(PagedRequestFieldValue),
          ) as PagedRequestFieldValue;
          result.value.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PagedRequestField deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PagedRequestFieldBuilder();
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

