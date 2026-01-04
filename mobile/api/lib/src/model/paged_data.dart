//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/paged_data_data_inner.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'paged_data.g.dart';

/// PagedData
///
/// Properties:
/// * [data] 
/// * [totalCount] 
@BuiltValue()
abstract class PagedData implements Built<PagedData, PagedDataBuilder> {
  @BuiltValueField(wireName: r'data')
  BuiltList<PagedDataDataInner> get data;

  @BuiltValueField(wireName: r'totalCount')
  int get totalCount;

  PagedData._();

  factory PagedData([void updates(PagedDataBuilder b)]) = _$PagedData;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PagedDataBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedData> get serializer => _$PagedDataSerializer();
}

class _$PagedDataSerializer implements PrimitiveSerializer<PagedData> {
  @override
  final Iterable<Type> types = const [PagedData, _$PagedData];

  @override
  final String wireName = r'PagedData';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedData object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'data';
    yield serializers.serialize(
      object.data,
      specifiedType: const FullType(BuiltList, [FullType(PagedDataDataInner)]),
    );
    yield r'totalCount';
    yield serializers.serialize(
      object.totalCount,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedData object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required PagedDataBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'data':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(PagedDataDataInner)]),
          ) as BuiltList<PagedDataDataInner>;
          result.data.replace(valueDes);
          break;
        case r'totalCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.totalCount = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  PagedData deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PagedDataBuilder();
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

