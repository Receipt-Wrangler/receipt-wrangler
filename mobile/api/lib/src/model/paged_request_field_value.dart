//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'dart:core';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';
import 'package:one_of/one_of.dart';

part 'paged_request_field_value.g.dart';

/// Field value
@BuiltValue()
abstract class PagedRequestFieldValue implements Built<PagedRequestFieldValue, PagedRequestFieldValueBuilder> {
  /// One Of [BuiltList<String>], [BuiltList<int>], [String], [int]
  OneOf get oneOf;

  PagedRequestFieldValue._();

  factory PagedRequestFieldValue([void updates(PagedRequestFieldValueBuilder b)]) = _$PagedRequestFieldValue;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PagedRequestFieldValueBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedRequestFieldValue> get serializer => _$PagedRequestFieldValueSerializer();
}

class _$PagedRequestFieldValueSerializer implements PrimitiveSerializer<PagedRequestFieldValue> {
  @override
  final Iterable<Type> types = const [PagedRequestFieldValue, _$PagedRequestFieldValue];

  @override
  final String wireName = r'PagedRequestFieldValue';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedRequestFieldValue object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedRequestFieldValue object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final oneOf = object.oneOf;
    return serializers.serialize(oneOf.value, specifiedType: FullType(oneOf.valueType))!;
  }

  @override
  PagedRequestFieldValue deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PagedRequestFieldValueBuilder();
    Object? oneOfDataSrc;
    final targetType = const FullType(OneOf, [FullType(String), FullType(int), FullType(BuiltList, [FullType(String)]), FullType(BuiltList, [FullType(int)]), ]);
    oneOfDataSrc = serialized;
    result.oneOf = serializers.deserialize(oneOfDataSrc, specifiedType: targetType) as OneOf;
    return result.build();
  }
}

